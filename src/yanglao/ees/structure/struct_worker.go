package structure

type WorkerClass struct {
	Class string `orm:"pk"`
}

type Worker struct {
	Idx     int `orm:"pk;auto"`
	Name    string
	ChinaId string
	Phone   string
	Class   string
}
