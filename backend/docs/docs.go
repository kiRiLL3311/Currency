// Package docs Code generated for FOREX Rates API. DO NOT EDIT by hand unless regenerating.
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
        "/rates": {
            "get": {
                "security": [{"BearerAuth": []}],
                "description": "Returns all stored USD-based exchange rates, including previous_rate and change_percent when available.",
                "produces": ["application/json"],
                "tags": ["Rates"],
                "summary": "List exchange rates",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/models.Rate"}
                        }
                    },
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}},
                    "500": {"description": "Internal server error", "schema": {"type": "string"}}
                }
            }
        },
        "/convert": {
            "get": {
                "security": [{"BearerAuth": []}],
                "description": "Converts an amount between two currencies using stored USD-based rates. Note: response field rate is the converted total; converted is the unit rate.",
                "produces": ["application/json"],
                "tags": ["Converter"],
                "summary": "Convert currency",
                "parameters": [
                    {"type": "string", "example": "USD", "name": "from", "in": "query", "required": true},
                    {"type": "string", "example": "EUR", "name": "to", "in": "query", "required": true},
                    {"type": "number", "example": 1000, "name": "amount", "in": "query", "required": true}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/models.ConvertResponse"}
                    },
                    "400": {"description": "Invalid amount", "schema": {"type": "string"}},
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}},
                    "500": {"description": "Internal server error", "schema": {"type": "string"}}
                }
            }
        },
        "/rates/sync": {
            "get": {
                "description": "Public endpoint that refreshes rates from Open ER API and previous-day sources.",
                "produces": ["text/plain"],
                "tags": ["Rates"],
                "summary": "Sync rates from upstream",
                "parameters": [
                    {"type": "string", "default": "USD", "example": "USD", "name": "base", "in": "query"}
                ],
                "responses": {
                    "200": {"description": "Rates synchronized successfully", "schema": {"type": "string"}},
                    "500": {"description": "Internal server error", "schema": {"type": "string"}}
                }
            }
        },
        "/health": {
            "get": {
                "description": "Database health check",
                "produces": ["text/plain"],
                "tags": ["Health"],
                "summary": "Health",
                "responses": {
                    "200": {"description": "OK", "schema": {"type": "string"}},
                    "500": {"description": "Database is not connected", "schema": {"type": "string"}}
                }
            }
        }
    },
    "definitions": {
        "models.Rate": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "base_currency": {"type": "string"},
                "target_currency": {"type": "string"},
                "rate": {"type": "number"},
                "previous_rate": {"type": "number"},
                "change_percent": {"type": "number"}
            }
        },
        "models.ConvertResponse": {
            "type": "object",
            "properties": {
                "from": {"type": "string"},
                "to": {"type": "string"},
                "amount": {"type": "number"},
                "rate": {"type": "number", "description": "Converted amount (service field naming)"},
                "converted": {"type": "number", "description": "Unit exchange rate (service field naming)"}
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
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "FOREX Rates API",
	Description:      "Currency rates and converter microservice.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
