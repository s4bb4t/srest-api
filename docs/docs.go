// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "s4bb4t",
            "email": "s4bb4t@yandex.ru"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Gets users by accepting a url query payload containing filters.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search in login or username or email",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sotrOrder asc or desc or none (asc, decs - sotrOrder by email, none - sotrOrder by id)",
                        "name": "sotrOrder",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "block status",
                        "name": "isBlocked",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit of users for query",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Retrieve successful. Returns users.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_admin.Users"
                        }
                    },
                    "401": {
                        "description": "User context not found.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Not enough rights.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves user's profile by id.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Retrieve user's profile",
                "responses": {
                    "200": {
                        "description": "Retrieve successful. Returns user.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "Missing or wrong id.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "User context not found.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Not enough rights.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates user by id by accepting a JSON payload containing user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Update user's fields",
                "parameters": [
                    {
                        "description": "Any user data",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.PutUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Update successful. Returns user ok.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "Login or email already used.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "User context not found.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Not enough rights.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Removes user by id in url.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Remove user",
                "responses": {
                    "200": {
                        "description": "Remove successful. Returns ok.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Missing or wrong id.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "User context not found.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Not enough rights.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}/block": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Blocks user by id in url.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Block user",
                "responses": {
                    "200": {
                        "description": "Block successful. Returns user ok.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "No such field.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}/rights": {
            "post": {
                "description": "Updates user by id by accepting a JSON payload containing user's rights.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Update user's rights",
                "parameters": [
                    {
                        "description": "Complete user data",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_admin.UpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Update successful. Returns ok.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "No such field.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}/unlock": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Unlocks user by id in url.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Unlock user",
                "responses": {
                    "200": {
                        "description": "Unlock successful. Returns user ok.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "No such field.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "Recieve a user's refresh token in JSON format.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Refresh user's access token",
                "parameters": [
                    {
                        "description": "User's refresh token",
                        "name": "RefreshToken",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_user.RefreshToken"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Authentication successful. Returns a JWT token.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_user.Tokens"
                        }
                    },
                    "400": {
                        "description": "failed to deserialize json request.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials: token is expired - must auth again.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/signin": {
            "post": {
                "description": "Authenticates a user by accepting their login credentials (login and password) in JSON format.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Authenticate user",
                "parameters": [
                    {
                        "description": "User login credentials",
                        "name": "AuthData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.AuthData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Authentication successful. Returns a JWT token.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_user.Tokens"
                        }
                    },
                    "400": {
                        "description": "Invalid input.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Handles the registration of a new user by accepting a JSON payload containing user data.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Any user data for registration",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.PutUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Registration successful. Returns user data.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "Invalid input.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "User already exists.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/todos": {
            "get": {
                "description": "Gets all tasks and returns a JSON containing task data.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todo"
                ],
                "summary": "Get all tasks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "all, completed, or inWork",
                        "name": "filter",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Retrieved successfully. Returns status code OK.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.MetaResponse"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Handles the creation of a new task by accepting a JSON payload containing task data.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todo"
                ],
                "summary": "Create a new task",
                "parameters": [
                    {
                        "description": "Complete task data for creation",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.TodoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Creation successful. Returns task with status code OK.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.Todo"
                        }
                    },
                    "400": {
                        "description": "failed to deserialize json request.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/todos/{id}": {
            "get": {
                "description": "Gets a task by ID in the URL and returns a JSON containing task data.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todo"
                ],
                "summary": "Get task",
                "responses": {
                    "200": {
                        "description": "Retrieved successfully. Returns task and status code OK.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.Todo"
                        }
                    },
                    "400": {
                        "description": "Missing or wrong id.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such task.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Handles the update of a task by accepting a JSON payload containing task data.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todo"
                ],
                "summary": "Update task",
                "parameters": [
                    {
                        "description": "Complete task data for update",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.TodoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Update successful. Returns task with status code OK.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.Todo"
                        }
                    },
                    "400": {
                        "description": "Invalid IsDone field.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such task.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a task by ID in the URL.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todo"
                ],
                "summary": "Delete task",
                "responses": {
                    "200": {
                        "description": "Deletion successful. Returns status code OK.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Missing or wrong id.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such task.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/profile": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves the full profile of the currently authenticated user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user profile",
                "responses": {
                    "200": {
                        "description": "Returns the user profile data.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates the user profile with new data provided in the JSON payload.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user profile",
                "parameters": [
                    {
                        "description": "Updated user data",
                        "name": "Userdata",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Profile successfully updated.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                        }
                    },
                    "400": {
                        "description": "Login or email already used.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No such user.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.Meta": {
            "type": "object",
            "properties": {
                "totalAmount": {
                    "type": "integer"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.MetaResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.Todo"
                    }
                },
                "info": {
                    "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.TodoInfo"
                },
                "meta": {
                    "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.Meta"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.Todo": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "isDone": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.TodoInfo": {
            "type": "object",
            "properties": {
                "all": {
                    "type": "integer"
                },
                "completed": {
                    "type": "integer"
                },
                "inWork": {
                    "type": "integer"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.TodoRequest": {
            "type": "object",
            "properties": {
                "isDone": {
                    "type": "string",
                    "enum": [
                        "true",
                        "false",
                        ""
                    ]
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_userConfig.AuthData": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_userConfig.PutUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "isAdmin": {
                    "type": "boolean"
                },
                "isBlocked": {
                    "type": "boolean"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_userConfig.User": {
            "type": "object",
            "required": [
                "email",
                "login",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "internal_http-server_handlers_admin.UpdateRequest": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "value": {}
            }
        },
        "internal_http-server_handlers_admin.Users": {
            "type": "object",
            "properties": {
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                    }
                }
            }
        },
        "internal_http-server_handlers_user.RefreshToken": {
            "type": "object",
            "properties": {
                "refresh": {
                    "type": "string"
                }
            }
        },
        "internal_http-server_handlers_user.Tokens": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                },
                "refresh": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "JWT Bearer token required for accessing protected routes. Format: Bearer \u003ctoken\u003e",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v0.2.0",
	Host:             "51.250.113.72:8082",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "sAPI",
	Description:      "This is a RESTful API service for EasyDev. It provides various user management functionalities such as user registration, authentication, profile updates, and admin operations.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
