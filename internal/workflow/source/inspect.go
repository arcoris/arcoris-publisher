// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"context"
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	portfs "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
)

// inspector owns the mutable state for one source inspection run.
type inspector struct {
	// deps contains the infrastructure ports used by this run.
	deps Dependencies

	// opts contains behavior toggles resolved by Service.New.
	opts Options

	// request is the immutable caller input for this run.
	request Request

	// repositoryDir is the cleaned absolute repository root.
	repositoryDir string

	// stagingDir is the cleaned absolute staging root inside repositoryDir.
	stagingDir string

	// warnings accumulates non-fatal diagnostics such as warn-mode dirty source.
	warnings issueCollector
}

// moduleSourcePlan is the small ModulePlan surface required by source
// inspection.
type moduleSourcePlan interface {
	Name() manifest.ModuleName
	SourceDir() manifest.SourceDir
	ModuleRoot() manifest.RelativePath
	PublishEntries() []manifest.PublishEntry
}

// inspect validates roots, captures Git state, and inspects every planned
// module in publication order.
func (i *inspector) inspect(ctx context.Context) (Snapshot, error) {
	if err := (Service{deps: i.deps, opts: i.opts}).validateDependencies(); err != nil {
		return Snapshot{}, err
	}

	if err := i.validateRequest(ctx); err != nil {
		return Snapshot{}, err
	}

	repository, err := i.inspectRepository(ctx)
	if err != nil {
		return Snapshot{}, err
	}

	modules, err := i.inspectModules(ctx)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		repository: repository,
		modules:    modules,
		warnings:   i.warnings.Issues(),
	}, nil
}

// validateRequest normalizes user-provided roots and checks the filesystem
// invariants needed before any Git or module inspection can be trusted.
func (i *inspector) validateRequest(ctx context.Context) error {
	issues := newIssueCollector()

	i.validatePlan(&issues)
	repositoryDir := i.cleanRequestPath(&issues, "repositoryDir", i.request.RepositoryDir)
	stagingDir := i.cleanRequestPath(&issues, "stagingDir", i.request.StagingDir)

	i.validateRootDir(
		ctx,
		&issues,
		repositoryDir,
		IssueRepositoryMissing,
		IssueRepositoryNotDirectory,
		"repositoryDir",
		"source repository",
	)
	i.validateRootDir(
		ctx,
		&issues,
		stagingDir,
		IssueStagingMissing,
		IssueStagingNotDirectory,
		"stagingDir",
		"staging root",
	)
	i.validateStagingRoot(&issues, repositoryDir, stagingDir)

	if err := issues.Err(); err != nil {
		return err
	}

	i.repositoryDir = repositoryDir
	i.stagingDir = stagingDir

	return nil
}

// validatePlan rejects empty plans before filesystem checks report confusing
// downstream module errors.
func (i *inspector) validatePlan(issues *issueCollector) {
	if !i.request.Plan.Empty() {
		return
	}

	issues.Add(
		IssueInvalidRequest,
		"",
		"plan",
		"plan must contain at least one module",
	)
}

// cleanRequestPath converts a request path into a cleaned absolute path.
func (i *inspector) cleanRequestPath(
	issues *issueCollector,
	path string,
	value string,
) string {
	cleaned, err := pathutil.CleanAbs(value)
	if err == nil {
		return cleaned
	}

	issues.AddMessage(IssueInvalidRequest, "", path, err.Error())
	return ""
}

// validateRootDir checks one required root directory when its path was parsed.
func (i *inspector) validateRootDir(
	ctx context.Context,
	issues *issueCollector,
	dir string,
	missing IssueCode,
	notDir IssueCode,
	path string,
	label string,
) {
	if dir == "" {
		return
	}

	if err := i.checkDir(ctx, dir, missing, notDir, path, label); err != nil {
		issues.Append(validationIssues(err))
	}
}

// validateStagingRoot ensures all planned module SourceDir values stay rooted
// under the inspected repository checkout.
func (i *inspector) validateStagingRoot(
	issues *issueCollector,
	repositoryDir string,
	stagingDir string,
) {
	if repositoryDir == "" || stagingDir == "" {
		return
	}

	if err := pathutil.EnsureInside(repositoryDir, stagingDir); err != nil {
		issues.Add(
			IssueStagingOutsideRepo,
			"",
			"stagingDir",
			"staging root must be inside repository: %v",
			err,
		)
	}
}

