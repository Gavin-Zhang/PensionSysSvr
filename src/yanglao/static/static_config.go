package static

import (
	"io/ioutil"

	"yanglao/utils"

	us "yanglao/utils/unmarshal"
)

var MyConfig *Static_Config = &Static_Config{}

type Static_Config struct {
	MinRoleListCount uint32 `json:"MinRoleListCount"`
}

func (p *Static_Config) Init(path string, output bool) {
	err := p.loadLocalInfo(path)

	if output {
		utils.OutputInfo("Static_Config", err)
	}
}

func (p *Static_Config) loadLocalInfo(path string) error {
	reader, err := ioutil.ReadFile(path)
	if err == nil {
		err = us.JSONUnmarshal(string(reader[:]), p)
	}
	return err
}
