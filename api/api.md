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

**接口描述**: 获取词汇数据，支持分页功能，分页参数完全可选

**请求方式**: `GET`

**接口地址**: `/api/words`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pageNum | int | 否 | 页码，从1开始 |
| pageSize | int | 否 | 每页显示数量，最大100 |
| sort | string | 否 | 排序字段 |

**请求示例**:
```bash
GET /api/words
GET /api/words?pageNum=1&pageSize=12
GET /api/words?pageNum=2&pageSize=24&sort=english
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "id": "12345678-1234-1234-1234-123456789012",
        "english": "abandon",
        "chinese": "放弃，抛弃",
        "category": "verb",
        "difficulty": "medium",
        "createdAt": 1663245600000,
        "updatedAt": 1663245600000
      },
      {
        "id": "87654321-4321-4321-4321-210987654321",
        "english": "ability",
        "chinese": "能力，才能",
        "category": "noun",
        "difficulty": "easy",
        "createdAt": 1663245600000,
        "updatedAt": 1663245600000
      }
    ],
    "total": 3673
  },
  "msg": ""
}
```

**响应字段说明**:
- `items`: 词汇数组，每个词汇包含完整的字段信息（id, english, chinese, category, difficulty, createdAt, updatedAt）
- `total`: 总词汇数（前端可据此计算总页数和页码信息）

**注意**: 当不提供分页参数时，返回所有词汇数据；提供分页参数时，返回指定页的数据

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

**接口描述**: 根据关键词搜索词汇，支持英文和中文搜索，支持分页

**请求方式**: `GET`

**接口地址**: `/api/search`

**请求参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| q | string | 是 | - | 搜索关键词，至少2个字符 |
| pageNum | int | 否 | 1 | 页码，从1开始 |
| pageSize | int | 否 | 12 | 每页显示数量，最大100 |

**请求示例**:
```bash
GET /api/search?q=abandon
GET /api/search?q=放弃
GET /api/search?q=abandon&pageNum=1&pageSize=12
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "id": "12345678-1234-1234-1234-123456789012",
        "english": "abandon",
        "chinese": "放弃，抛弃",
        "category": "verb",
        "difficulty": "medium",
        "createdAt": 1663245600000,
        "updatedAt": 1663245600000
      },
      {
        "id": "87654321-4321-4321-4321-210987654321",
        "english": "ability",
        "chinese": "能力，才能",
        "category": "noun",
        "difficulty": "easy",
        "createdAt": 1663245600000,
        "updatedAt": 1663245600000
      }
    ],
    "total": 25
  },
  "msg": ""
}
```

**响应字段说明**:
- `items`: 当前页的匹配词汇数组，每个词汇包含完整的字段信息
- `total`: 匹配结果总数（用于计算总页数）

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
    "totalPages": 306,
    "pageSize": 12
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

### 5. 获取用户学习进度

**接口描述**: 获取用户的学习进度统计信息，包括已认识单词数量、总单词数、学习进度率等

**请求方式**: `GET`

**接口地址**: `/api/user/{userID}/progress`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userID | string | 是 | 用户ID |

**请求示例**:
```bash
GET /api/user/12345678-1234-1234-1234-123456789012/progress
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "userID": "12345678-1234-1234-1234-123456789012",
    "knownWords": 156,
    "totalWords": 3673,
    "progressRate": 4.25,
    "recentActivity": 156
  },
  "msg": ""
}
```

**响应字段说明**:
- `userID`: 用户ID
- `knownWords`: 已认识单词数量
- `totalWords`: 总单词数
- `progressRate`: 学习进度率（百分比）
- `recentActivity`: 最近活动数

---

### 6. 标记单词为已认识

**接口描述**: 将指定单词标记为已认识，记录标记时间戳

**请求方式**: `POST`

**接口地址**: `/api/word/mark`

**请求参数** (JSON body):
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| wordID | string | 是 | 单词ID |
| userID | string | 是 | 用户ID |

**请求示例**:
```bash
POST /api/word/mark
Content-Type: application/json

{
  "wordID": "12345678-1234-1234-1234-123456789012",
  "userID": "87654321-4321-4321-4321-210987654321"
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "wordID": "12345678-1234-1234-1234-123456789012",
    "isMarked": true,
    "markCount": 1,
    "message": "单词已标记为认识"
  },
  "msg": ""
}
```

**响应字段说明**:
- `wordID`: 单词ID
- `isMarked`: 是否已标记
- `markCount`: 标记次数（新设计总是1）
- `message`: 操作结果消息

