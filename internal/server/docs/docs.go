// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/events": {
            "get": {
                "description": "GetEvents",
                "tags": [
                    "Events"
                ],
                "summary": "GetEvents",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "disabled",
                        "name": "disabled",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "code for relation",
                        "name": "code",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "country for relation",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 25,
                        "description": "limit",
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
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-array_github_com_worldline-go_calendar_pkg_models_Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            },
            "post": {
                "description": "AddEvents",
                "tags": [
                    "Events"
                ],
                "summary": "AddEvents",
                "parameters": [
                    {
                        "description": "Event",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_worldline-go_calendar_pkg_models.Event"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-array_string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/events/{id}": {
            "get": {
                "description": "GetEvent",
                "tags": [
                    "Events"
                ],
                "summary": "GetEvent",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Event ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-github_com_worldline-go_calendar_pkg_models_Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            },
            "delete": {
                "description": "RemoveEvent",
                "tags": [
                    "Events"
                ],
                "summary": "RemoveEvent",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Event ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/ics": {
            "get": {
                "description": "GetICS",
                "tags": [
                    "iCal"
                ],
                "summary": "GetICS",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "code for relation",
                        "name": "code",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "country for relation",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "specific year events",
                        "name": "year",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            },
            "post": {
                "description": "AddICS",
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "iCal"
                ],
                "summary": "AddICS",
                "parameters": [
                    {
                        "type": "file",
                        "description": "ICS file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "code for relation",
                        "name": "code",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "country for relation",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "timezone like Europe/Amsterdam",
                        "name": "tz",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/relations": {
            "get": {
                "description": "GetRelations",
                "tags": [
                    "Relations"
                ],
                "summary": "GetRelations",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "code for relation",
                        "name": "code",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "country for relation",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 25,
                        "description": "limit",
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
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-array_github_com_worldline-go_calendar_pkg_models_Relation"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            },
            "post": {
                "description": "AddRelations",
                "tags": [
                    "Relations"
                ],
                "summary": "AddRelations",
                "parameters": [
                    {
                        "description": "Relation",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_worldline-go_calendar_pkg_models.Relation"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-array_string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/relations/{id}": {
            "get": {
                "description": "GetRelation",
                "tags": [
                    "Relations"
                ],
                "summary": "GetRelation",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Relation ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-github_com_worldline-go_calendar_pkg_models_Relation"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            },
            "delete": {
                "description": "RemoveRelation",
                "tags": [
                    "Relations"
                ],
                "summary": "RemoveRelation",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Relation ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/workday": {
            "get": {
                "description": "GetEvents for specific date",
                "tags": [
                    "Search"
                ],
                "summary": "WorkDay",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "code for relation",
                        "name": "code",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "country for relation",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "date specific event",
                        "name": "date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "duration like 1d, 2w, 1h, 1m",
                        "name": "duration",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.Response-array_github_com_worldline-go_calendar_pkg_models_Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_worldline-go_calendar_pkg_models.Event": {
            "type": "object",
            "properties": {
                "all_day": {
                    "type": "boolean"
                },
                "date_from": {
                    "type": "string"
                },
                "date_to": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "disabled": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "rrule": {
                    "type": "string"
                },
                "tz": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                }
            }
        },
        "github_com_worldline-go_calendar_pkg_models.Relation": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "country": {
                    "type": "string"
                },
                "event_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                }
            }
        },
        "rest.Message": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "params": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "rest.Meta": {
            "type": "object",
            "properties": {
                "limit": {
                    "description": "Limit is the limit used within the request.\nIf not defined in the query parameters, this should be the default value used in the service endpoint.",
                    "type": "integer"
                },
                "offset": {
                    "description": "Offset is the offset used within the request.",
                    "type": "integer"
                },
                "total_item_count": {
                    "description": "TotalItemCount is the total number of entities that match the query.",
                    "type": "integer"
                }
            }
        },
        "rest.Response-array_github_com_worldline-go_calendar_pkg_models_Event": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/rest.Message"
                },
                "meta": {
                    "$ref": "#/definitions/rest.Meta"
                },
                "payload": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_worldline-go_calendar_pkg_models.Event"
                    }
                }
            }
        },
        "rest.Response-array_github_com_worldline-go_calendar_pkg_models_Relation": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/rest.Message"
                },
                "meta": {
                    "$ref": "#/definitions/rest.Meta"
                },
                "payload": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_worldline-go_calendar_pkg_models.Relation"
                    }
                }
            }
        },
        "rest.Response-array_string": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/rest.Message"
                },
                "meta": {
                    "$ref": "#/definitions/rest.Meta"
                },
                "payload": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "rest.Response-github_com_worldline-go_calendar_pkg_models_Event": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/rest.Message"
                },
                "meta": {
                    "$ref": "#/definitions/rest.Meta"
                },
                "payload": {
                    "$ref": "#/definitions/github_com_worldline-go_calendar_pkg_models.Event"
                }
            }
        },
        "rest.Response-github_com_worldline-go_calendar_pkg_models_Relation": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/rest.Message"
                },
                "meta": {
                    "$ref": "#/definitions/rest.Meta"
                },
                "payload": {
                    "$ref": "#/definitions/github_com_worldline-go_calendar_pkg_models.Relation"
                }
            }
        },
        "rest.ResponseMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/rest.Message"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/calendar/v1",
	Schemes:          []string{},
	Title:            "calendar API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
