package static

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/cihub/seelog"
)

// Init 初始包内各个管理类
func Init() {
	defer func() {
		if err := recover(); err != nil {
			filepath := fmt.Sprintf("staic_stack_%s", time.Now().Format("2006-01-02"))
			logfile, openErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if openErr != nil {
				panic(openErr)
			}
			defer logfile.Close()
			logger := log.New(logfile, "", log.Ldate|log.Ltime)

			stack := debug.Stack()
			logger.Println("===============================================================")
			logger.Println(fmt.Sprint(err))
			logger.Println(string(stack))

			os.Exit(0)
		}
	}()

	seelog.Info("===================================================================")
	Server.Init("data/json/server", true)
	Db.Init("data/json/db", true)
	MyConfig.Init("data/json/config", true)
	HttpConfig.Init("data/json/http", true)
	seelog.Info("===================================================================")
}
