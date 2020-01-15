package user

var ErrorCodeMap = map[uint32]string{
	0:   Ok,
	1:   Unknown,
	100: ErrorParams,
	101: ErrorDuplicateRequest,
	102: ErrorDataBaseError,
	103: ErrorDataNotFound,
	104: ErrorUsersig,
	105: ErrorRPC,
}
