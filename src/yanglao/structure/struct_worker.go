package structure

type WorkerClass struct {
	Class string `orm:"pk"`
}

type Worker struct {
	Idx     int `orm:"pk;auto"`
	Name    string
	ChinaId string
	Class   string
	Phone   string
}
