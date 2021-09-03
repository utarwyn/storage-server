// Copyright 2021 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	portPtr       = 8043
	basePath      = "/svr/http"
	pathPrefix    = "/"
	enableLogging = false
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

func configure() {
	flag.IntVar(&portPtr, "port", LookupEnvOrInt("PORT", portPtr), "listening port")
	flag.StringVar(&basePath, "base-path", LookupEnvOrString("BASE_PATH", basePath), "directory where all files are stored")
	flag.BoolVar(&enableLogging, "enable-logging", LookupEnvOrBool("ENABLE_LOGGING", enableLogging), "enable log request")
	flag.Parse()
}
