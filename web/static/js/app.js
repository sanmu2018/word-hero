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
    // Get current page from URL or default to 1
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

    // Load known words from API (with localStorage fallback)
    loadKnownWordsFromAPI();

    // Load word timestamps
    loadWordTimestamps();

    // Update page info display
    updatePageInfo();

    // Load initial page data via API call
    loadPage(currentPage);
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

    accentUS.addEventListener('click', () => selectAccent('us', true));
    accentUK.addEventListener('click', () => selectAccent('uk', true));

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
        const parsedPageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 12;
        pageSize = isNaN(parsedPageSize) || parsedPageSize <= 0 ? 12 : parsedPageSize;
    } else {
        // Validate provided pageSize
        pageSize = isNaN(pageSize) || pageSize <= 0 ? 12 : pageSize;
    }

    fetch(`/api/words?page=${pageNumber}&pageSize=${pageSize}`)
        .then(response => response.json())
        .then(response => {
            if (response.code === 0) {
                updatePageContent(response.data);
                updateNavigationState(response.data);
                updatePageInfo();
                scrollToTop();
            } else {
                showError(response.msg);
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
            showError('Network error: ' + error.message);
        });
}

function updatePageContent(data) {
    const vocabularyGrid = document.getElementById('vocabularyGrid');

    // Clear existing content
    vocabularyGrid.innerHTML = '';

    // Validate data structure - new format has items array and total count
    if (!data || !data.items || !Array.isArray(data.items) || typeof data.total !== 'number') {
        console.error('Invalid data structure:', data);
        showError('Invalid data received from server');
        return;
    }

    // Get page size to calculate display numbers
    const pageSizeSelect = document.getElementById('pageSizeSelect');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;
    const startIndex = (currentPage - 1) * pageSize + 1;

    // Store original word order for this page
    originalWordOrder = data.items.map((word, index) => ({
        word: word,
        displayNumber: startIndex + index
    }));

    // Add vocabulary cards
    data.items.forEach((word, index) => {
        const card = createVocabularyCard(word, startIndex + index);
        vocabularyGrid.appendChild(card);
    });

    // Apply current visibility settings to new content
    applyVisibilitySettings();

    // Update known words status from API
    updateKnownWordsStatus();
}

