package cmds

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/mitchellh/ioprogress"
)

// GetRequestBody http get body
func GetRequestBody(reqURL string) []byte {
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

// KillPID kill 进程id
func KillPID(pid int) (int, error) {
	err := syscall.Kill(pid, syscall.SIGKILL)
	return pid, err
}

// StartCmd 启动命令
func StartCmd(cmd string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	var err error
	c := exec.Command("/bin/sh", "-c", cmd)
	c.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, nil, nil, err
	}
	return c, stdout, stderr, err
}

// RemoveDir 删除文件夹
func RemoveDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// DownloadFile 下载文件，带进度
func DownloadFile(url, dest string) string {
	file := path.Base(url)

	log.Printf("从 %s 下载 %s\n", url, file)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	out, err := os.Create(path.String())
	if err != nil {
		fmt.Println(path.String())
		panic(err)
	}

	defer out.Close()

	headResp, err := http.Head(url)
	if err != nil {
		panic(err)
	}

	defer headResp.Body.Close()
	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bar := ioprogress.DrawTextFormatBar(20)
	progressR := &ioprogress.Reader{
		Reader: resp.Body,
		Size:   int64(size),
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, func(progress, total int64) string {
			return fmt.Sprintf(
				"正在下载: %s %s",
				bar(progress, total),
				ioprogress.DrawTextFormatBytes(progress, total))
		}),
	}

	io.Copy(out, progressR)

	return file
}
