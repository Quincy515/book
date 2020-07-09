`go run main.go help` 查看帮助

```shell
Usage:
   [command]

Available Commands:
  help        Help about any command
  sql         sql 转化和处理
  word        单词格式转换

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```

word

`go run main.go word -s custer_go -m 3`

输出

```shell
输出结果: CusterGo
```

sql

`go run main.go sql struct --username=? --password=? --db=? --table=?`
