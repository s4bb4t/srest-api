
---

# sAPI RESTful API Server

`sAPI` is a robust and secure RESTful API server built using the Go programming language. This server is designed to handle user authentication, profile management, and administrative tasks. The API follows best practices for RESTful API design, ensuring scalability, security, and maintainability.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Server](#running-the-server)
- [API Endpoints](#api-endpoints)
  - [User Endpoints](#user-endpoints)
  - [Admin Endpoints](#admin-endpoints)
- [Middleware](#middleware)
- [Error Handling](#error-handling)
- [Logging](#logging)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Authentication**: Secure registration and login endpoints.
- **Profile Management**: Update and retrieve user profiles.
- **Admin Functions**: Manage user permissions, delete users, and retrieve detailed user information.
- **JWT Authentication**: Secure API routes using JSON Web Tokens.
- **Structured Logging**: Detailed logging for easier debugging and monitoring.
- **Middleware Integration**: Comprehensive middleware for request handling, error recovery, and logging.

## Getting Started

### Prerequisites

Before running the server, ensure that you have the following installed:

- [Go](https://golang.org/dl/) (version 1.18 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) (for the database)

### Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/sabbatD/srest-api.git
cd srest-api
```

### Configuration

Configure the application by editing the `config.yaml` file located in the `internal/config` directory. The configuration file includes settings for:

- **Database Connection**: `DbString` for PostgreSQL connection.
- **Server Address**: `Address` for the server's binding address.
- **Timeouts**: `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` for HTTP server.
- **Environment**: `Env` for setting the application environment (e.g., `development`, `production`).

### Running the Server

To start the server, use the following command:

```bash
go run main.go
```

The server will start on the address specified in the configuration file. Logs will indicate that the server is running and ready to accept connections.

## API Endpoints

### User Endpoints

| Method | Endpoint           | Description                                   |
|--------|--------------------|-----------------------------------------------|
| POST   | `/signup`          | Register a new user.                          |
| POST   | `/signin`          | Authenticate a user and issue a JWT.          |
| GET    | `/u/profile`       | Retrieve the profile of the authenticated user.|
| POST   | `/u/profile/update`| Update the authenticated user's profile.      |

### Admin Endpoints

| Method | Endpoint                                   | Description                                                   |
|--------|--------------------------------------------|---------------------------------------------------------------|
| POST   | `/admin/users/rights/{field}`              | Update user rights (block/unblock, make admin/user).           |
| POST   | `/admin/users/registrate/new`              | Register a new user (admin only).                              |
| POST   | `/admin/users/user?username={username}/update`| Update all fields of a user's profile (admin only).            |
| GET    | `/admin/users/profile/user?={username}`    | Retrieve detailed information for a specific user (admin only).|
| GET    | `/admin/users/all`                         | Retrieve an array of all users with detailed information.      |
| DELETE | `/admin/users/remove?username={username}`  | Delete a user by username (admin only).                        |

### Middleware

- **RequestID**: Adds a unique request ID to each request for tracking purposes.
- **Logger**: Logs incoming requests and their corresponding responses.
- **Recoverer**: Recovers from panics and returns a 500 internal server error.
- **URLFormat**: Handles URLs with different formats.

### Error Handling

The API returns structured JSON responses for all errors, with a clear message and status code. Common error codes include:

- `400 Bad Request`: Invalid request data.
- `401 Unauthorized`: Authentication failed or missing token.
- `403 Forbidden`: Insufficient permissions.
- `404 Not Found`: Resource not found.
- `500 Internal Server Error`: General server error.

### Logging

`sAPI` uses structured logging via the `slog` package. Logs include detailed information such as the operation name, request ID, and any relevant data or errors. The log level can be configured via the `Env` setting in the configuration file:

- **Info**: General operational messages.
- **Debug**: Detailed debugging information.
- **Error**: Errors that occurred during processing.

### Contributing

Contributions are welcome! Please follow the standard GitHub flow for contributing:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a pull request.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

This `README.md` file should provide a comprehensive and clear guide for developers using or contributing to the `sAPI` project. It covers everything from installation and configuration to API usage and contribution guidelines.