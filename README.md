# sAPI: RestAPI Server by s4bb4t for Timur Seledkov

Welcome to **sAPI**, a RESTful API server developed by s4bb4t. This project serves as a backend framework aimed at supporting frontend education and experimentation.

## Author
**s4bb4t**
- **GitHub**: [sabbatD](https://github.com/sabbatD)
- **Telegram**: [sabbathtm](https://t.me/sabbathtm)
- **Email**: mentalrapewhr@gmail.com

---

## API Documentation

### Endpoints Overview

The API is organized into several routes that support user registration, authentication, and administrative tasks. Below is a detailed guide to each available endpoint.

### POST Requests

#### User Endpoints

##### **Sign Up**
- **Endpoint**: `POST http://localhost:8080/signup`
- **Description**: Registers a new user in the system.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Response**:
  ```json
  {
    "status": "string",
    "error": "string (optional)"
  }
  ```
  - If successful, the response will include the username and password used for registration.

##### **Sign In**
- **Endpoint**: `POST http://localhost:8080/signin`
- **Description**: Authenticates a user and, if the user has admin rights, returns a JWT token in the cookies.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response**: No JSON response is returned. However, a JWT token is set in the cookies if the user is an admin.

#### Admin Endpoints

**Note**: All admin routes require a valid JWT token extracted from cookies. Only signed-in users with admin rights can make these requests. Users who are blocked cannot perform admin actions, even if they possess admin rights.

##### **Update Rights**
- **Endpoint**: `POST http://localhost:8080/admin/rights/{field}`
- **Description**: Modifies a specified attribute for a user.
- **URL Path Parameter**:
  - `{field}`: Can be one of the following values:
    - `block` – Block the user.
    - `unblock` – Unblock the user.
    - `makeadmin` – Grant admin rights.
    - `makeuser` – Revoke admin rights.
- **Request Body**:
  ```json
  {
    "username": "string"
  }
  ```
- **Response**:
  ```json
  {
    "status": "string",
    "error": "string (optional)"
  }
  ```

##### **Remove User**
- **Endpoint**: `POST http://localhost:8080/admin/remove`
- **Description**: Deletes the specified user from the system.
- **Request Body**:
  ```json
  {
    "username": "string"
  }
  ```
- **Response**:
  ```json
  {
    "status": "string",
    "error": "string (optional)"
  }
  ```

### GET Requests

#### Admin Endpoints

##### **List All Users**
- **Endpoint**: `GET http://localhost:8080/admin/allusers`
- **Description**: Retrieves an array of all users, including all associated information.
- **Response**:
  ```json
  {
    "status": "string",
    "error": "string (optional)",
    "users": [
      {
        "username": "string",
        "password": "string",
        "email": "string",
        "date": "string",
        "blocked": "boolean",
        "admin": "boolean"
      }
    ]
  }
  ```

---

### Final Notes

This API is designed to be simple yet powerful, providing essential features for user management and administration. Whether you’re signing up users or managing admin rights, **sAPI** is a solid foundation for your frontend development.

For any questions or contributions, feel free to reach out via the contact methods listed above.