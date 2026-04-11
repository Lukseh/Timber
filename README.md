<p style="text-align: center;"><img src="https://raw.githubusercontent.com/Lukseh/Timber/refs/heads/main/timber.svg" height="300px"></p>
<h1 style="text-align: center;">Timber</h1>

Timber is Golang build tool

## Instalation

To install Timber you can use one of the following methods:

- Download the latest release from the [releases page](https://github.com/Lukseh/Timber/releases)

- Use `go install` command:

```bash
go install github.com/Lukseh/Timber@latest
```

## Gopherfile

Example Gopherfile:

```yaml
name: timber // Project name
version: 0.0.1 // Project version
go: 1.26.2 // Go version (minimum required)
description: Golang build tool // Project description

build: // Build scripts
  release: // Script name, release is default one
    entryfile: main.go // main file of the project, default is main.go
    outname: Timber.exe // Output binary name, default is project name
    options: -ldflags="-s -w" // Additional options to pass to the go build command, default is empty
    watch: false // Watch for file changes and re-build the command, default is false // Requires [gow](https://github.com/mitranim/gow)
    dockerize: false // Whenever to build the project in a docker container, default is false // Requires [docker](https://www.docker.com/)
    dockerfile: Dockerfile // Dockerfile to use when dockerize is true, default is Dockerfile, for other scripts defaults to Dockerfile.scriptname, e.g. Dockerfile.dev
```

### Logo of Timber uses graphics made by pch.vector downloaded on Freepik [HERE](https://www.freepik.com/free-vector/variety-wood-logs-trunks-flat-pictures-set_13146704.htm#fromView=search&page=1&position=0&uuid=a90f7664-1af3-4071-ad51-99d41841461c&query=timber")
