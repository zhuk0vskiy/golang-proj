<template>
  <div class="app">
    <!-- Header -->
    <header class="header">
      <div class="header-left">
        <button class="toggle-btn">▼</button>
        <h1 class="title">Audio Studio</h1>
      </div>
      <div class="header-right">
        <button class="logout-btn">
          Выйти
        </button>
        <div class="profile-pic"></div>
      </div>
    </header>

    <div class="main-container">
      <!-- Sidebar -->
      <aside class="sidebar">
        <nav class="menu">
          <button class="menu-item menu-item--active">
            <svg class="menu-icon" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
            </svg>
            Создать
          </button>
          <button class="menu-item">
            <svg class="menu-icon" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            Удалить
          </button>
          <button class="menu-item">
            <svg class="menu-icon" viewBox="0 0 20 20" fill="currentColor">
              <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
              <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
            </svg>
            Просмотреть
          </button>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="content">
        <div class="form-container">
          <div class="form-group">
            <div class="select-wrapper">
              <select class="form-select">
                <option disabled selected>Выберите студию</option>
              </select>
              <div class="select-arrow">
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </div>
            </div>

            <input
              type="text"
              placeholder="Введите дату"
              class="form-input"
            />

            <div class="time-inputs">
              <input
                type="text"
                placeholder="Начало"
                class="form-input"
              />
              <input
                type="text"
                placeholder="Конец"
                class="form-input"
              />
            </div>

            <button class="submit-btn">
              Продолжить
            </button>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AudioStudio',
  data() {
    return {
      selectedStudio: null,
      date: '',
      startTime: '',
      endTime: ''
    }
  }
}
</script>

<style>
/* Reset CSS */
*, *::before, *::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

/* Root styles */
:root {
  --primary-color: #3b82f6;
  --primary-hover: #2563eb;
  --background-color: #f9fafb;
  --sidebar-color: #f3f4f6;
  --border-color: #d1d5db;
  --text-color: #111827;
  --header-height: 64px;
}

/* Base styles */
html, body {
  height: 100%;
  width: 100%;
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  font-size: 16px;
  color: var(--text-color);
  overflow: hidden;
}

/* Component styles */
.app {
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  background-color: var(--background-color);
  overflow: hidden;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1.5rem;
  background-color: #ffffff;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  height: var(--header-height);
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.toggle-btn {
  padding: 0.5rem;
  background: none;
  border: none;
  cursor: pointer;
}

.title {
  font-size: 1.25rem;
  font-weight: 500;
}

.logout-btn {
  padding: 0.5rem 1rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.logout-btn:hover {
  background-color: var(--primary-hover);
}

.profile-pic {
  width: 2.5rem;
  height: 2.5rem;
  background-color: #e5e7eb;
  border-radius: 9999px;
}

.main-container {
  display: flex;
  height: calc(100vh - var(--header-height));
  overflow: hidden;
}

.sidebar {
  width: 16rem;
  background-color: var(--sidebar-color);
  padding: 1rem;
  flex-shrink: 0;
  overflow-y: auto;
  border-right: 1px solid var(--border-color);
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  background: none;
  cursor: pointer;
  transition: background-color 0.2s;
  text-align: left;
  color: var(--text-color);
}

.menu-item:hover {
  background-color: #e5e7eb;
}

.menu-item--active {
  background-color: #eff6ff;
  color: var(--primary-color);
}

.menu-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
}

.content {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.form-container {
  max-width: 28rem;
  width: 100%;
  margin: 2rem auto;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.select-wrapper {
  position: relative;
}

.form-select {
  width: 100%;
  padding: 0.625rem;
  background-color: white;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  appearance: none;
  cursor: pointer;
  font-size: 1rem;
}

.select-arrow {
  position: absolute;
  right: 0.5rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  width: 1rem;
  height: 1rem;
  color: #9ca3af;
}

.form-input {
  width: 100%;
  padding: 0.625rem;
  background-color: white;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  font-size: 1rem;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.time-inputs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.submit-btn {
  width: 100%;
  padding: 0.625rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 1rem;
  font-weight: 500;
}

.submit-btn:hover {
  background-color: var(--primary-hover);
}

/* Scrollbar styles */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background-color: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background-color: #888;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background-color: #666;
}
</style>