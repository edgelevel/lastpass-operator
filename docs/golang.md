## golang

* [go](https://golang.org/doc) documentation
* [How To Build and Install Go Programs](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs)
* [Convert JSON into a Go type definition](https://mholt.github.io/json-to-go)

```bash
# verify env
go env | grep GOPATH

# setup workspace
mkdir -p $HOME/go

# add to .bashrc or .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:$(go env GOPATH)/bin
```

```bash
# download source
mkdir -p $GOPATH/src/github.com/edgelevel && cd $_
git clone git@github.com:edgelevel/lastpass-operator.git

# resolve dependencies
go mod tidy

# add dependencies
go get github.com/USER/DEP1 github.com/USER/DEP2
# example
go get github.com/spf13/cobra
go get github.com/codeskyblue/go-sh

# verify dependencies
go list -m all

# upgrade all dependencies to the latest or minor patch release
go get -t -u ./...

# init cli
cobra init . --pkg-name lastpass-operator

cd $GOPATH/src/github.com/edgelevel/lastpass-operator

# run
go run main.go

# compile
go build

# compile and build executable
go install

# test
go test ./...
```