---

### 7. 取消单词标记

**接口描述**: 取消单词的已认识标记，将其标记为未知

**请求方式**: `DELETE`

**接口地址**: `/api/word/mark`

**请求参数** (JSON body):
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| wordID | string | 是 | 单词ID |
| userID | string | 是 | 用户ID |

**请求示例**:
```bash
DELETE /api/word/mark
Content-Type: application/json

{
  "wordID": "12345678-1234-1234-1234-123456789012",
  "userID": "87654321-4321-4321-4321-210987654321"
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "wordID": "12345678-1234-1234-1234-123456789012",
    "isMarked": false,
    "markCount": 0,
    "message": "单词标记已移除"
  },
  "msg": ""
}
```

---

### 8. 检查单词标记状态

**接口描述**: 检查指定单词的标记状态

**请求方式**: `GET`

**接口地址**: `/api/word/{wordID}/mark-status`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| wordID | string | 是 | 单词ID |

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userID | string | 是 | 用户ID |

**请求示例**:
```bash
GET /api/word/12345678-1234-1234-1234-123456789012/mark-status?userID=87654321-4321-4321-4321-210987654321
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "wordID": "12345678-1234-1234-1234-123456789012",
    "isMarked": true,
    "markCount": 1,
    "markedAt": 1663245600000
  },
  "msg": ""
}
```

**响应字段说明**:
- `wordID`: 单词ID
- `isMarked`: 是否已标记为认识
- `markCount`: 标记次数
- `markedAt`: 标记时间戳（毫秒）

---

### 9. 获取用户已认识单词列表

**接口描述**: 获取用户已认识单词的分页列表

**请求方式**: `GET`

**接口地址**: `/api/user/{userID}/known-words`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userID | string | 是 | 用户ID |

**查询参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从1开始 |
| pageSize | int | 否 | 24 | 每页显示数量，最大100 |

**请求示例**:
```bash
GET /api/user/87654321-4321-4321-4321-210987654321/known-words?page=1&pageSize=24
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "wordIDs": [
      "12345678-1234-1234-1234-123456789012",
      "87654321-4321-4321-4321-210987654321"
    ],
    "totalCount": 156
  },
  "msg": ""
}
```

**响应字段说明**:
- `wordIDs`: 已认识单词ID数组
- `totalCount`: 已认识单词总数

---

### 10. 批量忘光单词

**接口描述**: 批量忘光指定的已认识单词

**请求方式**: `DELETE`

**接口地址**: `/api/user/{userID}/forget-words`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userID | string | 是 | 用户ID |

**请求参数** (JSON body):
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| wordIDs | string[] | 是 | 要忘光的单词ID数组 |

**请求示例**:
```bash
DELETE /api/user/87654321-4321-4321-4321-210987654321/forget-words
Content-Type: application/json

{
  "wordIDs": [
    "12345678-1234-1234-1234-123456789012",
    "87654321-4321-4321-4321-210987654321"
  ]
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "wordIDs": [
      "12345678-1234-1234-1234-123456789012",
      "87654321-4321-4321-4321-210987654321"
    ],
    "forgottenCount": 2,
    "message": "已忘光 2 个已认识单词"
  },
  "msg": ""
}
```

**响应字段说明**:
- `wordIDs`: 请求忘光的单词ID数组
- `forgottenCount`: 实际忘光的单词数量
- `message`: 操作结果消息

---

### 11. 忘光所有单词

**接口描述**: 忘光用户所有已认识单词（需要确认）

**请求方式**: `DELETE`

**接口地址**: `/api/user/{userID}/forget-all`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userID | string | 是 | 用户ID |

**请求参数** (JSON body):
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| confirm | boolean | 是 | 确认标志，必须为true |

**请求示例**:
```bash
DELETE /api/user/87654321-4321-4321-4321-210987654321/forget-all
Content-Type: application/json

{
  "confirm": true
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "forgottenCount": 156,
    "message": "已忘光全部 156 个已认识单词"
  },
  "msg": ""
}
```

**响应字段说明**:
- `forgottenCount`: 忘光的单词数量
- `message`: 操作结果消息

---

### 12. 获取用户单词统计

**接口描述**: 获取用户的详细单词学习统计信息

**请求方式**: `GET`

**接口地址**: `/api/user/{userID}/word-stats`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userID | string | 是 | 用户ID |

