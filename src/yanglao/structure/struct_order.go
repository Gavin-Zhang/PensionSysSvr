package structure

import (
	"time"
)

// 消费类型
type ConsumptionType struct {
	Id   int `orm:"pk;auto"`
	Type string
}

// 支付类型
type PaymentType struct {
	Id   int `orm:"pk;auto"`
	Type string
}

type Order struct {
	Idx             string    `orm:"pk"` //工单编号
	Name            string    //老人名字
	Phone           string    //联系电话
	ChinaId         string    // 身份证
	Waiter          string    `orm:"null"` //服务人员
	WaiterPhone     string    `orm:"null"` //服务人员联系电话
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
	AssignTime      time.Time `orm:"null;description(分派时间)"`
	BeginTime       time.Time `orm:"null;description(服务开始时间)"`
	EndTime         time.Time `orm:"null;description(服务结束时间)"`
	PaymentTime     time.Time `orm:"null;type(date);description(支付日期)"`
	FinishTime      time.Time `orm:"null;description(订单结束时间)"`
}

// 工单状态
const (
	ORDER_STATUS_WAIT_ASSIGN  = "待分派"
	ORDER_STATUS_WAIT_SERVICE = "待服务"
	ORDER_STATUS_WAIT_PAY     = "待支付"
	ORDER_STATUS_FINISHED     = "已完成"
)

// 支付状态
const (
	ORDER_PAY_STATUS_WAIT = "待支付"
	ORDER_PAY_STATUS_OVER = "已支付"
)
