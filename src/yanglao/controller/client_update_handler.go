package controller

//func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
//	cors(&w, r)

//	if checkSession(w, r) == nil {
//		seelog.Error("UpdateClientHandler not find player by cookie")
//		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
//		return
//	}

//	if !checkNotEmptyParams(r, []string{"idx"}) {
//		seelog.Error("UpdateClientHandler checkNotEmptyParams fail")
//		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
//		return
//	}

//	client := structure.Client{
//		Idx:       r.FormValue("idx"),
//		ChinaId:   r.FormValue("chinaid"),
//		Name:      r.FormValue("name"),
//		Phone:     r.FormValue("phone"),
//		Addr:      r.FormValue("addr"),
//		Community: r.FormValue("ascription"),
//		Type:      r.FormValue("type"),
//		Healthy:   r.FormValue("healthdescription"),
//		Remarks:   r.FormValue("other")}

//	temp, _ := json.Marshal(r.FormValue("contacts"))
//	client.Contacts = string(temp)
//	temp, _ = json.Marshal(r.FormValue("slow"))
//	client.SlowIll = string(temp)

//}
