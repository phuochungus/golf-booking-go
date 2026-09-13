package response

const (
	ErrorCodeSuccess        = 0
	ErrorCodeBadRequest     = 40000
	ErrorCodeUnauthorized   = 40001
	ErrorCodeInternalServer = 50000
	ErrorCodeAdminExisted   = 20001
	ErrorCodeParamInvalid   = 20003
	ErrorCodeUserExisted    = 20004
)

var msg = map[int]string{
	ErrorCodeSuccess:      "Success",
	ErrorCodeParamInvalid: "Invalid parameters",
	ErrorCodeUserExisted:  "User already exists",
	ErrorCodeBadRequest:   "Bad request",
}
