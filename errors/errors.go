package errors

import "fmt"

// 定义错误类型
type HufuError struct {
	Code    int
	Message string
}

// 实现 error 接口
func (e *HufuError) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}

// 定义具体的错误
var (
	ErrInsufficientBalance = &HufuError{Code: 1001, Message: "余额不足"}
	ErrWalletNotFound      = &HufuError{Code: 1002, Message: "钱包未找到"}
	ErrTransferFailed      = &HufuError{Code: 1003, Message: "转账失败"}
	ErrNoProxyWallets      = &HufuError{Code: 1004, Message: "没有可用的代理钱包"}
	ErrPrivateKeyNotFound  = &HufuError{Code: 1005, Message: "私钥未找到"}
	ErrPrivateKeyInvalid   = &HufuError{Code: 1006, Message: "私钥无效"}
	ErrPrivateKeyApply     = &HufuError{Code: 1007, Message: "私钥申请失败"}
)

// 创建新的错误实例的函数
func NewHufuError(code int, message string) *HufuError {
	return &HufuError{
		Code:    code,
		Message: message,
	}
}
