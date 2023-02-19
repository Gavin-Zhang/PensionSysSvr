package structure

type SubsidyRecord struct {
	Idx       int    `orm:"pk;auto"`
	Date      string `orm:"description(日期)"`
	OrderIdx  string `orm:"description(工单编号)"`
	ClientIdx string `orm:"description(老人编号)"`
	Name      string `orm:"description(老人名字)"`
	Duration  int    `orm:"description(服务时长-分钟)"`
}

type SubsidyMonthRecord struct {
	Idx      int    `orm:"pk;auto"`
	Date     string `orm:"description(日期yyyymm)"`
	Duration int    `orm:"description(服务时长-分钟)"`
}
