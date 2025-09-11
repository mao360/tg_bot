# Telegram Bot

Простой телеграм бот с двумя режимами работы.

## Функциональность

1. **Режим 1: Перевернуть строку** - выводит введенную пользователем строку наоборот
2. **Режим 2: Hello + строка** - выводит "Hello" + введенная строка

## Установка и запуск

### Вариант 1: Запуск с Docker (рекомендуется)

1. Убедитесь, что у вас установлен Docker и Docker Compose

2. Соберите и запустите контейнер:
   ```bash
   docker-compose up --build
   ```

3. Для запуска в фоновом режиме:
   ```bash
   docker-compose up -d --build
   ```

4. Для остановки:
   ```bash
   docker-compose down
   ```

### Вариант 2: Запуск без Docker

1. Убедитесь, что у вас установлен Go (версия 1.23 или выше)

2. Установите зависимости:
   ```bash
   go mod tidy
   ```

3. Скомпилируйте проект:
   ```bash
   go build -o bot.exe main.go
   ```

4. Запустите бота:
   ```bash
   ./bot.exe
   ```

## Использование

1. Найдите вашего бота в Telegram по токену
2. Отправьте команду `/start`
3. Выберите режим работы с помощью кнопок
4. Введите текст для обработки
5. Получите результат и вернитесь к выбору режима

## Docker команды

### Основные команды:
```bash
# Собрать образ
docker build -t tg_bot .

# Запустить контейнер
docker run -d --name tg_bot tg_bot

# Остановить контейнер
docker stop tg_bot

# Удалить контейнер
docker rm tg_bot

# Посмотреть логи
docker logs tg_bot

# Посмотреть логи в реальном времени
docker logs -f tg_bot
```

### Docker Compose команды:
```bash
# Собрать и запустить
docker-compose up --build

# Запустить в фоне
docker-compose up -d --build

# Остановить
docker-compose down

# Посмотреть логи
docker-compose logs -f

# Пересобрать без кэша
docker-compose build --no-cache
```

## Структура проекта

- `main.go` - основной файл с кодом бота
- `go.mod` - файл зависимостей Go
- `Dockerfile` - конфигурация Docker образа
- `docker-compose.yml` - конфигурация Docker Compose
- `.dockerignore` - файлы, исключаемые из Docker контекста
- `bot.exe` - скомпилированный исполняемый файл (создается после сборки)