function createVocabularyCard(word, displayNumber) {
    // Validate word object
    if (!word || typeof word !== 'object') {
        console.error('Invalid word object:', word);
        const card = document.createElement('div');
        card.className = 'vocabulary-card error';
        card.innerHTML = '<div class="error-message">Invalid word data</div>';
        return card;
    }

    // Ensure required fields exist (API returns lowercase field names)
    const english = word.english || '';
    const chinese = word.chinese || '';

    const card = document.createElement('div');
    card.className = 'vocabulary-card';
    card.setAttribute('data-word', english);
    card.setAttribute('data-word-id', word.id || '');
    
    // Process Chinese translation
    const chineseTranslation = processTranslation(chinese);

    const isKnown = knownWords.has(english) || (window.apiKnownWordIds && word.id && window.apiKnownWordIds.has(word.id));

    card.innerHTML = `
        <button class="action-btn" onclick="toggleWordMenu(event, '${escapeJsString(english)}')" title="更多操作">
            <i class="fas fa-ellipsis-v"></i>
        </button>
        <div class="card-content">
            <div class="english-word">
                ${escapeHtml(english)}
                <button class="speaker-btn" onclick="speakWord('${escapeJsString(english)}', this)" title="播放发音">
                    <i class="fas fa-volume-up"></i>
                </button>
            </div>
            <div class="chinese-meaning">
                <span class="translation-text">${escapeHtml(chineseTranslation.display)}</span>
                ${chineseTranslation.hasTooltip ? `<div class="translation-tooltip">${escapeHtml(chinese)}</div>` : ''}
            </div>
        </div>
        <div class="word-menu" id="menu-${escapeJsString(english)}" style="display: none;">
            <div class="menu-item known-action" data-word="${escapeJsString(english)}" data-action="known">
                <i class="fas fa-check"></i>
                标为认识
            </div>
            <div class="menu-item unknown-action" data-word="${escapeJsString(english)}" data-action="unknown" style="display: none;">
                <i class="fas fa-times"></i>
                标为不认识
            </div>
            <div class="menu-item detail-action" data-word="${escapeJsString(english)}" data-action="detail">
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

    // Add click event listeners for inline menu items
    const knownAction = card.querySelector('.known-action');
    const unknownAction = card.querySelector('.unknown-action');
    const detailAction = card.querySelector('.detail-action');

    if (knownAction) {
        knownAction.addEventListener('click', function(e) {
            e.stopPropagation();
            const word = this.dataset.word;
            const wordId = card.dataset.wordId;
            if (word && wordId) {
                markWordAsKnownWithAPI(word, wordId);
            } else if (word) {
                markWordAsKnown(word);
            }
            // Hide the menu
            const menu = card.querySelector('.word-menu');
            if (menu) menu.style.display = 'none';
        });
    }

    if (unknownAction) {
        unknownAction.addEventListener('click', function(e) {
            e.stopPropagation();
            const word = this.dataset.word;
            const wordId = card.dataset.wordId;
            if (word && wordId) {
                markWordAsUnknownWithAPI(word, wordId);
            } else if (word) {
                markWordAsUnknown(word);
            }
            // Hide the menu
            const menu = card.querySelector('.word-menu');
            if (menu) menu.style.display = 'none';
        });
    }

    if (detailAction) {
        detailAction.addEventListener('click', function(e) {
            e.stopPropagation();
            const word = this.dataset.word;
            showWordDetail(word);
            // Hide the menu
            const menu = card.querySelector('.word-menu');
            if (menu) menu.style.display = 'none';
        });
    }
    
    // Apply known word state if needed
    if (isKnown) {
        card.classList.add('known-word');
    }
    
    return card;
}

function updateNavigationState(data) {
    // Validate data structure
    if (!data || typeof data !== 'object' || typeof data.total !== 'number') {
        console.error('Invalid navigation data:', data);
        return;
    }

    // Get page size to calculate total pages
    const pageSizeSelect = document.getElementById('pageSizeSelect');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;

    // Calculate total pages from total count
    totalPages = Math.ceil(data.total / pageSize);

    // Ensure currentPage is valid
    if (currentPage < 1) currentPage = 1;
    if (currentPage > totalPages) currentPage = totalPages;

    // Calculate navigation states
    const hasPrev = currentPage > 1;
    const hasNext = currentPage < totalPages;

    // Update buttons
    const firstBtn = document.getElementById('firstBtn');
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');
    const lastBtn = document.getElementById('lastBtn');

    if (firstBtn) firstBtn.disabled = !hasPrev;
    if (prevBtn) prevBtn.disabled = !hasPrev;
    if (nextBtn) nextBtn.disabled = !hasNext;
    if (lastBtn) lastBtn.disabled = !hasNext;

    // Update page input
    const pageInput = document.getElementById('pageInput');
    if (pageInput) {
        // Only set valid numeric values
        if (currentPage && currentPage > 0) {
            pageInput.value = currentPage;
        }
        if (totalPages && totalPages > 0) {
            pageInput.max = totalPages;
        }
    }
}

function updatePageInfo() {
    // Page info display removed, only update page input
    const pageInput = document.getElementById('pageInput');
    if (pageInput) {
        // Safely set values only if they are valid numbers
        if (currentPage && currentPage > 0) {
            pageInput.value = currentPage;
        }
        if (totalPages && totalPages > 0) {
            pageInput.max = totalPages;
        }
    }
}


// Modal functions
function showStats() {
    fetch('/api/stats')
        .then(response => response.json())
        .then(data => {
            if (data.code === 0) {
                displayStats(data.data);
                showModal('statsModal');
            } else {
                showError(data.msg);
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
    if (!text) return '';
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
        padding: 15px 25px;
        border-radius: 8px;
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
        display: flex;
        align-items: center;
        gap: 12px;
        z-index: 3000;
        opacity: 0;
        transform: translateX(100%);
        transition: all 0.3s ease;
        font-size: 16px;
        font-weight: 500;
        min-width: 250px;
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
        animation: successPulse 0.6s ease-in-out;
    }

    @keyframes successPulse {
        0% { transform: translateX(0) scale(1); }
        50% { transform: translateX(0) scale(1.05); }
        100% { transform: translateX(0) scale(1); }
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
            if (data.code === 0) {
                displayModalSearchResults(data.data, query);
            } else {
                showSearchError(data.msg);
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
                <div class="search-result-word">${escapeHtml(word.english)}</div>
                <div class="search-result-meaning">${escapeHtml(word.chinese)}</div>
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
            const targetLang = selectedAccent === 'us' ? 'en-US' : 'en-GB';
            if (voice.lang === targetLang) {
                selectedVoice = voice;
                break;
            }
        }

        // If no exact match, try partial match
        if (!selectedVoice) {
            for (const voice of voices) {
                const targetLang = selectedAccent === 'us' ? 'en-US' : 'en-GB';
                if (voice.lang.startsWith(targetLang.split('-')[0])) {
                    selectedVoice = voice;
                    break;
                }
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
let currentActionWordId = null; // Track current word ID for action modal

function toggleWordMenu(event, word) {
    event.stopPropagation();

    // Get word ID from the card
    const card = event.target.closest('.vocabulary-card');
    const wordId = card ? card.dataset.wordId : null;

    // Show modal instead of inline menu
    showWordActionModal(word, wordId);
}

function updateInlineMenuButtons(card, word) {
    const knownAction = card.querySelector('.known-action');
    const unknownAction = card.querySelector('.unknown-action');

    // Check if word is known using either local knownWords set or API word IDs
    let isKnown = knownWords.has(word);

    // Also check against API known word IDs if available
    if (!isKnown && window.apiKnownWordIds && card.dataset.wordId) {
        isKnown = window.apiKnownWordIds.has(card.dataset.wordId);
    }

    if (knownAction) {
        knownAction.style.display = isKnown ? 'none' : 'block';
    }
    if (unknownAction) {
        unknownAction.style.display = isKnown ? 'block' : 'none';
    }
}

function updateActionModalButtons(word) {
    const knownBtn = document.querySelector('.known-btn');
    const unknownBtn = document.querySelector('.unknown-btn');

    // Check if word is known using either local knownWords set or API word IDs
    let isKnown = knownWords.has(word);

    // Also check against API known word IDs if available and we have a current action word ID
    if (!isKnown && window.apiKnownWordIds && currentActionWordId) {
        isKnown = window.apiKnownWordIds.has(currentActionWordId);
    }

    if (isKnown) {
        knownBtn.style.display = 'none';
        unknownBtn.style.display = 'flex';
    } else {
        knownBtn.style.display = 'flex';
        unknownBtn.style.display = 'none';
    }
}

function markCurrentWordAsKnown() {

    // Try to get word from selected card first
    if (selectedCard) {
        const word = selectedCard.getAttribute('data-word');
        const wordId = selectedCard.getAttribute('data-word-id');

        if (word && wordId) {
            // Use the new API-based marking function
            markWordAsKnownWithAPI(word, wordId);
            return;
        } else if (word) {
            // Fallback to localStorage if no word ID
            markWordAsKnown(word);
            return;
        }
    }

    // Fallback to action modal context
    if (currentActionWord && currentActionWordId) {
        markWordAsKnownWithAPI(currentActionWord, currentActionWordId);
        closeWordActionModal();
    } else if (currentActionWord) {
        markWordAsKnown(currentActionWord);
        closeWordActionModal();
    } else {
        showToast('请先选择一个单词', 'error');
    }
}

function markCurrentWordAsUnknown() {
    // Try to get word from selected card first
    if (selectedCard) {
        const word = selectedCard.getAttribute('data-word');
        const wordId = selectedCard.getAttribute('data-word-id');

        if (word && wordId) {
            // Use the new API-based marking function
            markWordAsUnknownWithAPI(word, wordId);
            return;
        } else if (word) {
            // Fallback to localStorage if no word ID
            markWordAsUnknown(word);
            return;
        }
    }

    // Fallback to action modal context
    if (currentActionWord && currentActionWordId) {
        markWordAsUnknownWithAPI(currentActionWord, currentActionWordId);
        closeWordActionModal();
    } else if (currentActionWord) {
        markWordAsUnknown(currentActionWord);
        closeWordActionModal();
    } else {
        showToast('请先选择一个单词', 'error');
    }
}

function showCurrentWordDetail() {
    if (currentActionWord) {
        closeWordActionModal();
        showWordDetail(currentActionWord);
    }
}

function showWordActionModal(word, wordId = null) {
    const modal = document.getElementById('wordActionModal');
    const title = document.getElementById('wordActionTitle');

    // Set the word for action functions
    currentActionWord = word;
    currentActionWordId = wordId;

    // Update modal title
    title.textContent = `单词操作: ${word}`;

    // Update button states based on word status
    updateActionModalButtons(word);

    // Show modal
    modal.style.display = 'block';

}

function closeWordActionModal() {
    const modal = document.getElementById('wordActionModal');
    modal.style.display = 'none';
    currentActionWord = null;
}


function markWordAsKnown(word) {
    // First get the word ID from the card data attribute
    const card = document.querySelector(`.vocabulary-card[data-word="${word}"]`);
    const wordId = card ? card.dataset.wordId : null;

    if (!wordId) {
        showToast('无法获取单词ID', 'error');
        return;
    }

    // Call API to mark word as known
    fetch('/api/word-tags/mark', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAuthToken()}`
        },
        body: JSON.stringify({
            wordId: wordId
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0) {
            knownWords.add(word);
            wordTimestamps.set(word, new Date().toISOString());
            saveKnownWords();
            saveWordTimestamps();
            updateWordCardAppearance(word, true);
            updateWordMarkStatus(wordId, true, data.data.markCount);
            showToast(`"${word}" 已标记为认识`, 'success');
            closeWordActionModal();
        } else {
            showToast(data.msg || '标记失败', 'error');
        }
    })
    .catch(error => {
        console.error('Error marking word as known:', error);
        showToast('网络错误，请重试', 'error');
    });
}

