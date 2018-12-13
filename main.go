// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary apt-golang-s3 implements the APT method interface in order to
// allow hosting of APT packages in Amazon S3. For more information about
// the APT method interface see, http://www.fifi.org/doc/libapt-pkg-doc/method.html/ch2.html#s2.3.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/crashlytics/apt-golang-s3/method"
)

const (
	version = "1.0.0"
)

var (
	showVersion = flag.Bool("version", false, "Print version and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("apt-golang-s3 %s (Go version: %s)\n", version, runtime.Version())
		os.Exit(0)
	}

	method.New().Run()
}
