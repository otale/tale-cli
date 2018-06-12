# tale-cli

tale命令行帮助程序 🌟

[![License](http://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/otale/tale-cli/master/LICENSE)
[![Travis branch](https://img.shields.io/travis/otale/tale-cli/master.svg)](https://travis-ci.org/otale/tale-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/otale/tale-cli)](https://goreportcard.com/report/github.com/otale/tale-cli)
[![GoDoc](https://godoc.org/github.com/otale/tale-cli?status.svg)](https://godoc.org/github.com/otale/tale-cli)

## 特性

- 支持启动、停止、升级 Tale 博客
- 支持 Linux、MacOSX 系统
- 支持旧版本迁移
- 支持博客备份

## 预览

[![tale-cli](https://i.loli.net/2017/09/10/59b5241331c47.png)](https://asciinema.org/a/137112)

## 安装

**linux_64位**

```bash
cd tale
wget http://static.biezhi.me/tale-cli
chmod +x tale-cli
```

## 帮助

```bash
# tale-cli

NAME:
   tale - tale的命令行帮助程序

USAGE:
   tale-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
   https://github.com/biezhi <biezhi.me@gmail.com>

COMMANDS:
     start    启动tale
     stop     停止当前tale实例
     reload   重新启动当前tale实例
     log      查看当前tale日志
     status   查看当前tale状态
     upgrade  升级当前的tale版本
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## 操作指令

```bash
./tale-cli start  : 启动tale
./tale-cli stop   : 停止tale实例
./tale-cli reload : 重启tale实例
./tale-cli log    : 查看当前tale日志
./tale-cli status : 查看当前tale运行状态
./tale-cli upgrade: 升级tale版本，会自动帮你备份
```

> 备份的文件在当前目录下形如 `tale_backup_201709109281.zip`

