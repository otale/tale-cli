package cmds

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	statusURL = "https://otale.github.io/status/version.json"
	pidFile   = "tale.pid"
)

// StartAction 启动 tale 博客
func StartAction() error {
	dat, err := ioutil.ReadFile(pidFile)
	if err == nil {
		pid := strings.TrimSpace(string(dat))
		if pid != "" {
			pidInt, err := strconv.Atoi(strings.TrimSpace(string(dat)))
			if err != nil {
				return err
			}
			_, err = os.FindProcess(pidInt)
			if err != nil {
				return nil
			}
			log.Println("博客已经启动成功，请停止后重启.")
			return nil
		}
	}

	os.Remove(pidFile)
	shell := "nohup java -Xms256m -Xmx256m -Dfile.encoding=UTF-8 -jar tale-letast.jar > /dev/null 2>&1 & echo $! > " + pidFile
	_, _, _, err = StartCmd(shell)
	if err != nil {
		return err
	}
	log.Println("博客程序已经启动成功，可使用 log 命令查看日志")
	return err
}

// StopAction 停止 tale 博客
func StopAction() error {
	dat, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Println("博客程序已经停止")
		return nil
	}
	log.Println("pid:", strings.TrimSpace(string(dat)))

	pid, err := strconv.Atoi(strings.TrimSpace(string(dat)))
	if err != nil {
		return err
	}
	KillPID(pid)
	err = os.Remove(pidFile)
	if err != nil {
		return err
	}
	log.Println("博客程序已经停止")
	return nil
}

// RestartAction 重启 tale 博客
func RestartAction() error {
	err := StopAction()
	if err == nil {
		StartAction()
	}
	return err
}

// StatusAction 查看博客运行状态
func StatusAction() error {
	dat, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Println("博客已经停止运行")
		return nil
	}
	pid := strings.TrimSpace(string(dat))
	if pid == "" {
		log.Println("博客已经停止运行")
		return nil
	}
	pidInt, err := strconv.Atoi(strings.TrimSpace(string(dat)))
	if err != nil {
		return err
	}
	_, err = os.FindProcess(pidInt)
	if err != nil {
		return nil
	}
	log.Println("Tale 博客正在运行")
	return nil
}

// LogAction 输出日志
func LogAction() error {
	_, stdout, stderr, err := StartCmd("tail -f logs/tale.log")
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	select {}
}

// UpgradeAction 升级博客
func UpgradeAction() error {
	return nil
}

// BackupAction 备份博客，SQL和当前全部状态
func BackupAction() error {
	return nil
}
