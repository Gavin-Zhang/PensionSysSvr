package structure

import (
	"time"
)

type Record struct {
	Idx       string    `orm:"pk;description(工单编号)"`
	ClientIdx string    `orm:"description(老人编号)"`
	Name      string    `orm:"description(老人名字)"`
	Phone     string    `orm:"description(联系电话)"`
	Addr      string    `orm:"description(地址)"`
	ChinaId   string    `orm:"description(身份证)"`
	Service   string    `orm:"description(服务描述)"`
	Remarks   string    `orm:"type(text);null;description(备注)"`
	Created   time.Time `orm:"auto_now_add;type(date);"`
	BeginTime time.Time `orm:"null;description(服务开始时间)"`
	EndTime   time.Time `orm:"null;description(服务结束时间)"`
	Handler   string    `orm:"null;description(处理人)"`
}

type RecordWorker struct {
	Idx       int64  `orm:"pk;auto;description(编号)"`
	RecordIdx string `orm:"description(工单编号)"`
	Name      string `orm:"description(服务人员名字)" json:"name"`
	Phone     string `orm:"description(服务人员电话)" json:"phone"`
	ChinaId   string `orm:"description(服务人员身份证)" json:"chinaid"`
	Class     string `orm:"description(服务人员类型)" json:"class"`
}

type RecordPhoto struct {
	Name       string
	Path       string
	Size       int64
	CreateTime time.Time
}
