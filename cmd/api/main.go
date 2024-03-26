package main

import (
	"log"

	"github.com/gihyodocker/taskapp/pkg/app/api/cmd/config"
	"github.com/gihyodocker/taskapp/pkg/app/api/cmd/server"
	"github.com/gihyodocker/taskapp/pkg/cli"
)

func main() {
	// コマンドラインアプリケーションのインスタンスを作成
	c := cli.NewCLI("taskapp-api", "The API application of taskapp")
	//　サブコマンドの定義
	c.AddCommands(
		//　APIサーバの起動コマンド
		server.NewCommand(),
		config.NewCommand(),
	)
	// コマンドラインアプリケーションを実行
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
