// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Taufik Januar",
            "email": "taufikjanuar35@gmail.com"
        },
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health-check": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health Check"
                ],
                "summary": "Checking Health Services",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helper.BaseResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "helper.BaseResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {},
                "messages": {
                    "type": "string"
                },
                "meta": {
                    "$ref": "#/definitions/helper.Meta"
                }
            }
        },
        "helper.Meta": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "total_data": {
                    "type": "integer"
                },
                "total_page": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/v1",
	Schemes:          []string{"http"},
	Title:            "Singkatin Revamp API",
	Description:      "Revamped URL Shortener API - Shortener Services",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
