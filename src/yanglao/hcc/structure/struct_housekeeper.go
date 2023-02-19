package structure

type HouseKeeper struct {
	Id        int `orm:"pk;auto"`
	Name      string
	Phone     string `orm:"description(联系电话)"`
	Community string `orm:"description(所属社区/村部)"`
}
