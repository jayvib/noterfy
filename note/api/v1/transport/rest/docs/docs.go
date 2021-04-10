// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Jayson Vibandor",
            "email": "jayson.vibandor@gmail.com"
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
        "/note": {
            "post": {
                "description": "Creating a new note. The client can assign the note ID with a UUID value but the service will return a conflict error when the note with the ID provided is already exists.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new note.",
                "parameters": [
                    {
                        "description": "A body containing the new note",
                        "name": "CreateRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/rest.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully created a new note",
                        "schema": {
                            "$ref": "#/definitions/rest.CreateResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict error due to the new note with an ID already exists in the service",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "499": {
                        "description": "Cancel error when the request was aborted",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        },
        "/note/{id}": {
            "delete": {
                "description": "Deletes an existing note.",
                "summary": "Deletes an existing note.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the note",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful deleting a note",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Note's ID parameter is not provided in the path",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "499": {
                        "description": "Cancel error when the request was aborted",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "note.Note": {
            "type": "object",
            "properties": {
                "content": {
                    "description": "Content is the content of the note",
                    "type": "string"
                },
                "created_time": {
                    "description": "CreatedTime is the timestamp when the note was created.",
                    "type": "string"
                },
                "id": {
                    "description": "ID is a unique identifier UUID of the note.",
                    "type": "string"
                },
                "is_favorite": {
                    "description": "IsFavorite is a flag when then the note is marked as favorite",
                    "type": "boolean"
                },
                "title": {
                    "description": "Title is the title of the note",
                    "type": "string"
                },
                "updated_time": {
                    "description": "UpdateTime is the timestamp when the note last updated.",
                    "type": "string"
                }
            }
        },
        "rest.CreateRequest": {
            "type": "object",
            "properties": {
                "note": {
                    "$ref": "#/definitions/note.Note"
                }
            }
        },
        "rest.CreateResponse": {
            "type": "object",
            "properties": {
                "note": {
                    "$ref": "#/definitions/note.Note"
                }
            }
        },
        "rest.ResponseError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.2.1",
	Host:        "localhost:8080",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "Noterfy Note Service",
	Description: "Noterfy Note Service.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}