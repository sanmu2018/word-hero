// Global variables
let currentPage = 1;
let totalPages = 1;
let isLoading = false;
let wordsVisible = true;
let translationsVisible = true;
// Speech synthesis is managed by the browser API directly
let selectedCard = null;
let selectedAccent = 'uk'; // Default to UK accent
let knownWords = new Set(); // Track known words
let wordTimestamps = new Map(); // Track word marking timestamps
let originalWordOrder = []; // Store original word order for restore functionality

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Initialize speech synthesis voices
    if ('speechSynthesis' in window) {
        // Load voices (some browsers need this)
        window.speechSynthesis.getVoices();
        // Listen for voices changed event
        window.speechSynthesis.onvoiceschanged = function() {
            // Voices loaded
        };
    }
    
    initializeApp();
});

function initializeApp() {
    // Get current page from URL or initial data
    const urlParams = new URLSearchParams(window.location.search);
    currentPage = parseInt(urlParams.get('page')) || 1;
    
    // Initialize from template data if available
    if (window.initialPageData) {
        currentPage = window.initialPageData.currentPage;
        totalPages = window.initialPageData.totalPages;
        
        // Set page size selector
        const pageSizeSelect = document.getElementById('pageSizeSelect');
        if (pageSizeSelect && window.initialPageData.pageSize) {
            pageSizeSelect.value = window.initialPageData.pageSize;
        }
    }
    
    // Set up event listeners
    setupEventListeners();
    
    // Add click events to initial template cards
    setupTemplateCardEvents();
    
    // Load known words from localStorage
    loadKnownWords();
    
    // Apply known word states to template cards
    applyKnownWordStates();
    
    // Load word timestamps
    loadWordTimestamps();
    
    // Store original word order for template-rendered cards
    storeOriginalOrderFromTemplate();
    
    // Update page info display
    updatePageInfo();
}

function setupEventListeners() {
    // Navigation buttons
    document.getElementById('firstBtn').addEventListener('click', goToFirstPage);
    document.getElementById('prevBtn').addEventListener('click', goToPreviousPage);
    document.getElementById('nextBtn').addEventListener('click', goToNextPage);
    document.getElementById('lastBtn').addEventListener('click', goToLastPage);
    document.getElementById('goBtn').addEventListener('click', goToPage);
    
      
    // Search functionality (modal-based)
    const searchBtn = document.getElementById('searchBtn');
    searchBtn.addEventListener('click', showSearchModal);
    
    // Page input
    const pageInput = document.getElementById('pageInput');
    pageInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            goToPage();
        }
    });
    
    // Page size selector
    const pageSizeSelect = document.getElementById('pageSizeSelect');
    pageSizeSelect.addEventListener('change', handlePageSizeChange);
    
    // Visibility toggle buttons
    const toggleWordsBtn = document.getElementById('toggleWords');
    const toggleTranslationsBtn = document.getElementById('toggleTranslations');
    const shuffleBtn = document.getElementById('shuffleBtn');
    const restoreBtn = document.getElementById('restoreBtn');
    const resetKnownWordsBtn = document.getElementById('resetKnownWords');
    
    toggleWordsBtn.addEventListener('click', toggleWordsVisibility);
    toggleTranslationsBtn.addEventListener('click', toggleTranslationsVisibility);
    shuffleBtn.addEventListener('click', shuffleCards);
    restoreBtn.addEventListener('click', restoreOriginalOrder);
    resetKnownWordsBtn.addEventListener('click', resetKnownWords);
    
    // Accent selection
    const accentUS = document.getElementById('accentUS');
    const accentUK = document.getElementById('accentUK');
    accentUS.addEventListener('click', () => selectAccent('us'));
    accentUK.addEventListener('click', () => selectAccent('uk'));
    
    // Initialize accent selection to UK (without notification)
    selectAccent('uk', false);
    
    // Keyboard shortcuts
    document.addEventListener('keydown', handleKeyboardShortcuts);
    
    // Close search modal when clicking outside
    document.addEventListener('click', function(e) {
        const searchModal = document.getElementById('searchModal');
        if (searchModal.style.display === 'block' && e.target === searchModal) {
            closeSearchModal();
        }
    });
}

// Process translation to handle truncation and tooltip
function processTranslation(chineseText) {
    if (!chineseText) return { display: '', hasTooltip: false };
    
    // Always show tooltip for all translations to ensure users can see complete content
    return { display: chineseText, hasTooltip: true };
}

// Handle page size change
function handlePageSizeChange() {
    const pageSizeSelect = document.getElementById('pageSizeSelect');
    const newPageSize = parseInt(pageSizeSelect.value);
    
    // Update URL with new page size
    const url = new URL(window.location);
    url.searchParams.set('pageSize', newPageSize);
    url.searchParams.set('page', 1); // Reset to first page
    window.history.pushState({}, '', url);
    
    // Load first page with new page size
    loadPage(1, newPageSize);
}

// Navigation functions
function goToFirstPage() {
    if (currentPage > 1) {
        navigateToPage(1);
    }
}

function goToPreviousPage() {
    if (currentPage > 1) {
        navigateToPage(currentPage - 1);
    }
}

function goToNextPage() {
    if (currentPage < totalPages) {
        navigateToPage(currentPage + 1);
    }
}

function goToLastPage() {
    if (currentPage < totalPages) {
        navigateToPage(totalPages);
    }
}

function goToPage() {
    const pageInput = document.getElementById('pageInput');
    const pageNumber = parseInt(pageInput.value);
    
    if (pageNumber >= 1 && pageNumber <= totalPages && pageNumber !== currentPage) {
        navigateToPage(pageNumber);
    } else {
        pageInput.value = currentPage;
    }
}

