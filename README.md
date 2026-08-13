<div align="center">

# 🔄 DevSync

### A Role-Based Backend API for Team Collaboration & Task Management

Built with Go for speed, reliability, and clean architecture.

[![Go](https://img.shields.io/badge/Go-1.x-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![JWT](https://img.shields.io/badge/Auth-JWT-black?style=for-the-badge&logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
[![Swagger](https://img.shields.io/badge/API%20Docs-Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)
[![License](https://img.shields.io/badge/License-Learning%20Project-yellow?style=for-the-badge)](#-license)

</div>

---

## 📖 About

**DevSync** is a backend API that powers a team collaboration and task management platform. Managers can spin up projects and assign tasks, employees can track progress and collaborate through comments, and admins oversee the whole system through approval workflows for manager access and account deletion.

Built entirely in **Go**, DevSync focuses on clean role-based access control, a secure JWT auth flow, and a fully documented, Dockerized backend ready to plug into any frontend.

---

## 📑 Table of Contents

- [✨ Features](#-features)
- [🛠️ Tech Stack](#️-tech-stack)
- [📁 Project Structure](#-project-structure)
- [👥 User Roles](#-user-roles)
- [🔐 Authentication](#-authentication)
- [🛡️ Authorization](#️-authorization)
- [📡 API Endpoints](#-api-endpoints)
- [🗄️ Database Migrations](#️-database-migrations)
- [⚙️ Environment Variables](#️-environment-variables)
- [🚀 Getting Started](#-getting-started)
- [🐳 Docker Architecture](#-docker-architecture)
- [📘 Swagger Documentation](#-swagger-documentation)
- [🔄 Account Deletion Flow](#-account-deletion-flow)
- [📊 Task Activity Flow](#-task-activity-flow)
- [🎯 Project Goals](#-project-goals)
- [🔮 Future Improvements](#-future-improvements)
- [📦 Repository](#-repository)
- [📄 License](#-license)

---

## ✨ Features

<table>
<tr>
<td width="50%" valign="top">

**🔑 Core**
- User registration & login
- JWT-based authentication
- Role-based authorization
- Bcrypt password hashing
- User profile management

</td>
<td width="50%" valign="top">

**📋 Task & Project Management**
- Project management (CRUD)
- Task management (CRUD)
- Employee task assignment
- Task progress tracking
- Task comments & activity history

</td>
</tr>
<tr>
<td width="50%" valign="top">

**🧑‍💼 Workflows**
- Manager request system
- Admin approval / rejection of manager requests
- Account deletion request system
- Admin approval / rejection of deletion requests
- Deleted-account protection

</td>
<td width="50%" valign="top">

**⚙️ Platform**
- Manager / Employee / Admin dashboards
- Swagger API documentation
- MySQL database with migrations
- Docker & Docker Compose support

</td>
</tr>
</table>

---

## 🛠️ Tech Stack

| Technology | Purpose |
|---|---|
| 🐹 **Go** | Backend programming language |
| 🌐 **Gorilla Mux** | HTTP router |
| 🐬 **MySQL 8.4** | Relational database |
| 🔐 **JWT** | Authentication |
| 🔒 **bcrypt** | Password hashing |
| 📘 **Swagger** | API documentation |
| 🐳 **Docker** | Containerization |
| 🧩 **Docker Compose** | Multi-container setup |
| 🗄️ **golang-migrate** | Database migrations |
| 📄 **godotenv** | Environment variable loading |

---

## 📁 Project Structure

```text
DevSync/
│
├── config/
│   └── config.go
│
├── database/
│   ├── db.go
│   └── migrations/
│       ├── 001_create_users.up.sql / .down.sql
│       ├── 002_create_projects.up.sql / .down.sql
│       ├── 003_create_tasks.up.sql / .down.sql
│       ├── 004_create_comment.up.sql
│       ├── 005_create_activity.up.sql
│       ├── 006_add_created_at.up.sql
│       ├── 007_create_manager_requests.up.sql
│       ├── 008_update_user_roles.up.sql
│       ├── 009_create_account_deletion.up.sql
│       └── 010_add_account_deletion_field.up.sql
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── handlers/
│   ├── account_deletion.go
│   ├── activity.go
│   ├── auth.go
│   ├── comment.go
│   ├── dashboard.go
│   ├── health.go
│   ├── manager_requests.go
│   ├── project.go
│   ├── task.go
│   └── user.go
│
├── middleware/
│   ├── admin.go
│   ├── auth.go
│   ├── logger.go
│   └── manager.go
│
├── models/
│   ├── account_deletion_request.go
│   ├── activity.go
│   ├── comment.go
│   ├── dashboard.go
│   ├── manager_request.go
│   ├── project.go
│   ├── task.go
│   └── user.go
│
├── routes/
│   └── routes.go
│
├── seed/
│   └── admin.go
│
├── utils/
│   ├── jwt.go
│   ├── password.go
│   ├── response.go
│   ├── sanitize.go
│   └── validation.go
│
├── .dockerignore
├── .env
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

---

## 👥 User Roles

DevSync operates on **three roles**, each with a tailored set of permissions.

### 🛡️ Admin
Handles all administrative oversight:
- View the admin dashboard
- View / approve / reject manager requests
- View / approve / reject account deletion requests

### 👔 Manager
Owns projects and tasks end-to-end:
- Create, view, update, and delete **projects**
- Create, view, update, and delete **tasks**
- View task comments & task activities
- View the manager dashboard

### 👷 Employee
Executes and collaborates on assigned work:
- View profile & assigned tasks
- Update task status & progress
- Add and view task comments
- View task activities
- View the employee dashboard
- Request manager access
- Request account deletion

---

## 🔐 Authentication

DevSync uses **JWT-based authentication**. After a successful login, the API returns a signed JWT token that must be included in the `Authorization` header for all protected routes:

```text
Authorization: Bearer <JWT_TOKEN>
```

**Auth middleware pipeline:**

```text
1️⃣  Read the Authorization header
2️⃣  Validate the Bearer token format
3️⃣  Verify the JWT signature & claims
4️⃣  Extract the user ID and role
5️⃣  Check whether the account has been deleted
6️⃣  Attach authenticated user info to the request context
7️⃣  Forward the request to the next handler
```

---

## 🛡️ Authorization

Role-based middleware gates access to protected resources:

```text
              Auth Middleware
                    │
        ┌───────────┼───────────┐
        │                       │
 Admin Middleware      Manager Middleware
```

Employees are authorized through the authenticated user context and inline role checks on employee-specific endpoints.

---

## 📡 API Endpoints

### 🌍 Public

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/register` | Register a user |
| `POST` | `/login` | Login and receive JWT |

### 👤 User

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/me` | Get current user profile |
| `GET` | `/api/users` | Get users |

### 📨 Manager Requests

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/manager-requests` | Create manager request |
| `GET` | `/api/manager-request` | Get current manager request |
| `GET` | `/api/admin/manager-requests` | Get all manager requests |
| `PUT` | `/api/admin/manager-requests/{id}/approve` | Approve manager request |
| `PUT` | `/api/admin/manager-requests/{id}/reject` | Reject manager request |

### 📁 Projects

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/projects` | Create project |
| `GET` | `/api/projects` | Get manager projects |
| `GET` | `/api/projects/{id}` | Get project |
| `PUT` | `/api/projects/{id}` | Update project |
| `DELETE` | `/api/projects/{id}` | Delete project |

### ✅ Manager Tasks

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/tasks` | Create task |
| `GET` | `/api/tasks` | Get manager tasks |
| `GET` | `/api/tasks/{id}` | Get task |
| `PUT` | `/api/tasks/{id}` | Update task |
| `DELETE` | `/api/tasks/{id}` | Delete task |

### 🧑‍💻 Employee Tasks

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/my-tasks` | Get assigned tasks |
| `PUT` | `/api/my-tasks/{id}/progress` | Update task progress |

### 💬 Comments

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/tasks/{id}/comments` | Add a comment |
| `GET` | `/api/tasks/{id}/comments` | Get task comments |
| `DELETE` | `/api/comments/{id}` | Delete a comment |

> Comments include the name of the user who created them.

### 🕓 Activities

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/tasks/{id}/activities` | Get task activity history |

> Activity records include the name of the user who performed the action.

### 📊 Dashboards

| Role | Method | Endpoint |
|---|---|---|
| Manager | `GET` | `/api/dashboard/manager` |
| Employee | `GET` | `/api/employee/dashboard` |
| Admin | `GET` | `/api/admin/dashboard` |

### 🗑️ Account Deletion

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/account-deletion-request` | Create deletion request |
| `GET` | `/api/account-deletion-request` | Get current deletion request |
| `GET` | `/api/admin/account-deletion-requests` | Get deletion requests |
| `PUT` | `/api/admin/account-deletion-requests/{id}/approve` | Approve deletion |
| `PUT` | `/api/admin/account-deletion-requests/{id}/reject` | Reject deletion |

---

## 🗄️ Database Migrations

DevSync uses **`golang-migrate`** for schema versioning. Migration files live in:

```text
database/migrations/
```

| Version | Description |
|---|---|
| `001` | Users |
| `002` | Projects |
| `003` | Tasks |
| `004` | Comments |
| `005` | Activities |
| `006` | Created-at changes |
| `007` | Manager requests |
| `008` | User roles |
| `009` | Account deletion requests |
| `010` | Account deletion field |

**Apply migrations:**
```bash
migrate -path database/migrations -database "<MYSQL_DSN>" up
```

**Check current version:**
```bash
migrate -path database/migrations -database "<MYSQL_DSN>" version
```

---

## ⚙️ Environment Variables

Create a `.env` file in the project root:

```env
DB_HOST=localhost
DB_PORT=3310
DB_USER=root
DB_PASSWORD=<your-password>
DB_NAME=DevSync-api

JWT_SECRET=<your-jwt-secret>

PORT=8081

ADMIN_NAME=<admin-name>
ADMIN_EMAIL=<admin-email>
ADMIN_PASSWORD=<admin-password>
```


---

## 🚀 Getting Started

### ▶️ Option 1 — Run Without Docker

```bash
# 1. Download dependencies
go mod download

# 2. Run the application
go run main.go
```

The API will be live at:

```text
http://localhost:8081
```

### 🐳 Option 2 — Run With Docker

DevSync ships with both a backend container and a MySQL container.

```bash
# Build and start
docker compose up --build

# Run in detached mode
docker compose up -d --build

# Check running containers
docker compose ps

# Stop containers
docker compose down
```

Compose services:

```text
devsync-api      → the Go backend
devsync-mysql    → the MySQL 8.4 database
```

Inside Docker, the backend connects to MySQL via the service name:

```text
DB_HOST=mysql
DB_PORT=3306
```

---

## 🐳 Docker Architecture

```text
                 ┌─────────────────────┐
                 │      Client         │
                 │ Swagger / Frontend  │
                 └──────────┬──────────┘
                            │
                            │ HTTP
                            ▼
                 ┌─────────────────────┐
                 │    devsync-api      │
                 │      Go API         │
                 │     Port 8081       │
                 └──────────┬──────────┘
                            │
                            │ MySQL
                            ▼
                 ┌─────────────────────┐
                 │    devsync-mysql    │
                 │     MySQL 8.4       │
                 │      Port 3306      │
                 └─────────────────────┘
```

---

## 📘 Swagger Documentation

Once the server is running, explore and test every endpoint interactively at:

```text
http://localhost:8081/swagger/index.html
```

Use Swagger UI to:
- 🔍 Browse available endpoints
- 📖 Read full API documentation
- 🧪 Test requests live
- 🔑 Authorize using your JWT
- 🔬 Inspect request/response schemas

---

## 🔄 Account Deletion Flow

```text
        Employee / User
               │
               │  Create deletion request
               ▼
        Pending Request
               │
               ▼
             Admin
        ┌──────┴───────┐
        │              │
     Approve         Reject
        │              │
        ▼              ▼
  Account Deleted   Request Rejected
```

When an account is deleted:
- 🚫 The user can no longer log in
- 🛡️ Auth middleware checks the deleted status on every request
- 🔒 All protected access is blocked for the deleted account

---

## 📊 Task Activity Flow

Every task action is logged to the activities table:

```text
task_created
task_updated
task_deleted
progress_updated
```

Each activity record captures:

| Field | Description |
|---|---|
| 🆔 User ID | Who performed the action |
| 🙋 User Name | Display name of the actor |
| 📌 Task ID | Related task |
| ⚡ Action | The action type |
| 📝 Details | Additional context |
| 🕒 Created At | Timestamp of the action |

---

## 🎯 Project Goals

DevSync aims to be a solid backend foundation for a team collaboration platform, featuring:

- 🔐 Secure authentication
- 🛡️ Role-based access control
- 📁 Project & task management
- 🤝 Employee collaboration
- 📊 Activity tracking
- 🧑‍💼 Administrative workflows

---

## 🔮 Future Improvements

- [ ] Automated tests
- [ ] Refresh tokens
- [ ] Pagination
- [ ] Search and filtering
- [ ] Email notifications
- [ ] Real-time notifications
- [ ] WebSocket support
- [ ] CI/CD pipeline
- [ ] Production deployment
- [ ] Rate limiting
- [ ] Advanced analytics

---

## 📦 Repository

```text
https://github.com/Nehasirohi07/devSync
```

---

## 📄 License

This project is developed as part of backend/internship learning and project development.

<div align="center">

---

Made with ⚙️ Go and 💙 by **Neha Sirohi**

</div>
