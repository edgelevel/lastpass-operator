## golang

* [go](https://golang.org/doc) documentation
* [dep](https://golang.github.io/dep/docs/introduction.html) Dependency management for Go
* [Convert JSON into a Go type definition](https://mholt.github.io/json-to-go/)

```bash
# download source
mkdir -p $GOPATH/src/github.com/edgelevel && cd $_
git clone git@github.com:edgelevel/lastpass-operator.git

# first time only
dep init

# add dependencies
dep ensure -add github.com/USER/DEP1 github.com/USER/DEP2
# example
dep ensure -add github.com/spf13/cobra
dep ensure -add github.com/codeskyblue/go-sh

# verify and update all dependencies
dep status
dep check
dep ensure -update

# resolve dependencies
dep ensure

# init cli
cobra init . --pkg-name lastpass-operator

# run
go run main.go

# compile
go build $GOPATH/src/github.com/edgelevel/lastpass-operator

# compile and build executable
go install $GOPATH/src/github.com/edgelevel/lastpass-operator

# test
go test $GOPATH/src/github.com/edgelevel/lastpass-operator
```
