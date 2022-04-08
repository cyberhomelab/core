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
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestGetHash(t *testing.T) {
	message := "Hello World!"
	filePaths := []string{TestSourceFile, TestDestinationFile}

	// Create
	for _, filePath := range filePaths {
		assert.NilError(t, CreateFileWithMessage(filePath, message, DefaultMode, DefaultUserId, DefaultGroupId))
	}

	// Checks
	assert.NilError(t, CheckHash(TestSourceFile, TestDestinationFile))

	// Cleanup
	for _, filePath := range filePaths {
		assert.NilError(t, Remove(filePath))
	}

	// End
	t.Logf("GetHash() function works as expected.")
}

func TestCreateDirectoryHappyFlow(t *testing.T) {
	assert.NilError(t, CreateDirectory(TestSourceDirectory, DefaultMode, DefaultUserId, DefaultGroupId))
	assert.NilError(t, Remove(TestSourceDirectory))
}

func TestCreateDirectoryNegativeFlowCannotCreate(t *testing.T) {
	assert.ErrorContains(
		t,
		CreateDirectory(TestDirectoryCannotCreate, DefaultMode, DefaultUserId, DefaultGroupId),
		"couldn't create directory",
	)
}

func TestCreateDirectoryNegativeFlowCannotChangeOwner(t *testing.T) {
	assert.ErrorContains(
		t,
		CreateDirectory(TestSourceDirectory, DefaultMode, TestUserIdNotFound, DefaultGroupId),
		"couldn't change the owner",
	)
	assert.NilError(t, Remove(TestSourceDirectory))
}

func TestCountLinesInFileHappyFlow(t *testing.T) {
	nrOfLines, err := CountLinesInFile("../README.md")
	assert.NilError(t, err)
	assert.Assert(t, nrOfLines > 30)
}

func TestCountLinesInFileNegativeFlow(t *testing.T) {
	_, err := CountLinesInFile(TestFileNotFound)
	assert.ErrorContains(t, err, "couldn't read the file")
}

func TestReadFileHappyFlow(t *testing.T) {
	output, err := ReadFile("../README.md")
	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(output, "core"))
}

func TestReadFileNegativeFlow(t *testing.T) {
	_, err := ReadFile(TestFileNotFound)
	assert.ErrorContains(t, err, "no such file or directory")
}
