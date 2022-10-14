package structure

type ServiceClass struct {
	Class string `orm:"pk"`
}

type Service struct {
	Id    int `orm:"pk;auto"`
	Class string
	Name  string
	Price string
}
