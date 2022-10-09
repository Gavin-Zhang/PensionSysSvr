package structure

import (
	"time"
)

const (
	Indexs_Role_Index string = "roleIdx"
)

type Indexs struct {
	Key   string `orm:"pk"`
	Value uint32
}

type Client struct {
	Idx          string `orm:"pk"` // WLB_0000001
	Name         string
	ChinaId      string
	Phone        string
	Addr         string
	Community    string
	Type         string
	Contacts     string `orm:"type(text)"`
	SlowIll      string
	Healthy      string `orm:"type(text)"`
	Remarks      string `orm:"type(text)"`
	Handler      string
	RegisterTime time.Time
	InDbTime     time.Time
}

type User struct {
	Account  string `orm:"pk;index"`
	PassWord string
	UserName string
	Phone    string
}
