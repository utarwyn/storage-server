// Copyright 2021 Maxime Malgorn. All rights reserved.
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
		if r.Method == "PUT" && !strings.HasSuffix(r.URL.Path, "/") {
			LogRequest(r)

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
							w.WriteHeader(201)
						} else {
							panic(err)
						}
					} else {
						panic(err)
					}
				} else {
					w.WriteHeader(400)
				}
			} else {
				w.WriteHeader(400)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
