## LastPass

* [LastPass](https://www.lastpass.com)
* [lastpass-cli](https://github.com/lastpass/lastpass-cli)
* `lpass` [examples](lpass-example.txt)

CLI
```bash
# login
echo PASSWORD | LPASS_DISABLE_PINENTRY=1 lpass login --trust USERNAME

# list
lpass ls

# retrieve passwords
lpass show <GROUP>/<NAME> --json --expand-multi

# logout
lpass logout --force
```

Docker Alpine (8.07MB)
```bash
# build image
docker build -t niqdev/lastpass-cli -f example/lastpass-alpine .

# temporary container
docker run --rm -it niqdev/lastpass-cli /bin/sh
echo PASSWORD | LPASS_DISABLE_PINENTRY=1 lpass login --trust USERNAME
echo PASSWORD | lpass show <GROUP>/<NAME> --json --expand-multi
```

Docker Ubuntu (673MB)
```bash
# build image
docker build -t niqdev/lastpass-cli -f example/lastpass-ubuntu .

# temporary container
docker run --rm --name lastpass-cli niqdev/lastpass-cli

# access container
docker exec -it lastpass-cli bash
lpass --version

# execute command inline
docker run --rm -it niqdev/lastpass-cli lpass --version
```
