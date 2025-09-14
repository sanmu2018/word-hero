package pke

import (
	"fmt"
	"testing"
)

func Example() {
	// 示例：使用错误码
	fmt.Println("成功码:", CodeSuccess)
	fmt.Println("用户不存在:", CodeUserNotFound)
	fmt.Println("单词不存在:", CodeWordNotFound)
	fmt.Println("标记失败:", CodeMarkFailed)

	// 示例：获取错误消息
	fmt.Println("错误消息:", GetErrorMessage(CodeUserNotFound))

	// 示例：判断错误类型
	fmt.Println("是否成功:", IsSuccess(CodeSuccess))
	fmt.Println("是否错误:", IsError(CodeUserNotFound))

	// 示例：解析错误码
	fmt.Println("模块编码:", GetModuleCode(CodeUserNotFound))
	fmt.Println("错误编码:", GetErrorCode(CodeUserNotFound))

	// 示例：格式化错误描述
	fmt.Println(FormatErrorCode(CodeUserNotFound))
	fmt.Println(FormatErrorCode(CodeWordNotFound))

	// Output:
	// 成功码: 0
	// 用户不存在: 100000101
	// 单词不存在: 100000301
	// 标记失败: 100000403
	// 错误消息: 用户不存在
	// 是否成功: true
	// 是否错误: true
	// 模块编码: 1000
	// 错误编码: 101
	// word-hero模块错误[1000-101]: 用户不存在
	// word-hero模块错误[1000-301]: 单词不存在
}

func TestErrorCodes(t *testing.T) {
	// 测试错误码常量
	if CodeSuccess != 0 {
		t.Errorf("Expected success code to be 0, got %d", CodeSuccess)
	}

	// 测试错误码格式
	expectedModule := 1000
	if GetModuleCode(CodeUserNotFound) != expectedModule {
		t.Errorf("Expected module code %d, got %d", expectedModule, GetModuleCode(CodeUserNotFound))
	}

	// 测试错误消息
	expectedMsg := "用户不存在"
	if GetErrorMessage(CodeUserNotFound) != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, GetErrorMessage(CodeUserNotFound))
	}

	// 测试成功判断
	if !IsSuccess(CodeSuccess) {
		t.Error("Expected CodeSuccess to be success")
	}

	if !IsError(CodeUserNotFound) {
		t.Error("Expected CodeUserNotFound to be error")
	}

	// 测试格式化
	formatted := FormatErrorCode(CodeUserNotFound)
	expectedFormatted := "word-hero模块错误[1000-101]: 用户不存在"
	if formatted != expectedFormatted {
		t.Errorf("Expected formatted message '%s', got '%s'", expectedFormatted, formatted)
	}
}

// BenchmarkErrorCodes 性能测试
func BenchmarkErrorCodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 测试错误码相关操作的性能
		GetErrorMessage(CodeUserNotFound)
		IsError(CodeUserNotFound)
		GetModuleCode(CodeUserNotFound)
		FormatErrorCode(CodeUserNotFound)
	}
}