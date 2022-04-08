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
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestRunCommandHappyFlow(t *testing.T) {
	output, err := RunCommand(5, "curl", "--head", "https://www.google.ro/")
	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(output, "HTTP/2 200"))
}

func TestRunCommandNegativeFlow(t *testing.T) {
	_, err := RunCommand(5, "notacommand")
	assert.ErrorContains(t, err, "executable file not found in $PATH")
}
