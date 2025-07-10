# 🚀 **GopherMart — система лояльности на Go**  
**Выпускной проект курса Продвинутый Go-разработчик от Яндекс.Практикум**  

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?logo=postgresql)](https://www.postgresql.org/)
[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?logo=github-actions)](https://github.com/features/actions)

Производственная система лояльности с REST API, реализующая:  
✅ Регистрацию и аутентификацию пользователей  
✅ Учет заказов и баллов лояльности  
✅ Интеграцию с внешним сервисом начислений  
✅ Транзакционные операции с балансом  

---

## 🔍 **Ключевые технологии**  

- **Язык**: Go 1.20+ (чистый код без ORM)  
- **База данных**: PostgreSQL (транзакции, миграции)  
- **Архитектура**:  
  - Многослойная структура (Handler → Service → Repository)  
  - Dependency Injection (Google Wire)  
  - SQL-запросы отдельно от бизнес-логики (sqlc)  
- **Безопасность**:  
  - Хеширование паролей (bcrypt)  
  - JWT-аутентификация  
- **Инфраструктура**:  
  - Миграции (tern)  
  - CI/CD (GitHub Actions)  
  - Конфигурация через env-переменные или аргументы командной строки

---

## 🛠 **Функциональность API**  

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/user/register` | Регистрация пользователя |
| POST | `/api/user/login` | Аутентификация |
| POST | `/api/user/orders` | Загрузка номера заказа |
| GET | `/api/user/orders` | Список заказов пользователя |
| GET | `/api/user/balance` | Текущий баланс |
| POST | `/api/user/balance/withdraw` | Списание баллов |
| GET | `/api/user/withdrawals` | История списаний |

---

## 🏗 **Архитектурные решения**  

### **1. Работа с данными**  
- **Строгая схема БД** с контролем целостности  
- **Транзакционность** критических операций (начисления/списания)  
- **Оптимизированные SQL-запросы** (отдельный слой sqlc)
- **Использование миграций** (tern)


### **2. Безопасность**  
- Пароли хранятся как **хеши bcrypt**  
- **Нет raw SQL** в бизнес-логике  
- Валидация входящих данных (номера заказов по алгоритму Луна)  

### **3. Масштабируемость**  
- **Конкурентная обработка** заказов  
- **Экспоненциальный бекофф** при запросах к accrual-сервису  
- **Кеширование** частых запросов  

---

## 🚀 **Запуск проекта**  

### Требования:  
- Go 1.20+  
- PostgreSQL 15+  

```bash
go build -o gophermart ./cmd/gophermart
./gophermart -d "postgres://user:pass@localhost:5432/gophermart?sslmode=disable"
```

---

<a href="https://t.me/neooshka" target="_blank">
  <img src="https://img.shields.io/badge/Telegram-@neooshka-26A5E4?logo=telegram&style=flat" alt="Telegram">
</a>
<a href="mailto:neeeo@mail.ru">
  <img src="https://img.shields.io/badge/Email-neeeo@mail.ru-D14836?logo=gmail&style=flat" alt="Email">
</a>

[![GitHub Stars](https://img.shields.io/github/stars/Neeeooshka/gopher-club?style=social)](https://github.com/Neeeooshka/gopher-club)