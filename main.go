/*
 * @Author: zzz
 * @Date: 2021-06-08 12:04:38
 * @LastEditTime: 2021-06-09 08:20:27
 * @LastEditors: zzz
 * @Description: 提供Windows Server API
 * @FilePath: \go-windows-server-api\main.go
 */

package main

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/kkzzhizhou/go-windows-server-api/hello"
)

type powerShell struct {
	powerShell string // 定义powershell
}

func new() *powerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &powerShell{
		powerShell: ps,
	}
}

func (p *powerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.GET("favicon.ico", func(c *gin.Context) {
		file, _ := favicon.ReadFile("icon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	// Ping test
	r.GET("/flushdns", func(c *gin.Context) {
		pwsh := new()
		stdout, stderr, err := pwsh.execute("[Console]::OutputEncoding = [Text.Encoding]::UTF8; Clear-DnsServerCache -Force")
		if err != nil {
			c.JSON(404, gin.H{"status": "flush failed", "reason": stderr})
		} else {
			fmt.Println(stdout)
			c.JSON(200, gin.H{"status": "flush finish", "content": stdout})
		}
	})

	return r
}

//go:embed icon.ico
var favicon embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)
	c := exec.Command("cmd", "/C", "Title", "Windows Server API")
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(hello.Greet())
	fmt.Println("作者: zzz")
	fmt.Println("说明: 在Windows Server上提供一些API接口，例如刷新DNS缓存等")
	fmt.Println("使用: http://ip:5000/flushdns // 刷新DNS缓存")
	r := setupRouter()
	r.Run(":5000")
}