**请求示例**:
```bash
GET /api/user/87654321-4321-4321-4321-210987654321/word-stats
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "userID": "87654321-4321-4321-4321-210987654321",
    "knownWordsCount": 156,
    "totalWordsCount": 3673,
    "progressRate": 4.25,
    "recentMarks": [
      {
        "wordID": "12345678-1234-1234-1234-123456789012",
        "knownAt": 1663245600000
      }
    ],
    "knownWordsByDate": {},
    "topCategories": {}
  },
  "msg": ""
}
```

**响应字段说明**:
- `userID`: 用户ID
- `knownWordsCount`: 已认识单词数量
- `totalWordsCount`: 总单词数量
- `progressRate`: 学习进度率
- `recentMarks`: 最近标记记录
- `knownWordsByDate`: 按日期分组的已认识单词
- `topCategories`: 热门分类统计

---

## 请求数据模型

### WordSearchRequest 对象（搜索请求）

```json
{
  "q": "string",           // 搜索关键词
  "pageNum": 1,            // 页码，从1开始（可选）
  "pageSize": 12,          // 每页大小（可选）
  "sort": "string"         // 排序字段（可选）
}
```

### BaseList 对象（基础分页请求）

```json
{
  "pageNum": 1,            // 页码，从1开始（可选）
  "pageSize": 12,          // 每页大小（可选）
  "sort": "string"         // 排序字段（可选）
}
```

## 数据模型

### Word 对象

```json
{
  "id": "string",           // 单词ID (UUID)
  "english": "string",      // 英文单词
  "chinese": "string",      // 中文释义
  "category": "string",      // 单词分类
  "difficulty": "string",    // 难度等级
  "createdAt": 1663245600000, // 创建时间戳
  "updatedAt": 1663245600000  // 更新时间戳
}
```

### WordTag 对象（单词标记）

```json
{
  "id": "string",           // 标记ID (UUID)
  "wordId": "string",       // 单词ID
  "userId": "string",       // 用户ID
  "known": 1663245600000,   // 认识时间戳（null表示不认识）
  "createdAt": 1663245600000, // 创建时间戳
  "updatedAt": 1663245600000  // 更新时间戳
}
```

### WordWithMarkStatus 对象（带标记状态的单词）

```json
{
  "word": {                 // Word对象
    "id": "string",
    "english": "string",
    "chinese": "string",
    "category": "string",
    "difficulty": "string",
    "createdAt": 1663245600000,
    "updatedAt": 1663245600000
  },
  "isMarked": true,         // 是否已标记为认识
  "markCount": 1,           // 标记次数
  "markedAt": 1663245600000 // 标记时间戳
}
```

### BaseListResp 对象（基础列表响应）

```json
{
  "items": [],              // 对象数组（根据接口类型可能是Word或其他对象）
  "total": 3673             // 总数量
}
```

### Page 对象（分页响应）

```json
{
  "items": [],              // 对象数组（根据接口类型可能是Word或WordWithMarkStatus）
  "total": 3673,            // 总数量
  "pageNumber": 1,          // 当前页码（可选）
  "pageSize": 12,            // 每页大小（可选）
  "totalPages": 153          // 总页数（可选）
}
```

## 使用示例

### JavaScript 示例