function markWordAsUnknown(word) {
    // First get the word ID from the card data attribute
    const card = document.querySelector(`.vocabulary-card[data-word="${word}"]`);
    const wordId = card ? card.dataset.wordId : null;

    if (!wordId) {
        showToast('无法获取单词ID', 'error');
        return;
    }

    // Call API to unmark word
    fetch('/api/word-tags/unmark', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAuthToken()}`
        },
        body: JSON.stringify({
            wordId: wordId
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0) {
            knownWords.delete(word);
            wordTimestamps.delete(word);
            saveKnownWords();
            saveWordTimestamps();
            updateWordCardAppearance(word, false);
            updateWordMarkStatus(wordId, false, data.data.markCount);
            showToast(`"${word}" 已标记为不认识`, 'info');
            closeWordActionModal();
        } else {
            showToast(data.msg || '取消标记失败', 'error');
        }
    })
    .catch(error => {
        console.error('Error unmarking word:', error);
        showToast('网络错误，请重试', 'error');
    });
}

function markWordAsKnownWithAPI(word, wordId) {
    const token = getAuthToken();
    if (!token) {
        showToast('请先登录后再进行操作', 'error');
        setTimeout(() => showLoginModal(), 1000);
        return;
    }

    // Call new API to mark word as known
    fetch('/api/word-tags/mark', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            wordId: wordId
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0) {
            knownWords.add(word);
            wordTimestamps.set(word, new Date().toISOString());
            saveKnownWords();
            saveWordTimestamps();
            updateWordCardAppearance(word, true);

            // Update API known words cache
            if (!window.apiKnownWordIds) {
                window.apiKnownWordIds = new Set();
            }
            window.apiKnownWordIds.add(wordId);

            showToast(`"${word}" 已标记为认识`, 'success');
        } else if (data.code === 401) {
            showToast('请先登录后再进行操作', 'error');
            setTimeout(() => showLoginModal(), 1000);
        } else {
            showToast(data.msg || '标记失败', 'error');
        }
    })
    .catch(error => {
        showToast('网络错误，请重试', 'error');
    });
}

function markWordAsUnknownWithAPI(word, wordId) {
    const token = getAuthToken();
    if (!token) {
        showToast('请先登录后再进行操作', 'error');
        setTimeout(() => showLoginModal(), 1000);
        return;
    }

    // Call new API to unmark word
    fetch('/api/word-tags/unmark', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            wordId: wordId
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0) {
            knownWords.delete(word);
            wordTimestamps.delete(word);
            saveKnownWords();
            saveWordTimestamps();
            updateWordCardAppearance(word, false);

            // Update API known words cache - remove the word ID
            if (window.apiKnownWordIds) {
                window.apiKnownWordIds.delete(wordId);
            }

            showToast(`"${word}" 已标记为不认识`, 'info');
        } else if (data.code === 401) {
            showToast('请先登录后再进行操作', 'error');
            setTimeout(() => showLoginModal(), 1000);
        } else {
            showToast(data.msg || '取消标记失败', 'error');
        }
    })
    .catch(error => {
        showToast('网络错误，请重试', 'error');
    });
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

function updateWordMarkStatus(wordId, isMarked, markCount) {
    // Mark badge functionality removed - no longer needed
    // This function is kept for compatibility but does nothing
}

// New function to check word status using the new API
function checkWordStatusWithAPI(wordId, callback) {
    const token = getAuthToken();
    if (!token) {
        // If not logged in, use localStorage data
        const isKnown = knownWords.has(wordId) || (window.apiKnownWordIds && window.apiKnownWordIds.has(wordId));
        callback(isKnown);
        return;
    }

    // Call new API to check word status
    fetch(`/api/word-tags/status?wordId=${encodeURIComponent(wordId)}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0 && data.data) {
            callback(data.data.isKnown);
        } else {
            // Fallback to localStorage data
            const isKnown = knownWords.has(wordId) || (window.apiKnownWordIds && window.apiKnownWordIds.has(wordId));
            callback(isKnown);
        }
    })
    .catch(error => {
        // Fallback to localStorage data
        const isKnown = knownWords.has(wordId) || (window.apiKnownWordIds && window.apiKnownWordIds.has(wordId));
        callback(isKnown);
    });
}

