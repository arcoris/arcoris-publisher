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

package porttest

import (
	"context"
	"fmt"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

// GitCall records one fake Git operation.
type GitCall struct {
	// Op is the operation name.
	Op string

	// RepoDir is the repository worktree path.
	RepoDir string

	// Ref is a branch, tag, refspec, or commit depending on Op.
	Ref string

	// ForceWithLease records whether a push used lease protection.
	ForceWithLease bool

	// ForceWithLeaseRef records exact lease ref when supplied.
	ForceWithLeaseRef string

	// ForceWithLeaseExpect records exact lease object when supplied.
	ForceWithLeaseExpect git.CommitHash
}

// Git is a deterministic in-memory Git port for workflow tests.
type Git struct {
	// Statuses returns status by worktree path.
	Statuses map[string]git.Status

	// StatusErrors forces Status to fail for a worktree path.
	StatusErrors map[string]error

	// FetchError forces Fetch to fail.
	FetchError error

	// PushError forces Push to fail.
	PushError error

	// PushTagError forces PushTag to fail.
	PushTagError error

	// Tags reports local tag existence by tag name.
	Tags map[git.TagName]bool

	// RemoteRefs reports remote ref existence by "remote\x00ref".
	RemoteRefs map[string]bool

	// Refs reports local ref existence by "repoDir\x00ref" or ref.
	Refs map[string]bool

	// RemoteRefHashes reports remote ref object hashes by "remote\x00ref".
	RemoteRefHashes map[string]git.CommitHash

	// TagExistsError forces TagExists to fail.
	TagExistsError error

	// RemoteRefExistsError forces RemoteRefExists to fail.
	RemoteRefExistsError error

	// RemoteRefHashError forces RemoteRefHash to fail.
	RemoteRefHashError error

	// DeleteRemoteRefError forces DeleteRemoteRef to fail.
	DeleteRemoteRefError error

	// DeleteTagError forces DeleteTag to fail.
	DeleteTagError error

	// CommitHash is returned by Commit when non-empty.
	CommitHash git.CommitHash

	// Calls records mutating and transport operations in order.
	Calls []GitCall
}

// NewGit returns a fake Git port with clean default status.
func NewGit() *Git {
	return &Git{
		Statuses:        map[string]git.Status{},
		StatusErrors:    map[string]error{},
		Tags:            map[git.TagName]bool{},
		Refs:            map[string]bool{},
		RemoteRefs:      map[string]bool{},
		RemoteRefHashes: map[string]git.CommitHash{},
		CommitHash:      "abcdef1234567890",
	}
}

// Head returns a stable synthetic commit hash.
func (g *Git) Head(context.Context, string) (git.CommitHash, error) {
	return "abcdef1234567890", nil
}

// CurrentBranch returns a stable branch name.
func (g *Git) CurrentBranch(context.Context, string) (git.BranchName, error) {
	return "main", nil
}

// Status returns the configured worktree status or a clean default.
func (g *Git) Status(_ context.Context, repoDir string) (git.Status, error) {
	if err := g.StatusErrors[repoDir]; err != nil {
		return git.Status{}, err
	}
	status, ok := g.Statuses[repoDir]
	if !ok {
		return git.Status{Clean: true}, nil
	}

	return cloneStatus(status), nil
}

// RefExists reports configured local refs.
func (g *Git) RefExists(_ context.Context, repoDir string, ref string) (bool, error) {
	return g.Refs[repoDir+"\x00"+ref] || g.Refs[ref], nil
}

// RemoteRefExists reports configured remote refs.
func (g *Git) RemoteRefExists(_ context.Context, repoDir string, remote string, ref string) (bool, error) {
	if g.RemoteRefExistsError != nil {
		return false, g.RemoteRefExistsError
	}
	if hash := g.remoteRefHash(repoDir, remote, ref); hash != "" {
		return true, nil
	}
	if g.RemoteRefs[remoteRefKeyForRepo(repoDir, remote, ref)] {
		return true, nil
	}
	return g.RemoteRefs[remoteRefKey(remote, ref)], nil
}

// RemoteRefHash reports configured remote ref hashes.
func (g *Git) RemoteRefHash(_ context.Context, repoDir string, remote string, ref string) (git.CommitHash, bool, error) {
	if g.RemoteRefHashError != nil {
		return "", false, g.RemoteRefHashError
	}
	hash := g.remoteRefHash(repoDir, remote, ref)
	return hash, hash != "", nil
}

// CommitMessage returns an empty synthetic message.
func (g *Git) CommitMessage(context.Context, string, string) (string, error) {
	return "", nil
}

// Checkout records a checkout operation.
func (g *Git) Checkout(_ context.Context, repoDir string, ref string, _ git.CheckoutOptions) error {
	g.record("checkout", repoDir, ref, false)
	return nil
}

// CreateBranch records a branch creation operation.
func (g *Git) CreateBranch(
	_ context.Context,
	repoDir string,
	branch git.BranchName,
	startPoint string,
	_ git.CreateBranchOptions,
) error {
	g.record("create-branch", repoDir, fmt.Sprintf("%s@%s", branch, startPoint), false)
	return nil
}

// ResetHard records a hard reset operation.
func (g *Git) ResetHard(_ context.Context, repoDir string, ref string) error {
	g.record("reset-hard", repoDir, ref, false)
	return nil
}

// Clean records a Git clean operation.
func (g *Git) Clean(_ context.Context, repoDir string, _ git.CleanOptions) error {
	g.record("clean", repoDir, "", false)
	return nil
}

// AddAll records a Git add operation.
func (g *Git) AddAll(_ context.Context, repoDir string) error {
	g.record("add", repoDir, "", false)
	return nil
}

// Commit records a commit operation and returns CommitHash.
func (g *Git) Commit(
	_ context.Context,
	repoDir string,
	message string,
	_ git.CommitOptions,
) (git.CommitHash, error) {
	g.record("commit", repoDir, message, false)
	return g.CommitHash, nil
}

// Clone records a clone operation.
func (g *Git) Clone(_ context.Context, remoteURL string, dir string, _ git.CloneOptions) error {
	g.record("clone", dir, remoteURL, false)
	return nil
}

// Fetch records a fetch operation.
func (g *Git) Fetch(_ context.Context, repoDir string, remote string, _ git.FetchOptions) error {
	g.record("fetch", repoDir, remote, false)
	return g.FetchError
}

// Push records a branch push operation.
func (g *Git) Push(
	_ context.Context,
	repoDir string,
	remote string,
	refspec git.RefSpec,
	opts git.PushOptions,
) error {
	g.recordPush("push", repoDir, string(refspec), opts)
	if g.PushError == nil {
		g.recordRemotePush(repoDir, remote, refspec)
	}
	return g.PushError
}

// DeleteRemoteRef records a remote ref deletion.
func (g *Git) DeleteRemoteRef(
	_ context.Context,
	repoDir string,
	remote string,
	ref string,
	opts git.PushOptions,
) error {
	g.recordPush("delete-remote-ref", repoDir, ref, opts)
	if g.DeleteRemoteRefError != nil {
		return g.DeleteRemoteRefError
	}
	delete(g.RemoteRefs, remoteRefKey(remote, ref))
	delete(g.RemoteRefs, remoteRefKeyForRepo(repoDir, remote, ref))
	delete(g.RemoteRefHashes, remoteRefKey(remote, ref))
	delete(g.RemoteRefHashes, remoteRefKeyForRepo(repoDir, remote, ref))
	return nil
}

// TagExists reports configured local tags.
func (g *Git) TagExists(_ context.Context, _ string, tag git.TagName) (bool, error) {
	if g.TagExistsError != nil {
		return false, g.TagExistsError
	}
	return g.Tags[tag], nil
}

// CreateTag records a tag creation operation.
func (g *Git) CreateTag(
	_ context.Context,
	repoDir string,
	tag git.TagName,
	target git.CommitHash,
	_ git.TagOptions,
) error {
	g.record("tag", repoDir, fmt.Sprintf("%s@%s", tag, target), false)
	g.Tags[tag] = true
	return nil
}

// PushTag records a tag push operation.
func (g *Git) PushTag(
	_ context.Context,
	repoDir string,
	remote string,
	tag git.TagName,
	opts git.PushOptions,
) error {
	g.recordPush("push-tag", repoDir, string(tag), opts)
	if g.PushTagError == nil {
		ref := "refs/tags/" + tag.String()
		g.RemoteRefs[remoteRefKeyForRepo(repoDir, remote, ref)] = true
		g.RemoteRefHashes[remoteRefKeyForRepo(repoDir, remote, ref)] = g.CommitHash
	}
	return g.PushTagError
}

// DeleteTag records a local tag deletion.
func (g *Git) DeleteTag(_ context.Context, repoDir string, tag git.TagName) error {
	g.record("delete-tag", repoDir, string(tag), false)
	if g.DeleteTagError != nil {
		return g.DeleteTagError
	}
	delete(g.Tags, tag)
	return nil
}

func (g *Git) record(op, repoDir, ref string, forceWithLease bool) {
	g.Calls = append(g.Calls, GitCall{
		Op:             op,
		RepoDir:        repoDir,
		Ref:            ref,
		ForceWithLease: forceWithLease,
	})
}

func (g *Git) recordPush(op, repoDir, ref string, opts git.PushOptions) {
	g.Calls = append(g.Calls, GitCall{
		Op:                   op,
		RepoDir:              repoDir,
		Ref:                  ref,
		ForceWithLease:       opts.ForceWithLease || opts.ForceWithLeaseRef != "",
		ForceWithLeaseRef:    opts.ForceWithLeaseRef,
		ForceWithLeaseExpect: opts.ForceWithLeaseExpect,
	})
}

func cloneStatus(status git.Status) git.Status {
	return git.Status{
		Clean:   status.Clean,
		Entries: append([]git.StatusEntry(nil), status.Entries...),
	}
}

// RemoteRefKey returns the deterministic key used by RemoteRefs.
func RemoteRefKey(remote string, ref string) string {
	return remoteRefKey(remote, ref)
}

// RemoteRefKeyForRepo returns a repo-scoped key used by RemoteRefs.
func RemoteRefKeyForRepo(repoDir string, remote string, ref string) string {
	return remoteRefKeyForRepo(repoDir, remote, ref)
}

func remoteRefKey(remote string, ref string) string {
	return remote + "\x00" + ref
}

func remoteRefKeyForRepo(repoDir string, remote string, ref string) string {
	return repoDir + "\x00" + remote + "\x00" + ref
}

func (g *Git) remoteRefHash(repoDir string, remote string, ref string) git.CommitHash {
	if hash := g.RemoteRefHashes[remoteRefKeyForRepo(repoDir, remote, ref)]; hash != "" {
		return hash
	}
	return g.RemoteRefHashes[remoteRefKey(remote, ref)]
}

func (g *Git) recordRemotePush(repoDir string, remote string, refspec git.RefSpec) {
	before, after, ok := strings.Cut(refspec.String(), ":")
	if !ok || after == "" || before == "" {
		return
	}
	hash := git.CommitHash(before)
	if before == "HEAD" {
		hash = g.CommitHash
	}
	key := remoteRefKeyForRepo(repoDir, remote, after)
	g.RemoteRefs[key] = true
	g.RemoteRefHashes[key] = hash
}
