package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"os/exec"
	"log"
	"strings"
	"strconv"
)

func main() {

	app := cli.NewApp()
	app.Name = "tale"
	app.Usage = "tale的命令行帮助程序"
	app.Author = "https://github.com/biezhi"
	app.Email = "biezhi.me@gmail.com"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "启动tale",
			Action: doStart,
		},
		{
			Name:   "stop",
			Usage:  "停止当前tale实例",
			Action: doStop,
		},
		{
			Name:  "reload",
			Usage: "重新启动当前tale实例",
			Action: func(ctx *cli.Context) {
				doStop(ctx)
				doStart(ctx)
			},
		},
		{
			Name:  "log",
			Usage: "查看当前tale日志",
			Action: func(ctx *cli.Context) {
				tailLog()
			},
		},
		{
			Name:  "status",
			Usage: "查看当前tale状态",
			Action: func(ctx *cli.Context) {
				pid := findPid()
				if pid < 0 {
					fmt.Println("Tale 实例没有运行.")
				} else {
					fmt.Printf("Tale start with pid: %d\n", pid)
				}
			},
		},
		{
			Name:  "upgrade",
			Usage: "升级当前的tale版本",
			Action: func(ctx *cli.Context) {

			},
		},
	}
	app.Run(os.Args)
	os.Exit(0)
}

func doStart(ctx *cli.Context) {
	pid := findPid()
	if pid > 0 {
		fmt.Println("Tale 已经启动.")
	} else {
		cmd := exec.Command("java", "-jar", "tale-1.3.0-alpha1.jar", "&")
		// 重定向标准输出到文件
		stdout, err := os.OpenFile("tale.log", os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalln(err)
		}
		defer stdout.Close()
		cmd.Stdout = stdout
		// 执行命令
		if err := cmd.Start(); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Tale 启动成功, 可以使用 tale log 命令查看日志.")
	}
}

func doStop(ctx *cli.Context) {
	pid := findPid()
	if pid > 0 {
		proc, err := os.FindProcess(pid)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("kill pid: %d\n", pid)
		proc.Kill()
		os.Remove("resources/tale.pid")
	}
}

//func findPid() int {
//	pidByte, _ := ioutil.ReadFile("resources/tale.pid")
//	if len(pidByte) == 0 {
//		return -1
//	}
//	pid := strings.TrimSuffix(string(pidByte), "\n")
//	intVal, _ := strconv.Atoi(pid)
//	return intVal
//}

func tailLog() {
	cmd := exec.Command("tail", "-f", "tale.log")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func findPid() int {
	pidByte, err := exec.Command("/bin/sh", "-c", `ps -eaf|grep "tale-1.3.0-alpha1.jar"|grep -v "grep"|awk '{print $2}'`).Output()
	if err != nil {
		log.Fatal(err)
		return -1
	}
	if len(pidByte) == 0 {
		return -1;
	}
	pid := string(pidByte)
	pid = strings.TrimSuffix(string(pidByte), "\n")
	if len(pid) == 0 {
		return -1
	}
	intVal, _ := strconv.Atoi(pid)
	return intVal
}