function navigateToPage(pageNumber) {
    if (isLoading) return;
    
    currentPage = pageNumber;
    
    // Update URL
    const url = new URL(window.location);
    url.searchParams.set('page', currentPage);
    window.history.pushState({}, '', url);
    
    // Load page data
    loadPage(currentPage);
}

function loadPage(pageNumber, pageSize) {
    if (isLoading) return;
    
    // Get page size from selector or use default
    if (!pageSize) {
        const pageSizeSelect = document.getElementById('pageSizeSelect');
        pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;
    }
    
    fetch(`/api/words?page=${pageNumber}&pageSize=${pageSize}`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                updatePageContent(data.data);
                updatePageInfo();
                scrollToTop();
            } else {
                showError('Failed to load page: ' + data.error);
            }
        })
        .catch(error => {
            showError('Network error: ' + error.message);
        });
}

function updatePageContent(data) {
    const vocabularyGrid = document.getElementById('vocabularyGrid');
    
    // Clear existing content
    vocabularyGrid.innerHTML = '';
    
    // Store original word order for this page
    originalWordOrder = data.words.map((word, index) => ({
        word: word,
        displayNumber: data.startIndex + index
    }));
    
    // Add vocabulary cards
    data.words.forEach((word, index) => {
        const card = createVocabularyCard(word, data.startIndex + index);
        vocabularyGrid.appendChild(card);
    });
    
    // Apply current visibility settings to new content
    applyVisibilitySettings();
    
    // Update navigation state
    updateNavigationState(data);
}

function createVocabularyCard(word, displayNumber) {
    const card = document.createElement('div');
    card.className = 'vocabulary-card';
    card.setAttribute('data-word', word.English);
    
    // Process Chinese translation
    const chineseTranslation = processTranslation(word.Chinese);
    
    const isKnown = knownWords.has(word.English);
    
    card.innerHTML = `
        <button class="action-btn" onclick="toggleWordMenu(event, '${escapeJsString(word.English)}')" title="更多操作">
            <i class="fas fa-ellipsis-v"></i>
        </button>
        <div class="card-content">
            <div class="english-word">
                ${escapeHtml(word.English)}
                <button class="speaker-btn" onclick="speakWord('${escapeJsString(word.English)}', this)" title="播放发音">
                    <i class="fas fa-volume-up"></i>
                </button>
            </div>
            <div class="chinese-meaning">
                <span class="translation-text">${escapeHtml(chineseTranslation.display)}</span>
                ${chineseTranslation.hasTooltip ? `<div class="translation-tooltip">${escapeHtml(word.Chinese)}</div>` : ''}
            </div>
        </div>
        <div class="word-menu" id="menu-${escapeJsString(word.English)}" style="display: none;">
            <div class="menu-item known-action" data-word="${escapeJsString(word.English)}" data-action="known">
                <i class="fas fa-check"></i>
                标为认识
            </div>
            <div class="menu-item unknown-action" data-word="${escapeJsString(word.English)}" data-action="unknown" style="display: none;">
                <i class="fas fa-times"></i>
                标为不认识
            </div>
            <div class="menu-item detail-action" data-word="${escapeJsString(word.English)}" data-action="detail">
                <i class="fas fa-search-plus"></i>
                细品
            </div>
        </div>
    `;
    
    // Add click event for card selection
    card.addEventListener('click', function(e) {
        // Don't select if clicking on speaker button or action button
        if (!e.target.closest('.speaker-btn') && !e.target.closest('.action-btn') && !e.target.closest('.word-menu')) {
            selectCard(this);
        }
    });
    
    // Apply known word state if needed
    if (isKnown) {
        card.classList.add('known-word');
    }
    
    return card;
}

function updateNavigationState(data) {
    currentPage = data.currentPage;
    totalPages = data.totalPages;
    
    // Update buttons
    document.getElementById('firstBtn').disabled = !data.hasPrev;
    document.getElementById('prevBtn').disabled = !data.hasPrev;
    document.getElementById('nextBtn').disabled = !data.hasNext;
    document.getElementById('lastBtn').disabled = !data.hasNext;
    
    // Update page input
    const pageInput = document.getElementById('pageInput');
    if (pageInput) {
        pageInput.value = currentPage;
        pageInput.max = totalPages;
    }
}

function updatePageInfo() {
    // Page info display removed, only update page input
    const pageInput = document.getElementById('pageInput');
    if (pageInput) {
        pageInput.value = currentPage;
        pageInput.max = totalPages;
    }
}


// Modal functions
function showStats() {
    fetch('/api/stats')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayStats(data.data);
                showModal('statsModal');
            } else {
                showError('Failed to load stats: ' + data.error);
            }
        })
        .catch(error => {
            showError('Stats error: ' + error.message);
        });
}

function displayStats(data) {
    const statsContent = document.getElementById('statsContent');
    statsContent.innerHTML = `
        <div class="stats-grid">
            <div class="stat-card">
                <h4>总词汇数</h4>
                <div class="stat-value">${data.totalWords}</div>
            </div>
            <div class="stat-card">
                <h4>总页数</h4>
                <div class="stat-value">${data.totalPages}</div>
            </div>
            <div class="stat-card">
                <h4>每页显示</h4>
                <div class="stat-value">${data.pageSize}</div>
            </div>
            <div class="stat-card">
                <h4>数据来源</h4>
                <div class="stat-value" style="font-size: 1rem;">${data.fileSource}</div>
            </div>
        </div>
    `;
}

