// Copyright 2023 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExposedFile struct {
	Path      string    `json:"path"`
	MD5       string    `json:"md5,omitempty"`
	SHA256    string    `json:"sha256,omitempty"`
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

// computeFileHashes reads the file once and returns hex digests by hash name.
func computeFileHashes(filePath string, names []string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hashes := make(map[string]hash.Hash, len(names))
	writers := make([]io.Writer, 0, len(names))
	for _, name := range names {
		var h hash.Hash
		switch name {
		case "md5":
			h = md5.New()
		case "sha256":
			h = sha256.New()
		default:
			continue
		}
		hashes[name] = h
		writers = append(writers, h)
	}
	if len(writers) > 0 {
		if _, err := io.Copy(io.MultiWriter(writers...), file); err != nil {
			return nil, err
		}
	}

	result := make(map[string]string, len(hashes))
	for name, h := range hashes {
		result[name] = hex.EncodeToString(h.Sum(nil))
	}
	return result, nil
}

func generateDirectoryDetails(directory string) (string, error) {
	var files []ExposedFile

	walkDir := filepath.Join(basePath, directory)
	walkErr := filepath.Walk(walkDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				file := ExposedFile{
					Path:      strings.TrimLeft(strings.ReplaceAll(strings.Replace(path, walkDir, "", 1), string(filepath.Separator), "/"), "/"),
					Size:      info.Size(),
					UpdatedAt: info.ModTime().UTC(),
				}
				hashes, hashErr := computeFileHashes(path, fileHashList)
				if hashErr != nil {
					return hashErr
				}
				file.MD5 = hashes["md5"]
				file.SHA256 = hashes["sha256"]
				files = append(files, file)
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
					if _, err := w.Write([]byte(details)); err == nil {
						LogRequest(r)
						return
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
