# <img src="./assets/go-ppic.png" width="32" height="32" valign="middle" title="go-ppic example"> <img src="./assets/hello-world.png" width="32" height="32" valign="middle" title="go-ppic hello-world example"> <img src="./assets/jackwilsdon.png" width="32" height="32" valign="middle" title="go-ppic jackwilsdon example"> go-ppic <a href="https://travis-ci.org/jackwilsdon/go-ppic" title="Build status"><img src="https://img.shields.io/travis/jackwilsdon/go-ppic.svg" valign="middle" title="Build status"></a> <a href="https://goreportcard.com/report/github.com/jackwilsdon/go-ppic" title="Go Report Card"><img src="https://goreportcard.com/badge/github.com/jackwilsdon/go-ppic" valign="middle" title="Go Report Card status"></a> <a href="https://godoc.org/github.com/jackwilsdon/go-ppic" title="GoDoc reference"><img src="https://godoc.org/github.com/jackwilsdon/go-ppic?status.svg" valign="middle" title="GoDoc reference"></a>

Profile picture generation service written in Go. A demo can be found at [ppic.now.sh](https://ppic.now.sh/hello).

`go-ppic` provides two commands; [`ppicd`](#ppicd) and [`ppic`](#ppic).

## ppicd

`ppicd` is a web server providing image generation.

### Installation

```Shell
go get -u github.com/jackwilsdon/go-ppic/cmd/ppicd
```

### Usage

```Text
  -d	enable pprof debug routes
  -h string
    	host to run the server on
  -p uint
    	port to run the server on (default 3000)
  -v	enable verbose output
  -z	enable gzip compression
```

After starting up the server, you should see something similar to the following output;

```Text
2006/01/02 15:04:05 Starting server on http://127.0.0.1:3000...
```

Visiting the URL that the server is running on will give you the image for an empty string. You can get the image for
the string "example" by visiting `/example` on the server (`http://127.0.0.1:3000/example` in this case).

### URL Parameters

The server accepts the following query parameters to change the response;

 * `?size=N` → specify the size of the image to return (must be a multiple of 8)
 * `?monochrome` → change the image to black and white

### Supported Extensions

By default the server will respond in PNG format, but it also supports the following file extensions;

 * `.gif`
 * `.jpeg`

## ppic

`ppic` is used to generate profile pictures on the command line, without having to run a web server. `ppic` outputs the generated image to stdout.

### Installation

```Shell
go get -u github.com/jackwilsdon/go-ppic/cmd/ppic
```

### Usage

```Text
usage: ppic username [size] > image.png
```

> `size` defaults to 512 if not provided

### Examples

```Shell
ppic jackwilsdon 1024 > profile.png
```
