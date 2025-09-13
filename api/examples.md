# Word Hero API 使用示例

本文档提供了各种编程语言和工具中使用 Word Hero API 的具体示例。

## 快速开始

### 基础设置
- API 基础URL: `http://localhost:8080`
- 所有响应都是 JSON 格式
- 默认端口: 8080

---

## JavaScript / Fetch API

### 获取词汇列表

```javascript
// 获取第一页词汇，每页24个
async function getVocabularyList() {
    try {
        const response = await fetch('/api/words?page=1&pageSize=24');
        const data = await response.json();

        if (data.success) {
            console.log('词汇列表:', data.data.words);
            console.log('分页信息:', {
                currentPage: data.data.current_page,
                totalPages: data.data.total_pages,
                totalWords: data.data.total_words,
                hasNext: data.data.has_next,
                hasPrev: data.data.has_prev
            });

            // 渲染词汇卡片
            renderVocabularyCards(data.data.words);
        } else {
            console.error('获取失败:', data.error);
        }
    } catch (error) {
        console.error('网络错误:', error);
    }
}

// 渲染词汇卡片函数
function renderVocabularyCards(words) {
    const container = document.getElementById('vocabulary-container');
    container.innerHTML = '';

    words.forEach(word => {
        const card = document.createElement('div');
        card.className = 'vocabulary-card';
        card.innerHTML = `
            <div class="english-word">${word.english}</div>
            <div class="chinese-meaning">${word.chinese}</div>
        `;
        container.appendChild(card);
    });
}

// 调用函数
getVocabularyList();
```

### 搜索功能

```javascript
// 搜索词汇
async function searchVocabulary(query) {
    if (!query || query.length < 2) {
        alert('请输入至少2个字符的搜索关键词');
        return;
    }

    try {
        const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
        const data = await response.json();

        if (data.success) {
            console.log(`找到 ${data.data.count} 个结果`);
            renderSearchResults(data.data.results, data.data.query);
        } else {
            console.error('搜索失败:', data.error);
        }
    } catch (error) {
        console.error('搜索错误:', error);
    }
}

// 渲染搜索结果
function renderSearchResults(results, query) {
    const container = document.getElementById('search-results');
    container.innerHTML = `
        <h3>搜索 "${query}" 的结果</h3>
        <p>共找到 ${results.length} 个结果</p>
    `;

    results.forEach(word => {
        const item = document.createElement('div');
        item.className = 'search-result-item';
        item.innerHTML = `
            <strong>${word.english}</strong> - ${word.chinese}
        `;
        container.appendChild(item);
    });
}

// 使用示例
searchVocabulary('abandon');
```

### 分页导航

```javascript
// 分页导航类
class PaginationManager {
    constructor() {
        this.currentPage = 1;
        this.pageSize = 24;
        this.totalPages = 1;
    }

    async loadPage(pageNumber) {
        if (pageNumber < 1 || pageNumber > this.totalPages) {
            return;
        }

        try {
            const response = await fetch(`/api/words?page=${pageNumber}&pageSize=${this.pageSize}`);
            const data = await response.json();

            if (data.success) {
                this.currentPage = data.data.current_page;
                this.totalPages = data.data.total_pages;
                this.renderPage(data.data);
                this.updateNavigationButtons();
            }
        } catch (error) {
            console.error('加载页面失败:', error);
        }
    }

    renderPage(pageData) {
        console.log('渲染页面:', pageData);
        // 渲染逻辑...
    }

    updateNavigationButtons() {
        const prevBtn = document.getElementById('prev-btn');
        const nextBtn = document.getElementById('next-btn');

        if (prevBtn) prevBtn.disabled = this.currentPage <= 1;
        if (nextBtn) nextBtn.disabled = this.currentPage >= this.totalPages;
    }

    goToNextPage() {
        if (this.currentPage < this.totalPages) {
            this.loadPage(this.currentPage + 1);
        }
    }

    goToPrevPage() {
        if (this.currentPage > 1) {
            this.loadPage(this.currentPage - 1);
        }
    }
}

// 使用示例
const pagination = new PaginationManager();
pagination.loadPage(1);
```

---

## Python / Requests

### 基础示例

