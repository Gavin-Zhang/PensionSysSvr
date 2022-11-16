package utils

import (
	"fmt"
	//"log"

	"github.com/cihub/seelog"
)

// CheckError 检测错误
func CheckError(err error) {
	if err != nil {
		fmt.Println("=========================")
		fmt.Println(err)
		fmt.Println("=========================")
		seelog.Error("=========================")
		seelog.Error(err)
		seelog.Error("=========================")
		panic(err)
	}
}

// OutputInfo 加载服务及静态全局信息时输出信息
func OutputInfo(mgrName string, err error) {
	initStr := fmt.Sprintf("Init %s", mgrName)
	formatStr := "%" + fmt.Sprintf("%d", 50-len(initStr)) + "s"
	if err == nil {
		formatStr = initStr + formatStr
		seelog.Info(fmt.Sprintf(formatStr, "[v]"))
	} else {
		seelog.Info(fmt.Sprintf(formatStr, "[x]"))
	}
	CheckError(err)
}
