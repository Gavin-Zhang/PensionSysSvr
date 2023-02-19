package structure

type ServiceClass struct {
	Class string `orm:"pk"`
}

type Service struct {
	Id           int `orm:"pk;auto"`
	Class        string
	Name         string
	ExpensePrice string `orm:"description(自费价格)"`
	SubsidyPrice string `orm:"description(补贴价格)"`
}