// checkDir validates existence and directory type for one filesystem path.
func (i *inspector) checkDir(
	ctx context.Context,
	dir string,
	missing IssueCode,
	notDir IssueCode,
	path string,
	label string,
) error {
	exists, err := i.deps.FS.Exists(ctx, dir)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if !exists {
		return singleIssueError(missing, "", path, label+" does not exist")
	}

	isDir, err := i.deps.FS.IsDir(ctx, dir)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if !isDir {
		return singleIssueError(notDir, "", path, label+" is not a directory")
	}

	return nil
}

// validationIssues converts nested validation errors into detached issues.
func validationIssues(err error) []Issue {
	if ve, ok := err.(*ValidationError); ok {
		return cloneIssues(ve.Issues)
	}

	return []Issue{{
		Code:    IssueInvalidRequest,
		Message: err.Error(),
	}}
}

// inspectRepository captures Git provenance and enforces dirty-check policy.
func (i *inspector) inspectRepository(ctx context.Context) (RepositorySnapshot, error) {
	repository, err := i.repositorySnapshot(ctx)
	if err != nil {
		return RepositorySnapshot{}, err
	}

	if err := i.enforceDetachedHeadPolicy(repository); err != nil {
		return RepositorySnapshot{}, err
	}

	if err := i.enforceDirtyPolicy(repository); err != nil {
		return RepositorySnapshot{}, err
	}

	return repository, nil
}

// repositorySnapshot reads the Git state used by later provenance and policy
// checks.
func (i *inspector) repositorySnapshot(ctx context.Context) (RepositorySnapshot, error) {
	head, err := i.deps.Git.Head(ctx, i.repositoryDir)
	if err != nil {
		return RepositorySnapshot{}, fmt.Errorf("git head: %w", err)
	}

	branch, err := i.deps.Git.CurrentBranch(ctx, i.repositoryDir)
	if err != nil {
		return RepositorySnapshot{}, fmt.Errorf("git current branch: %w", err)
	}

	status, err := i.deps.Git.Status(ctx, i.repositoryDir)
	if err != nil {
		return RepositorySnapshot{}, fmt.Errorf("git status: %w", err)
	}

	return RepositorySnapshot{
		repositoryDir: i.repositoryDir,
		stagingDir:    i.stagingDir,
		head:          head,
		branch:        branch,
		status:        cloneStatus(status),
	}, nil
}

// enforceDetachedHeadPolicy rejects detached source checkouts unless the caller
// explicitly opted into weaker branch provenance.
func (i *inspector) enforceDetachedHeadPolicy(repository RepositorySnapshot) error {
	if repository.Branch() != "" || i.opts.AllowDetachedHEAD {
		return nil
	}

	return singleIssueError(
		IssueDetachedHead,
		"",
		"git.branch",
		"source checkout is in detached HEAD state",
	)
}

// enforceDirtyPolicy applies the resolved manifest dirty-source policy.
func (i *inspector) enforceDirtyPolicy(repository RepositorySnapshot) error {
	if !repository.Dirty() {
		return nil
	}

	switch i.request.Plan.Source().DirtyPolicy() {
	case manifest.DirtyPolicyFail:
		return singleIssueError(
			IssueDirtySource,
			"",
			"git.status",
			"source checkout is dirty",
		)
	case manifest.DirtyPolicyWarn:
		i.warnings.Add(
			IssueDirtySource,
			"",
			"git.status",
			"source checkout is dirty",
		)
	case manifest.DirtyPolicyAllow:
	default:
		return singleIssueError(
			IssueInvalidRequest,
			"",
			"source.dirtyPolicy",
			"unsupported dirty policy",
		)
	}

	return nil
}