```python
import requests
import json

class WordHeroAPI:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url

    def get_words(self, page=1, page_size=24):
        """获取词汇列表"""
        url = f"{self.base_url}/api/words"
        params = {
            'page': page,
            'pageSize': page_size
        }

        try:
            response = requests.get(url, params=params)
            response.raise_for_status()
            data = response.json()

            if data['success']:
                return data['data']
            else:
                print(f"API 错误: {data['error']}")
                return None
        except requests.exceptions.RequestException as e:
            print(f"网络错误: {e}")
            return None

    def search_words(self, query):
        """搜索词汇"""
        url = f"{self.base_url}/api/search"
        params = {'q': query}

        try:
            response = requests.get(url, params=params)
            response.raise_for_status()
            data = response.json()

            if data['success']:
                return data['data']
            else:
                print(f"搜索错误: {data['error']}")
                return None
        except requests.exceptions.RequestException as e:
            print(f"搜索网络错误: {e}")
            return None

    def get_stats(self):
        """获取统计信息"""
        url = f"{self.base_url}/api/stats"

        try:
            response = requests.get(url)
            response.raise_for_status()
            data = response.json()

            if data['success']:
                return data['data']
            else:
                print(f"统计错误: {data['error']}")
                return None
        except requests.exceptions.RequestException as e:
            print(f"统计网络错误: {e}")
            return None

# 使用示例
if __name__ == "__main__":
    api = WordHeroAPI()

    # 获取第一页词汇
    words_data = api.get_words(page=1, page_size=10)
    if words_data:
        print(f"总词汇数: {words_data['total_words']}")
        print(f"当前页: {words_data['current_page']}/{words_data['total_pages']}")

        # 打印前5个词汇
        for i, word in enumerate(words_data['words'][:5]):
            print(f"{i+1}. {word['english']} - {word['chinese']}")

    # 搜索词汇
    search_results = api.search_words("abandon")
    if search_results:
        print(f"\n搜索 '{search_results['query']}' 找到 {search_results['count']} 个结果:")
        for word in search_results['results']:
            print(f"- {word['english']}: {word['chinese']}")

    # 获取统计信息
    stats = api.get_stats()
    if stats:
        print(f"\n应用统计:")
        print(f"- 词汇总数: {stats['total_words']}")
        print(f"- 总页数: {stats['totalPages']}")
        print(f"- 每页大小: {stats['pageSize']}")
```

### 高级用法

```python
import requests
from typing import List, Dict, Optional
import time

class AdvancedWordHeroAPI:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()

    def get_all_words(self, page_size=100) -> List[Dict]:
        """获取所有词汇（分页获取）"""
        all_words = []

        # 先获取第一页以确定总页数
        first_page = self.get_words(page=1, page_size=page_size)
        if not first_page:
            return all_words

        total_pages = first_page['total_pages']
        all_words.extend(first_page['words'])

        # 获取剩余页面
        for page in range(2, total_pages + 1):
            page_data = self.get_words(page=page, page_size=page_size)
            if page_data:
                all_words.extend(page_data['words'])
                print(f"已获取第 {page}/{total_pages} 页")
                time.sleep(0.1)  # 避免请求过快

        return all_words

    def batch_search(self, queries: List[str]) -> Dict[str, List[Dict]]:
        """批量搜索多个词汇"""
        results = {}

        for query in queries:
            search_result = self.search_words(query)
            if search_result:
                results[query] = search_result['results']
            else:
                results[query] = []

            time.sleep(0.1)  # 避免请求过快

        return results

    def export_words_to_csv(self, filename: str, page_size=100):
        """导出词汇到CSV文件"""
        import csv

        words = self.get_all_words(page_size)

        with open(filename, 'w', newline='', encoding='utf-8') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(['English', 'Chinese'])

            for word in words:
                writer.writerow([word['english'], word['chinese']])

        print(f"已导出 {len(words)} 个词汇到 {filename}")

# 使用示例
if __name__ == "__main__":
    api = AdvancedWordHeroAPI()

    # 批量搜索
    queries = ['abandon', 'ability', 'academic']
    batch_results = api.batch_search(queries)

    for query, results in batch_results.items():
        print(f"'{query}': {len(results)} 个结果")

    # 导出词汇
    api.export_words_to_csv('vocabulary.csv')
```

---

## Node.js / Axios

### 基础示例

