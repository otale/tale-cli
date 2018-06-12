package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/mitchellh/ioprogress"
)

func DownloadFile(url string, dest string) {

	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

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
}

func main() {

	DownloadFile("https://wordpress.org/wordpress-4.4.2.zip", "./")

}
