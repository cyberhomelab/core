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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func addToArchive(tarWriter *tar.Writer, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("couldn't open file %s -> %s", filePath, err)
	}
	defer file.Close()

	// ...
	info, err := os.Stat(filePath)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(filePath)
	}

	// Go through each file
	return filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, filePath))
			}

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarWriter, file)
			return err
		},
	)
}

func CreateArchive(archivePath string, filePaths []string) error {
	// Create an empty file that will be used by tar
	outFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("couldn't create archive %s -> %s", archivePath, err)
	}
	defer outFile.Close()

	// Creating the gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Creating the tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Adding all the files to the archive
	for _, filePath := range filePaths {
		err := addToArchive(tarWriter, filePath)
		if err != nil {
			return fmt.Errorf("couldn't add file %s to archive %s -> %s", filePath, archivePath, err)
		}
	}
	return nil
}
