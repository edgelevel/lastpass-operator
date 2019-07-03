# lastpass-operator

**Resources**

* [lastpass-cli](https://github.com/lastpass/lastpass-cli)
* [A complete guide to Kubernetes Operator SDK](https://banzaicloud.com/blog/operator-sdk)
* [operator-sdk](https://github.com/operator-framework/operator-sdk)

```
echo PWD | LPASS_DISABLE_PINENTRY=1 lpass login --trust EMAIL
lpass ls
lpass show example-pwd --json
lpass logout --force
```

```
# build image
docker build -t niqdev/lastpass-cli .

# temporary container
docker run --rm --name lastpass-cli  niqdev/lastpass-cli

# access container
docker exec -it lastpass-cli bash
lpass --version

# execute command inline
docker run --rm -it niqdev/lastpass-cli lpass --version
```

## Setup

Documentation
* [go](https://golang.org/doc)
* [dep](https://golang.github.io/dep/docs/introduction.html)

Install `go` and `dep`
```bash
# ubuntu
sudo snap install --classic go
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# TODO operator-sdk

# macos
brew install go
brew install dep
brew install operator-sdk

# cli
go get github.com/spf13/cobra/cobra
```

Setup a [workspace](https://golang.org/doc/code.html#Workspaces)
```bash
# add to .bashrc or .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:$(go env GOPATH)/bin

# clone sources
mkdir -p $GOPATH/src/github.com/niqdev && cd $_
git clone git@github.com:niqdev/lastpass-operator.git
```

Development
```bash
# first time only
dep init

# add dependencies
dep ensure -add github.com/USER/DEP1 github.com/USER/DEP2
# example
dep ensure -add github.com/spf13/cobra

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
go build $GOPATH/src/github.com/niqdev/lastpass-operator

# compile and build executable
go install $GOPATH/src/github.com/niqdev/lastpass-operator

# test
go test $GOPATH/src/github.com/niqdev/lastpass-operator
```
