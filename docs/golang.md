## golang

* [go](https://golang.org/doc) documentation
* [Convert JSON into a Go type definition](https://mholt.github.io/json-to-go/)

```bash
# download source
mkdir -p $GOPATH/src/github.com/edgelevel && cd $_
git clone git@github.com:edgelevel/lastpass-operator.git

# add dependencies
go get github.com/USER/DEP1 github.com/USER/DEP2
# example
go get github.com/spf13/cobra
go get github.com/codeskyblue/go-sh

# verify and update all dependencies
go get -u
go mod tidy

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
