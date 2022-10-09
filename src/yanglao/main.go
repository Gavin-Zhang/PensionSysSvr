package main

import (
	"gonet"
	_ "yanglao/service"
	_ "yanglao/service/db"
	_ "yanglao/single"

	"yanglao/static"

	"github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defer seelog.Flush()
	redirect()
	initSeelog("")

	static.Init()
	gonet.Run("Mainsvr", 10, true)
}
