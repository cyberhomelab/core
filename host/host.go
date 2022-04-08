/*
   Copyright (c) 2022 Cyber Home Lab authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package host

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	logging "cyberhomelab.com/core/logging"
)

var log = logging.NewLogger()

func RunCommand(timeout int, command string, args ...string) (string, error) {
	// Create a new context and add a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel() // Cleanup

	// Creating the command with the context
	cmd := exec.CommandContext(ctx, command, args...)

	// Get the output
	log.Infof("Running command -> %s %s", command, strings.Join(args, " "))
	out, err := cmd.Output()

	// Timeout
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command timed out")
	}

	// Error
	if err != nil {
		return "", fmt.Errorf("non-zero exit code -> %s", err)
	}

	// Return
	return string(out), nil
}