```javascript
// 获取所有词汇（不提供分页参数）
fetch('/api/words')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('所有词汇:', data.data.items);
      console.log('词汇总数:', data.data.total);
    } else {
      console.error('请求失败:', data.msg);
      alert(data.msg);
    }
  })
  .catch(error => {
    console.error('网络错误:', error);
  });

// 获取分页词汇
fetch('/api/words?pageNum=1&pageSize=12')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('词汇列表:', data.data.items);
      console.log('词汇总数:', data.data.total);

      // 计算分页信息
      const pageSize = 12;
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
      alert(data.msg);
    }
  })
  .catch(error => {
    console.error('网络错误:', error);
  });

// 搜索词汇
fetch('/api/search?q=abandon&pageNum=1&pageSize=12')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log(`找到 ${data.data.total} 个匹配结果:`, data.data.items);

      // 计算分页信息
      const pageSize = 12;
      const totalPages = Math.ceil(data.data.total / pageSize);
      const currentPage = 1;

      console.log('搜索分页信息:', {
        currentPage: currentPage,
        totalPages: totalPages,
        hasPrev: currentPage > 1,
        hasNext: currentPage < totalPages
      });
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

// === 新增：单词标记功能示例 ===

// 获取用户学习进度
const userID = '87654321-4321-4321-4321-210987654321';
fetch(`/api/user/${userID}/progress`)
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('学习进度:', data.data);
      console.log(`已认识 ${data.data.knownWords} 个单词，进度 ${data.data.progressRate}%`);
    } else {
      console.error('获取进度失败:', data.msg);
    }
  });

// 标记单词为已认识
const wordID = '12345678-1234-1234-1234-123456789012';
fetch('/api/word/mark', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    wordID: wordID,
    userID: userID
  })
})
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('标记成功:', data.data.message);
      // 更新UI显示已标记状态
      updateWordMarkStatus(wordID, true);
    } else {
      console.error('标记失败:', data.msg);
      alert(data.msg);
    }
  });

// 取消单词标记
fetch('/api/word/mark', {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    wordID: wordID,
    userID: userID
  })
})
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('取消标记成功:', data.data.message);
      // 更新UI显示未标记状态
      updateWordMarkStatus(wordID, false);
    } else {
      console.error('取消标记失败:', data.msg);
      alert(data.msg);
    }
  });

// 检查单词标记状态
fetch(`/api/word/${wordID}/mark-status?userID=${userID}`)
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('单词标记状态:', data.data);
      // 更新UI显示标记状态
      if (data.data.isMarked) {
        console.log(`该单词已于 ${new Date(data.data.markedAt).toLocaleString()} 标记为认识`);
      }
    } else {
      console.error('获取标记状态失败:', data.msg);
    }
  });

// 获取用户已认识单词列表
fetch(`/api/user/${userID}/known-words?page=1&pageSize=24`)
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('已认识单词:', data.data);
      console.log(`共 ${data.data.totalCount} 个已认识单词`);
      // 显示已认识单词列表
      displayKnownWords(data.data.wordIDs);
    } else {
      console.error('获取已认识单词失败:', data.msg);
    }
  });

// 批量忘光单词
const wordsToForget = ['word-id-1', 'word-id-2'];
fetch(`/api/user/${userID}/forget-words`, {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    wordIDs: wordsToForget
  })
})
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('批量忘光成功:', data.data.message);
      alert(data.data.message);
      // 刷新页面或更新UI
      location.reload();
    } else {
      console.error('批量忘光失败:', data.msg);
      alert(data.msg);
    }
  });

// 忘光所有单词（需要用户确认）
if (confirm('确定要忘光所有已认识的单词吗？此操作不可恢复！')) {
  fetch(`/api/user/${userID}/forget-all`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      confirm: true
    })
  })
    .then(response => response.json())
    .then(data => {
      if (data.code === 0) {
        console.log('忘光所有单词成功:', data.data.message);
        alert(data.data.message);
        // 刷新页面或更新UI
        location.reload();
      } else {
        console.error('忘光所有单词失败:', data.msg);
        alert(data.msg);
      }
    });
}

// 获取用户单词统计
fetch(`/api/user/${userID}/word-stats`)
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('用户单词统计:', data.data);
      console.log(`学习进度: ${data.data.progressRate}%`);
      console.log(`最近标记: ${data.data.recentMarks.length} 个单词`);
    } else {
      console.error('获取单词统计失败:', data.msg);
    }
  });

// 辅助函数：更新单词标记状态
function updateWordMarkStatus(wordID, isMarked) {
  const wordElement = document.querySelector(`[data-word-id="${wordID}"]`);
  if (wordElement) {
    wordElement.classList.toggle('marked', isMarked);
    const markButton = wordElement.querySelector('.mark-button');
    if (markButton) {
      markButton.textContent = isMarked ? '已认识' : '标记为认识';
      markButton.classList.toggle('marked', isMarked);
    }
  }
}

// 辅助函数：显示已认识单词列表
function displayKnownWords(wordIDs) {
  const container = document.getElementById('known-words-container');
  if (container) {
    container.innerHTML = `<h3>已认识单词 (${wordIDs.length})</h3>`;
    // 这里可以进一步获取单词详情并显示
  }
}
```

### curl 示例

