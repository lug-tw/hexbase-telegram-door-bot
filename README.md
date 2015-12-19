# HEXBASE telegram door bot

A [telegram](https://telegram.org/) bot for the roll-up door of [HexBase](https://github.com/lug-tw/HexBase).

## Available commands

|command|  note         |
|-------|---------------|
|`/ping`|life check     |
|`/up`  |scroll up      |
|`/down`|scroll down    |
|`/stop`|stop scrolling |


## Development steps

```shell
# hack hack hack
vim bot.go


$GOPATH/bin/goimports -w bot.go
gofmt -w bot.go

# cross compile to Raspberry Pi architecture
GOARCH=arm go build -a

# install
scp bot rspi:~/dev/door/
```
