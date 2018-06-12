package cmds

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/archiver"
)

const (
	statusURL   = "https://otale.github.io/status/version.json"
	pidFile     = "tale.pid"
	dbFile      = "resources/tale.db"
	versionFile = "version.txt"
)

// Version tale版本信息
type Version struct {
	LatestVersion string   `json:"latest_version"`
	PublishTime   string   `json:"publish_time"`
	Hash          string   `json:"hash"`
	ChangeLogs    []string `json:"change_logs"`
	DownloadURL   string   `json:"download_url"`
}

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
	updated := false
	var ver string
	if _, err := os.Stat(versionFile); err == nil {
		data, _ := ioutil.ReadFile(versionFile)
		ver = string(data)
	} else {
		if _, err := os.Stat(dbFile); err == nil {
			updated = true
		}
	}

	if updated {
		log.Println("修复 SQL")
		cmd1 := "sqlite3 " + dbFile + " \"update t_comments set type = 'comment', status = 'approved';\""
		cmd2 := "sqlite3 " + dbFile + " \"update t_contents set categories = '默认分类' where categories is null;\""
		cmd3 := "sqlite3 " + dbFile + " \"update t_contents set allow_feed = 1 where allow_feed is null;\""
		cmd4 := "sqlite3 " + dbFile + " \"insert into t_options(name,value,description) values ('allow_comment_audit', 'true', '评论需要审核');\""
		log.Println(cmd1)
		_, _, _, err := StartCmd(cmd1)
		if err != nil {
			log.Println("修复 SQL 失败")
			return err
		}
		log.Println(cmd2)
		_, _, _, err = StartCmd(cmd2)
		if err != nil {
			log.Println("修复 SQL 失败")
			return err
		}
		log.Println(cmd3)
		_, _, _, err = StartCmd(cmd3)
		if err != nil {
			log.Println("修复 SQL 失败")
			return err
		}
		log.Println(cmd4)
		_, _, _, err = StartCmd(cmd4)
		if err != nil {
			log.Println("修复 SQL 失败")
			return err
		}
	}

	var version Version
	body := GetRequestBody(statusURL)
	if err := json.Unmarshal(body, &version); err != nil {
		log.Println("解析JSON失败")
		return err
	}
	log.Println(version)

	if ver != "" {
		verInt, _ := strconv.Atoi(ver)
		leteatVer, _ := strconv.Atoi(version.PublishTime)
		if verInt >= leteatVer {
			log.Println("您的版本已经是最新，无需升级 :)")
			return nil
		}
	}

	log.Println("最新版本:", version.LatestVersion)

	// 下载新文件
	newZipName := DownloadFile(version.DownloadURL, "./")

	// 备份
	datetime := time.Now().Format("20060102_150405") + ".zip"
	err := archiver.Zip.Make(datetime, []string{"resources", "lib", "tale-least.jar"})
	if err != nil {
		log.Println("备份失败")
		return err
	}
	// 删除旧文件
	_ = os.Remove("tale-latest.jar")
	_ = os.Remove("tale-least.jar")
	RemoveDir("lib")
	RemoveDir("resources/static")
	RemoveDir("resources/templates")

	// 解压新文件
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	err = archiver.Zip.Open(newZipName, dir)
	if err != nil {
		log.Println("解压 " + newZipName + " 失败")
		return err
	}
	log.Println("升级完毕")
	_ = os.Remove(newZipName)
	return nil
}

// BackupAction 备份博客，SQL和当前全部状态
func BackupAction() error {
	// 备份 SQL
	_, _, _, err := StartCmd("sqlite3 tale.db .dump > tale.sql")
	if err != nil {
		return err
	}
	// 备份目录

	return nil
}
