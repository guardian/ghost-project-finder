package main

import (
	"errors"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

type FileEntry struct {
	os.FileInfo

	FullPath string
}

func shouldExclude(path string, excludeList *[]string) bool {
	for _, entry := range *excludeList {
		if strings.HasPrefix(path, entry) {
			return true
		}
	}
	return false
}

func recursiveRead(forPath string, matcher *regexp.Regexp, outputCh chan *FileEntry, excludeList *[]string) error {
	//log.Print("DEBUG recursiveRead on ", forPath)
	entries, err := os.ReadDir(forPath)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			newPath := path.Join(forPath, e.Name())
			if !shouldExclude(newPath, excludeList) {
				subErr := recursiveRead(newPath, matcher, outputCh, excludeList)
				if subErr != nil && !errors.Is(subErr, os.ErrPermission) {
					log.Print("WARNING ", subErr)
				}
			}
		} else if matcher.MatchString(e.Name()) {
			//log.Printf("DEBUG recursiveRead found %s in %s", e.Name(), forPath)
			fullPath := path.Join(forPath, e.Name())
			fileInfo, statErr := os.Stat(fullPath)
			if statErr != nil {
				return statErr
			} else {
				output := FileEntry{
					FileInfo: fileInfo,
					FullPath: fullPath,
				}
				outputCh <- &output
			}
		}
	}

	return nil
}

func AsyncScanner(startPath *string, matcher *regexp.Regexp, excludeList *[]string) (chan *FileEntry, chan error) {
	outputCh := make(chan *FileEntry, 100)
	errCh := make(chan error, 1)

	go func() {
		err := recursiveRead(*startPath, matcher, outputCh, excludeList)
		if err != nil {
			log.Printf("ERROR could not read paths: %s", err)
			errCh <- err
		}
		outputCh <- nil
	}()

	return outputCh, errCh
}
