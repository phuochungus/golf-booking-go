package response

const (
	ErrorCodeSuccess      = 20001
	ErrorCodeParamInvalid = 20003
	ErrorCodeUserExisted  = 20004
)

var msg = map[int]string{
	ErrorCodeSuccess:      "Success",
	ErrorCodeParamInvalid: "Invalid parameters",
	ErrorCodeUserExisted:  "User already exists",
}
