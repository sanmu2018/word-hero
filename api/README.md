# Word Hero API 文档目录

本目录包含了 Word Hero 应用的完整 API 规范文档，为开发者提供详细的接口说明和使用指南。

## 📁 文件结构

```
api/
├── README.md              # 本文档
├── api.md                 # 完整的 API 接口文档
├── openapi.yaml           # OpenAPI 3.0 规范文件
├── examples.md            # 多语言使用示例
└── test-api.js            # API 测试脚本
```

## 📖 文档说明

### 1. `api.md` - 主要 API 文档
**用途**: 完整的 API 接口说明文档
**内容包括**:
- API 概述和基础信息
- 通用响应格式
- 错误码规范
- 4个主要接口的详细说明：
  - 获取分页词汇列表 (`/api/words`)
  - 获取指定页面词汇 (`/api/page/{pageNumber}`)
  - 词汇搜索 (`/api/search`)
  - 获取应用统计信息 (`/api/stats`)
- 数据模型定义
- 使用示例和注意事项

### 2. `openapi.yaml` - OpenAPI 规范
**用途**: 标准化的 API 规范描述
**支持工具**:
- Swagger UI
- Postman
- API 文档生成器
- 代码生成工具

**特点**:
- 符合 OpenAPI 3.0 标准
- 包含完整的请求/响应模式
- 支持多服务器配置
- 包含错误处理说明

### 3. `examples.md` - 多语言使用示例
**用途**: 提供各种编程语言的具体实现示例
**包含语言**:
- JavaScript / Fetch API
- Python / Requests
- Node.js / Axios
- cURL 命令行
- Java / OkHttp
- React Hooks

**内容特色**:
- 完整可运行的代码示例
- 错误处理机制
- 最佳实践建议
- 高级用法展示

### 4. `test-api.js` - API 测试脚本
**用途**: 自动化测试所有 API 接口
**测试内容包括**:
- 功能测试（正常流程）
- 边界值测试
- 错误处理测试
- 性能测试
- 数据验证测试

**使用方法**:
```javascript
// 浏览器控制台
runAllTests();

// Node.js 环境
node test-api.js
```

## 🚀 快速开始

### 1. 查看接口文档
```bash
# 查看主要 API 文档
cat api/api.md

# 或在浏览器中打开查看
```

### 2. 使用 OpenAPI 工具
```bash
# 使用 Swagger UI 查看 (需要安装 swagger-ui)
swagger-ui-dist/api/openapi.yaml

# 或导入到 Postman
# 文件: api/openapi.yaml
```

### 3. 运行测试脚本
```bash
# 确保服务器运行在 localhost:8080
# 在浏览器控制台运行
open http://localhost:8080
# 然后在开发者工具控制台输入: runAllTests()
```

### 4. 复制使用示例
```bash
# 查看 JavaScript 示例
cat api/examples.md | grep -A 20 "JavaScript"

# 查看 Python 示例
cat api/examples.md | grep -A 30 "Python"
```

## 📊 API 接口概览

| 接口 | 方法 | 路径 | 描述 |
|------|------|------|------|
| 获取词汇列表 | GET | `/api/words` | 分页获取词汇数据 |
| 获取指定页面 | GET | `/api/page/{pageNumber}` | 获取指定页码的词汇 |
| 搜索词汇 | GET | `/api/search` | 根据关键词搜索词汇 |
| 获取统计信息 | GET | `/api/stats` | 获取应用统计信息 |

## 🔧 技术规范

### 请求格式
- **方法**: GET
- **内容类型**: application/json
- **字符编码**: UTF-8
- **基础URL**: http://localhost:8080

### 响应格式
```json
{
  "success": true,
  "data": {},
  "error": ""
}
```

### 状态码
- `200`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误

## 📝 版本信息

- **当前版本**: v1.0.0
- **最后更新**: 2025-09-13
- **维护团队**: Word Hero 开发团队

## 🤝 贡献指南

如果您发现文档中的错误或有改进建议，请：

1. 检查现有的 API 实现
2. 确认文档准确性
3. 提交 Issue 或 Pull Request
4. 更新相关的示例代码

## 📞 联系方式

- **项目地址**: https://github.com/sanmu2018/word-hero
- **Issues**: https://github.com/sanmu2018/word-hero/issues
- **文档问题**: 请提交 Issue 或直接修改文档文件

## 📄 许可证

本文档遵循与主项目相同的许可证。