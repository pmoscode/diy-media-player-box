// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "pmoscode",
            "url": "https://pmoscode.de",
            "email": "info@pmoscode.de"
        },
        "license": {
            "name": "GNU General Public License v3.0",
            "url": "https://gitlab.com/pmoscode/diy-media-player-box/-/raw/master/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/audio-books": {
            "get": {
                "description": "Get all audiobooks stored",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book"
                ],
                "summary": "GetAllAudioBooks",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookFull"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new audiobook",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book"
                ],
                "summary": "AddAudioBook",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookFull"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/audio-books/pause": {
            "post": {
                "description": "Pauses the current playing audio track (if any) - Can be called again to resume playback",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book-debug"
                ],
                "summary": "PauseTrack",
                "responses": {
                    "200": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/audio-books/stop": {
            "post": {
                "description": "Stops the current playing audio track (if any)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book-debug"
                ],
                "summary": "StopTrack",
                "responses": {
                    "200": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/audio-books/{id}": {
            "delete": {
                "description": "Delete an existing audiobook",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book"
                ],
                "summary": "DeleteAudioBook",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id of audio-book",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookFull"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "description": "Update an existing audiobook",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book"
                ],
                "summary": "UpdateAudioBook",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id of audio-book",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "content of audio-book",
                        "name": "audio-book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookUi"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookFull"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/audio-books/{id}/track/{track}/play": {
            "post": {
                "description": "Plays an audio tracks of an existing audiobook",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-book-debug"
                ],
                "summary": "PlayTrack",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id of audio-book",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "track number to be played",
                        "name": "track",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/audio-books/{id}/tracks": {
            "post": {
                "description": "Uploads audio tracks to an existing audiobook",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-track"
                ],
                "summary": "UploadTracks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id of audio-book",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "filename",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "length",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "title",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "name": "track",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookFull"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete all audio tracks of an existing audiobook",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "audio-track"
                ],
                "summary": "DeleteAllTracks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id of audio-book",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.AudioBookFull"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cards/unassigned": {
            "get": {
                "description": "Get all unassigned cards (rfid card ids)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "card"
                ],
                "summary": "GetAllUnassignedCards",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.Card"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "schema.AudioBookFull": {
            "type": "object",
            "properties": {
                "card": {
                    "$ref": "#/definitions/schema.Card"
                },
                "id": {
                    "type": "integer"
                },
                "lastPlayed": {
                    "type": "string"
                },
                "timesPlayed": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "trackList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.AudioTrack"
                    }
                }
            }
        },
        "schema.AudioBookUi": {
            "type": "object",
            "properties": {
                "card": {
                    "$ref": "#/definitions/schema.Card"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "schema.AudioTrack": {
            "type": "object",
            "properties": {
                "filename": {
                    "type": "string"
                },
                "length": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "track": {
                    "type": "integer"
                }
            }
        },
        "schema.Card": {
            "type": "object",
            "properties": {
                "cardId": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:2020",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "DIY Music Box for Children",
	Description:      "This is the controller app for the DIY Music Box for Children.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