```javascript
const axios = require('axios');

class WordHeroClient {
    constructor(baseUrl = 'http://localhost:8080') {
        this.baseUrl = baseUrl;
        this.client = axios.create({
            baseURL: baseUrl,
            timeout: 10000,
            headers: {
                'Content-Type': 'application/json'
            }
        });
    }

    async getWords(page = 1, pageSize = 24) {
        try {
            const response = await this.client.get('/api/words', {
                params: { page, pageSize }
            });

            if (response.data.success) {
                return response.data.data;
            } else {
                throw new Error(response.data.error);
            }
        } catch (error) {
            console.error('获取词汇失败:', error.message);
            throw error;
        }
    }

    async searchWords(query) {
        try {
            const response = await this.client.get('/api/search', {
                params: { q: query }
            });

            if (response.data.success) {
                return response.data.data;
            } else {
                throw new Error(response.data.error);
            }
        } catch (error) {
            console.error('搜索失败:', error.message);
            throw error;
        }
    }

    async getStats() {
        try {
            const response = await this.client.get('/api/stats');

            if (response.data.success) {
                return response.data.data;
            } else {
                throw new Error(response.data.error);
            }
        } catch (error) {
            console.error('获取统计失败:', error.message);
            throw error;
        }
    }
}

// 使用示例
async function main() {
    const client = new WordHeroClient();

    try {
        // 获取词汇
        const wordsData = await client.getWords(1, 10);
        console.log(`词汇总数: ${wordsData.total_words}`);
        console.log(`当前页: ${wordsData.current_page}/${wordsData.total_pages}`);

        // 搜索
        const searchResults = await client.searchWords('abandon');
        console.log(`搜索结果: ${searchResults.count} 个`);

        // 统计
        const stats = await client.getStats();
        console.log('应用统计:', stats);

    } catch (error) {
        console.error('API 调用失败:', error.message);
    }
}

main();
```

---

## cURL 示例

### 基础用法

```bash
#!/bin/bash

# 获取第一页词汇
curl -X GET "http://localhost:8080/api/words?page=1&pageSize=24" \
  -H "Content-Type: application/json"

# 获取指定页面
curl -X GET "http://localhost:8080/api/page/2" \
  -H "Content-Type: application/json"

# 搜索词汇
curl -X GET "http://localhost:8080/api/search?q=abandon" \
  -H "Content-Type: application/json"

# 获取统计信息
curl -X GET "http://localhost:8080/api/stats" \
  -H "Content-Type: application/json"
```

### 高级用法

```bash
#!/bin/bash

# 输出格式化的 JSON
function call_api() {
    local endpoint="$1"
    shift
    local params="$@"

    curl -s -X GET "http://localhost:8080$endpoint?$params" \
      -H "Content-Type: application/json" | \
      python3 -m json.tool
}

# 获取词汇列表
echo "=== 获取词汇列表 ==="
call_api "/api/words" "page=1&pageSize=5"

# 搜索词汇
echo -e "\n=== 搜索词汇 ==="
call_api "/api/search" "q=abandon"

# 获取统计信息
echo -e "\n=== 统计信息 ==="
call_api "/api/stats"

# 测试错误情况
echo -e "\n=== 测试错误情况 ==="
call_api "/api/words" "page=999"
call_api "/api/search" "q="
```

---

## Java / OkHttp

### 基础示例

