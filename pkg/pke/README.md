# PKE - Package Error Code

Word Hero 项目的统一错误码管理模块

## 概述

PKE 模块提供了 Word Hero 项目的统一错误码定义和管理功能。所有前端返回的错误都在此模块中定义常量，确保错误码的统一性和可维护性。

## 错误码格式

错误码采用 9 位数字格式：`100000001`

- **前 4 位**：模块编码（从 1000 开始）
  - Word Hero 模块编码：`1000`

- **后 5 位**：具体错误码（从 00000 开始）

示例：
- `100000101` = Word Hero 模块(1000) + 用户错误(101)
- `100000301` = Word Hero 模块(1000) + 单词错误(301)

## 错误码分类

### 通用错误 (100000000-100000099)
- `100000001` - 系统内部错误
- `100000002` - 无效的请求参数
- `100000003` - 未授权访问
- `100000004` - 禁止访问
- `100000005` - 资源未找到
- `100000006` - 方法不允许
- `100000007` - 请求超时
- `100000008` - 请求过于频繁
- `100000009` - 服务不可用

### 用户相关错误 (100000100-100000199)
- `100000101` - 用户不存在
- `100000102` - 用户已存在
- `100000103` - 无效的密码
- `100000104` - 用户已被禁用
- `100000105` - 用户未验证
- `100000106` - 无效的令牌
- `100000107` - 令牌已过期

### 认证相关错误 (100000200-100000299)
- `100000201` - 认证失败
- `100000202` - 需要登录
- `100000203` - 权限不足
- `100000204` - 会话已过期
- `100000205` - 无效的凭据

### 单词相关错误 (100000300-100000399)
- `100000301` - 单词不存在
- `100000302` - 无效的单词ID
- `100000303` - 单词已存在
- `100000304` - 单词导入失败
- `100000305` - 无效的单词数据

### 单词标记相关错误 (100000400-100000499)
- `100000401` - 单词标记不存在
- `100000402` - 无效的单词标记
- `100000403` - 标记操作失败
- `100000404` - 取消标记失败
- `100000405` - 忘光操作失败
- `100000406` - 无效的标记请求
- `100000407` - 单词已标记
- `100000408` - 单词未标记

### 分页相关错误 (100000500-100000599)
- `100000501` - 无效的页码
- `100000502` - 无效的页面大小
- `100000503` - 页码超出范围
- `100000504` - 无效的分页参数

### 搜索相关错误 (100000600-100000699)
- `100000601` - 搜索失败
- `100000602` - 无效的搜索查询
- `100000603` - 搜索超时
- `100000604` - 无搜索结果

### 数据库相关错误 (100000700-100000799)
- `100000701` - 数据库错误
- `100000702` - 数据库连接错误
- `100000703` - 数据库查询错误
- `100000704` - 数据库事务错误
- `100000705` - 数据冲突

### 文件相关错误 (100000800-100000899)
- `100000801` - 文件不存在
- `100000802` - 文件读取错误
- `100000803` - 文件写入错误
- `100000804` - 无效的文件类型
- `100000805` - 文件过大

### 网络相关错误 (100000900-100000999)
- `100000901` - 网络错误
- `100000902` - 连接错误
- `100000903` - 超时错误
- `100000904` - 服务错误

## 使用方法

### 基本使用

```go
package main

import (
	"fmt"
	"github.com/sanmu2018/word-hero/pkg/pke"
)

func main() {
	// 使用预定义的错误码
	errorCode := pke.CodeUserNotFound
	errorMessage := pke.GetErrorMessage(errorCode)

	fmt.Printf("错误码: %d, 错误消息: %s\n", errorCode, errorMessage)

	// 判断错误类型
	if pke.IsError(errorCode) {
		fmt.Println("这是一个错误")
	}

	// 格式化错误描述
	formatted := pke.FormatErrorCode(errorCode)
	fmt.Println("格式化错误:", formatted)
}
```

### 在 API 响应中使用

```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sanmu2018/word-hero/pkg/pke"
)

func GetUserHandler(c *gin.Context) {
	userID := c.Param("id")

	user, err := userService.GetUser(userID)
	if err != nil {
		// 根据错误类型返回相应的错误码
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(200, gin.H{
				"code": pke.CodeUserNotFound,
				"data": nil,
				"msg":  pke.GetErrorMessage(pke.CodeUserNotFound),
			})
			return
		}

		// 其他系统错误
		c.JSON(200, gin.H{
			"code": pke.CodeSystemError,
			"data": nil,
			"msg":  pke.GetErrorMessage(pke.CodeSystemError),
		})
		return
	}

	// 成功响应
	c.JSON(200, gin.H{
		"code": pke.CodeSuccess,
		"data": user,
		"msg":  "",
	})
}
```

### 在服务层使用

```go
package service

import (
	"errors"
	"github.com/sanmu2018/word-hero/pkg/pke"
)

func (s *UserService) GetUser(userID string) (*User, error) {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) GetUserWithErrorCode(userID string) (*User, int) {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pke.CodeUserNotFound
		}
		return nil, pke.CodeDatabaseError
	}

	if user == nil {
		return nil, pke.CodeUserNotFound
	}

	return user, pke.CodeSuccess
}
```

## 工具函数

### GetErrorMessage(code int) string
根据错误码获取对应的错误消息

### IsSuccess(code int) bool
判断是否为成功码（0）

### IsError(code int) bool
判断是否为错误码（非0）

### GetModuleCode(code int) int
获取错误码的模块编码部分

### GetErrorCode(code int) int
获取错误码的具体错误码部分

### FormatErrorCode(code int) string
格式化错误码为可读描述

## 扩展错误码

如果需要添加新的错误码，请按照以下步骤：

1. 在适当的错误码范围内选择下一个可用的错误码
2. 在 `errors.go` 文件中添加新的常量定义
3. 在 `ErrorMessages` 映射中添加对应的错误消息
4. 更新相关文档和测试用例

### 示例：添加新的错误码

```go
// 在 errors.go 中添加
const (
    // 在单词相关错误范围内添加
    CodeWordTranslationMissing = 100000306 // 单词翻译缺失
)

// 在 ErrorMessages 映射中添加
var ErrorMessages = map[int]string{
    // ... 现有错误消息
    CodeWordTranslationMissing: "单词翻译缺失",
}
```

## 测试

运行测试：

```bash
go test ./pkg/pke/...
```

运行示例：

```bash
go test ./pkg/pke/ -run ExampleErrorCodes
```

## 性能考虑

错误码查询操作是 O(1) 时间复杂度的，性能开销很小。错误消息映射在初始化时创建，后续查询只是简单的哈希表查找。

## 版本历史

- v1.0.0 (2025-09-14)
  - 初始版本发布
  - 定义了完整的错误码体系
  - 提供了工具函数和示例代码