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

package cli

import "strings"

type command string

const (
	commandHelp    command = "help"
	commandPlan    command = "plan"
	commandVerify  command = "verify"
	commandPublish command = "publish"
	commandVersion command = "version"
)

func parseCommand(args []string) (command, []string, error) {
	if len(args) == 0 {
		return commandHelp, nil, nil
	}

	name := strings.TrimSpace(args[0])
	switch command(name) {
	case commandHelp, commandPlan, commandVerify, commandPublish, commandVersion:
		return command(name), args[1:], nil
	case "-h", "--help":
		return commandHelp, args[1:], nil
	default:
		return "", nil, &Error{Code: CodeInvalidCommand, Message: "unknown command: " + name}
	}
}

func isHelpRequest(args []string) bool {
	for _, arg := range args {
		switch strings.TrimSpace(arg) {
		case "-h", "--help":
			return true
		}
	}
	return false
}
