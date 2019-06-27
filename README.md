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
```
