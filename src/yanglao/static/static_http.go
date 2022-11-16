package static

import (
	"io/ioutil"

	"utils"

	us "utils/unmarshal"
)

// 服务器类表管理类
var HttpConfig *Static_Http = &Static_Http{}

type Static_Http struct {
	CookieName  string `json:"cookie_name"`
	CookieLife  int64  `json:"cookie_life"`
	Domain      string `json:"domain"`
	AllowOrigin string `json:"allow_origin"`
}

func (p *Static_Http) Init(path string, output bool) {
	err := p.loadLocalInfo(path)

	if output {
		utils.OutputInfo("Static_Http", err)
	}
}

func (p *Static_Http) loadLocalInfo(path string) error {
	reader, err := ioutil.ReadFile(path)
	if err == nil {
		err = us.JSONUnmarshal(string(reader[:]), p)
	}
	return err
}
