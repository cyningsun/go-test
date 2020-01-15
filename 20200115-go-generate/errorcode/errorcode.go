package user

type ErrorCode struct {
	Code uint32 `json:"err_code"`
	Msg  string `json:"err_msg"`
}

//go:generate codemap -type ErrorCode

var (
	// [0, 999] 通用
	Ok      = ErrorCode{Code: 0, Msg: "success"}
	Unknown = ErrorCode{Code: 1, Msg: "unknown"}

	// 已废弃
	ErrorParams           = ErrorCode{Code: 100, Msg: "params invalid"}
	ErrorDuplicateRequest = ErrorCode{Code: 101, Msg: "duplicate request"}
	ErrorDataBaseError    = ErrorCode{Code: 102, Msg: "database error"}
	ErrorDataNotFound     = ErrorCode{Code: 103, Msg: "data not found"}
	ErrorUsersig          = ErrorCode{Code: 104, Msg: "usersig cal error"}
	ErrorRPC              = ErrorCode{Code: 105, Msg: "rpc call error"}
)
