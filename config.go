// Copyright 2023 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

var (
	portPtr         = 8043
	basePath        = "/svr/http"
	pathPrefix      = "/"
	clientSecret    = ""
	enableLogging   = false
	cacheDirList    []string
	exposeDirList   []string
	allowOriginList []string
)

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, _ := strconv.Atoi(val)
		return v
	}
	return defaultVal
}

func LookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, _ := strconv.ParseBool(val)
		return v
	}
	return defaultVal
}

func Configure() {
	var cacheDirString string
	var exposeDirString string
	var allowOriginString string

	flag.IntVar(&portPtr, "port", LookupEnvOrInt("PORT", portPtr), "listening port")
	flag.StringVar(&basePath, "base-path", LookupEnvOrString("BASE_PATH", basePath), "directory where all files are stored")
	flag.StringVar(&clientSecret, "client-secret", LookupEnvOrString("CLIENT_SECRET", clientSecret), "secret key used to access privileged routes")
	flag.BoolVar(&enableLogging, "enable-logging", LookupEnvOrBool("ENABLE_LOGGING", enableLogging), "enable log request")
	flag.StringVar(&cacheDirString, "caching-directories", LookupEnvOrString("CACHING_DIRECTORIES", ""), "list of directories to cache")
	flag.StringVar(&exposeDirString, "expose-directories", LookupEnvOrString("EXPOSE_DIRECTORIES", ""), "list of directories to expose")
	flag.StringVar(&allowOriginString, "allow-origins", LookupEnvOrString("ALLOW_ORIGINS", ""), "list of origins to allow using CORS")
	flag.Parse()

	cacheDirList = strings.Split(cacheDirString, ";")
	exposeDirList = strings.Split(exposeDirString, ";")
	allowOriginList = strings.Split(allowOriginString, ";")
}
