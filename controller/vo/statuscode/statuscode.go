package statuscode

var (
	Success              = &Response{Code: "000000", Message: "成功"}
	BadRequest           = &Response{Code: "000001", Message: "参数错误"}
	InternalServerError  = &Response{Code: "000002", Message: "系统错误"}
	SystemNotInitialized = &Response{Code: "000100", Message: "系统未初始化"}
	SystemInitialized    = &Response{Code: "000101", Message: "系统已初始化"}
	NotLogin             = &Response{Code: "000400", Message: "请先登录"}
	InvalidLogin         = &Response{Code: "000401", Message: "登录失效"}
	PermissionDenied     = &Response{Code: "000402", Message: "权限不足"}

	DatabaseConnectFail = &Response{Code: "001001", Message: "连接数据库失败"}
	DatabaseCreateFail  = &Response{Code: "001002", Message: "数据库创建失败"}
	DatabaseMigrateFail = &Response{Code: "001003", Message: "数据库迁移失败"}

	InvalidUsernameOrPassword = &Response{Code: "002001", Message: "用户名不存在或密码错误"}
	UserNotExists             = &Response{Code: "002002", Message: "用户不存在"}
	UsernameExists            = &Response{Code: "002003", Message: "用户名已存在"}
	InvalidPassword           = &Response{Code: "002004", Message: "密码错误"}
	CannotDeleteCurrentUser   = &Response{Code: "002005", Message: "不能删除自己"}

	DirectoryNotExists             = &Response{Code: "003001", Message: "目录不存在"}
	DirectoryExists                = &Response{Code: "003002", Message: "目录已存在"}
	FileNotExists                  = &Response{Code: "003003", Message: "文件不存在"}
	FileExists                     = &Response{Code: "003004", Message: "文件已存在"}
	TargetFolderInsideSourceFolder = &Response{Code: "003005", Message: "目标文件夹是源文件夹的子文件夹"}

	StoreSpaceNotExists = &Response{Code: "004001", Message: "存储空间不存在"}
	StoreSpaceExists    = &Response{Code: "004002", Message: "存储空间已存在"}
	NoStoreSpace        = &Response{Code: "004003", Message: "请先添加存储空间"}
	NotEnoughFreeSpace  = &Response{Code: "004004", Message: "剩余空间不足"}

	UploadFailed            = &Response{Code: "005001", Message: "上传失败"}
	InvalidUploadStartPoint = &Response{Code: "005002", Message: "上传文件开始位置错误"}

	GroupNotExists = &Response{Code: "006001", Message: "群组不存在"}
	GroupExists    = &Response{Code: "006002", Message: "群组名已被占用"}
)
