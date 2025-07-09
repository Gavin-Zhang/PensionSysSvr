package static

import (
	"fmt"
	"io/ioutil"

	"yanglao/utils"

	us "yanglao/utils/unmarshal"
)

// 服务器类表管理类
var Server *Static_Server = &Static_Server{}

// ServerData 服务器信息
type ServerData struct {
	IP   string // 地址
	Port uint32 // 端口
}

type Static_Server struct {
	ServerIp   string `json:"ServerIp"`
	ServerPort uint32 `json:"ServerPort"`
}

func (p *Static_Server) Init(path string, output bool) {
	err := p.loadLocalInfo(path)

	if output {
		utils.OutputInfo("Static_Server", err)
	}
}

func (p *Static_Server) loadLocalInfo(path string) error {
	reader, err := ioutil.ReadFile(path)
	if err == nil {
		err = us.JSONUnmarshal(string(reader[:]), p)
	}
	return err
}

func (p *Static_Server) GetLocalAddr() string {
	return fmt.Sprintf("%s:%d", p.ServerIp, p.ServerPort)
}

func (p *Static_Server) GetLocalInfo() ServerData {
	return ServerData{
		IP:   p.ServerIp,
		Port: p.ServerPort}
}
