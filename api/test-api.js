/**
 * Word Hero API 测试脚本
 * 用于测试所有 API 接口的功能
 */

const BASE_URL = 'http://localhost:8080';

// API 端点配置
const API_ENDPOINTS = {
    words: '/api/words',
    page: '/api/page',
    search: '/api/search',
    stats: '/api/stats'
};

/**
 * 通用请求函数
 */
async function apiRequest(endpoint, params = {}) {
    const url = new URL(BASE_URL + endpoint);

    // 添加查询参数
    Object.keys(params).forEach(key => {
        if (params[key] !== undefined && params[key] !== null) {
            url.searchParams.append(key, params[key]);
        }
    });

    try {
        const response = await fetch(url);
        const data = await response.json();

        console.log(`\n=== ${endpoint} ===`);
        console.log('URL:', url.toString());
        console.log('Status:', response.status);
        console.log('Response:', JSON.stringify(data, null, 2));

        return data;
    } catch (error) {
        console.error(`\n=== ${endpoint} ERROR ===`);
        console.error('URL:', url.toString());
        console.error('Error:', error.message);
        return null;
    }
}

/**
 * 测试词汇分页接口
 */
async function testWordsPagination() {
    console.log('\n🧪 测试词汇分页接口');

    // 测试默认参数
    await apiRequest(API_ENDPOINTS.words);

    // 测试指定页码和每页数量
    await apiRequest(API_ENDPOINTS.words, { page: 1, pageSize: 12 });

    // 测试边界值
    await apiRequest(API_ENDPOINTS.words, { page: 999, pageSize: 24 });

    // 测试无效参数
    await apiRequest(API_ENDPOINTS.words, { page: 0, pageSize: 0 });
}

/**
 * 测试指定页面接口
 */
async function testPageByNumber() {
    console.log('\n🧪 测试指定页面接口');

    // 测试有效页码
    await apiRequest(API_ENDPOINTS.page + '/1');
    await apiRequest(API_ENDPOINTS.page + '/2');

    // 测试无效页码
    await apiRequest(API_ENDPOINTS.page + '/0');
    await apiRequest(API_ENDPOINTS.page + '/999');

    // 测试无效格式
    await apiRequest(API_ENDPOINTS.page + '/abc');
}

/**
 * 测试搜索接口
 */
async function testSearch() {
    console.log('\n🧪 测试搜索接口');

    // 测试英文搜索
    await apiRequest(API_ENDPOINTS.search, { q: 'abandon' });

    // 测试中文搜索
    await apiRequest(API_ENDPOINTS.search, { q: '放弃' });

    // 测试短关键词
    await apiRequest(API_ENDPOINTS.search, { q: 'a' });

    // 测试空查询
    await apiRequest(API_ENDPOINTS.search, { q: '' });

    // 测试不存在的词
    await apiRequest(API_ENDPOINTS.search, { q: 'nonexistentword' });
}

/**
 * 测试统计信息接口
 */
async function testStats() {
    console.log('\n🧪 测试统计信息接口');

    await apiRequest(API_ENDPOINTS.stats);
}

/**
 * 性能测试
 */
async function testPerformance() {
    console.log('\n🧪 性能测试');

    const testCases = [
        { page: 1, pageSize: 10 },
        { page: 50, pageSize: 50 },
        { page: 100, pageSize: 100 }
    ];

    for (const testCase of testCases) {
        const startTime = performance.now();
        await apiRequest(API_ENDPOINTS.words, testCase);
        const endTime = performance.now();
        console.log(`⏱️  响应时间: ${(endTime - startTime).toFixed(2)}ms`);
    }
}

/**
 * 数据验证测试
 */
async function testDataValidation() {
    console.log('\n🧪 数据验证测试');

    // 测试词汇数据结构
    const response = await apiRequest(API_ENDPOINTS.words, { page: 1, pageSize: 5 });
    if (response && response.success && response.data && response.data.words) {
        const words = response.data.words;

        console.log('\n📋 词汇数据验证:');
        console.log('词汇数量:', words.length);
        console.log('分页信息:', {
            currentPage: response.data.current_page,
            totalPages: response.data.total_pages,
            totalWords: response.data.total_words
        });

        // 验证每个词汇对象的结构
        words.forEach((word, index) => {
            console.log(`\n词汇 ${index + 1}:`);
            console.log('  英文:', word.english);
            console.log('  中文:', word.chinese);
            console.log('  结构完整:', !!word.english && !!word.chinese);
        });
    }
}

/**
 * 错误处理测试
 */
async function testErrorHandling() {
    console.log('\n🧪 错误处理测试');

    // 测试各种错误情况
    const errorCases = [
        { endpoint: API_ENDPOINTS.words, params: { page: -1 } },
        { endpoint: API_ENDPOINTS.words, params: { pageSize: 101 } },
        { endpoint: API_ENDPOINTS.search, params: {} },
        { endpoint: API_ENDPOINTS.page + '/invalid' }
    ];

    for (const testCase of errorCases) {
        const endpoint = testCase.endpoint + (testCase.params ? '' : '');
        await apiRequest(endpoint, testCase.params);
    }
}

/**
 * 主测试函数
 */
async function runAllTests() {
    console.log('🚀 开始 Word Hero API 测试');
    console.log('==========================================');

    try {
        await testWordsPagination();
        await testPageByNumber();
        await testSearch();
        await testStats();
        await testPerformance();
        await testDataValidation();
        await testErrorHandling();

        console.log('\n✅ 所有测试完成！');
    } catch (error) {
        console.error('\n❌ 测试过程中发生错误:', error);
    }
}

// 如果直接运行此脚本，则执行所有测试
if (typeof window === 'undefined') {
    // Node.js 环境
    const fetch = require('node-fetch');
    runAllTests();
} else {
    // 浏览器环境
    console.log('💡 在浏览器控制台中运行 runAllTests() 来执行所有测试');
    window.runAllTests = runAllTests;
}

// 导出测试函数供其他模块使用
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        apiRequest,
        testWordsPagination,
        testPageByNumber,
        testSearch,
        testStats,
        testPerformance,
        testDataValidation,
        testErrorHandling,
        runAllTests
    };
}