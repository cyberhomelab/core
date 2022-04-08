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
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func TestCreateArchive(t *testing.T) {
	// Happy flow
	filePaths := []string{TestFile1, TestFile2, TestFile3}
	filePathsWithDir := append(filePaths, TestSourceDirectory)
	var filePathWithDir string
	message := "Hello World in an Archive!"

	// Prerequisites
	assert.NilError(t, CreateDirectory(TestSourceDirectory, DefaultMode, DefaultUserId, DefaultGroupId))
	for _, filePath := range filePaths {
		assert.NilError(t, CreateFileWithMessage(filePath, message, DefaultMode, DefaultUserId, DefaultGroupId))
		filePathWithDir = filepath.Join(TestSourceDirectory, filepath.Base(filePath))
		assert.NilError(t, CreateFileWithMessage(filePathWithDir, message, DefaultMode, DefaultUserId, DefaultGroupId))
	}

	// Create archive
	assert.NilError(t, CreateArchive(TestArchivePath, filePathsWithDir))

	// Cleanup
	// for _, filePath := range filePathsWithDir {
	// 	assert.NilError(t, Remove(filePath))
	// }
	// Remove(TestArchivePath)

	// End
	t.Logf("CreateArchive() function works as expected.")
}
