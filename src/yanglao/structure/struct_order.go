package structure

import (
	"time"
)

type Order struct {
	Idx           string    `orm:"pk"` //工单编号
	Name          string    //老人名字
	Phone         string    //联系电话
	Waiter        string    //服务人员
	WaiterPhne    string    //服务人员联系电话
	Service       string    //服务项目
	ServiceTime   time.Time //预定服务时间
	Addr          string    //地址
	Community     string    //所属社区
	OrderStatus   string    //工单状态
	PaymentStatus string    //支付状态
	Handler       string    //下单人
	HandlerTime   time.Time //下单时间
	Remarks       string    `orm:"type(text)"` //备注
}
