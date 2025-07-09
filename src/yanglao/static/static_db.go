package static

import (
	"os"
	"yanglao/utils"

	us "yanglao/utils/unmarshal"
)

// 服务器类表管理类
var Db *Static_DB = &Static_DB{}

type StoreIdxsInfo struct {
	IdxName  string `json:"name"`
	YearFlag bool   `json:"year"`
}

type Static_DB struct {
	Conn      string   `json:"db"`
	Indexs    []string `json:"index"`
	IndexHead string   `json:"index_head"`
	OrderHead string   `json:"order_head"`

	EESConn      string   `json:"ees_db"`
	EESIndexs    []string `json:"ees_roleidx"`
	EESIndexHead string   `json:"ees_idxhead"`
	EESRecord    string   `json:"ess_record"`

	JLYConn   string   `json:"jly_db"`
	JLYIndexs []string `json:"jly_idxs"`

	StoreConn      string          `json:"store_db"`
	StoreClassRoot string          `json:"store_class_root"`
	StoreIdxs      []StoreIdxsInfo `json:"store_idx"`
	StoreIdxsMap   map[string]StoreIdxsInfo
	SupplierRoot   string `json:"supplier_root"`
}

func (p *Static_DB) Init(path string, output bool) {
	err := p.loadLocalInfo(path)

	if output {
		utils.OutputInfo("Static_DB", err)
	}

	p.buildStoreIdxsMap()
}

func (p *Static_DB) loadLocalInfo(path string) error {
	reader, err := os.ReadFile(path)
	if err == nil {
		err = us.JSONUnmarshal(string(reader[:]), p)
	}
	return err
}

func (p *Static_DB) buildStoreIdxsMap() {
	p.StoreIdxsMap = make(map[string]StoreIdxsInfo)
	for _, idx := range p.StoreIdxs {
		p.StoreIdxsMap[idx.IdxName] = idx
	}
}