function showHelp() {
    showModal('helpModal');
}

function showModal(modalId) {
    document.getElementById(modalId).style.display = 'block';
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

// Utility functions
// Loading functions removed for smoother pagination experience

function showError(message) {
    // Create a toast notification
    const toast = document.createElement('div');
    toast.className = 'toast error';
    toast.innerHTML = `
        <i class="fas fa-exclamation-circle"></i>
        <span>${message}</span>
    `;
    document.body.appendChild(toast);
    
    // Show toast
    setTimeout(() => toast.classList.add('show'), 100);
    
    // Hide and remove toast
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => document.body.removeChild(toast), 300);
    }, 3000);
}

function scrollToTop() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function escapeJsString(text) {
    return text.replace(/['"\\]/g, '\\$&').replace(/\n/g, '\\n').replace(/\r/g, '\\r');
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Keyboard shortcuts
function handleKeyboardShortcuts(e) {
    // Ignore if user is typing in an input field
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
        return;
    }
    
    switch(e.key) {
        case 'ArrowLeft':
            e.preventDefault();
            goToPreviousPage();
            break;
        case 'ArrowRight':
            e.preventDefault();
            goToNextPage();
            break;
        case 'f':
            if (e.ctrlKey || e.metaKey) {
                e.preventDefault();
                showSearchModal();
            }
            break;
        case 'Escape':
            closeModal('statsModal');
            closeModal('helpModal');
            closeSearchModal();
            break;
        case 'w':
            if (e.ctrlKey || e.metaKey) {
                e.preventDefault();
                toggleWordsVisibility();
            }
            break;
        case 't':
            if (e.ctrlKey || e.metaKey) {
                e.preventDefault();
                toggleTranslationsVisibility();
            }
            break;
        case 'r':
            if (e.ctrlKey || e.metaKey) {
                if (e.shiftKey) {
                    e.preventDefault();
                    resetKnownWords();
                } else {
                    e.preventDefault();
                    shuffleCards();
                }
            }
            break;
    }
}

// Add toast styles dynamically
const toastStyles = document.createElement('style');
toastStyles.textContent = `
    .toast {
        position: fixed;
        top: 20px;
        right: 20px;
        background: #dc3545;
        color: white;
        padding: 12px 20px;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        display: flex;
        align-items: center;
        gap: 10px;
        z-index: 3000;
        opacity: 0;
        transform: translateX(100%);
        transition: all 0.3s ease;
    }
    
    .toast.show {
        opacity: 1;
        transform: translateX(0);
    }
    
    .toast.error {
        background: #dc3545;
    }
    
    .toast.success {
        background: #28a745;
    }
`;
document.head.appendChild(toastStyles);

// Handle browser back/forward
window.addEventListener('popstate', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const page = parseInt(urlParams.get('page')) || 1;
    
    if (page !== currentPage) {
        currentPage = page;
        loadPage(currentPage);
    }
});

// Ensure tooltips work for dynamically generated content
document.addEventListener('DOMContentLoaded', function() {
    // Enable tooltips for existing elements
    enableTooltips();
    
    // Watch for dynamically added content
    const observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            if (mutation.addedNodes) {
                enableTooltips();
            }
        });
    });
    
    observer.observe(document.body, {
        childList: true,
        subtree: true
    });
});

function enableTooltips() {
    const chineseMeanings = document.querySelectorAll('.chinese-meaning');
    chineseMeanings.forEach(function(element) {
        // Make sure tooltip functionality is working
        element.addEventListener('mouseenter', function() {
            const tooltip = this.querySelector('.translation-tooltip');
            if (tooltip) {
                tooltip.style.opacity = '1';
                tooltip.style.visibility = 'visible';
            }
        });
        
        element.addEventListener('mouseleave', function() {
            const tooltip = this.querySelector('.translation-tooltip');
            if (tooltip) {
                tooltip.style.opacity = '0';
                tooltip.style.visibility = 'hidden';
            }
        });
    });
}

// Visibility toggle functions
function toggleWordsVisibility() {
    wordsVisible = !wordsVisible;
    const toggleBtn = document.getElementById('toggleWords');
    const englishWords = document.querySelectorAll('.english-word');
    
    if (wordsVisible) {
        // Show words
        englishWords.forEach(word => {
            word.classList.remove('hidden');
        });
        toggleBtn.classList.remove('words-hidden');
        toggleBtn.innerHTML = '<i class="fas fa-eye"></i> 单词';
    } else {
        // Hide words
        englishWords.forEach(word => {
            word.classList.add('hidden');
        });
        toggleBtn.classList.add('words-hidden');
        toggleBtn.innerHTML = '<i class="fas fa-eye-slash"></i> 单词';
    }
    
    // Update card classes for proper layout
    updateCardLayoutClasses();
}

function toggleTranslationsVisibility() {
    translationsVisible = !translationsVisible;
    const toggleBtn = document.getElementById('toggleTranslations');
    const chineseMeanings = document.querySelectorAll('.chinese-meaning');
    
    if (translationsVisible) {
        // Show translations
        chineseMeanings.forEach(meaning => {
            meaning.classList.remove('hidden');
        });
        toggleBtn.classList.remove('translations-hidden');
        toggleBtn.innerHTML = '<i class="fas fa-eye"></i> 翻译';
    } else {
        // Hide translations
        chineseMeanings.forEach(meaning => {
            meaning.classList.add('hidden');
        });
        toggleBtn.classList.add('translations-hidden');
        toggleBtn.innerHTML = '<i class="fas fa-eye-slash"></i> 翻译';
    }
    
    // Update card classes for proper layout
    updateCardLayoutClasses();
}

