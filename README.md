# sAPI RESTful API Documentation

Welcome to the sAPI RESTful API documentation! This guide provides a comprehensive overview of all available endpoints in the API. For every endpoint, ensure to include an `Authorization` header with a valid JWT token for authentication.

---

## Table of Contents

- [Authentication](#authentication)
- [User Endpoints](#user-endpoints)
  - [Sign Up](#sign-up)
  - [Sign In](#sign-in)
  - [Update Profile](#update-profile)
  - [Get Profile](#get-profile)
- [Admin Endpoints](#admin-endpoints)
  - [Update User Rights](#update-user-rights)
  - [Block User](#block-user)
  - [Unblock User](#unblock-user)
  - [Register New User](#register-new-user)
  - [Update User Profile](#update-user-profile)
  - [Get User Profile](#get-user-profile)
  - [Get All Users](#get-all-users)
  - [Delete User](#delete-user)
- [Errors](#error-handling)

---

## Base URL

The base URL for all endpoints is:

```
http://localhost:8080
```

## Authentication

All endpoints, except `/signup` and `/signin`, require authentication. You must include an `Authorization` header in your requests with a valid JWT token. The token should be prefixed with `Bearer `.

---

# User Endpoints


## Register a New User

- **Endpoint:** `POST /signup`
- **Description:** Registers a new user.
- **Request Body:**
  ```json
  {
    "login": "string",
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Response:**
  ```json
  {
    "status": "OK",
    "authdata": {
      "login": "string",
      "password": "string"
    }
  }
  ```

---

## Sign In

- **Endpoint:** `POST /signin`
- **Description:** Authenticates a user and returns a JWT token.
- **Request Body:**
  ```json
  {
    "login": "string",
    "password": "string"
  }
  ```
- **Response:**
  ```json
  {
    "status": "OK",
    "token": "string"
  }
  ```

---

## Update User Profile

- **Endpoint:** `PUT /user/profile`
- **Description:** Updates the logged-in user's profile.
- **Request Body:**
  ```json
  {
    "login": "string",
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Authentication Required**

---

## Get User Profile

- **Endpoint:** `GET /user/profile`
- **Description:** Retrieves the logged-in user's profile information.
- **Response:**
  ```json
  {
    "status": "OK",
    "user": {
      "id": 1,
      "login": "string",
      "username": "string",
      "email": "string",
      "date": "string",
      "block": false,
      "admin": false
    }
  }
  ```
- **Authentication Required**

---

# Admin Endpoints


## Update User Rights

- **Endpoint:** `POST /admin/users/user={id}/rights`
- **Description:** Updates the rights of a user.
- **Request Body:**
  ```json
  {
    "field": "string",
    "value": "string"
  }
  ```
- **Authentication Required**

---

## Block User

- **Endpoint:** `POST /admin/users/user={id}/block`
- **Description:** Blocks a user.
- **Authentication Required**

---

## Unblock User

- **Endpoint:** `POST /admin/users/user={id}/unblock`
- **Description:** Unblocks a user.
- **Authentication Required**

---

## Register a New User (Admin)

- **Endpoint:** `POST /admin/users/registrate/new`
- **Description:** Registers a new user as an admin.
- **Request Body:**
  ```json
  {
    "login": "string",
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Authentication Required**

---

## Update User Profile (Admin)

- **Endpoint:** `PUT /admin/users/profile/user={id}`
- **Description:** Updates a specific user's profile.
- **Request Body:**
  ```json
  {
    "login": "string",
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Authentication Required**

---

## Get User Profile (Admin)

- **Endpoint:** `GET /admin/users/profile/user={id}`
- **Description:** Retrieves the profile information of a specific user.
- **Response:**
  ```json
  {
    "status": "OK",
    "user": {
      "id": 1,
      "login": "string",
      "username": "string",
      "email": "string",
      "date": "string",
      "block": false,
      "admin": false
    }
  }
  ```
- **Authentication Required**

---

## Get All Users

- **Endpoint:** `GET /admin/users/all`
- **Description:** Retrieves a list of all users.
- **Response:**
  ```json
  {
    "status": "OK",
    "users": [
      {
        "id": 1,
        "login": "string",
        "username": "string",
        "email": "string",
        "date": "string",
        "block": false,
        "admin": false
      }
    ]
  }
  ```
- **Authentication Required**

---

## Delete User

- **Endpoint:** `DELETE /admin/users/user={id}`
- **Description:** Deletes a specific user.
- **Authentication Required**

---

# Error Handling

For any errors, you will receive a response with the status set to `"Error"` and an appropriate message describing the issue.

Example Error Response:
```json
{
  "status": "Error",
  "error": "Internal Server Error"
}
```