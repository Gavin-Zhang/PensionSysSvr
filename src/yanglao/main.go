package main

import (
	"yanglao/gonet"
	_ "yanglao/service"
	_ "yanglao/service/db"
	_ "yanglao/service/http"

	"yanglao/single"
	"yanglao/static"

	"github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defer seelog.Flush()
	redirect()
	initSeelog("")

	static.Init()
	single.Init()
	gonet.Run("Mainsvr", 10, true)
}