function updateCardLayoutClasses() {
    const vocabularyCards = document.querySelectorAll('.vocabulary-card');
    
    vocabularyCards.forEach(card => {
        card.classList.remove('words-only', 'translations-only');
        
        if (!wordsVisible && translationsVisible) {
            card.classList.add('translations-only');
        } else if (wordsVisible && !translationsVisible) {
            card.classList.add('words-only');
        }
    });
}

// Apply visibility settings to dynamically loaded content
function applyVisibilitySettings() {
    if (!wordsVisible) {
        document.querySelectorAll('.english-word').forEach(word => {
            word.classList.add('hidden');
        });
        document.getElementById('toggleWords').classList.add('words-hidden');
        document.getElementById('toggleWords').innerHTML = '<i class="fas fa-eye-slash"></i> 单词';
    }
    
    if (!translationsVisible) {
        document.querySelectorAll('.chinese-meaning').forEach(meaning => {
            meaning.classList.add('hidden');
        });
        document.getElementById('toggleTranslations').classList.add('translations-hidden');
        document.getElementById('toggleTranslations').innerHTML = '<i class="fas fa-eye-slash"></i> 翻译';
    }
    
    updateCardLayoutClasses();
}

// Shuffle functionality
function shuffleCards() {
    const vocabularyGrid = document.getElementById('vocabularyGrid');
    const cards = Array.from(vocabularyGrid.children);
    
    if (cards.length <= 1) {
        return; // No need to shuffle if there's only one or no cards
    }
    
    // Add shuffle animation to all cards
    cards.forEach(card => {
        card.classList.add('shuffling');
    });
    
    // Fisher-Yates shuffle algorithm
    for (let i = cards.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [cards[i], cards[j]] = [cards[j], cards[i]];
    }
    
    // Re-append cards in new order
    cards.forEach(card => {
        vocabularyGrid.appendChild(card);
    });
    
    // Remove animation class after animation completes
    setTimeout(() => {
        cards.forEach(card => {
            card.classList.remove('shuffling');
        });
    }, 500);
    
    // Show shuffle feedback
    showShuffleFeedback();
}

function showShuffleFeedback() {
    const shuffleBtn = document.getElementById('shuffleBtn');
    const originalHTML = shuffleBtn.innerHTML;
    
    // Show feedback
    shuffleBtn.innerHTML = '<i class="fas fa-check"></i> 已打乱';
    shuffleBtn.style.background = '#28a745';
    
    // Reset after 1 second
    setTimeout(() => {
        shuffleBtn.innerHTML = originalHTML;
        shuffleBtn.style.background = '';
    }, 1000);
}

// Search Modal Functions
function showSearchModal() {
    const modal = document.getElementById('searchModal');
    const searchInput = document.getElementById('searchModalInput');
    const searchBtn = document.getElementById('searchModalBtn');
    
    modal.style.display = 'block';
    
    // Focus on search input after a short delay
    setTimeout(() => {
        searchInput.focus();
        searchInput.value = '';
        resetSearchResults();
    }, 100);
    
    // Set up modal event listeners
    setupModalSearchListeners();
}

function closeSearchModal() {
    const modal = document.getElementById('searchModal');
    modal.style.display = 'none';
    resetSearchResults();
}

function setupModalSearchListeners() {
    const searchInput = document.getElementById('searchModalInput');
    const searchBtn = document.getElementById('searchModalBtn');
    
    // Remove existing listeners to prevent duplicates
    searchInput.removeEventListener('input', handleModalSearch);
    searchInput.removeEventListener('keypress', handleModalSearchKeypress);
    searchBtn.removeEventListener('click', handleModalSearchClick);
    
    // Add new listeners
    searchInput.addEventListener('input', debounce(handleModalSearch, 300));
    searchInput.addEventListener('keypress', handleModalSearchKeypress);
    searchBtn.addEventListener('click', handleModalSearchClick);
}

function handleModalSearchKeypress(e) {
    if (e.key === 'Enter') {
        handleModalSearchClick();
    }
}

function handleModalSearchClick() {
    const query = document.getElementById('searchModalInput').value.trim();
    if (query.length >= 2) {
        performModalSearch(query);
    }
}

function handleModalSearch() {
    const query = document.getElementById('searchModalInput').value.trim();
    
    if (query.length === 0) {
        resetSearchResults();
        return;
    }
    
    if (query.length < 2) {
        updateSearchInfo('输入至少2个字符开始搜索');
        return;
    }
    
    performModalSearch(query);
}

