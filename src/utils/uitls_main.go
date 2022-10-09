package utils

import (
	"fmt"
	"log"
)

// CheckError 检测错误
func CheckError(err error) {
	if err != nil {
		fmt.Println("=========================")
		fmt.Println(err)
		fmt.Println("=========================")
		panic(err)
	}
}

// OutputInfo 加载服务及静态全局信息时输出信息
func OutputInfo(mgrName string, err error) {
	initStr := fmt.Sprintf("Init %s", mgrName)
	//fmt.Print(initStr)
	formatStr := "%" + fmt.Sprintf("%d", 50-len(initStr)) + "s"
	if err == nil {
		formatStr = initStr + formatStr
		log.Println(fmt.Sprintf(formatStr, "[v]"))
	} else {
		log.Println(fmt.Sprintf(formatStr, "[x]"))
	}
	CheckError(err)
}
