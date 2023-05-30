package global

var (
	SUCCESS = Error{Code: 10000, Message: "操作成功"}
	FAIL    = Error{Code: 99999, Message: "操作失败"}
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetCode() int {
	return e.Code
}

func (e Error) GetMessage() string {
	return e.Message
}
