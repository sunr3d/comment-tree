let currentParentId = null; // Загружаем все корневые комментарии
let searchQuery = '';

// Загрузка корневых комментариев
async function loadComments() {
    try {
        console.log('Начинаем загрузку комментариев...');
        
        const params = new URLSearchParams({
            parent: 0 // Загружаем корневые комментарии
        });
        
        if (searchQuery) {
            params.append('search', searchQuery);
        }

        console.log('Отправляем запрос:', `/comments?${params}`);
        
        const response = await fetch(`/comments?${params}`);
        
        console.log('Получен ответ:', response.status, response.statusText);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Ошибка ответа:', errorText);
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('Получены данные:', data);
        console.log('Количество комментариев:', data.comments ? data.comments.length : 0);
        
        displayRootComments(data.comments);
    } catch (error) {
        console.error('Ошибка загрузки комментариев:', error);
        document.getElementById('commentsTree').innerHTML = 
            '<div class="error">Ошибка загрузки комментариев: ' + error.message + '</div>';
    }
}

// Отображение корневых комментариев с кнопками "Показать ответы"
function displayRootComments(comments) {
    const container = document.getElementById('commentsTree');
    
    if (!comments || comments.length === 0) {
        container.innerHTML = '<div class="loading">Комментариев пока нет</div>';
        return;
    }
    
    // Сортируем по времени создания
    const sortedComments = comments.sort((a, b) => {
        return new Date(a.created_at) - new Date(b.created_at);
    });
    
    container.innerHTML = sortedComments.map((comment, index) => {
        console.log('Создаем комментарий с ID:', comment.id);
        
        return `
            <div class="root-comment" data-comment-id="${comment.id}">
                <div class="comment-header">
                    <span class="comment-author">${escapeHtml(comment.author)}</span>
                    <span class="comment-date">${formatDate(comment.created_at)}</span>
                </div>
                <div class="comment-content">${escapeHtml(comment.content)}</div>
                <div class="comment-actions">
                    <button class="show-replies-btn" onclick="loadReplies(${comment.id}, ${index})">
                        Показать ответы
                    </button>
                    <button class="reply-btn" onclick="replyToComment(${comment.id})">Ответить</button>
                </div>
                <div class="replies-container" id="replies-${comment.id}" style="display: none;"></div>
            </div>
        `;
    }).join('');
}

// Загрузка ответов для конкретного комментария
async function loadReplies(parentId, index) {
    try {
        console.log('Загружаем ответы для parentId:', parentId, 'тип:', typeof parentId);
        
        const response = await fetch(`/comments?parent=${parentId}`);
        
        console.log('Ответ от сервера:', response.status, response.statusText);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Ошибка ответа:', errorText);
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('Получены ответы:', data);
        
        displayReplies(parentId, data.comments);
        
        // Скрываем кнопку "Показать ответы"
        const btn = document.querySelector(`[data-comment-id="${parentId}"] .show-replies-btn`);
        if (btn) {
            btn.style.display = 'none';
        } else {
            console.error('Кнопка не найдена для parentId:', parentId);
        }
        
    } catch (error) {
        console.error('Ошибка загрузки ответов:', error);
        showMessage('Ошибка загрузки ответов', 'error');
    }
}

// Отображение ответов с отступами по уровням
function displayReplies(parentId, replies) {
    const container = document.getElementById(`replies-${parentId}`);
    
    if (!replies || replies.length === 0) {
        container.innerHTML = '<div class="no-replies">Ответов пока нет</div>';
        container.style.display = 'block';
        return;
    }
    
    // Сортируем по времени создания
    const sortedReplies = replies.sort((a, b) => {
        return new Date(a.created_at) - new Date(b.created_at);
    });
    
    container.innerHTML = sortedReplies.map((reply, index) => `
        <div class="reply level-${reply.level || 0}" data-reply-id="${reply.id}">
            <div class="comment-header">
                <span class="comment-author">${escapeHtml(reply.author)}</span>
                <span class="comment-date">${formatDate(reply.created_at)}</span>
            </div>
            <div class="comment-content">${escapeHtml(reply.content)}</div>
            <div class="comment-actions">
                <button class="reply-btn" onclick="replyToComment(${reply.id})">Ответить</button>
            </div>
            <!-- Форма ответа для этого комментария -->
            <div class="reply-form" id="replyForm-${reply.id}" style="display: none;">
                <form onsubmit="submitReply(event, ${reply.id})">
                    <input type="text" id="replyAuthor-${reply.id}" placeholder="Ваше имя" required maxlength="50">
                    <textarea id="replyContent-${reply.id}" placeholder="Ваш ответ..." required maxlength="1000" rows="3"></textarea>
                    <div class="reply-actions">
                        <button type="submit">Отправить</button>
                        <button type="button" onclick="cancelReply(${reply.id})">Отмена</button>
                    </div>
                </form>
            </div>
        </div>
    `).join('');
    
    container.style.display = 'block';
}

