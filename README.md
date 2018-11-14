# apt-golang-s3

An s3 transport method for the `apt` package management system

## TL;DR
The `apt` program creates a child process using the transport method and writes to it's STDIN, the method communicates back to `apt` by writing to STDOUT.

## Building the go program:

```
$ go build -o apt-golang-s3 main.go
```

There is an included Dockerfile to setup an environment for building the binary in a sandboxed environment:

```
$ ls
Dockerfile  main.go  method  README.md

$ docker build -t apt-golang-s3 .

$ docker run -it --rm -v $(pwd):/app apt-golang-s3 bash

root@83823fffd369:/app# ls
Dockerfile  README.md  build-deb.sh  go.mod  go.sum  main.go  method

root@83823fffd369:/app# go build -o apt-golang-s3 main.go
go: finding github.com/pmezard/go-difflib v1.0.0
go: finding github.com/davecgh/go-spew v1.1.1
go: finding github.com/stretchr/testify v1.2.2
go: finding github.com/aws/aws-sdk-go v1.15.73
go: finding golang.org/x/text v0.3.0
go: finding golang.org/x/net v0.0.0-20181108082009-03003ca0c849
go: finding github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8
go: downloading github.com/aws/aws-sdk-go v1.15.73
go: downloading github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8

root@83823fffd369:/app# ls
Dockerfile  README.md  apt-golang-s3  build-deb.sh  go.mod  go.sum  main.go  method

root@83823fffd369:/app# exit
exit

$ ls
apt-golang-s3  build-deb.sh  Dockerfile  go.mod  go.sum  main.go  method  README.md
```

## Building a debian package:

For convenience, there is a small bash script in the repository that can build the binary and package it as a .deb.

```
$ ls
build-deb.sh  Dockerfile  go.mod  go.sum  main.go  method  README.md

$ docker build -t apt-golang-s3 .

$ docker run -it --rm -v $(pwd):/app apt-golang-s3 /app/build-deb.sh
go: finding github.com/stretchr/testify v1.2.2
go: finding github.com/davecgh/go-spew v1.1.1
go: finding github.com/pmezard/go-difflib v1.0.0
go: finding github.com/aws/aws-sdk-go v1.15.73
go: finding golang.org/x/text v0.3.0
go: finding golang.org/x/net v0.0.0-20181108082009-03003ca0c849
go: finding github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8
go: downloading github.com/aws/aws-sdk-go v1.15.73
go: downloading github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8
/var/lib/gems/2.3.0/gems/fpm-1.10.2/lib/fpm/util.rb:29: warning: Insecure world writable dir /go/bin in PATH, mode 040777
Debian packaging tools generally labels all files in /etc as config files, as mandated by policy, so fpm defaults to this behavior for deb packages. You can disable this default behavior with --deb-no-default-config-files flag {:level=>:warn}
Created package {:path=>"apt-golang-s3_1_amd64.deb"}

$ ls
apt-golang-s3  apt-golang-s3_1_amd64.deb  build-deb.sh  Dockerfile  go.mod  go.sum  main.go  method  README.md
```

## Installing in production:

The `apt-golang-s3` binary is an executable. To install it copy it to `/usr/lib/apt/methods/s3` on your computer.
The .deb file produced by `build-deb.sh` will install the method in the correct place.

## APT Repository Configuration:

AWS keys are specified in the apt sources list configuration as follows:

```
$ cat /etc/apt/sources.list.d/my-private-repo.list
deb s3://aws-access-key-id:aws-secret-access-key@s3.amazonaws.com/my-private-repo-bucket stable main
```
