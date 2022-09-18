package apps

type UserData struct {
	Imsi       int    `json:"imsi"`
	Erab       int    `json:"erab"`
	Imei       int    `json:"imei"`
	Apn        string `json:"apn"`
	Pdn_type   string `json:"pdn_type"`
	State      string `json:"state"`
	Msg_id     int    `json:"msgId"`
	Ul_data_KB int    `json:"ulDataKB"`
	Dl_data_kB int    `json:"dlDataKB"`
}
