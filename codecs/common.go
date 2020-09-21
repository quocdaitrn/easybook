package codecs

// Common error codes
const (
	Success = iota + 100000
	Fail
	InvalidParams
	FailedCreate
	FailedUpdate
	FailedDelete
	RecordNotExist
	SystemError
)

var code2text = map[int]string{
	Success:        "success",
	Fail:           "fail",
	InvalidParams:  "invalid parametes",
	FailedCreate:   "create record failed",
	FailedUpdate:   "update record failed",
	FailedDelete:   "delete record failed",
	RecordNotExist: "record doesn't exist",
	SystemError:    "system error",
}

// CommonResponse define
type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SetCode function
func (c *CommonResponse) SetCode(code int) {
	c.Code = code
	c.Message = code2text[code]
}
