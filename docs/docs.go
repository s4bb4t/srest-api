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
                        "description": "Search term",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "order asc or desc",
                        "name": "order",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "block status",
                        "name": "blocked",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit of users for query",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Retrieve successful. Returns users.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_admin.GetAllResponse"
                        }
                    },
                    "401": {
                        "description": "Retrieve failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}": {
            "get": {
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
                            "$ref": "#/definitions/internal_http-server_handlers_admin.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Retrieve failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            },
            "put": {
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
                        "description": "Complete user data",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Update successful. Returns user ok.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_admin.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Update failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            },
            "delete": {
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
                        "description": "Remove successful. Returns user ok.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    },
                    "401": {
                        "description": "Remove failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}/block": {
            "post": {
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
                            "$ref": "#/definitions/internal_http-server_handlers_admin.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Block failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
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
                        "description": "Update successful. Returns user ok.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    },
                    "401": {
                        "description": "Update failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        },
        "/admin/users/{id}/unlock": {
            "post": {
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
                            "$ref": "#/definitions/internal_http-server_handlers_admin.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Unlock failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        },
        "/signin": {
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
                            "$ref": "#/definitions/internal_http-server_handlers_user.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input. Returns error message for improper data structure.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        },
        "/signup": {
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
                        "description": "Complete user data for registration",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Registration successful. Returns user data.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_user.GetResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input. Returns error message for improper data structure.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
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
                        "description": "all, completed, or inwork",
                        "name": "filter",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Retrieved successfully. Returns status code OK.",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_todo.GetAllResponse"
                        }
                    },
                    "401": {
                        "description": "Retrieving failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
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
                            "$ref": "#/definitions/internal_http-server_handlers_todo.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Creation failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
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
                            "$ref": "#/definitions/internal_http-server_handlers_todo.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Retrieving failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
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
                            "$ref": "#/definitions/internal_http-server_handlers_todo.GetResponse"
                        }
                    },
                    "401": {
                        "description": "Update failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
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
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    },
                    "401": {
                        "description": "Deletion failed. Returns error message.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        },
        "/user/profile": {
            "get": {
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
                            "$ref": "#/definitions/internal_http-server_handlers_user.GetResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input. Returns error message for improper data structure.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            },
            "put": {
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
                            "$ref": "#/definitions/internal_http-server_handlers_user.GetResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input. Returns error message for improper data structure.",
                        "schema": {
                            "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_sabbatD_srest-api_internal_lib_api_response.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.Meta": {
            "type": "object",
            "properties": {
                "total_amount": {
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
                "isdone": {
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
                "inwork": {
                    "type": "integer"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_todoConfig.TodoRequest": {
            "type": "object",
            "properties": {
                "isdone": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_userConfig.AuthData": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser": {
            "type": "object",
            "properties": {
                "admin": {
                    "type": "boolean"
                },
                "block": {
                    "type": "boolean"
                },
                "date": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
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
        "github_com_sabbatD_srest-api_internal_lib_userConfig.User": {
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
        "internal_http-server_handlers_admin.GetAllResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                    }
                }
            }
        },
        "internal_http-server_handlers_admin.GetResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
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
        "internal_http-server_handlers_todo.GetAllResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "metaresponse": {
                    "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.MetaResponse"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "internal_http-server_handlers_todo.GetResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "todo": {
                    "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_todoConfig.Todo"
                }
            }
        },
        "internal_http-server_handlers_user.AuthResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "internal_http-server_handlers_user.GetResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/github_com_sabbatD_srest-api_internal_lib_userConfig.TableUser"
                }
            }
        }
    },
    "externalDocs": {
        "description": "Readme.md from github",
        "url": "https://github.com/sabbatD/srest-api"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v0.1.1",
	Host:             "51.250.113.72:8082/api/v1",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "sAPI",
	Description:      "This is RESTful-API service for EasyDev.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
