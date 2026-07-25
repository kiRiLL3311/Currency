// Package docs Swagger for FOREX News API.
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
        "/news": {
            "get": {
                "security": [{"BearerAuth": []}],
                "description": "Returns stored headlines for a region, optionally filtered by keyword.",
                "produces": ["application/json"],
                "tags": ["News"],
                "summary": "List news articles",
                "parameters": [
                    {
                        "type": "string",
                        "enum": ["US","EU","GB","JP","AU"],
                        "example": "US",
                        "name": "region",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "euro",
                        "name": "keyword",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/models.Article"}
                        }
                    },
                    "401": {"description": "Unauthorized", "schema": {"type": "string"}},
                    "500": {"description": "Internal server error", "schema": {"type": "string"}}
                }
            }
        },
        "/news/sync": {
            "get": {
                "description": "Pulls articles from NewsAPI and/or RSS for a region (and keyword). Public endpoint.",
                "produces": ["text/plain"],
                "tags": ["News"],
                "summary": "Sync news from upstream",
                "parameters": [
                    {
                        "type": "string",
                        "enum": ["US","EU","GB","JP","AU"],
                        "example": "US",
                        "name": "region",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "yen",
                        "name": "keyword",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {"description": "News synchronized successfully", "schema": {"type": "string"}},
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
                    "200": {"description": "News service OK", "schema": {"type": "string"}}
                }
            }
        }
    },
    "definitions": {
        "models.Article": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "title": {"type": "string"},
                "summary": {"type": "string"},
                "url": {"type": "string"},
                "source": {"type": "string"},
                "region": {"type": "string"},
                "currency_tags": {
                    "type": "array",
                    "items": {"type": "string"}
                },
                "published_at": {"type": "string"},
                "fetched_at": {"type": "string"}
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
	Host:             "localhost:8082",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "FOREX News API",
	Description:      "Market news microservice (region + keyword pull).",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
