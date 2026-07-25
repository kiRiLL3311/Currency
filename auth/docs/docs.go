// Package docs Swagger for FOREX Auth API.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/register": {
            "post": {
                "description": "Creates a new account. Does not return tokens — sign in afterwards.",
                "consumes": ["application/json"],
                "produces": ["text/plain"],
                "tags": ["Auth"],
                "summary": "Register a user",
                "parameters": [{
                    "description": "Registration payload",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/models.RegisterRequest"}
                }],
                "responses": {
                    "201": {"description": "User created", "schema": {"type": "string"}},
                    "400": {"description": "Invalid JSON", "schema": {"type": "string"}},
                    "409": {"description": "Username or email already exists", "schema": {"type": "string"}},
                    "500": {"description": "Internal Server Error", "schema": {"type": "string"}}
                }
            }
        },
        "/login": {
            "post": {
                "description": "Returns access and refresh tokens for an existing user.",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Login",
                "parameters": [{
                    "description": "Login payload",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/models.LoginRequest"}
                }],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.AuthResponse"}},
                    "400": {"description": "Invalid JSON", "schema": {"type": "string"}},
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}}
                }
            }
        },
        "/refresh": {
            "post": {
                "description": "Rotates refresh token and returns a new access + refresh pair.",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Refresh tokens",
                "parameters": [{
                    "description": "Refresh payload",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/models.RefreshRequest"}
                }],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.AuthResponse"}},
                    "400": {"description": "Invalid JSON", "schema": {"type": "string"}},
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}}
                }
            }
        },
        "/me": {
            "get": {
                "security": [{"BearerAuth": []}],
                "description": "Returns the authenticated user profile.",
                "produces": ["application/json"],
                "tags": ["Profile"],
                "summary": "Current user",
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.User"}},
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}},
                    "500": {"description": "Internal server error", "schema": {"type": "string"}}
                }
            },
            "patch": {
                "security": [{"BearerAuth": []}],
                "description": "Updates the user's news region (US, EU, GB, JP, AU).",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Profile"],
                "summary": "Update profile region",
                "parameters": [{
                    "description": "Profile update",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/models.UpdateProfileRequest"}
                }],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.User"}},
                    "400": {"description": "Invalid JSON or region", "schema": {"type": "string"}},
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}},
                    "500": {"description": "Internal server error", "schema": {"type": "string"}}
                }
            }
        },
        "/health": {
            "get": {
                "produces": ["text/plain"],
                "tags": ["Health"],
                "summary": "Health",
                "responses": {
                    "200": {"description": "Auth service OK", "schema": {"type": "string"}}
                }
            }
        }
    },
    "definitions": {
        "models.RegisterRequest": {
            "type": "object",
            "properties": {
                "username": {"type": "string"},
                "email": {"type": "string"},
                "password": {"type": "string"}
            }
        },
        "models.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {"type": "string"},
                "password": {"type": "string"}
            }
        },
        "models.RefreshRequest": {
            "type": "object",
            "properties": {
                "refresh_token": {"type": "string"}
            }
        },
        "models.UpdateProfileRequest": {
            "type": "object",
            "properties": {
                "region": {"type": "string", "example": "US"}
            }
        },
        "models.AuthResponse": {
            "type": "object",
            "properties": {
                "access_token": {"type": "string"},
                "refresh_token": {"type": "string"}
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "username": {"type": "string"},
                "email": {"type": "string"},
                "region": {"type": "string"},
                "created_at": {"type": "string"}
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "JWT access token. Example: Bearer \u003ctoken\u003e"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "FOREX Auth API",
	Description:      "Authentication and user profile microservice.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
