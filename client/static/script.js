// Глобальные переменные
const API_URL = 'http://localhost:8080/comments';
let allComments = [];

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
  loadComments();
  setupEventListeners();
});

/**
 * Загружает корневые комментарии с сервера
 */
function loadComments() {
  showLoading(true);
  
  fetch(`${API_URL}?parent=0`)
    .then(handleResponse)
    .then(comments => {
      allComments = comments || [];
      renderComments(allComments, document.getElementById('commentsList'));
      updateNoCommentsVisibility();
    })
    .catch(showError)
    .finally(() => showLoading(false));
}

/**
 * Обрабатывает ответ от сервера
 */
function handleResponse(response) {
  if (!response.ok) {
    return response.json().then(err => {
      throw new Error(err.error || `Ошибка ${response.status}`);
    }).catch(() => {
      throw new Error(`Ошибка ${response.status}`);
    });
  }
  return response.json().then(data => {
    return data || [];
  });
}

/**
 * Отображает комментарии в DOM
 */
function renderComments(comments, container, level = 0) {
  container.innerHTML = '';
  
  if (!comments || comments.length === 0) {
    if (level === 0) {
      container.innerHTML = '';
    } else {
      container.innerHTML = '<li class="no-comments" style="padding: 12px 16px; margin: 8px 0;">Нет ответов</li>';
    }
    return;
  }
  
  comments.forEach(comment => {
    const commentElement = createCommentElement(comment, level);
    container.appendChild(commentElement);
    
    // Рекурсивная отрисовка дочерних комментариев
    if (comment.children && comment.children.length > 0) {
      const childrenContainer = document.createElement('ul');
      childrenContainer.className = 'comments-list';
      commentElement.querySelector('.children').appendChild(childrenContainer);
      renderComments(comment.children, childrenContainer, level + 1);
    }
  });
}

/**
 * Создает элемент комментария
 */
function createCommentElement(comment, level) {
  const li = document.createElement('li');
  li.className = 'comment';
  li.dataset.id = comment.id;
  
  // Добавляем класс для визуального отступа
  if (level > 0) {
    li.classList.add('indent-level');
  }
  
  li.innerHTML = `
    <div class="comment-content">${escapeHTML(comment.content)}</div>
    <div class="comment-actions">
      <button class="comment-action reply-btn">
        <svg class="reply-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path d="M12 19l-7-7 7-7m-7 7h10a4 4 0 0 1 0 8h1"></path>
        </svg>
        Ответить
      </button>
      <button class="comment-action delete-action">
        <svg class="delete-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
        Удалить
      </button>
    </div>
    <div class="children"></div>
  `;
  
  return li;
}

/**
 * Экранирует HTML для безопасного отображения
 */
function escapeHTML(str) {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "<")
    .replace(/>/g, ">")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

/**
 * Настройка обработчиков событий
 */
function setupEventListeners() {
  // Добавление нового комментария
  document.getElementById('addCommentBtn').addEventListener('click', () => {
    const content = document.getElementById('newCommentContent').value.trim();
    if (!content) return;
    
    addComment(content, null);
  });
  
  // Поиск комментариев
  document.getElementById('searchBtn').addEventListener('click', performSearch);
  document.getElementById('searchInput').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') performSearch();
  });
  
  // Делегирование событий для комментариев
  document.getElementById('commentsList').addEventListener('click', (e) => {
    const commentElement = e.target.closest('.comment');
    if (!commentElement) return;
    
    const id = commentElement.dataset.id;
    
    if (e.target.closest('.reply-btn')) {
      toggleReplyForm(commentElement, id);
    }
    
    if (e.target.closest('.delete-action')) {
      deleteComment(id, commentElement);
    }
  });
}

/**
 * Выполняет поиск по ключевому слову
 */