// Ответ на комментарий (упрощенная версия)
function replyToComment(commentId) {
    console.log('Отвечаем на комментарий:', commentId);
    
    // Скрываем все другие формы ответов
    document.querySelectorAll('.reply-form').forEach(form => {
        form.style.display = 'none';
    });
    
    // Показываем форму ответа для конкретного комментария
    const replyForm = document.getElementById(`replyForm-${commentId}`);
    console.log('Найдена форма:', replyForm);
    
    if (replyForm) {
        replyForm.style.display = 'block';
        console.log('Форма показана');
    } else {
        console.error('Форма ответа не найдена для комментария:', commentId);
        // Создаем форму динамически
        createReplyForm(commentId);
    }
}

// Создание формы ответа динамически
function createReplyForm(commentId) {
    console.log('Создаем форму для комментария:', commentId);
    
    // Ищем контейнер для формы
    const commentElement = document.querySelector(`[data-comment-id="${commentId}"]`);
    if (!commentElement) {
        console.error('Комментарий не найден:', commentId);
        return;
    }
    
    // Создаем форму
    const formHTML = `
        <div class="reply-form" id="replyForm-${commentId}" style="display: block;">
            <form onsubmit="submitReply(event, ${commentId})">
                <input type="text" id="replyAuthor-${commentId}" placeholder="Ваше имя" required maxlength="50">
                <textarea id="replyContent-${commentId}" placeholder="Ваш ответ..." required maxlength="1000" rows="3"></textarea>
                <div class="reply-actions">
                    <button type="submit">Отправить</button>
                    <button type="button" onclick="cancelReply(${commentId})">Отмена</button>
                </div>
            </form>
        </div>
    `;
    
    // Добавляем форму в контейнер ответов
    const repliesContainer = commentElement.querySelector('.replies-container');
    if (repliesContainer) {
        repliesContainer.innerHTML = formHTML + repliesContainer.innerHTML;
        repliesContainer.style.display = 'block';
        console.log('Форма создана и добавлена');
    } else {
        console.error('Контейнер для ответов не найден');
    }
}

// Отмена ответа
function cancelReply(commentId) {
    console.log('Отменяем ответ для комментария:', commentId);
    const replyForm = document.getElementById(`replyForm-${commentId}`);
    if (replyForm) {
        replyForm.style.display = 'none';
        
        // Находим форму внутри div'а и сбрасываем её
        const form = replyForm.querySelector('form');
        if (form) {
            form.reset();
        }
    }
}

// Отправка ответа (упрощенная версия)
async function submitReply(event, commentId) {
    event.preventDefault();
    console.log('Отправляем ответ для комментария:', commentId);
    
    const author = document.getElementById(`replyAuthor-${commentId}`).value;
    const content = document.getElementById(`replyContent-${commentId}`).value;
    
    if (!author.trim() || !content.trim()) {
        showMessage('Заполните все поля', 'error');
        return;
    }
    
    try {
        console.log('Отправляем запрос:', {
            author: author.trim(),
            content: content.trim(),
            parent_id: commentId
        });
        
        const response = await fetch('/comments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                author: author.trim(),
                content: content.trim(),
                parent_id: commentId
            })
        });
        
        console.log('Ответ сервера:', response.status, response.statusText);
        
        if (!response.ok) {
            const error = await response.json();
            console.error('Ошибка сервера:', error);
            throw new Error(error.error || 'Ошибка создания ответа');
        }
        
        // Скрываем форму ответа
        cancelReply(commentId);
        
        // Перезагружаем все комментарии
        console.log('Перезагружаем комментарии...');
        await loadComments();
        
        showMessage('Ответ успешно добавлен', 'success');
    } catch (error) {
        console.error('Ошибка создания ответа:', error);
        showMessage(error.message, 'error');
    }
}

