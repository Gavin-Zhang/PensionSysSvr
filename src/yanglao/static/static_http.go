package static

import (
	"io/ioutil"

	"yanglao/utils"

	us "yanglao/utils/unmarshal"
)

// 服务器类表管理类
var HttpConfig *Static_Http = &Static_Http{}

type static_http_info struct {
	Domain      string `json:"domain"`
	AllowOrigin string `json:"allow_origin"`
}

type Static_Http struct {
	CookieName string             `json:"cookie_name"`
	CookieLife int64              `json:"cookie_life"`
	Info       []static_http_info `json:"info"`

	HCCPort   string `json:"hcc_port"`
	ESSPort   string `json:"ees_port"`
	JLYPort   string `json:"jly_port"`
	StorePort string `json:"store_port"`
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
