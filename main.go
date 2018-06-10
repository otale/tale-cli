package main

import (
	"fmt"
	"os"

	"github.com/otale/tale-cli/cmds"
	"github.com/robmerrell/comandante"
)

const banner = `
Tale 博客程序帮助工具

Github: https://github.com/otale/tale`

func main() {
	bin := comandante.New("tale-cli", banner)

	// list command
	startCmd := comandante.NewCommand("start", "启动 Tale 博客", cmds.StartAction)
	bin.RegisterCommand(startCmd)

	// list command
	stopCmd := comandante.NewCommand("stop", "停止 Tale 博客", cmds.StopAction)
	bin.RegisterCommand(stopCmd)

	// list command
	restartCmd := comandante.NewCommand("restart", "重启 Tale 博客", cmds.RestartAction)
	bin.RegisterCommand(restartCmd)

	logCmd := comandante.NewCommand("log", "查看 Tale 博客日志", cmds.LogAction)
	bin.RegisterCommand(logCmd)

	upgradeCmd := comandante.NewCommand("upgrade", "升级 Tale 博客", cmds.UpgradeAction)
	bin.RegisterCommand(upgradeCmd)

	backupCmd := comandante.NewCommand("backup", "备份 Tale 博客", cmds.BackupAction)
	bin.RegisterCommand(backupCmd)

	// run the commands
	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
