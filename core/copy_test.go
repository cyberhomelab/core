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

package core

import (
	"fmt"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func TestCopyFile(t *testing.T) {
	var err error
	var filesMatch bool
	message := "Hello World!"

	// Create source
	assert.NilError(t, CreateFileWithMessage(TestSourceFile, message, DefaultMode, DefaultUserId, DefaultGroupId))

	// Copy
	nrOfBytes, err := CopyFile(TestSourceFile, TestDestinationFile)
	assert.NilError(t, err)

	// Checks
	if nrOfBytes == 0 {
		t.Fatalf("0 bytes copied from %s to %s", TestSourceFile, TestDestinationFile)
	}
	filesMatch, err = CheckIfFilesMatch(TestSourceFile, TestDestinationFile)
	assert.NilError(t, err)
	assert.Assert(t, filesMatch)

	// Cleanup
	for _, filePath := range []string{TestSourceFile, TestDestinationFile} {
		assert.NilError(t, Remove(filePath))
	}

	// End
	t.Logf("%d bytes copied successfully from %s to %s", nrOfBytes, TestSourceFile, TestDestinationFile)
}

func TestCopyDirectory(t *testing.T) {
	var message string

	// Create source
	assert.NilError(t, CreateDirectory(TestSourceDirectory, DefaultMode, DefaultUserId, DefaultGroupId))
	for k := 1; ; k++ {
		if k == 10 {
			break
		}
		filePath := filepath.Join(TestSourceDirectory, fmt.Sprint(k))
		message = fmt.Sprintf("Hello %d!", k)
		assert.NilError(t, CreateFileWithMessage(filePath, message, DefaultMode, DefaultUserId, DefaultGroupId))
	}

	// Copy
	assert.NilError(t, CopyDirectory(TestSourceDirectory, TestDestinationDirectory))

	// Checks
	assert.NilError(t, CheckIfDirectoriesMatch(TestSourceDirectory, TestDestinationDirectory))

	// Cleanup
	for _, directoryPath := range []string{TestSourceDirectory, TestDestinationDirectory} {
		assert.NilError(t, Remove(directoryPath))
	}

	// End
	t.Logf("%s copied successfully to %s", TestSourceDirectory, TestDestinationDirectory)
}

func TestCopy(t *testing.T) {
	var err error
	var filesMatch bool
	message := "Hello World!"

	// Create source
	assert.NilError(t, CreateFileWithMessage(TestSourceFile, message, DefaultMode, DefaultUserId, DefaultGroupId))

	// Copy
	assert.NilError(t, Copy(TestSourceFile, TestDestinationFile))

	// Checks
	filesMatch, err = CheckIfFilesMatch(TestSourceFile, TestDestinationFile)
	assert.NilError(t, err)
	assert.Assert(t, filesMatch)

	// Cleanup
	for _, filePath := range []string{TestSourceFile, TestDestinationFile} {
		assert.NilError(t, Remove(filePath))
	}

	// End
	t.Logf("Copy() function works as expected.")
}
