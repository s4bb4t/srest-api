# sAPI - RESTful API Server

`sAPI` is a robust RESTful API server designed for managing user profiles and administering user access. It provides a suite of endpoints for user registration, authentication, profile management, and administrative control over user accounts. The server is built using the `Go` programming language and leverages the `chi` router for handling HTTP requests.

## Table of Contents

- [Overview](#overview)
- [API Endpoints](#api-endpoints)
  - [User Endpoints](#user-endpoints)
    - [Sign Up](#1-sign-up)
    - [Sign In](#2-sign-in)
    - [Update User Profile](#3-update-user-profile)
    - [Get User Profile](#4-get-user-profile)
  - [Admin Endpoints](#admin-endpoints)
    - [Update User Rights](#1-update-user-rights)
    - [Register a New User](#2-register-a-new-user)
    - [Update Any User Profile](#3-update-any-user-profile)
    - [Get User Information](#4-get-user-information)
    - [Get All Users](#5-get-all-users)
    - [Delete a User](#6-delete-a-user)
- [Running the Server](#running-the-server)
- [Authorization](#authorization)
- [Error Handling](#error-handling)

## Overview

`sAPI` provides endpoints to manage user accounts and perform administrative tasks. The server uses JWT for authorization, and all endpoints return JSON responses.

## API Endpoints

### User Endpoints

#### 1. Sign Up

- **URL**: `http://localhost:8080/signup`
- **Method**: `POST`
- **Description**: Registers a new user with a `username`, `password`, and `email`.

**Request Body**:

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**Response**:

```json
{
  "status": "OK",
  "authdata": {
    "username": "string",
    "password": "string"
  }
}
```

---

#### 2. Sign In

- **URL**: `http://localhost:8080/signin`
- **Method**: `POST`
- **Description**: Authenticates a user and issues a JWT token in the `token` cookie.

**Request Body**:

```json
{
  "username": "string",
  "password": "string"
}
```

**Response**:

```json
{
  "status": "OK"
}
```

**Cookies**:

- **Name**: `token`
- **Value**: JWT token
- **Expires**: 12 hours
- **HttpOnly**: true

---

#### 3. Update User Profile

- **URL**: `http://localhost:8080/u/profile/update`
- **Method**: `POST`
- **Description**: Updates the authenticated user's profile information (`username`, `password`, `email`).

**Authorization**: JWT required (set as a cookie)

**Request Body**:

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**Response**:

```json
{
  "status": "OK"
}
```

---

#### 4. Get User Profile

- **URL**: `http://localhost:8080/u/profile`
- **Method**: `GET`
- **Description**: Retrieves the authenticated user's profile information.

**Authorization**: JWT required (set as a cookie)

**Response**:

```json
{
  "status": "OK",
  "user": {
    "username": "string",
    "password": "string",
    "email":    "string",
    "date":     "string",
    "blocked":  "boolean",
    "admin":    "boolean"  
  }
}
```

### Admin Endpoints

#### 1. Update User Rights

- **URL**: `http://localhost:8080/admin/users/rights/{field}`
- **Method**: `POST`
- **Description**: Updates the rights of a user by setting them as blocked, unblocked, admin, or a regular user.

**Authorization**: JWT required (admin access)

**URL Parameters**:

- `{field}`: `block`, `unblock`, `makeadmin`, `makeuser`

**Request Body**:

```json
{
  "username": "string"
}
```

**Response**:

```json
{
  "status": "OK"
}
```
 
---

#### 2. Register a New User

- **URL**: `http://localhost:8080/admin/users/registrate/new`
- **Method**: `POST`
- **Description**: Registers a new user via an admin request. Useful for bulk or managed registrations.

**Authorization**: JWT required (admin access)

**Request Body**:

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**Response**:

```json
{
  "status": "OK",
  "authdata": {
    "username": "string",
    "password": "string"
  }
}
```

---

#### 3. Update Any User Profile

- **URL**: `http://localhost:8080/admin/users/user={username}/update`
- **Method**: `POST`
- **Description**: Updates the profile information of any user by an admin.

**Authorization**: JWT required (admin access)

**Request Body**:

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**Response**:

```json
{
  "status": "OK"
}
```

---

#### 4. Get User Information

- **URL**: `http://localhost:8080/admin/users/profile/user={username}`
- **Method**: `GET`
- **Description**: Retrieves the profile information of any user by an admin.

**Authorization**: JWT required (admin access)

**Response**:

```json
{
  "status": "OK",
  "user": {
    "username": "string",
    "email": "string",
    "is_admin": "boolean",
    "blocked": "boolean"
  }
}
```

---

#### 5. Get All Users

- **URL**: `http://localhost:8080/admin/users/all`
- **Method**: `GET`
- **Description**: Retrieves a list of all registered users in the system.

**Authorization**: JWT required (admin access)

**Response**:

```json
{
  "status": "OK",
  "users": [
    {
      "username": "string",
      "email": "string",
      "is_admin": "boolean",
      "blocked": "boolean"
    }
  ]
}
```

---

#### 6. Delete a User

- **URL**: `http://localhost:8080/admin/users/remove/user={username}`
- **Method**: `DELETE`
- **Description**: Deletes a user account by username.

**Authorization**: JWT required (admin access)

**Response**:

```json
{
  "status": "OK"
}
```

## Running the Server

To run the server, use the following address:

- **URL**: `http://localhost:8080`

The server will start and listen for requests on `localhost` port `8080`.

## Authorization

All endpoints that modify user data or access sensitive information require JWT-based authorization. The JWT token is issued upon successful login and must be provided in the cookie named `token`.

## Error Handling

Errors in requests, such as invalid JSON bodies, incorrect endpoints, or internal server errors, will return a response with the following format:

**Response**:

```json
{
  "status": "Error",
  "error": "Error message"
}
```