function getAuthToken() {
    return localStorage.getItem('authToken');
}

function loadKnownWordsFromAPI() {
    // Load from localStorage first for immediate feedback
    loadKnownWords();

    const token = getAuthToken();
    if (!token) {
        // Not logged in, use localStorage data only
        return;
    }

    // Try to load from API for accurate data using new endpoint
    fetch('/api/word-tags/known-words', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0 && data.data && data.data.wordIds) {
            // Store known word IDs from API
            window.apiKnownWordIds = new Set(data.data.wordIds);

            // Update the UI based on API data
            updateKnownWordsStatus();
        } else if (data.code === 401) {
            // Token invalid, clear it and continue with localStorage
            localStorage.removeItem('authToken');
        }
    })
    .catch(error => {
        // Continue with localStorage data
    });
}

// Update known words status on the current page
function updateKnownWordsStatus() {
    if (!window.apiKnownWordIds) return;

    const cards = document.querySelectorAll('.vocabulary-card');
    cards.forEach(card => {
        const wordId = card.dataset.wordId;
        if (wordId && window.apiKnownWordIds.has(wordId)) {
            // Add known-word class and animation
            card.classList.add('known-word');

            // Add marking animation for visual feedback
            card.classList.add('marking-known');

            // Remove animation class after animation completes
            setTimeout(() => {
                card.classList.remove('marking-known');
            }, 600);

            // Also add to local knownWords set for consistency
            const wordText = card.querySelector('.english-word').textContent;
            knownWords.add(wordText);
        }
    });

    // Update inline menu buttons for all cards to show correct action buttons
    cards.forEach(card => {
        const wordText = card.querySelector('.english-word').textContent;
        updateInlineMenuButtons(card, wordText);
    });

    // Save to localStorage for consistency
    saveKnownWords();
}

