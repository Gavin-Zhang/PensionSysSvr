package structure

import (
	"time"
)

type ConsumptionType struct {
	Type string `orm:"pk"`
}

type Order struct {
	Idx             string    `orm:"pk"` //工单编号
	Name            string    //老人名字
	Phone           string    //联系电话
	Waiter          string    `orm:"null"` //服务人员
	WaiterPhne      string    `orm:"null"` //服务人员联系电话
	Service         string    `orm:"description(服务项目)"`
	ConsumptionType string    `orm:"description(消费类型)"` //[自费，政府购买，积分，赠送...]
	Charge          int16     `orm:"description(服务费用)"`
	Fare            int16     `orm:"description(车费)"`
	HighRise        float32   `orm:"description(步梯高层费用)"`
	ServiceTime     time.Time //预定服务时间
	Addr            string    //地址
	Community       string    //所属社区
	OrderStatus     string    //工单状态
	PaymentStatus   string    //支付状态
	PaymentType     string    `orm:"null"` //支付方式
	Handler         string    //下单人
	Remarks         string    `orm:"type(text);null"` //备注
	Created         time.Time `orm:"auto_now_add"`
}

// 工单状态
const (
	ORDER_STATUS_CREATED  = "已下单"
	ORDER_STATUS_DISPENSE = "已分配"
	ORDER_STATUS_FINISHED = "已完成"
)

// 支付状态
const (
	ORDER_PAY_STATUS_WAIT = "待支付"
	ORDER_PAY_STATUS_OVER = "已支付"
)
