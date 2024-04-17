# Golang Authentication System Boilerplate

This Golang project provides a comprehensive boilerplate for an authentication system utilizing JWTs (JSON Web Tokens), sessions, cookies, and Redis. It is designed to offer a robust starting point for any web application requiring user authentication.

## Features

- **JWT Authentication:** Securely authenticate users using JSON Web Tokens.
- **Session Handling:** Manage user sessions to keep users logged in across visits.
- **Cookie Management:** Utilize cookies to store session data for recurring visits.
- **Redis Integration:** Leverage Redis to store and manage session data efficiently.
- **Scalability:** Ready-to-use for scaling in production with Redis handling session data.
- **Security Practices:** Implements basic security measures to protect user data.

## Prerequisites

Before you begin, ensure you have installed the following:
- [Go](https://golang.org/dl/) (version 1.15 or later recommended)
- [Redis](https://redis.io/download), running on localhost.

## Installation

Clone the repository to your desired location:

```bash
git clone https://github.com/Meenachinmay/golang-cookies-auth
cd golang-cookies-auth
```
# Configuration
Configure your application settings and database connection by editing the .env file in the root directory:

```bash
# Rename .env.example to .env and update the values accordingly

DB_URL=postgres://username:password@localhost:5432/database_name
REDIS_ADDR=localhost:6379
REDIS_PW=password
REDIS_DB=0

go run cmd/main.go
```

# Endpoints
The application defines the following endpoints:

- *POST /signup: Register a new user.
- *POST /login: Authenticate a user and returns a JWT.
- *GET /profile: Access a user's profile using a valid JWT.
- *POST /logout: Logout a user and clear the session.

## Example Requests
# Signup
```bash
curl -X POST http://localhost:8080/signup \
-H 'Content-Type: application/json' \
-d '{
    "username": "newuser",
    "password": "password123"
}'
```

# Login
```bash
curl -X POST http://localhost:8080/login \
-H 'Content-Type: application/json' \
-d '{
    "username": "newuser",
    "password": "password123"
}'
```

# Contributing
Contributions are welcome! Feel free to submit pull requests.

# License
Distributed under the MIT License. See LICENSE for more information.

# Acknowledgements
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Go Redis](https://github.com/redis/go-redis)
- [JSON Web Token Go](https://github.com/dgrijalva/jwt-go)