// Создание комментария
document.getElementById('commentForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const author = document.getElementById('author').value;
    const content = document.getElementById('content').value;
    
    if (!author.trim() || !content.trim()) {
        showMessage('Заполните все поля', 'error');
        return;
    }
    
    try {
        const response = await fetch('/comments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                author: author.trim(),
                content: content.trim(),
                parent_id: null // Корневой комментарий
            })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка создания комментария');
        }
        
        // Очистка формы
        document.getElementById('commentForm').reset();
        
        // Перезагрузка комментариев
        await loadComments();
        
        showMessage('Комментарий успешно создан', 'success');
    } catch (error) {
        console.error('Ошибка создания комментария:', error);
        showMessage(error.message, 'error');
    }
});

// Удаление комментария
async function deleteComment(commentIndex) {
    if (!confirm('Вы уверены, что хотите удалить этот комментарий?')) {
        return;
    }
    
    showMessage('Функция удаления будет реализована позже', 'error');
}

// Поиск комментариев
function searchComments() {
    searchQuery = document.getElementById('searchInput').value.trim();
    console.log('Поиск по запросу:', searchQuery);
    
    if (searchQuery) {
        // Если есть поиск, ищем по всем комментариям
        loadCommentsWithSearch();
    } else {
        // Если поиск пустой, загружаем корневые комментарии
        loadComments();
    }
}

// Загрузка комментариев с поиском
async function loadCommentsWithSearch() {
    try {
        console.log('Загружаем комментарии с поиском:', searchQuery);
        
        // Сначала ищем по корневым комментариям
        const rootParams = new URLSearchParams({
            parent: 0,
            search: searchQuery
        });

        const rootResponse = await fetch(`/comments?${rootParams}`);
        const rootData = await rootResponse.json();
        
        console.log('Результаты поиска по корневым:', rootData);
        
        // Затем ищем по всем комментариям
        const allParams = new URLSearchParams({
            parent: 1,
            search: searchQuery
        });

        const allResponse = await fetch(`/comments?${allParams}`);
        const allData = await allResponse.json();
        
        console.log('Результаты поиска по всем:', allData);
        
        // Объединяем результаты
        const allResults = [...(rootData.comments || []), ...(allData.comments || [])];
        
        // Убираем дубликаты по ID
        const uniqueResults = allResults.filter((comment, index, self) => 
            index === self.findIndex(c => c.id === comment.id)
        );
        
        console.log('Объединенные результаты:', uniqueResults);
        
        // Показываем результаты поиска
        displaySearchResults(uniqueResults);
        
    } catch (error) {
        console.error('Ошибка поиска:', error);
        document.getElementById('commentsTree').innerHTML = 
            '<div class="error">Ошибка поиска</div>';
    }
}

// Отображение результатов поиска
function displaySearchResults(comments) {
    const container = document.getElementById('commentsTree');
    
    if (!comments || comments.length === 0) {
        container.innerHTML = '<div class="loading">По вашему запросу ничего не найдено</div>';
        return;
    }
    
    // Сортируем по времени создания
    const sortedComments = comments.sort((a, b) => {
        return new Date(a.created_at) - new Date(b.created_at);
    });
    
    container.innerHTML = sortedComments.map((comment, index) => `
        <div class="search-result level-${comment.level || 0}" data-comment-id="${comment.id}">
            <div class="comment-header">
                <span class="comment-author">${escapeHtml(comment.author)}</span>
                <span class="comment-date">${formatDate(comment.created_at)}</span>
                <span class="search-level">Уровень: ${comment.level || 0}</span>
            </div>
            <div class="comment-content">${escapeHtml(comment.content)}</div>
            <div class="comment-actions">
                <button class="reply-btn" onclick="replyToComment(${comment.id})">Ответить</button>
            </div>
        </div>
    `).join('');
}

// Очистка поиска
function clearSearch() {
    document.getElementById('searchInput').value = '';
    searchQuery = '';
    console.log('Очищаем поиск, загружаем корневые комментарии');
    loadComments();
}

// Вспомогательные функции
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU');
}

function showMessage(message, type) {
    const container = document.querySelector('.container');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = message;
    container.insertBefore(messageDiv, container.firstChild);
    
    setTimeout(() => {
        messageDiv.remove();
    }, 3000);
}

// Загрузка при старте
document.addEventListener('DOMContentLoaded', loadComments);