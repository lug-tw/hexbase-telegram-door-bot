# HEXBASE telegram door bot

A [telegram](https://telegram.org/) bot for the roll-up door of [HexBase](https://github.com/lug-tw/HexBase).

## Available commands

|command|  note         |
|-------|---------------|
|`/ping`|life check     |
|`/up`  |scroll up      |
|`/down`|scroll down    |
|`/stop`|stop scrolling |


# Contribute

```shell
# setup go env and download goimports
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
which goimports || go install golang.org/x/tools/cmd/goimports

# fetch latest code
go get -ud github.com/lug-tw/hexbase-telegram-door-bot

# hack, hack, hack
cd $GOPATH/src/github.com/lug-tw/hexbase-telegram-door-bot
your_favorite_editor bot.go

# reformat the code, build and upload
goimports -w bot.go

# cross compile to Raspberry Pi architecture
GOARCH=arm go build -a

# install
scp hexbase-telegram-door-bot rspi:~/dev/door/
```

## Hardware

- Our Raspberry Pi and the bread board

![img](https://i.imgur.com/yo0Fa0L.jpg)

![img](https://i.imgur.com/xrI2j9K.jpg)
