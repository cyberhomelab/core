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
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("couldn't read the file -> %s", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("couldn't calculate the hash -> %s", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func CountLinesInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("couldn't read the file -> %s", err)
	}
	defer file.Close()
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func ReadFile(filePath string) (string, error) {
	lineLimit := 10000
	nrOfLines, err := CountLinesInFile(filePath)
	if err != nil {
		return "", fmt.Errorf("couldn't determine the number of lines -> %s", err)
	}
	if nrOfLines > lineLimit {
		return "", fmt.Errorf("couldn't read the content, file %s exceeds the number of lines (%d) allowed", filePath, nrOfLines)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("couldn't read the file -> %s", err)
	}
	return string(content), nil
}

func WriteToFile(filePath string, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func Remove(fileOrDirPath string) error {
	fileOrDirStat, err := os.Stat(fileOrDirPath)
	if err != nil {
		return fmt.Errorf("couldn't run os.Stat() -> %s", err)
	}
	if fileOrDirStat.IsDir() {
		err := os.RemoveAll(fileOrDirPath)
		if err != nil {
			return fmt.Errorf("couldn't delete directory %s -> %s", fileOrDirStat, err)
		}
		return nil
	}
	if fileOrDirStat.Mode().IsRegular() {
		err := os.Remove(fileOrDirPath)
		if err != nil {
			return fmt.Errorf("couldn't delete file %s -> %s", fileOrDirStat, err)
		}
		return nil
	}
	return fmt.Errorf("%s can't be deleted because it is not a regular file or a directory", fileOrDirPath)
}

func CreateFile(filePath string, permissions fs.FileMode, userId int, groupId int) error {
	currentFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("couldn't create file %s -> %s", filePath, err)
	}
	currentFile.Close()
	err = os.Chmod(filePath, permissions)
	if err != nil {
		return fmt.Errorf(
			"couldn't change permissions (%s) for file %s -> %s",
			permissions, filePath, err)
	}
	err = os.Chown(filePath, userId, groupId)
	if err != nil {
		return fmt.Errorf(
			"couldn't change the owner (%d:%d) for file %s -> %s",
			userId, groupId, filePath, err)
	}
	return nil
}

func CreateFileWithMessage(filePath string, message string, mode fs.FileMode, userId int, groupId int) error {
	err := CreateFile(filePath, mode, userId, groupId)
	if err != nil {
		return fmt.Errorf("couldn't create file -> %s", err)
	}
	err = WriteToFile(filePath, message)
	if err != nil {
		return fmt.Errorf("couldn't write to file -> %s", err)
	}
	return nil
}

func CreateDirectory(directoryPath string, permissions fs.FileMode, userId int, groupId int) error {
	err := os.Mkdir(directoryPath, permissions)
	if err != nil {
		return fmt.Errorf(
			"couldn't create directory %s with permissions %s -> %s",
			directoryPath, permissions, err)
	}
	err = os.Chown(directoryPath, userId, groupId)
	if err != nil {
		return fmt.Errorf(
			"couldn't change the owner (%d:%d) for directory %s -> %s",
			userId, groupId, directoryPath, err)
	}
	return nil
}

func GetNumberOfFiles(directoryPath string) (int, error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return 0, fmt.Errorf("couldn't get the number of files in directory %s -> %s", directoryPath, err)
	}
	return len(files), nil
}

func GetDirectorySize(directoryPath string) (int64, error) {
	var size int64
	err := filepath.Walk(directoryPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("couldn't get the size for directory %s -> %s", directoryPath, err)
	}
	return size, nil
}