function resetKnownWords() {
    // Show reset options modal instead of confirmation dialog
    showResetOptionsModal();
}

function showResetOptionsModal() {
    const modal = document.getElementById('resetOptionsModal');

    // Update the modal with current page information
    const currentPageElement = document.getElementById('currentPageNumber');
    const totalPagesElement = document.getElementById('totalPagesNumber');

    if (currentPageElement) {
        currentPageElement.textContent = currentPage;
    }
    if (totalPagesElement) {
        totalPagesElement.textContent = totalPages;
    }

    modal.style.display = 'block';
}

function closeResetOptionsModal() {
    const modal = document.getElementById('resetOptionsModal');
    modal.style.display = 'none';
}

function resetCurrentPage() {
    const token = getAuthToken();
    if (!token) {
        showToast('请先登录后再进行操作', 'error');
        setTimeout(() => showLoginModal(), 1000);
        return;
    }

    // Get current page and page size
    const pageSizeSelect = document.getElementById('pageSizeSelect');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;

    // Collect word IDs from current page
    const currentCards = document.querySelectorAll('.vocabulary-card');
    const wordIds = [];

    currentCards.forEach(card => {
        const wordId = card.getAttribute('data-word-id');
        if (wordId) {
            wordIds.push(wordId);
        }
    });

    if (wordIds.length === 0) {
        showToast('当前页面没有可忘光的单词', 'info');
        return;
    }

    // Call new API to forget specific words
    fetch('/api/word-tags/forget-words', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            wordIds: wordIds
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0) {
            // Update local knownWords cache
            const currentCards = document.querySelectorAll('.vocabulary-card');
            currentCards.forEach(card => {
                const word = card.getAttribute('data-word');
                if (word && knownWords.has(word)) {
                    knownWords.delete(word);
                    wordTimestamps.delete(word);
                    card.classList.remove('known-word');
                }
            });

            // Update API known words cache
            if (window.apiKnownWordIds && data.data && data.data.wordIds) {
                data.data.wordIds.forEach(wordId => {
                    window.apiKnownWordIds.delete(wordId);
                });
            }

            // Save changes
            saveKnownWords();
            saveWordTimestamps();

            // Close modal
            closeResetOptionsModal();

            // Show feedback
            showToast(data.data.message || `已忘光当前页面的 ${data.data.forgottenCount} 个已认识单词`, 'success');

            // Refresh current page
            loadPage(currentPage, pageSize);
        } else if (data.code === 401) {
            showToast('请先登录后再进行操作', 'error');
            setTimeout(() => showLoginModal(), 1000);
        } else {
            showToast(data.msg || '忘光失败', 'error');
        }
    })
    .catch(error => {
        showToast('网络错误，请重试', 'error');
    });
}

