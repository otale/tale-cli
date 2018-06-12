package cmds

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
