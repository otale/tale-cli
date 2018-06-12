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

[![tale-cli](https://wx3.sinaimg.cn/large/0079jBTLgy1fs8x1zwwwgj31j60wa7co.jpg)](https://asciinema.org/a/186825)

## 安装

目前二进制文件存储在仓库 `bin` 目录，等发布版本后会提供一个在线下载地址。

## 帮助

```bash
» ./tale-cli

Tale 博客程序帮助工具

Github: https://github.com/otale/tale

Usage:
	tale-cli command [arguments]

Available commands:
backup   备份 Tale 博客
log      查看 Tale 博客日志
restart  重启 Tale 博客
start    启动 Tale 博客
status   查看 Tale 博客运行状态
stop     停止 Tale 博客
upgrade  升级 Tale 博客
```

> 备份的文件在当前目录下形如 `backup_201709109281.zip`

