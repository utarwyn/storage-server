// Copyright 2021 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"net/http"
	"strings"
)

func isURLPathCached(urlPath string) bool {
	for _, directoryInList := range cacheDirList {
		if directoryInList == "*" || (directoryInList != "" && strings.HasPrefix(urlPath, directoryInList)) {
			return true
		}
	}
	return false
}

func CachingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if isURLPathCached(r.URL.Path) {
			header := w.Header()
			header.Set("vary", "Accept-Encoding")
			header.Set("cache-control", "public, max-age=31536000")
		}
	})
}
