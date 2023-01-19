package main

import (
	_ "yanglao/ees/service/db"
	_ "yanglao/ees/service/http"
	_ "yanglao/hcc/service/db"
	_ "yanglao/hcc/service/http"
	_ "yanglao/service"

	"yanglao/gonet"
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
