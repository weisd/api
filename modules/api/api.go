package api

type Respons struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResOk(data interface{}) Respons {
	return Respons{Code: 200, Status: "ok", Message: "ok", Data: data}
}

func ResErr(code int, msg string) Respons {
	return Respons{Code: code, Status: "err", Message: msg, Data: nil}
}
