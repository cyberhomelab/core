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
	"os"
	"syscall"
)

func CheckIfIsFile(filePath string) (bool, error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if !fileStat.Mode().IsRegular() {
		return false, fmt.Errorf("%s is not a regular file", filePath)
	}
	return true, nil
}

func CheckHash(sourceFilePath string, destinationFilePath string) error {
	sourceHash, err := GetHash(sourceFilePath)
	if err != nil {
		return fmt.Errorf("couldn't get the hash for file %s -> %s", sourceFilePath, err)
	}

	destinationHash, err := GetHash(destinationFilePath)
	if err != nil {
		return fmt.Errorf("couldn't get the hash for file %s -> %s", sourceFilePath, err)
	}

	if sourceHash != destinationHash {
		return fmt.Errorf(
			"hash missmatch between %s (%s) and %s (%s)",
			sourceFilePath, sourceHash, destinationFilePath, destinationHash)
	}

	return nil
}

func CheckOwner(sourceFilePath string, destinationFilePath string) error {
	sourceFileStat, err := os.Stat(sourceFilePath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	destinationFileStat, err := os.Stat(destinationFilePath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	sourceFileStatSys := sourceFileStat.Sys().(*syscall.Stat_t)
	destinationFileStatSys := destinationFileStat.Sys().(*syscall.Stat_t)
	if sourceFileStatSys.Uid != destinationFileStatSys.Uid {
		return fmt.Errorf(
			"file %s (Uid %d) has a different Uid in comparison with %s (Uid %d)",
			sourceFilePath, sourceFileStatSys.Uid,
			destinationFilePath, destinationFileStatSys.Uid)
	}
	if sourceFileStatSys.Gid != destinationFileStatSys.Gid {
		return fmt.Errorf(
			"file %s (Gid %d) has a different Gid in comparison with %s (Gid %d)",
			sourceFilePath, sourceFileStatSys.Gid,
			destinationFilePath, destinationFileStatSys.Gid)
	}
	return nil
}

func CheckPermissions(sourceFilePath string, destinationFilePath string) error {
	sourceFileStat, err := os.Stat(sourceFilePath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	destinationFileStat, err := os.Stat(destinationFilePath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if sourceFileStat.Mode() != destinationFileStat.Mode() {
		return fmt.Errorf(
			"permissions missmatch between %s (%s) and %s (%s)",
			sourceFilePath, sourceFileStat.Mode(),
			destinationFilePath, destinationFileStat.Mode())
	}
	return nil
}

func CheckIfFilesMatch(sourceFilePath string, destinationFilePath string) (bool, error) {
	var err error

	for _, filePath := range []string{sourceFilePath, destinationFilePath} {
		_, err := CheckIfIsFile(filePath)
		if err != nil {
			return false, fmt.Errorf("couldn't run CheckIfIsFile() -> %s", err)
		}
	}

	err = CheckHash(sourceFilePath, destinationFilePath)
	if err != nil {
		return false, fmt.Errorf("an error received from CheckHash() -> %s", err)
	}

	err = CheckOwner(sourceFilePath, destinationFilePath)
	if err != nil {
		return false, fmt.Errorf("an error received from CheckOwner() -> %s", err)
	}

	err = CheckPermissions(sourceFilePath, destinationFilePath)
	if err != nil {
		return false, fmt.Errorf("an error received from CheckPermissions() -> %s", err)
	}

	return true, nil
}

func CheckIfDirectoriesMatch(sourceDirectoryPath string, destinationDirectoryPath string) error {
	// Prechecks
	sourceStat, err := os.Stat(sourceDirectoryPath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	destinationStat, err := os.Stat(destinationDirectoryPath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if !sourceStat.IsDir() || !destinationStat.IsDir() {
		return fmt.Errorf("%s or %s is not a directory", sourceDirectoryPath, destinationDirectoryPath)
	}

	// Nr. of files
	nrOfFilesInSource, err := GetNumberOfFiles(sourceDirectoryPath)
	if err != nil {
		return fmt.Errorf("couldn't get the numer of files in directory %s -> %s", sourceDirectoryPath, err)
	}
	nrOfFilesInDestination, err := GetNumberOfFiles(destinationDirectoryPath)
	if err != nil {
		return fmt.Errorf("couldn't get the numer of files in directory %s -> %s", destinationDirectoryPath, err)
	}
	if nrOfFilesInSource != nrOfFilesInDestination {
		return fmt.Errorf("the number of files differ between %s and %s", sourceDirectoryPath, destinationDirectoryPath)
	}

	// Size
	sourceDirectoryBytes, err := GetDirectorySize(sourceDirectoryPath)
	if err != nil {
		return err
	}
	destinationDirectoryBytes, err := GetDirectorySize(destinationDirectoryPath)
	if err != nil {
		return err
	}
	if sourceDirectoryBytes != destinationDirectoryBytes {
		return fmt.Errorf("not the same amount of data are available in %s and %s", sourceDirectoryPath, destinationDirectoryPath)
	}

	// End
	return nil
}
