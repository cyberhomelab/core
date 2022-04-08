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
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func TestCheckIfIsFile_HappyFlow(t *testing.T) {
	// Define
	filePath := filepath.Join(os.TempDir(), "TestCheckIfIsFile_HappyFlow_file")
	dirPath := filepath.Join(os.TempDir(), "TestCheckIfIsFile_HappyFlow_dir")

	// Create
	assert.NilError(t, CreateFile(filePath, DefaultMode, DefaultUserId, DefaultGroupId))
	assert.NilError(t, CreateDirectory(dirPath, DefaultMode, DefaultUserId, DefaultGroupId))

	// Checks
	isFile, err := CheckIfIsFile(filePath)
	assert.NilError(t, err)
	assert.Assert(t, isFile)

	// Cleanup
	assert.NilError(t, Remove(filePath))
	assert.NilError(t, Remove(dirPath))
}

func TestCheckIfIsFile_NegativeFlow(t *testing.T) {
	// Define
	var isFile bool
	var err error

	// Test Case 1
	dirPath := filepath.Join(os.TempDir(), "TestCheckIfIsFile_NegativeFlow_dir")
	assert.NilError(t, CreateDirectory(dirPath, DefaultMode, DefaultUserId, DefaultGroupId))

	isFile, err = CheckIfIsFile(dirPath)
	assert.ErrorContains(t, err, "is not a regular file")
	assert.Assert(t, !isFile)

	assert.NilError(t, Remove(dirPath))

	// Test Case 2
	filePath := filepath.Join(os.TempDir(), "TestCheckIfIsFile_NegativeFlow_fileNotFound")

	isFile, err = CheckIfIsFile(filePath)
	assert.ErrorContains(t, err, "couldn't run os.Stat()")

	assert.Assert(t, !isFile)
}

func TestCheckHash_HappyFlow(t *testing.T) {
	// Define
	filePath1 := filepath.Join(os.TempDir(), "TestCheckHash_HappyFlow_1")
	filePath2 := filepath.Join(os.TempDir(), "TestCheckHash_HappyFlow_2")

	// Create
	assert.NilError(t, CreateFileWithMessage(filePath1, "Hello, 1!", DefaultMode, DefaultUserId, DefaultGroupId))
	assert.NilError(t, CreateFileWithMessage(filePath1, "Hello, 1!", DefaultMode, DefaultUserId, DefaultGroupId))

	// Check
	assert.NilError(t, CheckHash(filePath1, filePath2))

	// Cleanup
	for _, filePath := range []string{filePath1, filePath2} {
		assert.NilError(t, Remove(filePath))
	}
}

func TestCheckHash_NegativeFlow(t *testing.T) {
	// Define
	filePath1 := filepath.Join(os.TempDir(), "TestCheckHash_HappyFlow_1")
	filePath2 := filepath.Join(os.TempDir(), "TestCheckHash_HappyFlow_2")

	// Create
	assert.NilError(t, CreateFileWithMessage(filePath1, "Hello, 1!", DefaultMode, DefaultUserId, DefaultGroupId))
	assert.NilError(t, CreateFileWithMessage(filePath1, "Hello, 2!", DefaultMode, DefaultUserId, DefaultGroupId))

	// Check
	assert.ErrorContains(t, CheckHash(filePath1, filePath2), "hash missmatch between")

	// Cleanup
	for _, filePath := range []string{filePath1, filePath2} {
		assert.NilError(t, Remove(filePath))
	}
}

func TestCheckOwner_HappyFlow(t *testing.T) {
	// Define
	filePath1 := filepath.Join(os.TempDir(), "TestCheckOwner_HappyFlow_1")
	filePath2 := filepath.Join(os.TempDir(), "TestCheckOwner_HappyFlow_2")

	// Create
	assert.NilError(t, CreateFile(filePath1, DefaultMode, DefaultUserId, DefaultGroupId))
	assert.NilError(t, CreateFile(filePath2, DefaultMode, DefaultUserId, DefaultGroupId))

	// Check
	assert.NilError(t, CheckOwner(filePath1, filePath2))

	// Cleanup
	for _, filePath := range []string{filePath1, filePath2} {
		assert.NilError(t, Remove(filePath))
	}
}

func TestCheckOwner_NegativeFlow(t *testing.T) {
	// Define
	filePath := filepath.Join(os.TempDir(), "TestCheckOwner_NegativeFlow")

	// Check
	assert.ErrorContains(t, CreateFile(filePath, DefaultMode, OtherUserId, DefaultGroupId), "couldn't change the owner")

	// Cleanup
	assert.NilError(t, Remove(filePath))
}
