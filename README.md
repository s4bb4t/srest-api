### Documentation for `sAPI`

# sAPI - RESTful API Server

`sAPI` is a robust RESTful API server designed for managing user profiles and administering user access. It provides a suite of endpoints for user registration, authentication, profile management, and administrative control over user accounts. The server is built using the `Go` programming language and leverages the `chi` router for handling HTTP requests.

## Table of Contents

- [Overview](#overview)
- [API Endpoints](#api-endpoints)
  - [User Endpoints](#user-endpoints)
    - [Sign Up](#sign-up)
    - [Sign In](#sign-in)
    - [Update User Profile](#update-user-profile)
    - [Get User Profile](#get-user-profile)
  - [Admin Endpoints](#admin-endpoints)
    - [Update User Rights](#update-user-rights)
    - [Register a New User](#register-a-new-user)
    - [Update Any User Profile](#update-any-user-profile)
    - [Get User Information](#get-user-information)
    - [Get All Users](#get-all-users)
    - [Delete a User](#delete-a-user)
- [Running the Server](#running-the-server)
- [Authorization](#authorization)
- [Error Handling](#error-handling)

## Overview

`sAPI` provides endpoints to manage user accounts and administrative tasks. The server uses JWT for authorization, and all endpoints return JSON responses.

## API Endpoints

### User Endpoints

#### 1. Sign Up

**URL**: `/signup`

**Method**: `POST`

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

#### 2. Sign In

**URL**: `/signin`

**Method**: `POST`

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

#### 3. Update User Profile

**URL**: `/u/profile/update`

**Method**: `POST`

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

#### 4. Get User Profile

**URL**: `/u/profile`

**Method**: `GET`

**Authorization**: JWT required (set as a cookie)

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

### Admin Endpoints

#### 1. Update User Rights

**URL**: `/admin/users/rights/{field}`

**Method**: `POST`

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

#### 2. Register a New User

**URL**: `/admin/users/registrate/new`

**Method**: `POST`

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

#### 3. Update Any User Profile

**URL**: `/admin/users/user?username={username}/update`

**Method**: `POST`

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

#### 4. Get User Information

**URL**: `/admin/users/profile/user?username={username}`

**Method**: `GET`

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

#### 5. Get All Users

**URL**: `/admin/users/all`

**Method**: `GET`

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

#### 6. Delete a User

**URL**: `/admin/users/remove?username={username}`

**Method**: `DELETE`

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

This structure is used to communicate issues that may arise during request processing.