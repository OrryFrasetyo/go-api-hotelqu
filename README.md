# HotelQu API

## Overview
HotelQu API is a REST API-based hotel management system built with Go and the Gin framework. This API is designed to manage various aspects of hotel operations, including:

- 👥 Employee Management
- 📅 Shift and Schedule
- ⏰ Presence System
- 🏢 Department Management
- 👔 Position Management

## Technology Used
- **Backend:** Go (Golang) with Gin Framework
- **Database:** MySQL
- **Autentikasi:** JWT (JSON Web Token)
- **Arsitektur:** Domain-Driven Design (DDD)

## Pre-requisite
Before you begin, make sure your system has:

- **Go** version 1.24.2 atau latest ([Download Go](https://golang.org/dl/))
- **MySQL** Server ([Download MySQL](https://dev.mysql.com/downloads/))
- **Git** for clone repository

## Installation Guide

### 1. Cloning repository
```bash
git clone https://github.com/OrryFrasetyo/go-api-hotelqu.git
cd go-api-hotelqu
```

### 2. Configuration Database
1. Make new database mysql:
```sql
CREATE DATABASE hotelqu_db;
```

2. Adjust the database configuration in `models/setup.go` with your MySQL settings:
- **Host:** 127.0.0.1
- **Port:** 3306
- **User:** root
- **Password:** (sesuaikan)
- **Database:** hotelqu_db

### 3. Installation Dependency
```bash
go mod tidy
```

### 4. Run server application
```bash
go run main.go
```

The application will run on `http://localhost:3000`

## Main Features
- ✅ Login
- ✅ Register
- ✅ Department Management
- ✅ Position Management
- ✅ Shift Management
- ✅ Presence System
- ✅ Schedule Employee
- ✅ Profile Employee Management

## Project Structure
```
api-hotelqu-go-ddd/
├── controllers/     # Handler HTTP requests
├── middlewares/    # Middleware (auth, logging, etc)
├── models/         # Domain models & database
├── utils/          # Utilitas (JWT, helper functions)
└── uploads/        # File upload directory
```

---

