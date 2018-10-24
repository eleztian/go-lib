package errno

var (
	OK           = &Errno{Code: 0, Message: "OK"}
	ErrParam     = &Errno{Code: 50000001, Message: "Parameter error."}
	ErrToken     = &Errno{Code: 50000002, Message: "Error occurred while signing the JSON web token."}
	TokenExpire  = &Errno{Code: 50000003, Message: "Token expire."}
	RejectAccess = &Errno{Code: 50000004, Message: "Reject access."}
	RejectUpdate = &Errno{Code: 50000005, Message: "Reject update."}
	ErrPassword  = &Errno{Code: 50000006, Message: "Password was incorrect."}

	NotFound         = &Errno{Code: 50000100, Message: "Not found."}
	NotFoundUser     = &Errno{Code: 50000101, Message: "User not found."}
	NotFoundProduct  = &Errno{Code: 50000102, Message: "Production not found."}
	NotFoundPreFunc  = &Errno{Code: 50000103, Message: "Pre define func not found."}
	NotFoundCusFunc  = &Errno{Code: 50000104, Message: "Customer func not found."}
	NoThingToChange  = &Errno{Code: 50000105, Message: "Nothing to change."} // nothing to change
	NotFoundFirmware = &Errno{Code: 50000106, Message: "Firmware not found."}
	NotFoundBatch    = &Errno{Code: 50000107, Message: "Batch not found."}
	NotFoundCategory = &Errno{Code: 50000108, Message: "Category not found."}
	NotFoundUserKey  = &Errno{Code: 50000109, Message: "User key not found."}

	ExistProduction = &Errno{Code: 50000200, Message: "Production already exists."}
	ExistFunc       = &Errno{Code: 50000201, Message: "Function already exists."}
	ExistFirmware   = &Errno{Code: 50000202, Message: "Firmware already exists."}

	ErrInternalServer = &Errno{Code: 50001000, Message: "Internal server error."}
	ErrDatabase       = &Errno{Code: 50001001, Message: "Database error."}
	ErrConnectCenter  = &Errno{Code: 50001002, Message: "Request VDMP's connectCenter error."}
	ErrShadow         = &Errno{Code: 50001003, Message: "Request VDMP's shadow error."}
	ErrDeCompress     = &Errno{Code: 50001003, Message: "Decompress error."}
)
