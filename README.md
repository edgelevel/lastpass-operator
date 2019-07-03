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

# macos
brew install go
brew install dep
```

Setup a [workspace](https://golang.org/doc/code.html#Workspaces)
```bash
# add to .bashrc or .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:$(go env GOPATH)/bin

# clone sources
mkdir -p $GOPATH/src/github.com/niqdev && cd $_
git clone git@github.com:niqdev/lastpass-operator.git

# first time only
dep init
```

Development
```bash
# compile
go build $GOPATH/src/github.com/niqdev/lastpass-operator

# compile and build executable
go install $GOPATH/src/github.com/niqdev/lastpass-operator

# test
go test $GOPATH/src/github.com/niqdev/lastpass-operator
```
