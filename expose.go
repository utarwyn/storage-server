// Copyright 2021 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExposedFile struct {
	Path      string    `json:"path"`
	MD5       string    `json:"md5"`
	Size      int64     `json:"size"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func extractExposeDirectory(path string) *string {
	if strings.HasSuffix(path, "/") {
		return &path
	} else {
		return nil
	}
}

func isDirectoryExposed(directory string) bool {
	for _, directoryInList := range exposeDirList {
		if directoryInList == directory {
			return true
		}
	}
	return false
}

func getFileMD5Hash(filePath string) (string, error) {
	var md5Ptr string
	file, err := os.Open(filePath)
	if err != nil {
		return md5Ptr, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5Ptr, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5Ptr = hex.EncodeToString(hashInBytes)
	return md5Ptr, nil
}

func generateDirectoryDetails(directory string) (string, error) {
	var files []ExposedFile

	walkErr := filepath.Walk(filepath.Join(basePath, directory),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				md5Hash, hashErr := getFileMD5Hash(path)
				if hashErr != nil {
					return hashErr
				}
				files = append(files, ExposedFile{
					Path:      strings.TrimLeft(strings.ReplaceAll(strings.TrimLeft(path, basePath), "\\", "/"), directory),
					MD5:       md5Hash,
					Size:      info.Size(),
					UpdatedAt: info.ModTime().UTC(),
				})
			}
			return nil
		})

	if walkErr != nil {
		return "", walkErr
	}

	if encoded, err := json.Marshal(files); err == nil {
		return string(encoded), nil
	} else {
		return "", err
	}
}

func ExposeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expose middleware is only accessible when authorized
		if IsAuthorized(r) {
			directory := extractExposeDirectory(r.URL.Path)
			if directory != nil && isDirectoryExposed(*directory) {
				if details, err := generateDirectoryDetails(*directory); err == nil {
					w.Header().Set("content-type", "application/json; charset=utf-8")
					w.Write([]byte(details))
					LogRequest(r)
					return
				} else {
					panic(err)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