function performModalSearch(query) {
    const resultsContainer = document.getElementById('searchResultsContainer');
    const resultsList = document.getElementById('searchModalResults');
    const searchInfo = document.getElementById('searchResultCount');
    
    // Show loading state
    resultsContainer.style.display = 'block';
    resultsList.innerHTML = `
        <div class="search-loading">
            <i class="fas fa-spinner fa-spin"></i>
            <span>搜索中...</span>
        </div>
    `;
    searchInfo.textContent = '正在搜索...';
    
    fetch(`/api/search?q=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayModalSearchResults(data.data, query);
            } else {
                showSearchError('搜索失败: ' + data.error);
            }
        })
        .catch(error => {
            showSearchError('网络错误: ' + error.message);
        });
}

function displayModalSearchResults(data, query) {
    const resultsList = document.getElementById('searchModalResults');
    const searchInfo = document.getElementById('searchResultCount');
    
    resultsList.innerHTML = '';
    
    if (!data.results || !Array.isArray(data.results) || data.results.length === 0) {
        resultsList.innerHTML = `
            <div class="search-results-empty">
                <i class="fas fa-search"></i>
                <div>没有找到与 "<strong>${escapeHtml(query)}</strong>" 相关的词汇</div>
            </div>
        `;
        searchInfo.textContent = '没有找到相关词汇';
    } else {
        data.results.forEach(word => {
            const resultItem = document.createElement('div');
            resultItem.className = 'search-result-item';
            resultItem.innerHTML = `
                <div class="search-result-word">${escapeHtml(word.English)}</div>
                <div class="search-result-meaning">${escapeHtml(word.Chinese)}</div>
            `;
            resultItem.addEventListener('click', () => {
                closeSearchModal();
            });
            resultsList.appendChild(resultItem);
        });
        searchInfo.textContent = `找到 ${data.results.length} 个相关词汇`;
    }
}

function showSearchError(message) {
    const resultsList = document.getElementById('searchModalResults');
    const searchInfo = document.getElementById('searchResultCount');
    
    resultsList.innerHTML = `
        <div class="search-results-empty">
            <i class="fas fa-exclamation-triangle"></i>
            <div>${escapeHtml(message)}</div>
        </div>
    `;
    searchInfo.textContent = '搜索出现错误';
}

function updateSearchInfo(message) {
    const searchInfo = document.getElementById('searchResultCount');
    searchInfo.textContent = message;
}

function resetSearchResults() {
    const resultsContainer = document.getElementById('searchResultsContainer');
    const resultsList = document.getElementById('searchModalResults');
    const searchInfo = document.getElementById('searchResultCount');
    
    resultsContainer.style.display = 'none';
    resultsList.innerHTML = '';
    searchInfo.textContent = '输入至少2个字符开始搜索';
}

// Accent selection function
function selectAccent(accent, showNotification = true) {
    selectedAccent = accent;
    
    // Update button states
    const accentUS = document.getElementById('accentUS');
    const accentUK = document.getElementById('accentUK');
    
    if (accent === 'us') {
        accentUS.classList.add('active');
        accentUK.classList.remove('active');
    } else {
        accentUS.classList.remove('active');
        accentUK.classList.add('active');
    }
    
    // Show feedback only if requested
    if (showNotification) {
        const accentName = accent === 'us' ? '美式英语' : '英式英语';
        showToast(`已切换到${accentName}发音`, 'success');
    }
}

// Card selection function
function selectCard(card) {
    // Remove selection from previous card
    if (selectedCard) {
        selectedCard.classList.remove('selected');
    }
    
    // Select new card
    card.classList.add('selected');
    selectedCard = card;
}

// Text-to-speech function using browser speech synthesis
function speakWord(word, button) {
    // Toggle speaking state
    if (button.classList.contains('speaking')) {
        // Stop speaking
        stopSpeech(button);
        return;
    }
    
    // Check if browser supports speech synthesis
    if (!('speechSynthesis' in window)) {
        showError('您的浏览器不支持语音播放功能');
        return;
    }
    
    // Show speaking state
    button.classList.add('speaking');
    button.innerHTML = '<i class="fas fa-stop"></i>';
    button.title = '停止播放';
    
    // Speak the word
    speakWordBrowser(word, button);
}

// Stop speech synthesis
function stopSpeech(button) {
    if ('speechSynthesis' in window) {
        window.speechSynthesis.cancel();
        
        // Small delay to ensure cancellation is processed
        setTimeout(() => {
            resetButtonState(button);
        }, 100);
    }
}

// Reset button to normal state
function resetButtonState(button) {
    button.classList.remove('speaking');
    button.innerHTML = '<i class="fas fa-volume-up"></i>';
    button.title = '播放发音';
}

// Speak word using browser speech synthesis
function speakWordBrowser(word, button) {
    // Cancel any existing speech first
    window.speechSynthesis.cancel();
    
    // Small delay to ensure cancellation
    setTimeout(() => {
        const utterance = new SpeechSynthesisUtterance(word);
        
        // Configure speech settings
        utterance.rate = 0.8;
        utterance.pitch = 1.0;
        utterance.volume = 1.0;
        utterance.lang = selectedAccent === 'us' ? 'en-US' : 'en-GB';
        
        // Get voices and try to find a good one
        const voices = window.speechSynthesis.getVoices();
        let selectedVoice = null;
        
        // Try to find a voice matching the selected accent
        for (const voice of voices) {
            const lang = selectedAccent === 'us' ? 'en-US' : 'en-GB';
            if (voice.lang === lang || voice.lang.startsWith('en')) {
                selectedVoice = voice;
                break;
            }
        }
        
        if (selectedVoice) {
            utterance.voice = selectedVoice;
        }
        
        // Set up event handlers
        utterance.onend = function() {
            resetButtonState(button);
        };
        
        utterance.onerror = function(event) {
            resetButtonState(button);
            
            // Only show error for non-interrupted errors
            if (event.error !== 'interrupted') {
                showError('语音播放失败，请重试');
            }
        };
        
        // Try to speak
        try {
            window.speechSynthesis.speak(utterance);
            
            // Backup timeout in case onend doesn't fire
            setTimeout(() => {
                if (button.classList.contains('speaking')) {
                    resetButtonState(button);
                }
            }, 5000);
            
        } catch (error) {
            resetButtonState(button);
            showError('语音播放失败，请重试');
        }
    }, 100);
}




// Setup click events for template-rendered cards
function setupTemplateCardEvents() {
    const templateCards = document.querySelectorAll('.vocabulary-card');
    templateCards.forEach(card => {
        if (!card.hasAttribute('data-event-setup')) {
            card.addEventListener('click', function(e) {
                // Don't select if clicking on speaker button or action button
                if (!e.target.closest('.speaker-btn') && !e.target.closest('.action-btn')) {
                    selectCard(this);
                }
            });
            card.setAttribute('data-event-setup', 'true');
        }
    });
}

// Store original word order from template-rendered cards
function storeOriginalOrderFromTemplate() {
    const templateCards = document.querySelectorAll('.vocabulary-card');
    if (templateCards.length > 0) {
        originalWordOrder = Array.from(templateCards).map((card, index) => {
            const englishWord = card.getAttribute('data-word');
            const chineseText = card.querySelector('.translation-text').textContent;
            return {
                word: {
                    English: englishWord,
                    Chinese: chineseText
                },
                displayNumber: window.initialPageData ? window.initialPageData.startIndex + index : index + 1
            };
        });
        console.log('Stored original word order from template:', originalWordOrder.length, 'words');
    }
}

// Toast notification function
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <i class="fas fa-${type === 'success' ? 'check' : 'info'}-circle"></i>
        <span>${message}</span>
    `;
    document.body.appendChild(toast);
    
    // Show toast
    setTimeout(() => toast.classList.add('show'), 100);
    
    // Hide and remove toast
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => document.body.removeChild(toast), 300);
    }, 2000);
}

