# Word Hero API Documentation

## 概述

Word Hero 是一个基于 Go 的雅思词汇学习 Web 应用，提供完整的 RESTful API 接口用于词汇数据的查询、分页和搜索功能。

**基础信息**
- **基础URL**: `http://localhost:8080`
- **API版本**: v1
- **数据格式**: JSON
- **字符编码**: UTF-8

## 通用响应格式

所有 API 接口都使用统一的响应格式：

```json
{
  "code": 0,
  "data": {},
  "msg": ""
}
```

**响应字段说明：**
- `code`: 数字，表示请求状态码，0 表示成功，非零表示失败
- `data`: 对象，包含具体的响应数据（成功时存在）
- `msg`: 字符串，状态信息或错误信息（失败时必须存在）

**状态码说明：**
- `0`: 成功
- `150321309`: 通用错误码
- 其他非零码: 具体业务错误码

### 分页接口响应格式

分页接口统一使用以下简化格式：

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 1234
  },
  "msg": ""
}
```

**分页响应字段说明：**
- `items`: 数组，包含当前页的数据项
- `total`: 数字，表示数据总数（前端可据此计算总页数）

## 错误码规范

| HTTP状态码 | 错误类型 | 说明 |
|-----------|---------|------|
| 200 | OK | 请求成功 |
| 400 | Bad Request | 请求参数错误 |
| 500 | Internal Server Error | 服务器内部错误 |

---

## API 接口详情

### 1. 获取分页词汇列表

**接口描述**: 获取指定页码的词汇数据，支持自定义每页显示数量

**请求方式**: `GET`

**接口地址**: `/api/words`

**请求参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从1开始 |
| pageSize | int | 否 | 24 | 每页显示数量，最大100 |

**请求示例**:
```bash
GET /api/words?page=1&pageSize=24
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "english": "abandon",
        "chinese": "放弃，抛弃"
      },
      {
        "english": "ability",
        "chinese": "能力，才能"
      }
    ],
    "total": 3673
  },
  "msg": ""
}
```

**响应字段说明**:
- `items`: 词汇数组
- `total`: 总词汇数（前端可据此计算总页数和页码信息）

**错误响应**:
```json
{
  "code": 150321309,
  "data": null,
  "msg": "page number 999 is out of range (1-153)"
}
```

---

### 2. 获取指定页面词汇

**接口描述**: 根据URL路径中的页码获取指定页面的词汇数据

**请求方式**: `GET`

**接口地址**: `/api/page/{pageNumber}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pageNumber | int | 是 | 页码，从1开始 |

**请求示例**:
```bash
GET /api/page/2
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "english": "academic",
        "chinese": "学术的，学院的"
      }
    ],
    "total": 3673
  },
  "msg": ""
}
```

**错误响应**:
```json
{
  "code": 150321309,
  "data": null,
  "msg": "Invalid page number"
}
```

---

### 3. 词汇搜索

**接口描述**: 根据关键词搜索词汇，支持英文和中文搜索

**请求方式**: `GET`

**接口地址**: `/api/search`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| q | string | 是 | 搜索关键词，至少2个字符 |

**请求示例**:
```bash
GET /api/search?q=abandon
GET /api/search?q=放弃
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "query": "abandon",
    "results": [
      {
        "english": "abandon",
        "chinese": "放弃，抛弃"
      }
    ],
    "count": 1
  },
  "msg": ""
}
```

**响应字段说明**:
- `query`: 搜索关键词
- `results`: 匹配的词汇数组
- `count`: 匹配结果数量

**错误响应**:
```json
{
  "code": 150321309,
  "data": null,
  "msg": "Search query is required"
}
```

---

### 4. 获取应用统计信息

**接口描述**: 获取应用的统计信息，包括词汇总数、分页信息等

**请求方式**: `GET`

**接口地址**: `/api/stats`

**请求参数**: 无

**请求示例**:
```bash
GET /api/stats
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "totalWords": 3673,
    "file_source": "words/IELTS.xlsx",
    "totalPages": 153,
    "pageSize": 24
  },
  "msg": ""
}
```

**响应字段说明**:
- `totalWords`: 词汇总数
- `file_source`: 数据源文件路径
- `totalPages`: 总页数
- `pageSize`: 当前每页显示数量

---

## 数据模型

### Word 对象

```json
{
  "english": "string",    // 英文单词
  "chinese": "string"     // 中文释义
}
```

### Page 对象（分页响应）

```json
{
  "items": [],            // Word对象数组
  "total": 3673           // 总词汇数
}
```

## 使用示例

### JavaScript 示例

```javascript
// 获取第一页词汇
fetch('/api/words?page=1&pageSize=24')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('词汇列表:', data.data.items);
      console.log('词汇总数:', data.data.total);

      // 计算分页信息
      const pageSize = 24;
      const totalPages = Math.ceil(data.data.total / pageSize);
      const currentPage = 1;

      console.log('分页信息:', {
        currentPage: currentPage,
        totalPages: totalPages,
        hasPrev: currentPage > 1,
        hasNext: currentPage < totalPages
      });
    } else {
      console.error('请求失败:', data.msg);
      // 显示错误信息给用户
      alert(data.msg);
    }
  })
  .catch(error => {
    console.error('网络错误:', error);
  });

// 搜索词汇
fetch('/api/search?q=abandon')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log(`找到 ${data.data.count} 个匹配结果:`, data.data.results);
    } else {
      console.error('搜索失败:', data.msg);
      alert(data.msg);
    }
  });

// 获取统计信息
fetch('/api/stats')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('应用统计:', data.data);
    } else {
      console.error('获取统计失败:', data.msg);
    }
  });
```

### curl 示例

```bash
# 获取第一页词汇
curl "http://localhost:8080/api/words?page=1&pageSize=24"

# 获取指定页面
curl "http://localhost:8080/api/page/2"

# 搜索词汇
curl "http://localhost:8080/api/search?q=abandon"

# 获取统计信息
curl "http://localhost:8080/api/stats"
```

## 注意事项

1. **分页参数**: `pageSize` 参数的有效范围为 1-100，超出范围会自动调整
2. **搜索要求**: 搜索关键词至少需要2个字符
3. **字符编码**: 所有接口均使用 UTF-8 编码
4. **跨域**: 默认不支持跨域请求，如需跨域请配置 CORS
5. **性能**: 大量请求时建议适当增加缓存机制
6. **数据源**: 词汇数据来源于 `words/IELTS.xlsx` 文件，包含3673个雅思词汇
7. **分页计算**: 前端需根据 `data.total` 和当前页大小计算总页数，API 不再返回分页元数据
8. **错误处理**: 当 `code` 不为 0 时，应将 `msg` 字段内容显示给用户作为错误提示

## 更新日志

### v1.0.0 (2025-09-13)
- 初始版本发布
- 支持词汇分页查询
- 支持词汇搜索功能
- 支持应用统计信息获取
- 统一API响应格式

### v1.1.0 (2025-09-13)
- **响应格式升级**: 从 success/error 格式改为 code/msg 格式
- **分页响应简化**: 分页接口只返回 items 数组和 total 总数
- **错误处理优化**: 所有错误码使用 150321309，错误信息在 msg 字段中
- **前端计算优化**: 分页信息由前端根据 total 和页大小计算，减少服务端计算

---

**联系方式**: 如有问题请通过 GitHub Issues 反馈

**项目地址**: https://github.com/sanmu2018/word-hero