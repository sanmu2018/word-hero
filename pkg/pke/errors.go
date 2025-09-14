package pke

import "fmt"

// 错误码常量定义
// 成功码是0表示成功
// 错误码格式：100000001（9位）
// 前四位：从1000开始的模块编码，word-hero模块编码是1000
// 后5位：具体错误码，从0开始

const (
	// 成功码
	CodeSuccess = 0

	// ===== 模块编码: 1000 (word-hero模块) =====

	// 通用错误 (100000000-100000099)
	CodeSystemError        = 100000001 // 系统内部错误
	CodeInvalidRequest     = 100000002 // 无效的请求参数
	CodeUnauthorized       = 100000003 // 未授权访问
	CodeForbidden          = 100000004 // 禁止访问
	CodeNotFound           = 100000005 // 资源未找到
	CodeMethodNotAllowed   = 100000006 // 方法不允许
	CodeRequestTimeout     = 100000007 // 请求超时
	CodeTooManyRequests    = 100000008 // 请求过于频繁
	CodeServiceUnavailable = 100000009 // 服务不可用

	// 用户相关错误 (100000100-100000199)
	CodeUserNotFound      = 100000101 // 用户不存在
	CodeUserAlreadyExists = 100000102 // 用户已存在
	CodeInvalidPassword   = 100000103 // 无效的密码
	CodeUserDisabled      = 100000104 // 用户已被禁用
	CodeUserNotVerified   = 100000105 // 用户未验证
	CodeInvalidToken      = 100000106 // 无效的令牌
	CodeTokenExpired      = 100000107 // 令牌已过期

	// 认证相关错误 (100000200-100000299)
	CodeAuthFailed        = 100000201 // 认证失败
	CodeLoginRequired     = 100000202 // 需要登录
	CodePermissionDenied  = 100000203 // 权限不足
	CodeSessionExpired    = 100000204 // 会话已过期
	CodeInvalidCredentials = 100000205 // 无效的凭据

	// 单词相关错误 (100000300-100000399)
	CodeWordNotFound      = 100000301 // 单词不存在
	CodeInvalidWordID     = 100000302 // 无效的单词ID
	CodeWordAlreadyExists = 100000303 // 单词已存在
	CodeWordImportFailed  = 100000304 // 单词导入失败
	CodeInvalidWordData   = 100000305 // 无效的单词数据

	// 单词标记相关错误 (100000400-100000499)
	CodeWordTagNotFound    = 100000401 // 单词标记不存在
	CodeInvalidWordTag     = 100000402 // 无效的单词标记
	CodeMarkFailed         = 100000403 // 标记操作失败
	CodeUnmarkFailed       = 100000404 // 取消标记失败
	CodeForgetFailed       = 100000405 // 忘光操作失败
	CodeInvalidMarkRequest = 100000406 // 无效的标记请求
	CodeAlreadyMarked      = 100000407 // 单词已标记
	CodeNotMarked          = 100000408 // 单词未标记

	// 分页相关错误 (100000500-100000599)
	CodeInvalidPage       = 100000501 // 无效的页码
	CodeInvalidPageSize   = 100000502 // 无效的页面大小
	CodePageOutOfRange    = 100000503 // 页码超出范围
	CodeInvalidPagination = 100000504 // 无效的分页参数

	// 搜索相关错误 (100000600-100000699)
	CodeSearchFailed      = 100000601 // 搜索失败
	CodeInvalidSearchQuery = 100000602 // 无效的搜索查询
	CodeSearchTimeout     = 100000603 // 搜索超时
	CodeNoSearchResults   = 100000604 // 无搜索结果

	// 数据库相关错误 (100000700-100000799)
	CodeDatabaseError     = 100000701 // 数据库错误
	CodeDBConnectionError = 100000702 // 数据库连接错误
	CodeDBQueryError      = 100000703 // 数据库查询错误
	CodeDBTransactionError = 100000704 // 数据库事务错误
	CodeDataConflict      = 100000705 // 数据冲突

	// 文件相关错误 (100000800-100000899)
	CodeFileNotFound      = 100000801 // 文件不存在
	CodeFileReadError     = 100000802 // 文件读取错误
	CodeFileWriteError    = 100000803 // 文件写入错误
	CodeInvalidFileType   = 100000804 // 无效的文件类型
	CodeFileTooLarge      = 100000805 // 文件过大

	// 网络相关错误 (100000900-100000999)
	CodeNetworkError      = 100000901 // 网络错误
	CodeConnectionError   = 100000902 // 连接错误
	CodeTimeoutError      = 100000903 // 超时错误
	CodeServiceError      = 100000904 // 服务错误

	// 配置相关错误 (100001000-100001099)
	CodeConfigError       = 100001001 // 配置错误
	CodeMissingConfig     = 100001002 // 缺少配置
	CodeInvalidConfig     = 100001003 // 无效的配置

	// 业务逻辑错误 (100001100-100001199)
	CodeBusinessLogicError = 100001101 // 业务逻辑错误
	CodeInvalidOperation   = 100001102 // 无效的操作
	CodeOperationFailed    = 100001103 // 操作失败
	CodeInvalidState       = 100001104 // 无效状态
	CodeResourceBusy      = 100001105 // 资源繁忙

	// 数据验证错误 (100001200-100001299)
	CodeValidationError    = 100001201 // 数据验证错误
	CodeRequiredField      = 100001202 // 必填字段缺失
	CodeInvalidFormat      = 100001203 // 格式无效
	CodeOutOfRange        = 100001204 // 超出范围
	CodeDuplicateEntry     = 100001205 // 重复条目

	// 外部服务错误 (100001300-100001399)
	CodeExternalServiceError = 100001301 // 外部服务错误
	CodeExternalServiceDown  = 100001302 // 外部服务不可用
	CodeRateLimitExceeded   = 100001303 // 速率限制超出
	CodeInvalidResponse    = 100001304 // 无效响应

	// 缓存相关错误 (100001400-100001499)
	CodeCacheError        = 100001401 // 缓存错误
	CodeCacheMiss         = 100001402 // 缓存未命中
	CodeCacheExpired      = 100001403 // 缓存过期
)