// Word action functions
let currentActionWord = null; // Track current word for action modal

function toggleWordMenu(event, word) {
    event.stopPropagation();
    
    currentActionWord = word;
    
    // Update modal title
    const title = document.getElementById('wordActionTitle');
    if (title) {
        title.textContent = `单词 - ${word}`;
    }
    
    // Show modal
    const modal = document.getElementById('wordActionModal');
    if (modal) {
        modal.style.display = 'block';
        
        // Update button states based on current word status
        updateActionModalButtons(word);
    } else {
        console.error('Word action modal not found');
    }
}

function updateActionModalButtons(word) {
    const knownBtn = document.querySelector('.known-btn');
    const unknownBtn = document.querySelector('.unknown-btn');
    
    if (knownWords.has(word)) {
        knownBtn.style.display = 'none';
        unknownBtn.style.display = 'flex';
    } else {
        knownBtn.style.display = 'flex';
        unknownBtn.style.display = 'none';
    }
}

function markCurrentWordAsKnown() {
    if (currentActionWord) {
        markWordAsKnown(currentActionWord);
        closeWordActionModal();
    }
}

function markCurrentWordAsUnknown() {
    if (currentActionWord) {
        markWordAsUnknown(currentActionWord);
        closeWordActionModal();
    }
}

function showCurrentWordDetail() {
    if (currentActionWord) {
        closeWordActionModal();
        showWordDetail(currentActionWord);
    }
}

function closeWordActionModal() {
    const modal = document.getElementById('wordActionModal');
    modal.style.display = 'none';
    currentActionWord = null;
}


function markWordAsKnown(word) {
    knownWords.add(word);
    wordTimestamps.set(word, new Date().toISOString());
    saveKnownWords();
    saveWordTimestamps();
    updateWordCardAppearance(word, true);
    showToast(`"${word}" 已标记为认识`, 'success');
    closeWordActionModal();
}

function markWordAsUnknown(word) {
    knownWords.delete(word);
    wordTimestamps.delete(word);
    saveKnownWords();
    saveWordTimestamps();
    updateWordCardAppearance(word, false);
    showToast(`"${word}" 已标记为不认识`, 'info');
    closeWordActionModal();
}

function updateWordCardAppearance(word, isKnown) {
    const card = document.querySelector(`.vocabulary-card[data-word="${word}"]`);
    if (card) {
        if (isKnown) {
            card.classList.add('known-word');
        } else {
            card.classList.remove('known-word');
        }
    }
}

function resetKnownWords() {
    // Show reset options modal instead of confirmation dialog
    showResetOptionsModal();
}

function showResetOptionsModal() {
    const modal = document.getElementById('resetOptionsModal');
    
    // Update the modal with current page information
    document.getElementById('currentPageNumber').textContent = currentPage;
    document.getElementById('totalPagesNumber').textContent = totalPages;
    document.getElementById('totalWordsNumber').textContent = window.initialPageData ? window.initialPageData.totalWords : '3673';
    
    modal.style.display = 'block';
}

function closeResetOptionsModal() {
    const modal = document.getElementById('resetOptionsModal');
    modal.style.display = 'none';
}

function resetCurrentPage() {
    // Get current page words
    const currentCards = document.querySelectorAll('.vocabulary-card');
    let resetCount = 0;
    
    currentCards.forEach(card => {
        const word = card.getAttribute('data-word');
        if (word && knownWords.has(word)) {
            knownWords.delete(word);
            wordTimestamps.delete(word);
            card.classList.remove('known-word');
            resetCount++;
        }
    });
    
    // Save changes
    saveKnownWords();
    saveWordTimestamps();
    
    // Close modal
    closeResetOptionsModal();
    
    // Show feedback
    showToast(`已重置当前页面的 ${resetCount} 个单词标记`, 'success');
    
    // Restore original order for current page
    if (originalWordOrder.length > 0) {
        restoreOriginalOrder();
    }
}

