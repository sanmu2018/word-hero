// Global variables
let currentPage = 1;
let totalPages = 1;
let isLoading = false;
let wordsVisible = true;
let translationsVisible = true;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
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
    
      
    // Search functionality
    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    
    searchInput.addEventListener('input', debounce(handleSearch, 300));
    searchInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            handleSearch();
        }
    });
    searchBtn.addEventListener('click', handleSearch);
    
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
    
    toggleWordsBtn.addEventListener('click', toggleWordsVisibility);
    toggleTranslationsBtn.addEventListener('click', toggleTranslationsVisibility);
    shuffleBtn.addEventListener('click', shuffleCards);
    
    // Keyboard shortcuts
    document.addEventListener('keydown', handleKeyboardShortcuts);
    
    // Close search results when clicking outside
    document.addEventListener('click', function(e) {
        const searchContainer = document.querySelector('.search-container');
        if (!searchContainer.contains(e.target)) {
            closeSearchResults();
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
    
    showLoading();
    
    // Get page size from selector or use default
    if (!pageSize) {
        const pageSizeSelect = document.getElementById('pageSizeSelect');
        pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 24;
    }
    
    fetch(`/api/words?page=${pageNumber}&pageSize=${pageSize}`)
        .then(response => response.json())
        .then(data => {
            hideLoading();
            
            if (data.success) {
                updatePageContent(data.data);
                updatePageInfo();
                scrollToTop();
            } else {
                showError('Failed to load page: ' + data.error);
            }
        })
        .catch(error => {
            hideLoading();
            showError('Network error: ' + error.message);
        });
}

function updatePageContent(data) {
    const vocabularyGrid = document.getElementById('vocabularyGrid');
    
    // Clear existing content
    vocabularyGrid.innerHTML = '';
    
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
    
    // Process Chinese translation
    const chineseTranslation = processTranslation(word.Chinese);
    
    card.innerHTML = `
        <div class="card-content">
            <div class="english-word">${escapeHtml(word.English)}</div>
            <div class="chinese-meaning">
                <span class="translation-text">${escapeHtml(chineseTranslation.display)}</span>
                ${chineseTranslation.hasTooltip ? `<div class="translation-tooltip">${escapeHtml(word.Chinese)}</div>` : ''}
            </div>
        </div>
    `;
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

// Search functionality
function handleSearch() {
    const query = document.getElementById('searchInput').value.trim();
    
    if (query.length === 0) {
        closeSearchResults();
        return;
    }
    
    if (query.length < 2) {
        return; // Don't search for very short queries
    }
    
    performSearch(query);
}

function performSearch(query) {
    fetch(`/api/search?q=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displaySearchResults(data.data);
            } else {
                showError('Search failed: ' + data.error);
            }
        })
        .catch(error => {
            showError('Search error: ' + error.message);
        });
}

function displaySearchResults(data) {
    const searchResults = document.getElementById('searchResults');
    const searchResultsList = document.getElementById('searchResultsList');
    
    // Clear previous results
    searchResultsList.innerHTML = '';
    
    // Check if results exist and is an array
    if (!data.results || !Array.isArray(data.results) || data.results.length === 0) {
        searchResultsList.innerHTML = '<div class="search-result-item">没有找到相关词汇</div>';
    } else {
        data.results.forEach(word => {
            const resultItem = document.createElement('div');
            resultItem.className = 'search-result-item';
            resultItem.innerHTML = `
                <div class="search-result-word">${escapeHtml(word.English)}</div>
                <div class="search-result-meaning">${escapeHtml(word.Chinese)}</div>
            `;
            resultItem.addEventListener('click', () => {
                closeSearchResults();
                document.getElementById('searchInput').value = '';
            });
            searchResultsList.appendChild(resultItem);
        });
    }
    
    // Show search results
    searchResults.style.display = 'block';
}

function closeSearchResults() {
    document.getElementById('searchResults').style.display = 'none';
    document.getElementById('searchInput').value = '';
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
function showLoading() {
    isLoading = true;
    document.getElementById('loading').style.display = 'flex';
}

function hideLoading() {
    isLoading = false;
    document.getElementById('loading').style.display = 'none';
}

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
                document.getElementById('searchInput').focus();
            }
            break;
        case 'Escape':
            closeModal('statsModal');
            closeModal('helpModal');
            closeSearchResults();
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
                e.preventDefault();
                shuffleCards();
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