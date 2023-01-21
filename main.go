// Copyright 2023 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func LogRequest(r *http.Request) {
	if enableLogging {
		log.Println(r.Method, r.URL.Path)
	}
}

func ApplyCorsPolicy(origin string, w http.ResponseWriter) {
	if origin != "" && len(allowOriginList) > 0 {
		for _, allowedOrigin := range allowOriginList {
			if allowedOrigin == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
	}
}

func handleRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// forward to https if still using http
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			if enableLogging {
				log.Println(301, r.Method, r.URL.Path)
			}
			return
		}

		// do not display directory listings
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		// apply CORS policy based on origin header
		ApplyCorsPolicy(r.Header.Get("Origin"), w)

		LogRequest(r)
		h.ServeHTTP(w, r)
	})
}

func main() {
	Configure()

	// Setup file system
	var fileSystem http.FileSystem = http.Dir(basePath)

	// Create file server object
	handler := handleRequest(http.FileServer(fileSystem))

	// Middlewares
	handler = UploadMiddleware(handler)
	handler = DeleteMiddleware(handler)
	handler = ExposeMiddleware(handler)
	handler = CachingMiddleware(handler)
	handler = GzipMiddleware(handler)

	// Start the server
	http.Handle(pathPrefix, handler)

	// Log errors if found
	port := ":" + strconv.FormatInt(int64(portPtr), 10)
	log.Printf("Serving directory %v", basePath)
	log.Printf("Listening at 0.0.0.0%v %v...", port, pathPrefix)
	log.Fatalln(http.ListenAndServe(port, nil))
}