function resetAllPages() {
    // Confirm before resetting all pages
    if (confirm('确定要重置所有页面的单词标记状态吗？这将清除所有已认识单词的记录。')) {
        const totalKnown = knownWords.size;
        
        // Clear all known words and timestamps
        knownWords.clear();
        wordTimestamps.clear();
        
        // Save to localStorage
        saveKnownWords();
        saveWordTimestamps();
        
        // Close modal
        closeResetOptionsModal();
        
        // Reload current page to refresh display
        const pageSizeSelect = document.getElementById('pageSizeSelect');
        const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;
        loadPage(currentPage, pageSize);
        
        // Show feedback
        showToast(`已重置全部 ${totalKnown} 个已标记单词`, 'success');
    }
}

function showWordDetail(word) {
    const modal = document.getElementById('wordDetailModal');
    const title = document.getElementById('wordDetailTitle');
    const content = document.getElementById('wordDetailContent');
    
    // Find the word data
    const card = document.querySelector(`.vocabulary-card[data-word="${word}"]`);
    const chineseText = card ? card.querySelector('.translation-text').textContent : '';
    const isKnown = knownWords.has(word);
    
    title.textContent = '单词详情';
    content.innerHTML = `
        <div class="word-detail-header">
            <div class="word-detail-word">${escapeHtml(word)}</div>
            <div class="word-detail-pronunciation">
                <button class="speaker-btn" onclick="speakWord('${escapeHtml(word)}', this)" title="播放发音" style="position: relative; top: auto; right: auto; opacity: 1; transform: scale(1); margin-left: 10px;">
                    <i class="fas fa-volume-up"></i>
                </button>
                /${getPhonetic(word)}/
            </div>
            <div class="word-detail-meaning">${escapeHtml(chineseText)}</div>
            <div class="word-detail-actions">
                ${isKnown ? 
                    `<button class="word-detail-btn danger" onclick="markWordAsUnknown('${escapeHtml(word)}'); updateWordDetailModal('${escapeHtml(word)}')">标为不认识</button>` :
                    `<button class="word-detail-btn success" onclick="markWordAsKnown('${escapeHtml(word)}'); updateWordDetailModal('${escapeHtml(word)}')">标为认识</button>`
                }
                <button class="word-detail-btn primary" onclick="speakWord('${escapeHtml(word)}', this)">朗读单词</button>
            </div>
        </div>
        <div class="word-detail-section">
            <h3>学习进度</h3>
            <div class="progress-info">
                <div class="progress-item">
                    <span>状态：</span>
                    <span class="${isKnown ? 'status-known' : 'status-unknown'}">${isKnown ? '已认识' : '学习中'}</span>
                </div>
                <div class="progress-item">
                    <span>标记时间：</span>
                    <span>${isKnown ? new Date().toLocaleString('zh-CN') : '未标记'}</span>
                </div>
            </div>
        </div>
    `;
    
    modal.style.display = 'block';
    closeAllMenus();
}

function updateWordDetailModal(word) {
    showWordDetail(word);
}

function closeWordDetailModal() {
    const modal = document.getElementById('wordDetailModal');
    modal.style.display = 'none';
}

function showStatisticsModal() {
    const modal = document.getElementById('statisticsModal');
    modal.style.display = 'block';
    updateStatisticsData();
}

function closeStatisticsModal() {
    const modal = document.getElementById('statisticsModal');
    modal.style.display = 'none';
}

function updateStatisticsData() {
    // Update overview statistics
    const totalKnownWords = knownWords.size;
    const dailyStats = calculateDailyWordCount();
    const totalLearningDays = dailyStats.length;
    const averageWordsPerDay = totalLearningDays > 0 ? Math.round(totalKnownWords / totalLearningDays) : 0;

    document.getElementById('totalKnownWords').textContent = totalKnownWords;
    document.getElementById('totalLearningDays').textContent = totalLearningDays;
    document.getElementById('averageWordsPerDay').textContent = averageWordsPerDay;

    // Update chart
    updateDailyProgressChart(dailyStats);
}

function calculateDailyWordCount() {
    const dailyCount = new Map();
    
    wordTimestamps.forEach((timestamp) => {
        const date = new Date(timestamp).toISOString().split('T')[0];
        dailyCount.set(date, (dailyCount.get(date) || 0) + 1);
    });

    // Get last 30 days of data
    const result = [];
    const today = new Date();
    for (let i = 29; i >= 0; i--) {
        const date = new Date(today);
        date.setDate(date.getDate() - i);
        const dateStr = date.toISOString().split('T')[0];
        const count = dailyCount.get(dateStr) || 0;
        
        result.push({
            date: dateStr,
            count: count,
            displayDate: `${date.getMonth() + 1}/${date.getDate()}`
        });
    }

    return result;
}

function updateDailyProgressChart(dailyStats) {
    const chartDom = document.getElementById('dailyProgressChart');
    if (!chartDom) return;

    // Initialize chart
    const myChart = echarts.init(chartDom);
    
    const option = {
        title: {
            text: '每日学习进度',
            left: 'center',
            top: 10
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'cross'
            },
            formatter: function(params) {
                const data = params[0];
                return `${data.name}<br/>新增单词: ${data.value}个`;
            }
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: dailyStats.map(item => item.displayDate),
            axisLabel: {
                rotate: 45
            }
        },
        yAxis: {
            type: 'value',
            name: '单词数量',
            minInterval: 1
        },
        series: [
            {
                name: '新增单词',
                type: 'line',
                smooth: true,
                data: dailyStats.map(item => item.count),
                itemStyle: {
                    color: '#28a745'
                },
                areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        {
                            offset: 0,
                            color: 'rgba(40, 167, 69, 0.3)'
                        },
                        {
                            offset: 1,
                            color: 'rgba(40, 167, 69, 0.1)'
                        }
                    ])
                },
                emphasis: {
                    focus: 'series'
                }
            }
        ]
    };

    myChart.setOption(option);

    // Handle window resize
    window.addEventListener('resize', function() {
        myChart.resize();
    });
}

