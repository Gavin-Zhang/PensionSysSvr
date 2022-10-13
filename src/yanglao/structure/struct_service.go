package structure

type ServiceClass struct {
	Class string `orm:"pk"`
}

type Service struct {
	Id    string `orm:"pk"`
	Class string
	Name  string
	Price string
}
