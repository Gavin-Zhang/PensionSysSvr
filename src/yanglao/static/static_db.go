package static

import (
	"io/ioutil"

	"bitgame/utils"

	us "bitgame/utils/unmarshal"
)

// 服务器类表管理类
var Db *Static_DB = &Static_DB{}

type Static_DB struct {
	Conn      string   `json:"db"`
	Indexs    []string `json:"index"`
	IndexHead string   `json:"index_head"`
}

func (p *Static_DB) Init(path string, output bool) {
	err := p.loadLocalInfo(path)

	if output {
		utils.OutputInfo("Static_DB", err)
	}
}

func (p *Static_DB) loadLocalInfo(path string) error {
	reader, err := ioutil.ReadFile(path)
	if err == nil {
		err = us.JSONUnmarshal(string(reader[:]), p)
	}
	return err
}
