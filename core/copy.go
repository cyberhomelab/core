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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func CopyFile(sourceFilePath string, destinationPath string) (int64, error) {
	// Checks
	sourceFileStat, err := os.Stat(sourceFilePath)
	if err != nil {
		return 0, fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if sourceFileStat.IsDir() {
		return 0, fmt.Errorf("%s is a directory, not a file", sourceFilePath)
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", sourceFileStat)
	}

	// Open the file
	source, err := os.Open(sourceFilePath)
	if err != nil {
		return 0, fmt.Errorf("couldn't run os.Open() -> %s", err)
	}
	defer source.Close()

	// Preparing the destination
	destinationStat, err := os.Stat(destinationPath)
	if err != nil && !os.IsNotExist(err) {
		return 0, fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if !os.IsNotExist(err) {
		if destinationStat.IsDir() {
			destinationPath = filepath.Join(destinationPath, sourceFileStat.Name())
		}
	}
	destination, err := os.Create(destinationPath)
	if err != nil {
		return 0, fmt.Errorf("couldn't run os.Create() -> %s", err)
	}
	defer destination.Close()

	// Coping file to the destination
	nrOfBytes, err := io.Copy(destination, source)
	if err != nil {
		return 0, fmt.Errorf("couldn't copy file to destination -> %s", err)
	}

	// Preserving the permissions
	err = os.Chmod(destinationPath, sourceFileStat.Mode())
	if err != nil {
		return 0, fmt.Errorf(
			"couldn't preserve the permissions in the destination location %s -> %s",
			destinationPath, err)
	}

	// Preserving the ownership
	sourceFileStatSys := sourceFileStat.Sys().(*syscall.Stat_t)
	err = os.Chown(destinationPath, int(sourceFileStatSys.Uid), int(sourceFileStatSys.Gid))
	if err != nil {
		return 0, fmt.Errorf(
			"couldn't preserve the ownership in the destination location %s -> %s",
			destinationPath, err)
	}

	// Finish
	return nrOfBytes, nil
}

func CopyDirectory(sourceDirectoryPath string, destinationDirectoryPath string) error {
	// Cleanup
	sourceDirectoryPath = filepath.Clean(sourceDirectoryPath)
	destinationDirectoryPath = filepath.Clean(destinationDirectoryPath)

	// Checks
	sourceStat, err := os.Stat(sourceDirectoryPath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if !sourceStat.IsDir() {
		return fmt.Errorf("%s is not a directory", sourceDirectoryPath)
	}
	_, err = os.Stat(destinationDirectoryPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if err == nil {
		return fmt.Errorf(
			"destination directory %s already exists, please delete it or use a different path",
			destinationDirectoryPath)
	}

	// Create the destination directory
	err = os.MkdirAll(destinationDirectoryPath, sourceStat.Mode())
	if err != nil {
		return fmt.Errorf("couldn't create a directory under the following path %s -> %s", destinationDirectoryPath, err)
	}

	// Get all the files in source directory
	entries, err := ioutil.ReadDir(sourceDirectoryPath)
	if err != nil {
		return fmt.Errorf("couldn't get all the files under the following path %s -> %s", sourceDirectoryPath, err)
	}

	// Copy files
	for _, entry := range entries {
		currentSourcePath := filepath.Join(sourceDirectoryPath, entry.Name())
		currentDestinationPath := filepath.Join(destinationDirectoryPath, entry.Name())

		if entry.IsDir() {
			err = CopyDirectory(currentSourcePath, currentDestinationPath)
			if err != nil {
				return fmt.Errorf("couldn't run CopyDirectory() -> %s", err)
			}
		} else {
			// Skip symlinks
			// TODO: Copy symlinks
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			_, err = CopyFile(currentSourcePath, currentDestinationPath)
			if err != nil {
				return fmt.Errorf("couldn't run CopyFile() -> %s", err)
			}
		}
	}

	// Finish
	return nil
}

func Copy(sourcePath string, destinationPath string) error {
	sourceStat, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if sourceStat.IsDir() {
		return CopyDirectory(sourcePath, destinationPath)
	}
	if sourceStat.Mode().IsRegular() {
		_, err := CopyFile(sourcePath, destinationPath)
		return err
	}
	return fmt.Errorf("coudln't copy file %s because it has an unsuported type for the Copy() function", sourcePath)
}
