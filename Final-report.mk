#  Library Management System (Go + Echo + MySQL + Redis + RabbitMQ)

A production-ready **Library Management System API**, built with modern backend technologies:

- **Golang + Echo Framework**
- **MySQL (GORM ORM)**
- **Redis** (Caching + Token Bucket Rate Limiting)
- **RabbitMQ** (Event-Driven Reservation Queue)
- **Docker & Docker Compose**
- **Viper** for configuration management

This project follows a clean, modular structure suitable for scalable production systems and professional backend development.

##  Overview

The Library Management System provides backend functionality for:

- User registration & authentication  
- Managing books (CRUD, search, favorites)  
- Loan operations (reserve, borrow, return, cancel)  
- Real-time & async workflows (RabbitMQ)
- Redis caching and login rate limiting  
- Clean architecture for maintainability  

---

##  Key Features

- ✔ JWT-based authentication  
- ✔ Clean Architecture (Handlers → Services → Repositories)
- ✔ MySQL migrations included  
- ✔ Redis caching for book performance  
- ✔ Redis-backed Token Bucket rate limiting  
- ✔ RabbitMQ asynchronous reservation queue  
- ✔ Dockerized deployment  
- ✔ Configurable environment using Viper  


---

#  Modules

---

##  Users Module

Handles:

- User signup  
- Login  
- Password hashing (bcrypt)  
- Email uniqueness  
- Role management (`member`, `student`)  
- JWT generation  

Repository:

- `Create()`
- `FindByEmail()`

Service:

- `Signup()`
- `Login()`

Handler:

- `/users/signup`
- `/users/login`

---

##  Books Module

Features:

- CRUD for books  
- Search (title/author)  
- Favorite system  
- Redis caching:
  - Cache for single book: `book:<id>`
  - Cache for full list: `books:all`

Cache invalidation:

- On create/update/delete

Service:

- `CreateBook()`
- `ListBooks()`
- `GetBookByID()`
- `SearchBooks()`
- `AddToFavorites()`

---

##  Loans Module

Implements the complete loan lifecycle:

- Reserve book (async)
- Confirm borrow
- Return book
- Cancel reservation

All operations:

- Are transactional (GORM transactions)
- Update book stock  
- Update reservation status  
- Track timestamps (`reserved_at`, `borrowed_at`, `due_date`, etc.)

---

---

# API summary

---

## Users

POST /users/signup
POST /users/login

## Books

POST   /books
GET    /books
GET    /books/:id
PUT    /books/:id
DELETE /books/:id
GET    /books/search
POST   /books/:id/favorite/:user_id
GET    /books/favorites/:user_id

## Loans (JWT Required)

POST /api/loans/reserve
POST /api/loans/:id/confirm
POST /api/loans/:id/return
POST /api/loans/:id/cancel
GET  /api/loans/user/:userID

