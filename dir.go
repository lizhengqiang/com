// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"errors"
	"os"
	"path"
	"strings"
)

// IsDir checks whether the path is a directory.
// It returns false when it's a file or does not exist.
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

func statDir(dirPath, recPath string, includeDir bool) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	statList := make([]string, 0)
	for _, fi := range fis {
		relPath := strings.TrimPrefix(path.Join(recPath, fi.Name()), "/")
		curPath := strings.TrimPrefix(path.Join(dirPath, fi.Name()), "/")
		if fi.IsDir() {
			recPath = relPath
			if includeDir {
				statList = append(statList, relPath+"/")
			}
			s, err := statDir(curPath, recPath, includeDir)
			if err != nil {
				return nil, err
			}
			statList = append(statList, s...)
		} else {
			statList = append(statList, relPath)
		}
	}
	return statList, nil
}

// StatDir gathers information of given directory by depth-first.
// It returns slice of file list and includes subdirectories if enabled;
// it returns error and nil slice when error occurs in underlying functions,
// or given path is not a directory or does not exist.
//
// Slice does not include given path itself.
// If subdirectories is enabled, they will have suffix '/'.
func StatDir(dirPath string, includeDir ...bool) ([]string, error) {
	if !IsDir(dirPath) {
		return nil, errors.New("not a directory or does not exist: " + dirPath)
	}

	isIncludeDir := false
	if len(includeDir) >= 1 {
		isIncludeDir = includeDir[0]
	}

	return statDir(dirPath, "", isIncludeDir)
}

// CopyDir copy files recursively from source to target directory.
// It returns error when error occurs in underlying functions.
func CopyDir(srcPath, destPath string) error {
	// Check if target directory exists.
	if IsExist(destPath) {
		return errors.New("file or directory alreay exists: " + destPath)
	}

	err := os.Mkdir(destPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Gather directory info.
	infos, err := StatDir(srcPath, true)
	if err != nil {
		return err
	}

	for _, info := range infos {
		curPath := path.Join(destPath, info)
		if strings.HasSuffix(info, "/") {
			err = os.Mkdir(curPath, os.ModePerm)
		} else {
			_, err = Copy(path.Join(srcPath, info), curPath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
