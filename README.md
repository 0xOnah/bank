# 🏦 Bank API

The **Bank API** is a robust and scalable backend system for managing user accounts, transactions, and money transfers.  
It follows **Clean Architecture** principles to ensure maintainability, testability, and clear separation of concerns.

---

## ✨ Features

- **User Account Management** – Create, retrieve, update, and delete bank accounts securely.  
- **Money Transfers** – Support for deposits, withdrawals, and internal transfers between accounts.  
- **Transaction History** – View all past transactions with pagination and filtering.  
- **Role-Based Access Control (RBAC)** – Differentiate between regular users and administrators.  
- **Event-Driven Ready** – Designed to support future integrations with external services.

---

## 🛠️ Tech Stack

- **Language:** Go (1.23+)
- **Frameworks / Tools:**
  - **gRPC** – High-performance RPC for internal service-to-service communication.
  - **HTTP + REST** – Public-facing API for clients and external integrations.
  - **sqlc** – Compile-time type-safe SQL query generation.
  - **PostgreSQL** – Reliable relational database for storing account and transaction data.
  - **Docker** – Containerized development and deployment environment.
  - **GitHub Actions** – CI/CD pipeline with automated tests and linting.
  - **GIN Framework**

---

## 🏗️ Architecture


- **Transport Layer** – Handles HTTP routes and gRPC endpoints.  
- **Service Layer** – Contains business logic and orchestrates workflows.  
- **Repository Layer** – Interfaces with the database using `sqlc`.  
- **Database** – PostgreSQL with transaction-safe operations and migrations.

---

## ✅ Testing & Quality

- **Unit Tests** – Service layer tested with mocks (`GoMock`).
- **Integration Tests** – Run in CI against a real PostgreSQL container (via GitHub Actions).
- **Linting & Static Analysis** – Enforced with `golangci-lint`.
- **Logging & Observability** – Structured logs for debugging and monitoring.

---

## 🚀 Getting Started

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/your-username/bank-api.git
cd bank-api
```

### 2️⃣ Run Locally with Docker

```bash
docker-compose up --build
```

### 3️⃣ Run Tests

```bash 
make test
```



## 📦 CI/CD

### This project uses GitHub Actions for:

- **Automated tests on each push or pull request**
- **Linting and static code analysis**
- **Building Docker images (optional)**