// ErrorMessages 错误码对应的中文错误消息
var ErrorMessages = map[int]string{
	CodeSuccess:              "成功",

	// 通用错误
	CodeSystemError:          "系统内部错误",
	CodeInvalidRequest:       "无效的请求参数",
	CodeUnauthorized:         "未授权访问",
	CodeForbidden:            "禁止访问",
	CodeNotFound:             "资源未找到",
	CodeMethodNotAllowed:     "方法不允许",
	CodeRequestTimeout:       "请求超时",
	CodeTooManyRequests:      "请求过于频繁",
	CodeServiceUnavailable:   "服务不可用",

	// 用户相关错误
	CodeUserNotFound:         "用户不存在",
	CodeUserAlreadyExists:    "用户已存在",
	CodeInvalidPassword:      "无效的密码",
	CodeUserDisabled:         "用户已被禁用",
	CodeUserNotVerified:      "用户未验证",
	CodeInvalidToken:         "无效的令牌",
	CodeTokenExpired:         "令牌已过期",

	// 认证相关错误
	CodeAuthFailed:          "认证失败",
	CodeLoginRequired:        "需要登录",
	CodePermissionDenied:     "权限不足",
	CodeSessionExpired:       "会话已过期",
	CodeInvalidCredentials:   "无效的凭据",

	// 单词相关错误
	CodeWordNotFound:         "单词不存在",
	CodeInvalidWordID:        "无效的单词ID",
	CodeWordAlreadyExists:    "单词已存在",
	CodeWordImportFailed:     "单词导入失败",
	CodeInvalidWordData:      "无效的单词数据",

	// 单词标记相关错误
	CodeWordTagNotFound:      "单词标记不存在",
	CodeInvalidWordTag:       "无效的单词标记",
	CodeMarkFailed:           "标记操作失败",
	CodeUnmarkFailed:         "取消标记失败",
	CodeForgetFailed:         "忘光操作失败",
	CodeInvalidMarkRequest:   "无效的标记请求",
	CodeAlreadyMarked:        "单词已标记",
	CodeNotMarked:            "单词未标记",

	// 分页相关错误
	CodeInvalidPage:          "无效的页码",
	CodeInvalidPageSize:      "无效的页面大小",
	CodePageOutOfRange:       "页码超出范围",
	CodeInvalidPagination:     "无效的分页参数",

	// 搜索相关错误
	CodeSearchFailed:         "搜索失败",
	CodeInvalidSearchQuery:   "无效的搜索查询",
	CodeSearchTimeout:        "搜索超时",
	CodeNoSearchResults:      "无搜索结果",

	// 数据库相关错误
	CodeDatabaseError:        "数据库错误",
	CodeDBConnectionError:    "数据库连接错误",
	CodeDBQueryError:         "数据库查询错误",
	CodeDBTransactionError:   "数据库事务错误",
	CodeDataConflict:         "数据冲突",

	// 文件相关错误
	CodeFileNotFound:         "文件不存在",
	CodeFileReadError:        "文件读取错误",
	CodeFileWriteError:       "文件写入错误",
	CodeInvalidFileType:      "无效的文件类型",
	CodeFileTooLarge:         "文件过大",

	// 网络相关错误
	CodeNetworkError:         "网络错误",
	CodeConnectionError:      "连接错误",
	CodeTimeoutError:         "超时错误",
	CodeServiceError:         "服务错误",

	// 配置相关错误
	CodeConfigError:          "配置错误",
	CodeMissingConfig:        "缺少配置",
	CodeInvalidConfig:        "无效的配置",

	// 业务逻辑错误
	CodeBusinessLogicError:  "业务逻辑错误",
	CodeInvalidOperation:     "无效的操作",
	CodeOperationFailed:      "操作失败",
	CodeInvalidState:         "无效状态",
	CodeResourceBusy:         "资源繁忙",

	// 数据验证错误
	CodeValidationError:      "数据验证错误",
	CodeRequiredField:        "必填字段缺失",
	CodeInvalidFormat:        "格式无效",
	CodeOutOfRange:          "超出范围",
	CodeDuplicateEntry:       "重复条目",

	// 外部服务错误
	CodeExternalServiceError: "外部服务错误",
	CodeExternalServiceDown:  "外部服务不可用",
	CodeRateLimitExceeded:    "速率限制超出",
	CodeInvalidResponse:      "无效响应",

	// 缓存相关错误
	CodeCacheError:           "缓存错误",
	CodeCacheMiss:            "缓存未命中",
	CodeCacheExpired:         "缓存过期",
}

// GetErrorMessage 根据错误码获取错误消息
func GetErrorMessage(code int) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "未知错误"
}

// IsSuccess 判断是否为成功码
func IsSuccess(code int) bool {
	return code == CodeSuccess
}

// IsError 判断是否为错误码
func IsError(code int) bool {
	return code != CodeSuccess
}

// GetModuleCode 获取模块编码
func GetModuleCode(code int) int {
	return code / 100000
}

// GetErrorCode 获取具体错误码
func GetErrorCode(code int) int {
	return code % 100000
}

// FormatErrorCode 格式化错误码描述
func FormatErrorCode(code int) string {
	if IsSuccess(code) {
		return "成功"
	}

	moduleCode := GetModuleCode(code)
	errorCode := GetErrorCode(code)

	moduleName := "未知模块"
	if moduleCode == 1000 {
		moduleName = "word-hero"
	}

	return fmt.Sprintf("%s模块错误[%d-%d]: %s",
		moduleName, moduleCode, errorCode, GetErrorMessage(code))
}