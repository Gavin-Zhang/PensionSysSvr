package structure

import (
	"time"
)

const (
	Indexs_Role_Index string = "roleIdx"
)

type Indexs struct {
	Key   string `orm:"pk"`
	Value uint32 `orm:"default(1)"`
}

type Client struct {
	Idx          string    `orm:"pk;description(服务对象编号)"` // EES_0000001
	Name         string    `orm:"description(服务对象姓名)"`
	ChinaId      string    `orm:"description(服务对象身份证)"`
	Phone        string    `orm:"description(服务对象电话号)"`
	Addr         string    `orm:"description(地址)"`
	Contacts     string    `orm:"type(text);description(联系人)"`
	SlowIll      string    `orm:"description(慢病)"`
	Healthy      string    `orm:"type(text);description(健康描述)"`
	Remarks      string    `orm:"type(text);description(备注)"`
	Handler      string    `orm:"description(登记人)"`
	Avatar       string    `orm:"null;description(头像)"`
	RegisterTime time.Time `orm:"auto_now_add;type(date);description(登记时间)"`
}
