# HEXBASE telegram door bot

# Development steps

```shell
vim bot.go
$GOPATH/bin/goimports -w bot.go
gofmt -w bot.go
GOARCH=arm go build -a
scp bot rspi:~/dev/door/
```