function resetAllPages() {
    const token = getAuthToken();
    if (!token) {
        showToast('请先登录后再进行操作', 'error');
        setTimeout(() => showLoginModal(), 1000);
        return;
    }

    // Confirm before forgetting all words
    if (confirm('确定要忘光所有已认识的单词吗？这将清除所有单词的已认识标记。')) {
        // Call new API to forget all words
        fetch('/api/word-tags/forget-all', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                confirm: true
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.code === 0) {
                // Clear local caches
                knownWords.clear();
                wordTimestamps.clear();
                if (window.apiKnownWordIds) {
                    window.apiKnownWordIds.clear();
                }

                // Save to localStorage
                saveKnownWords();
                saveWordTimestamps();

                // Close modal
                closeResetOptionsModal();

                showToast(data.data.message || `已忘光全部 ${data.data.forgottenCount} 个已认识单词`, 'success');

                // Reload current page to refresh display
                const pageSizeSelect = document.getElementById('pageSizeSelect');
                const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;
                loadPage(currentPage, pageSize);
            } else {
                showToast(data.msg || '操作失败，请重试', 'error');
            }
        })
        .catch(error => {
            console.error('Error forgetting all words:', error);
            showToast('网络错误，请重试', 'error');
        });
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

    // Close inline menus when clicking outside
    if (!e.target.closest('.vocabulary-card') || !e.target.closest('.action-btn')) {
        // Close all inline menus
        const allMenus = document.querySelectorAll('.word-menu');
        allMenus.forEach(menu => {
            if (menu.style.display === 'block') {
                menu.style.display = 'none';
            }
        });
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
    
    if (cards.length <= 1 || originalWordOrder.length === 0) {
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

// Authentication functionality
class AuthManager {
    constructor() {
        this.token = localStorage.getItem('authToken') || null;
        this.user = JSON.parse(localStorage.getItem('currentUser')) || null;
        this.init();
    }

    init() {
        this.updateAuthUI();
        this.bindAuthEvents();
        this.checkAuthStatus();
    }

    bindAuthEvents() {
        // Login form submission
        document.getElementById('loginForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.login();
        });

        // Register form submission
        document.getElementById('registerForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.register();
        });

        // Profile form submission
        document.getElementById('profileForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.updateProfile();
        });

        // Change password form submission
        document.getElementById('changePasswordForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.changePassword();
        });
    }

    async login() {
        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password })
            });

            const result = await response.json();

            if (result.code === 0) {
                this.token = result.data.token;
                this.user = result.data.user;
                localStorage.setItem('authToken', this.token);
                localStorage.setItem('currentUser', JSON.stringify(this.user));
                this.updateAuthUI();
                this.showNotification('登录成功！', 'success');
                closeModal('loginModal');
            } else {
                this.showNotification(result.msg || '登录失败', 'error');
            }
        } catch (error) {
            this.showNotification('网络错误，请重试', 'error');
        }
    }

    async register() {
        const formData = {
            username: document.getElementById('registerUsername').value,
            email: document.getElementById('registerEmail').value,
            full_name: document.getElementById('registerFullName').value,
            password: document.getElementById('registerPassword').value,
            confirm_password: document.getElementById('registerConfirmPassword').value
        };

        if (formData.password !== formData.confirm_password) {
            this.showNotification('两次输入的密码不一致', 'error');
            return;
        }

        try {
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            const result = await response.json();

            if (result.code === 0) {
                this.token = result.data.token;
                this.user = result.data.user;
                localStorage.setItem('authToken', this.token);
                localStorage.setItem('currentUser', JSON.stringify(this.user));
                this.updateAuthUI();
                this.showNotification('注册成功！', 'success');
                closeModal('registerModal');
            } else {
                this.showNotification(result.msg || '注册失败', 'error');
            }
        } catch (error) {
            this.showNotification('网络错误，请重试', 'error');
        }
    }

    async updateProfile() {
        const formData = {
            full_name: document.getElementById('profileFullName').value,
            bio: document.getElementById('profileBio').value
        };

        try {
            const response = await fetch('/api/auth/profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.token}`
                },
                body: JSON.stringify(formData)
            });

            const result = await response.json();

            if (result.code === 0) {
                this.user = result.data;
                localStorage.setItem('currentUser', JSON.stringify(this.user));
                this.updateAuthUI();
                this.showNotification('个人资料更新成功！', 'success');
                closeModal('profileModal');
            } else {
                this.showNotification(result.msg || '更新失败', 'error');
            }
        } catch (error) {
            this.showNotification('网络错误，请重试', 'error');
        }
    }

    async changePassword() {
        const formData = {
            current_password: document.getElementById('currentPassword').value,
            new_password: document.getElementById('newPassword').value,
            confirm_new_password: document.getElementById('confirmNewPassword').value
        };

        if (formData.new_password !== formData.confirm_new_password) {
            this.showNotification('两次输入的新密码不一致', 'error');
            return;
        }

        try {
            const response = await fetch('/api/auth/change-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.token}`
                },
                body: JSON.stringify(formData)
            });

            const result = await response.json();

            if (result.code === 0) {
                this.showNotification('密码修改成功！', 'success');
                closeModal('changePasswordModal');
                document.getElementById('changePasswordForm').reset();
            } else {
                this.showNotification(result.msg || '密码修改失败', 'error');
            }
        } catch (error) {
            this.showNotification('网络错误，请重试', 'error');
        }
    }

    logout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('authToken');
        localStorage.removeItem('currentUser');
        this.updateAuthUI();
        this.showNotification('已退出登录', 'info');
    }

    updateAuthUI() {
        const loginSection = document.getElementById('loginSection');
        const userSection = document.getElementById('userSection');
        const userDisplay = document.getElementById('userDisplay');

        if (this.user) {
            loginSection.style.display = 'none';
            userSection.style.display = 'flex';
            userDisplay.textContent = this.user.full_name || this.user.username;
        } else {
            loginSection.style.display = 'flex';
            userSection.style.display = 'none';
        }
    }

    async checkAuthStatus() {
        if (!this.token) return;

        try {
            const response = await fetch('/api/auth/me', {
                headers: {
                    'Authorization': `Bearer ${this.token}`
                }
            });

            if (response.ok) {
                const result = await response.json();
                if (result.code === 0) {
                    this.user = result.data;
                    localStorage.setItem('currentUser', JSON.stringify(this.user));
                    this.updateAuthUI();
                } else {
                    // Token is invalid, clear it
                    this.logout();
                }
            } else {
                this.logout();
            }
        } catch (error) {
            this.logout();
        }
    }

    showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        document.body.appendChild(notification);

        setTimeout(() => {
            notification.remove();
        }, 3000);
    }
}

// Global authentication functions
function showLoginModal() {
    showModal('loginModal');
}

function showRegisterModal() {
    showModal('registerModal');
}

function showProfileModal() {
    if (authManager.user) {
        document.getElementById('profileUsername').value = authManager.user.username;
        document.getElementById('profileEmail').value = authManager.user.email;
        document.getElementById('profileFullName').value = authManager.user.full_name || '';
        document.getElementById('profileBio').value = authManager.user.bio || '';
        showModal('profileModal');
    }
}

function showChangePasswordModal() {
    showModal('changePasswordModal');
}

function logout() {
    if (confirm('确定要退出登录吗？')) {
        authManager.logout();
    }
}

// Initialize authentication manager
const authManager = new AuthManager();