# Archive Service

[![Typing SVG](https://readme-typing-svg.herokuapp.com?color=%2336BCF7&lines=Archive+Service)](https://git.io/typing-svg)


Микросервис для создания ZIP-архивов с REST API интерфейсом.

## 📌 Функционал

- Создание задач на архивацию файлов
- Добавление URL файлов в задачу (.pdf, .jpeg/jpg)
- Получение статуса задачи
- Скачивание готового архива
- Ограничение: 3 одновременно обрабатываемых задачи
- Ограничение: максимум 3 файла на архив

## 🚀 Запуск проекта

### Требования
- Go 1.21+

### Установка

```bash
git clone https://github.com/BabichevDima/2025-07-30-archive-service.git
```

```bash
cd 2025-07-30-archive-service
```

### Запуск

```bash
go run cmd/archive-service/main.go
```

### 📚 Документация

🔗 http://localhost:8080/swagger/index.html