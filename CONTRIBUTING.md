# Contribute

gitscanner uses GitHub to manage reviews of pull requests:

- If you have a trivial fix or improvement, go ahead and create a pull request.
- If you plan to do something more involved, discuss your ideas on the relevant GitHub issue or the discussions page.

## Steps to contribute

For now, you need to add your fork as a remote on the original **\$GOPATH**/src/github.com/MrWong99/gitscanner clone, so:

```bash

$ go get github.com/MrWong99/gitscanner
$ cd $GOPATH/src/github.com/gMrWong99/gitscanner # GOPATH is $HOME/go by default.

$ git remote add <FORK_NAME> <FORK_URL>
```

Notice: `go get` return `package github.com/MrWong99/gitscanner: no Go files in /go/src/github.com/MrWong99/gitscanner` is normal.

Also you need to be able to use `gcc` or `g++` in your cli. So if you are using Windows consider using `Mingw`, `Cygwin` or a Unix VM (e.g. using `WSL 2`)
to build the software.

### Dependency management

We use [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies on external packages.
This requires a working Go environment with version 1.17 or greater and git installed.

To add or update a new dependency, use the `go get` command:

```bash
# Pick the latest tagged release.
go get example.com/some/module/pkg

# Pick a specific version.
go get example.com/some/module/pkg@vX.Y.Z
```

Tidy up the `go.mod` and `go.sum` files:

```bash
go mod tidy
go mod vendor
git add go.mod go.sum vendor
git commit
```

You have to commit the changes to `go.mod` and `go.sum` before submitting the pull request.

## Coding Standards

### go imports
imports should follow `std libs`, `other libs` format

Example
```
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/MrWong99/gitscanner/checks"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/gorilla/mux"
)
```