function performSearch() {
  const keyword = document.getElementById('searchInput').value.trim().toLowerCase();
  if (!keyword) {
    resetSearch();
    return;
  }
  
  // Скрываем все комментарии
  document.querySelectorAll('.comment').forEach(comment => {
    comment.classList.add('hidden');
  });
  
  // Показываем только подходящие и их родителей
  let foundCount = 0;
  document.querySelectorAll('.comment').forEach(comment => {
    const content = comment.querySelector('.comment-content').textContent.toLowerCase();
    if (content.includes(keyword)) {
      foundCount++;
      let current = comment;
      while (current) {
        current.classList.remove('hidden');
        current = current.parentElement.closest('.comment');
      }
    }
  });
  
  if (foundCount === 0) {
    showError(new Error(`По запросу "${keyword}" ничего не найдено`));
  }
}

/**
 * Сбрасывает результаты поиска
 */
function resetSearch() {
  document.querySelectorAll('.comment').forEach(comment => {
    comment.classList.remove('hidden');
  });
  document.getElementById('searchInput').value = '';
}

/**
 * Показывает/скрывает форму ответа
 */
function toggleReplyForm(commentElement, parentId) {
  const formContainer = commentElement.querySelector('.children');
  let replyForm = formContainer.querySelector('.reply-form');
  
  if (replyForm) {
    replyForm.remove();
    return;
  }
  
  replyForm = document.createElement('div');
  replyForm.className = 'reply-form';
  replyForm.innerHTML = `
    <textarea placeholder="Напишите ваш ответ..."></textarea>
    <button class="submit-reply">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="22" y1="2" x2="11" y2="13"></line>
        <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
      </svg>
      Отправить ответ
    </button>
  `;
  
  formContainer.appendChild(replyForm);
  
  // Обработчик отправки ответа
  replyForm.querySelector('.submit-reply').addEventListener('click', () => {
    const content = replyForm.querySelector('textarea').value.trim();
    if (!content) return;
    
    const numericParentId = parentId !== null ? parseInt(parentId, 10) : null;
    addComment(content, numericParentId);
    replyForm.remove();
  });
}

/**
 * Добавляет новый комментарий
 */
function addComment(content, parentId) {
  const payload = { content };
  
  if (parentId !== null && !isNaN(parentId) && Number.isInteger(parentId)) {
    payload.id = parentId;
  }
  
  fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
  .then(handleResponse)
  .then(() => {
    resetSearch();
    loadComments();
    document.getElementById('newCommentContent').value = '';
  })
  .catch(showError);
}

/**
 * Удаляет комментарий
 */
function deleteComment(id, commentElement) {
  if (!confirm('Вы уверены, что хотите удалить этот комментарий и все ответы к нему?')) {
    return;
  }
  
  fetch(`${API_URL}/${id}`, { method: 'DELETE' })
    .then(handleResponse)
    .then(() => {
      commentElement.remove();
      updateNoCommentsVisibility();
    })
    .catch(showError);
}

/**
 * Обновляет видимость сообщения "Нет комментариев"
 */
function updateNoCommentsVisibility() {
  const noComments = document.getElementById('noComments');
  const commentsList = document.getElementById('commentsList');
  
  const hasComments = commentsList.children.length > 0;
  
  if (noComments) {
    noComments.style.display = hasComments ? 'none' : 'block';
  }
}

/**
 * Отображает сообщение об ошибке
 */
function showError(error) {
  const errorElement = document.getElementById('errorMessage');
  errorElement.textContent = error.message;
  errorElement.style.display = 'block';
  
  // Скрываем сообщение об ошибке через 5 секунд
  setTimeout(() => {
    errorElement.style.display = 'none';
  }, 5000);
  
  console.error('Ошибка:', error);
}

/**
 * Показывает индикатор загрузки
 */
function showLoading(isLoading) {
  const btn = document.getElementById('addCommentBtn');
  const originalText = 'Добавить комментарий';
  
  if (isLoading) {
    btn.innerHTML = `<div class="loading-spinner"></div> Загрузка...`;
    btn.disabled = true;
  } else {
    btn.innerHTML = `
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="22" y1="2" x2="11" y2="13"></line>
        <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
      </svg>
      Добавить комментарий
    `;
    btn.disabled = false;
  }
}