// inspectModules inspects all planned modules and reports every module failure
// together instead of stopping at the first bad module.
func (i *inspector) inspectModules(ctx context.Context) ([]ModuleSnapshot, error) {
	modules := i.request.Plan.Modules()
	out := make([]ModuleSnapshot, 0, len(modules))
	issues := newIssueCollector()

	for index, modulePlan := range modules {
		module, moduleIssues := i.inspectModule(ctx, index, modulePlan)
		if len(moduleIssues) > 0 {
			issues.Append(moduleIssues)
			continue
		}

		out = append(out, module)
	}

	if err := issues.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// inspectModule validates module roots and then inspects every explicit publish
// entry declared by the plan.
func (i *inspector) inspectModule(
	ctx context.Context,
	index int,
	modulePlan moduleSourcePlan,
) (ModuleSnapshot, []Issue) {
	name := modulePlan.Name()
	basePath := fmt.Sprintf("modules[%d]", index)
	moduleDir := resolveModuleDir(i.stagingDir, modulePlan.SourceDir())
	moduleRootDir := resolveModuleRootDir(moduleDir, modulePlan.ModuleRoot())

	issues := newIssueCollector()
	i.validateModuleRoots(ctx, &issues, name, basePath, moduleDir, moduleRootDir)
	if issues.Len() > 0 {
		return ModuleSnapshot{}, issues.Issues()
	}

	entries, entryHashes := i.inspectEntries(
		ctx,
		&issues,
		name,
		basePath,
		moduleRootDir,
		modulePlan.PublishEntries(),
	)
	if issues.Len() > 0 {
		return ModuleSnapshot{}, issues.Issues()
	}

	return ModuleSnapshot{
		name:          name,
		sourceDir:     moduleDir,
		moduleRootDir: moduleRootDir,
		entries:       entries,
		hash:          combineHashes("module", entryHashes),
	}, nil
}

// validateModuleRoots checks both the module source directory and the module
// root used to resolve publish entries.
func (i *inspector) validateModuleRoots(
	ctx context.Context,
	issues *issueCollector,
	name manifest.ModuleName,
	basePath string,
	moduleDir string,
	moduleRootDir string,
) {
	if err := pathutil.EnsureInside(i.stagingDir, moduleDir); err != nil {
		issues.AddMessage(IssueEntryPathEscape, name, basePath+".sourceDir", err.Error())
	}

	if err := pathutil.EnsureInside(moduleDir, moduleRootDir); err != nil {
		issues.AddMessage(IssueEntryPathEscape, name, basePath+".moduleRoot", err.Error())
	}

	i.validateModuleDir(
		ctx,
		issues,
		name,
		moduleDir,
		IssueModuleSourceMissing,
		IssueModuleSourceNotDir,
		basePath+".sourceDir",
		"module source directory",
	)
	i.validateModuleDir(
		ctx,
		issues,
		name,
		moduleRootDir,
		IssueModuleRootMissing,
		IssueModuleRootNotDir,
		basePath+".moduleRoot",
		"module root directory",
	)
}

// validateModuleDir checks one module-owned directory and annotates nested
// issues with the module name.
func (i *inspector) validateModuleDir(
	ctx context.Context,
	issues *issueCollector,
	name manifest.ModuleName,
	dir string,
	missing IssueCode,
	notDir IssueCode,
	path string,
	label string,
) {
	if err := i.checkDir(ctx, dir, missing, notDir, path, label); err != nil {
		issues.Append(withModule(validationIssues(err), name))
	}
}

// inspectEntries inspects explicit publish entries and gathers their content
// hashes for the module hash.
func (i *inspector) inspectEntries(
	ctx context.Context,
	issues *issueCollector,
	name manifest.ModuleName,
	basePath string,
	moduleRootDir string,
	entries []manifest.PublishEntry,
) ([]EntrySnapshot, []Hash) {
	entrySnapshots := make([]EntrySnapshot, 0, len(entries))
	entryHashes := make([]Hash, 0, len(entries))

	for index, entry := range entries {
		path := fmt.Sprintf("%s.publish.entries[%d]", basePath, index)
		entrySnapshot, entryIssues := i.inspectEntry(ctx, name, path, moduleRootDir, entry)
		if len(entryIssues) > 0 {
			issues.Append(entryIssues)
			continue
		}

		entrySnapshots = append(entrySnapshots, entrySnapshot)
		entryHashes = append(entryHashes, entrySnapshot.Hash())
	}

	return entrySnapshots, entryHashes
}

// withModule annotates detached nested issues with their owning module name.
func withModule(in []Issue, name manifest.ModuleName) []Issue {
	out := cloneIssues(in)
	for index := range out {
		out[index].Module = name
	}

	return out
}

// inspectEntry validates one explicit publish source and hashes it when
// present.
func (i *inspector) inspectEntry(
	ctx context.Context,
	module manifest.ModuleName,
	path string,
	moduleRootDir string,
	entry manifest.PublishEntry,
) (EntrySnapshot, []Issue) {
	sourcePath := resolveEntrySource(moduleRootDir, entry)
	if err := pathutil.EnsureInside(moduleRootDir, sourcePath); err != nil {
		return EntrySnapshot{}, []Issue{entryIssue(
			IssueEntryPathEscape,
			module,
			path+".from",
			err.Error(),
		)}
	}

	exists, err := i.deps.FS.Exists(ctx, sourcePath)
	if err != nil {
		return EntrySnapshot{}, []Issue{entryIssue(
			IssueInvalidRequest,
			module,
			path+".from",
			err.Error(),
		)}
	}
	if !exists {
		return missingEntrySnapshot(module, path, sourcePath, entry)
	}

	isDir, err := i.deps.FS.IsDir(ctx, sourcePath)
	if err != nil {
		return EntrySnapshot{}, []Issue{entryIssue(
			IssueInvalidRequest,
			module,
			path+".from",
			err.Error(),
		)}
	}

	if issue, ok := entryTypeIssue(module, path, entry, isDir); ok {
		return EntrySnapshot{}, []Issue{issue}
	}

	hash, err := i.hashEntry(ctx, entry, sourcePath, isDir)
	if err != nil {
		return EntrySnapshot{}, []Issue{entryIssue(
			IssueEntryHashFailed,
			module,
			path+".from",
			err.Error(),
		)}
	}

	return EntrySnapshot{
		entry:      entry,
		sourcePath: sourcePath,
		targetPath: entry.To(),
		present:    true,
		hash:       hash,
	}, nil
}

// missingEntrySnapshot accepts absent optional entries and rejects absent
// required entries.
func missingEntrySnapshot(
	module manifest.ModuleName,
	path string,
	sourcePath string,
	entry manifest.PublishEntry,
) (EntrySnapshot, []Issue) {
	if entry.Optional() {
		return EntrySnapshot{
			entry:      entry,
			sourcePath: sourcePath,
			targetPath: entry.To(),
			present:    false,
		}, nil
	}

	return EntrySnapshot{}, []Issue{entryIssue(
		IssueEntryMissing,
		module,
		path+".from",
		"required source entry does not exist",
	)}
}

// entryTypeIssue rejects file/directory kind mismatches.
func entryTypeIssue(
	module manifest.ModuleName,
	path string,
	entry manifest.PublishEntry,
	isDir bool,
) (Issue, bool) {
	if entry.Kind() == manifest.PublishEntryFile && isDir {
		return entryIssue(
			IssueEntryTypeMismatch,
			module,
			path+".from",
			"file publish entry points to a directory",
		), true
	}

	if entry.Kind() == manifest.PublishEntryDirectory && !isDir {
		return entryIssue(
			IssueEntryTypeMismatch,
			module,
			path+".from",
			"directory publish entry points to a file",
		), true
	}

	return Issue{}, false
}

// entryIssue creates one module-scoped entry diagnostic.
func entryIssue(
	code IssueCode,
	module manifest.ModuleName,
	path string,
	message string,
) Issue {
	return Issue{
		Code:    code,
		Module:  module,
		Path:    path,
		Message: message,
	}
}

// hashEntry computes the stable content identity for one present entry.
func (i *inspector) hashEntry(
	ctx context.Context,
	entry manifest.PublishEntry,
	sourcePath string,
	isDir bool,
) (Hash, error) {
	if i.opts.DisableHashes {
		return "", nil
	}

	if isDir {
		return i.hashDirectoryEntry(ctx, entry, sourcePath)
	}

	return i.hashFileEntry(ctx, entry, sourcePath)
}

// hashDirectoryEntry combines the filesystem tree hash with entry routing.
func (i *inspector) hashDirectoryEntry(
	ctx context.Context,
	entry manifest.PublishEntry,
	sourcePath string,
) (Hash, error) {
	treeHash, err := i.deps.FS.TreeHash(ctx, sourcePath, portfs.TreeHashOptions{
		IncludeFileMode: true,
		SymlinkPolicy:   portfs.SymlinkReject,
	})
	if err != nil {
		return "", err
	}

	return hashBytes(
		"dir",
		entry.From().String(),
		entry.To().String(),
		treeHash.String(),
	), nil
}

// hashFileEntry hashes file content and entry routing.
func (i *inspector) hashFileEntry(
	ctx context.Context,
	entry manifest.PublishEntry,
	sourcePath string,
) (Hash, error) {
	data, err := i.deps.FS.ReadFile(ctx, sourcePath)
	if err != nil {
		return "", err
	}

	contentHash := hashBytes("file-content", string(data))
	return hashBytes(
		"file",
		entry.From().String(),
		entry.To().String(),
		contentHash.String(),
	), nil
}
