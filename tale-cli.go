package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"os/exec"
	"log"
	"strings"
	"strconv"
	"time"
)

const (
	taleZipName     = "tale-least.zip"
	taleDownloadUrl = "http://static.biezhi.me/" + taleZipName
	jarFileName     = "tale-least.jar"
)

func main() {

	app := cli.NewApp()
	app.Name = "tale"
	app.Usage = "tale的命令行帮助程序"
	app.Author = "https://github.com/biezhi"
	app.Email = "biezhi.me@gmail.com"
	app.Version = "0.0.3"

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
			Name:   "upgrade",
			Usage:  "升级当前的tale版本",
			Action: doUpgrade,
		},
	}
	app.Run(os.Args)
	os.Exit(0)
}

// start tale instance
func doStart(ctx *cli.Context) {
	pid := findPid()
	if pid > 0 {
		fmt.Println("Tale 已经启动.")
	} else {
		// `nohup java -Xms128m -Xmx128m -Dfile.encoding=UTF-8 -jar tale-1.3.0-alpha1.jar >tale.log &`
		cmd := exec.Command("/bin/sh", "-c", `nohup java -Xms128m -Xmx128m -Dfile.encoding=UTF-8 -jar `+jarFileName+` >tale.log &`)
		cmd.Dir = "."
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
		fmt.Println("Tale 启动成功, 可以使用 ./tale-cli log 命令查看日志.")
	}
}

// stop tale instance
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

// tail -f tale.log
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

// 升级tale版本
func doUpgrade(ctx *cli.Context) {
	dir, err := os.Open("./resources")
	if err != nil {
		log.Fatal(err)
	}
	var files = []*os.File{dir}
	dest := "tale_backup_" + time.Now().Format("20060102150405") + ".zip"
	err = Compress(files, dest)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("备份成功.")
	fmt.Println("开始下载最新版tale安装包, 客官请稍等...")
	os.Remove(taleZipName)
	//// 下载tale.zip

	rand := "?t=_" + time.Now().Format("20060102150405")
	DownloadFile(taleDownloadUrl+rand, "./")

	Unzip(taleZipName+rand, "./")
	fmt.Println("解压完成")
	fmt.Println("正在升级...")

	// delete 除了 resources 目录下的所有
	// cd resources && delete 除了 app.properties、static、
	// templates/admin、templates/install.html、templates/comm
	RemoveContents("lib")

	os.Rename("./tale/lib", "./lib")
	os.Remove(jarFileName)

	os.Rename("./tale/"+jarFileName, "./"+jarFileName)

	RemoveContents("./resources/static")
	os.Rename("./tale/resources/static", "./resources/static")

	RemoveContents("./resources/templates/admin")
	os.Rename("./tale/resources/templates/admin", "./resources/templates/admin")

	RemoveContents("./resources/templates/comm")
	os.Rename("./tale/resources/templates/comm", "./resources/templates/comm")

	os.Remove("./resources/templates/install.html")
	os.Rename("./tale/resources/templates/install.html", "./resources/templates/install.html")

	RemoveContents("tale")
	os.Remove("tale")

	fmt.Println("Tale 升级成功, 请手动启动.")

}

// find tale.jar process id
func findPid() int {
	jarFileName := "tale.jar"
	pidByte, err := exec.Command("/bin/sh", "-c", `ps -eaf|grep "`+jarFileName+`"|grep -v "grep"|awk '{print $2}'`).Output()
	if err != nil {
		log.Fatal(err)
		return -1
	}
	if len(pidByte) == 0 {
		return -1
	}
	pid := string(pidByte)
	pid = strings.TrimSuffix(string(pidByte), "\n")
	if len(pid) == 0 {
		return -1
	}
	intVal, _ := strconv.Atoi(pid)
	return intVal
}
