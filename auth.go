// Copyright 2023 Maxime Malgorn. All rights reserved.
// Use of this source code is governed by a MIT-style.
// The license can be found in the LICENSE file.

package main

import "net/http"

func IsAuthorized(r *http.Request) bool {
	return clientSecret != "" && r.Header.Get("x-client-secret") == clientSecret
}