```java
import okhttp3.*;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import java.io.IOException;
import java.util.List;
import java.util.Map;

public class WordHeroClient {
    private final OkHttpClient client;
    private final Gson gson;
    private final String baseUrl;

    public WordHeroClient(String baseUrl) {
        this.client = new OkHttpClient();
        this.gson = new Gson();
        this.baseUrl = baseUrl;
    }

    public JsonObject getWords(int page, int pageSize) throws IOException {
        HttpUrl url = HttpUrl.parse(baseUrl + "/api/words")
                .newBuilder()
                .addQueryParameter("page", String.valueOf(page))
                .addQueryParameter("pageSize", String.valueOf(pageSize))
                .build();

        Request request = new Request.Builder()
                .url(url)
                .build();

        try (Response response = client.newCall(request).execute()) {
            String responseBody = response.body().string();
            return gson.fromJson(responseBody, JsonObject.class);
        }
    }

    public JsonObject searchWords(String query) throws IOException {
        HttpUrl url = HttpUrl.parse(baseUrl + "/api/search")
                .newBuilder()
                .addQueryParameter("q", query)
                .build();

        Request request = new Request.Builder()
                .url(url)
                .build();

        try (Response response = client.newCall(request).execute()) {
            String responseBody = response.body().string();
            return gson.fromJson(responseBody, JsonObject.class);
        }
    }

    public JsonObject getStats() throws IOException {
        HttpUrl url = HttpUrl.parse(baseUrl + "/api/stats").build();

        Request request = new Request.Builder()
                .url(url)
                .build();

        try (Response response = client.newCall(request).execute()) {
            String responseBody = response.body().string();
            return gson.fromJson(responseBody, JsonObject.class);
        }
    }

    public static void main(String[] args) {
        WordHeroClient client = new WordHeroClient("http://localhost:8080");

        try {
            // 获取词汇
            JsonObject wordsResponse = client.getWords(1, 10);
            if (wordsResponse.get("success").getAsBoolean()) {
                JsonObject data = wordsResponse.getAsJsonObject("data");
                System.out.println("词汇总数: " + data.get("total_words").getAsInt());
                System.out.println("当前页: " + data.get("current_page").getAsInt());
            }

            // 搜索词汇
            JsonObject searchResponse = client.searchWords("abandon");
            if (searchResponse.get("success").getAsBoolean()) {
                JsonObject searchData = searchResponse.getAsJsonObject("data");
                System.out.println("搜索结果: " + searchData.get("count").getAsInt() + " 个");
            }

            // 获取统计
            JsonObject statsResponse = client.getStats();
            if (statsResponse.get("success").getAsBoolean()) {
                JsonObject stats = statsResponse.getAsJsonObject("data");
                System.out.println("总页数: " + stats.get("totalPages").getAsInt());
            }

        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
```

---

## React Hooks 示例

```javascript
import { useState, useEffect } from 'react';

// 使用自定义 Hook 获取词汇数据
export function useVocabulary(page = 1, pageSize = 24) {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            setError(null);

            try {
                const response = await fetch(`/api/words?page=${page}&pageSize=${pageSize}`);
                const result = await response.json();

                if (result.success) {
                    setData(result.data);
                } else {
                    setError(result.error);
                }
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [page, pageSize]);

    return { data, loading, error };
}

// 搜索 Hook
export function useSearch(query) {
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        if (!query || query.length < 2) {
            setResults([]);
            return;
        }

        const search = async () => {
            setLoading(true);
            setError(null);

            try {
                const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
                const result = await response.json();

                if (result.success) {
                    setResults(result.data.results);
                } else {
                    setError(result.error);
                }
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        const debounceTimer = setTimeout(search, 300);
        return () => clearTimeout(debounceTimer);
    }, [query]);

    return { results, loading, error };
}

// 组件示例
function VocabularyList() {
    const [currentPage, setCurrentPage] = useState(1);
    const [searchQuery, setSearchQuery] = useState('');

    const { data: vocabularyData, loading: vocabLoading } = useVocabulary(currentPage);
    const { results: searchResults, loading: searchLoading } = useSearch(searchQuery);

    if (vocabLoading) {
        return <div>加载中...</div>;
    }

    if (!vocabularyData) {
        return <div>加载失败</div>;
    }

    return (
        <div>
            {/* 搜索框 */}
            <input
                type="text"
                placeholder="搜索词汇..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
            />

            {/* 搜索结果 */}
            {searchQuery && searchLoading && <div>搜索中...</div>}
            {searchQuery && !searchLoading && searchResults.length > 0 && (
                <div>
                    <h3>搜索结果</h3>
                    {searchResults.map((word, index) => (
                        <div key={index}>
                            {word.english} - {word.chinese}
                        </div>
                    ))}
                </div>
            )}

            {/* 词汇列表 */}
            <div>
                {vocabularyData.words.map((word, index) => (
                    <div key={index}>
                        <strong>{word.english}</strong> - {word.chinese}
                    </div>
                ))}
            </div>

            {/* 分页控制 */}
            <div>
                <button
                    onClick={() => setCurrentPage(currentPage - 1)}
                    disabled={currentPage <= 1}
                >
                    上一页
                </button>
                <span>
                    第 {currentPage} 页，共 {vocabularyData.total_pages} 页
                </span>
                <button
                    onClick={() => setCurrentPage(currentPage + 1)}
                    disabled={currentPage >= vocabularyData.total_pages}
                >
                    下一页
                </button>
            </div>
        </div>
    );
}
```

这些示例展示了在不同环境中使用 Word Hero API 的方法，你可以根据自己的需求选择合适的实现方式。