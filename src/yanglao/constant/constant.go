package constant

import (
	"errors"
)

const (
	LoginResult_Success uint16 = 0
	LoginResult_NotExit        = 1
)

// 注册结果码
const (
	RegisterResult_Success   uint16 = 0
	RegisterResult_RepeatID         = 1
	RegisterResult_CheckErr         = 2
	RegisterResult_InsertErr        = 3
)

const (
	ResponseCode_Success    int = 0
	ResponseCode_Fail           = 1000
	ResponseCode_CookieErr      = 1001
	ResponseCode_ProgramErr     = 1002
	ResponseCode_ParamErr       = 1003
	ResponseCode_Custom         = 1999
)

const (
	Power_ALL   uint32 = 0
	Power_HCC   uint32 = 1 << 0
	Power_EES   uint32 = 1 << 1
	Power_JLY   uint32 = 1 << 2
	Power_Store uint32 = 1 << 3
)

var (
	Error_RepeatID    = errors.New("身份证号重复")
	Error_Program     = errors.New("程序错误")
	Error_NoRows      = errors.New("没有找到对应记录")
	Error_NullPointer = errors.New("空指针")
)
