package constant

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
)
