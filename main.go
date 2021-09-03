// Copyright 2021 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handleReq(h http.Handler) http.Handler {
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

		if enableLogging {
			log.Println(r.Method, r.URL.Path)
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	configure()

	// Setup file system
	var fileSystem http.FileSystem = http.Dir(basePath)

	// Create file server object
	handler := handleReq(http.FileServer(fileSystem))

	// GZIP handler
	fileServer := handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fileServer.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzPool.Get().(*gzip.Writer)
			defer gzPool.Put(gz)

			gz.Reset(w)
			defer gz.Close()
			fileServer.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
		}
	})

	// Start the server
	http.Handle(pathPrefix, handler)

	// Log errors if found
	port := ":" + strconv.FormatInt(int64(portPtr), 10)
	log.Printf("Listening at 0.0.0.0%v %v...", port, pathPrefix)
	log.Fatalln(http.ListenAndServe(port, nil))
}