function getPhonetic(word) {
    // Simple phonetic approximation - in real app this would come from dictionary API
    const phonetics = {
        'hello': 'həˈloʊ',
        'world': 'wɜːrld',
        'computer': 'kəmˈpjuːtər',
        'language': 'ˈlæŋɡwɪdʒ',
        'english': 'ˈɪŋɡlɪʃ',
        'chinese': 'ˌtʃaɪˈniːz',
        'student': 'ˈstuːdnt',
        'teacher': 'ˈtiːtʃər',
        'school': 'skuːl',
        'university': 'ˌjuːnɪˈvɜːrsəti'
    };
    return phonetics[word.toLowerCase()] || 'pronunciation';
}

function loadKnownWords() {
    const saved = localStorage.getItem('knownWords');
    if (saved) {
        try {
            const parsed = JSON.parse(saved);
            knownWords = new Set(parsed);
        } catch (e) {
            console.error('Error loading known words:', e);
            knownWords = new Set();
        }
    }
}

function saveKnownWords() {
    try {
        localStorage.setItem('knownWords', JSON.stringify(Array.from(knownWords)));
    } catch (e) {
        console.error('Error saving known words:', e);
    }
}

function loadWordTimestamps() {
    const saved = localStorage.getItem('wordTimestamps');
    if (saved) {
        try {
            const parsed = JSON.parse(saved);
            wordTimestamps = new Map(parsed);
        } catch (e) {
            console.error('Error loading word timestamps:', e);
            wordTimestamps = new Map();
        }
    }
}

function saveWordTimestamps() {
    try {
        localStorage.setItem('wordTimestamps', JSON.stringify(Array.from(wordTimestamps)));
    } catch (e) {
        console.error('Error saving word timestamps:', e);
    }
}

function applyKnownWordStates() {
    knownWords.forEach(word => {
        updateWordCardAppearance(word, true);
    });
}

// Setup click events for menu items
document.addEventListener('click', function(e) {
    // Close action modal when clicking outside
    const actionModal = document.getElementById('wordActionModal');
    if (actionModal && actionModal.style.display === 'block' && e.target === actionModal) {
        closeWordActionModal();
    }
    
    // Prevent card selection when clicking action button
    if (e.target.closest('.action-btn')) {
        e.stopPropagation();
    }
});

// Add CSS for word detail modal
const detailStyles = document.createElement('style');
detailStyles.textContent = `
    .word-detail-section {
        margin-bottom: 25px;
    }
    
    .word-detail-section h3 {
        color: #2c3e50;
        margin-bottom: 15px;
        font-size: 1.2rem;
        border-bottom: 1px solid #f0f0f0;
        padding-bottom: 8px;
    }
    
    .progress-info {
        background: #f8f9fa;
        padding: 15px;
        border-radius: 8px;
    }
    
    .progress-item {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;
        font-size: 0.95rem;
    }
    
    .progress-item:last-child {
        margin-bottom: 0;
    }
    
    .status-known {
        color: #28a745;
        font-weight: 600;
    }
    
    .status-unknown {
        color: #6c757d;
        font-weight: 600;
    }
    
    .speaker-btn {
        position: relative !important;
        top: auto !important;
        right: auto !important;
        opacity: 1 !important;
        transform: scale(1) !important;
        display: inline-flex !important;
        margin-left: 10px;
    }
`;
document.head.appendChild(detailStyles);

// Update all speaker button states (fallback for backward compatibility)
function updateSpeakerButtons(isSpeaking) {
    const speakerButtons = document.querySelectorAll('.speaker-btn');
    speakerButtons.forEach(button => {
        updateSpeakerButton(button, isSpeaking);
    });
}

// Restore original order functionality
function restoreOriginalOrder() {
    const vocabularyGrid = document.getElementById('vocabularyGrid');
    const cards = Array.from(vocabularyGrid.children);
    
    console.log('Restore function called:');
    console.log('- Current cards:', cards.length);
    console.log('- Original order length:', originalWordOrder.length);
    
    if (cards.length <= 1 || originalWordOrder.length === 0) {
        console.log('Restore aborted: not enough cards or no original order');
        return; // No need to restore if there's only one card or no original order
    }
    
    // Add restore animation to all cards
    cards.forEach(card => {
        card.classList.add('restoring');
    });
    
    // Clear current cards
    vocabularyGrid.innerHTML = '';
    
    // Re-add cards in original order
    originalWordOrder.forEach(item => {
        const card = createVocabularyCard(item.word, item.displayNumber);
        vocabularyGrid.appendChild(card);
    });
    
    // Remove animation class after animation completes
    setTimeout(() => {
        const newCards = Array.from(vocabularyGrid.children);
        newCards.forEach(card => {
            card.classList.remove('restoring');
        });
    }, 500);
    
    // Show restore feedback
    showRestoreFeedback();
}

function showRestoreFeedback() {
    const restoreBtn = document.getElementById('restoreBtn');
    const originalHTML = restoreBtn.innerHTML;
    
    // Show feedback
    restoreBtn.innerHTML = '<i class="fas fa-check"></i> 已恢复';
    restoreBtn.style.background = '#28a745';
    
    // Reset after 1 second
    setTimeout(() => {
        restoreBtn.innerHTML = originalHTML;
        restoreBtn.style.background = '';
    }, 1000);
}