```bash
# 获取所有词汇（不提供分页参数）
curl "http://localhost:8080/api/words"

# 获取分页词汇
curl "http://localhost:8080/api/words?pageNum=1&pageSize=12"
curl "http://localhost:8080/api/words?pageNum=2&pageSize=24&sort=english"

# 获取指定页面
curl "http://localhost:8080/api/page/2"

# 搜索词汇
curl "http://localhost:8080/api/search?q=abandon"
curl "http://localhost:8080/api/search?q=abandon&pageNum=1&pageSize=12"

# 获取统计信息
curl "http://localhost:8080/api/stats"

# === 新增：单词标记功能示例 ===

# 获取用户学习进度
curl "http://localhost:8080/api/user/87654321-4321-4321-4321-210987654321/progress"

# 标记单词为已认识
curl -X POST "http://localhost:8080/api/word/mark" \
  -H "Content-Type: application/json" \
  -d '{
    "wordID": "12345678-1234-1234-1234-123456789012",
    "userID": "87654321-4321-4321-4321-210987654321"
  }'

# 取消单词标记
curl -X DELETE "http://localhost:8080/api/word/mark" \
  -H "Content-Type: application/json" \
  -d '{
    "wordID": "12345678-1234-1234-1234-123456789012",
    "userID": "87654321-4321-4321-4321-210987654321"
  }'

# 检查单词标记状态
curl "http://localhost:8080/api/word/12345678-1234-1234-1234-123456789012/mark-status?userID=87654321-4321-4321-4321-210987654321"

# 获取用户已认识单词列表
curl "http://localhost:8080/api/user/87654321-4321-4321-4321-210987654321/known-words?page=1&pageSize=24"

# 批量忘光单词
curl -X DELETE "http://localhost:8080/api/user/87654321-4321-4321-4321-210987654321/forget-words" \
  -H "Content-Type: application/json" \
  -d '{
    "wordIDs": ["12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321"]
  }'

# 忘光所有单词
curl -X DELETE "http://localhost:8080/api/user/87654321-4321-4321-4321-210987654321/forget-all" \
  -H "Content-Type: application/json" \
  -d '{
    "confirm": true
  }'

# 获取用户单词统计
curl "http://localhost:8080/api/user/87654321-4321-4321-4321-210987654321/word-stats"
```

## 注意事项

1. **分页参数**: 分页参数完全可选，不提供时返回所有数据；`pageSize` 参数的有效范围为 1-100，超出范围会自动调整
2. **搜索要求**: 搜索关键词至少需要2个字符，搜索结果支持分页显示
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

### v1.2.0 (2025-09-14)
- **单词标记系统重构**: 完全重新设计 word_tags 表结构，从复杂的 JSONB 设计改为简单的时间戳设计
- **用户数据隔离**: 所有单词标记操作都基于用户ID，确保数据安全和隔离
- **新增API接口**:
  - 用户学习进度查询 (/api/user/{userID}/progress)
  - 单词标记/取消标记 (/api/word/mark)
  - 单词标记状态查询 (/api/word/{wordID}/mark-status)
  - 已认识单词列表查询 (/api/user/{userID}/known-words)
  - 批量忘光单词 (/api/user/{userID}/forget-words)
  - 忘光所有单词 (/api/user/{userID}/forget-all)
  - 用户单词统计查询 (/api/user/{userID}/word-stats)
- **数据模型优化**:
  - WordTag 结构简化，使用 Known 时间戳字段替代复杂的 UserTags JSONB
  - 新增 WordWithMarkStatus 结构支持标记状态查询
  - 所有数据模型完善字段定义和时间戳管理
- **安全性增强**:
  - 所有标记操作必须传入有效的用户ID
  - 批量忘光操作需要明确确认
  - 用户数据完全隔离，防止跨用户数据访问
- **API文档完善**: 完整的接口文档和使用示例，包括 JavaScript 和 curl 示例

### v1.3.0 (2025-09-23)
- **API接口优化**: 更新 `/api/words` 接口响应格式，使用 `items` 字段替代 `words` 字段
- **分页参数重构**: `/api/words` 接口分页参数改为完全可选，使用 `BaseList` 结构体统一参数格式
- **搜索接口统一**: `/api/search` 接口响应格式统一为 `BaseListResp`，使用 `items` 和 `total` 字段
- **搜索分页功能**: `/api/search` 接口新增分页支持，支持 `pageNum` 和 `pageSize` 参数
- **词汇数据结构完善**: API 响应中的词汇对象现在包含完整的字段信息（id, english, chinese, category, difficulty, createdAt, updatedAt）
- **统计信息更新**: 更新 `/api/stats` 接口的响应数据，反映新的分页参数格式
- **文档同步更新**: 更新所有相关的 API 文档、示例代码和数据模型说明
- **新增数据模型**: 添加 `BaseListResp` 和 `BaseList` 对象定义，用于标准化的请求和响应格式

---

**联系方式**: 如有问题请通过 GitHub Issues 反馈

**项目地址**: https://github.com/sanmu2018/word-hero