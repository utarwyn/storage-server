// Copyright 2023 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && !strings.HasSuffix(r.URL.Path, "/") {
			LogRequest(r)

			// Check for authorization
			if !IsAuthorized(r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Parse input file
			if r.ParseMultipartForm(32<<20) == nil {
				if input, _, err := r.FormFile("file"); err == nil {
					defer input.Close()
					destFilepath := filepath.Join(basePath, r.URL.Path)

					// Create all nested directories if needed (not needed for now)
					//os.MkdirAll(filepath.Dir(destFilepath), os.ModePerm)

					// Copy file into the disk
					if dest, err := os.OpenFile(destFilepath, os.O_WRONLY|os.O_CREATE, os.ModePerm); err == nil {
						defer dest.Close()
						if _, err := io.Copy(dest, input); err == nil {
							w.WriteHeader(http.StatusNoContent)
						} else {
							panic(err)
						}
					} else {
						panic(err)
					}
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func DeleteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && !strings.HasSuffix(r.URL.Path, "/") {
			LogRequest(r)

			// Check for authorization
			if !IsAuthorized(r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Delete file if exists, otherwise write a not found error
			destFilepath := filepath.Join(basePath, r.URL.Path)
			if _, err := os.Stat(destFilepath); err == nil {
				if os.Remove(destFilepath) == nil {
					w.WriteHeader(http.StatusNoContent)
				} else {
					panic(err)
				}
			} else if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				panic(err)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
