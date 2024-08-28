## LastPass

* [Website](https://www.lastpass.com)
* [lastpass-cli](https://github.com/lastpass/lastpass-cli)
* `lpass` [examples](../example/lpass-examples.txt)

CLI
```bash
# login
echo <PASSWORD> | LPASS_DISABLE_PINENTRY=1 lpass login --trust <USERNAME>

# list
lpass ls

# retrieve passwords
lpass show <GROUP>/<NAME> --json --expand-multi

# logout
lpass logout --force
```

Docker Alpine (20.1MB)
```bash
# temporary base container
docker run --rm -it alpine /bin/sh

# build image
docker build -t edgelevel/lastpass-cli -f example/lastpass-alpine .

# temporary container
docker run --rm -it edgelevel/lastpass-cli /bin/sh
echo <PASSWORD> | LPASS_DISABLE_PINENTRY=1 lpass login --trust <USERNAME>
echo <PASSWORD> | lpass show <GROUP>/<NAME> --json --expand-multi
```

Docker Ubuntu (192MB)
```bash
# build image
docker build -t edgelevel/lastpass-cli -f example/lastpass-ubuntu .

# temporary container
docker run --rm --name lastpass-cli edgelevel/lastpass-cli

# access container
docker exec -it lastpass-cli bash
lpass --version

# execute command inline
docker run --rm -it edgelevel/lastpass-cli lpass --version
```
