// Code generated by swaggo/swag. DO NOT EDIT
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
        "/auth/login": {
            "post": {
                "description": "Login user with given credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Data to login user",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.AuthInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login success",
                        "schema": {
                            "$ref": "#/definitions/dtos.RegisterDetails"
                        }
                    },
                    "400": {
                        "description": "Error logging in",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register new user with given credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Data to register user",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.AuthInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Register success",
                        "schema": {
                            "$ref": "#/definitions/dtos.RegisterDetails"
                        }
                    },
                    "400": {
                        "description": "Failed to register user",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/auth/whoami": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get current user id and name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authenticated user details",
                "responses": {
                    "200": {
                        "description": "User details",
                        "schema": {
                            "$ref": "#/definitions/dtos.WhoAmI"
                        }
                    },
                    "400": {
                        "description": "Error getting user data",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament": {
            "get": {
                "description": "Get tournament details by its id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Tournament details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tournament id",
                        "name": "tournamentId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament",
                        "schema": {
                            "$ref": "#/definitions/models.Tournament"
                        }
                    },
                    "400": {
                        "description": "Tournament not found",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/contests": {
            "get": {
                "description": "Get tournament contests",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Tournament contests",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tournament id",
                        "name": "tournamentId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Contest bracket",
                        "schema": {
                            "$ref": "#/definitions/dtos.Contest"
                        }
                    },
                    "400": {
                        "description": "Failed to return tournament contests",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new tournament for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Create new tournament",
                "parameters": [
                    {
                        "description": "Data to create tournament",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.CreateTournament"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament created",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournament creation",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete tournaments for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Delete tournaments",
                "parameters": [
                    {
                        "description": "Data to delete tournaments",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.TournamentIds"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournaments deleted",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournaments deletion",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/edit": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Edit tournament for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Edit tournament",
                "parameters": [
                    {
                        "description": "Data to edit tournament",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.EditTournament"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament edited",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournament edition",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/tiktoks": {
            "get": {
                "description": "Get tournament tiktoks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Tournament tiktoks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tournament id",
                        "name": "tournamentId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament tiktoks",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Tiktok"
                            }
                        }
                    },
                    "400": {
                        "description": "Tournament not found",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/user/photo": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Change user photo for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Change user photo",
                "parameters": [
                    {
                        "description": "Data to change photo",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.ChangePhotoURL"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Photo edited",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during photo change",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/user/tournaments": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get tournaments for user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get tournaments for user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "page size",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sort page by name",
                        "name": "sort_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sort page by size",
                        "name": "sort_size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournaments of user",
                        "schema": {
                            "$ref": "#/definitions/dtos.TournamentsResponse"
                        }
                    },
                    "400": {
                        "description": "Couldn't get tournaments for specific user",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/winner/tournament": {
            "put": {
                "description": "Increment wins and increment times_played",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Update tournament winner statistics",
                "parameters": [
                    {
                        "description": "Data to update tournament winner",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.TournamentWinner"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Winner updated",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during winner updating",
                        "schema": {
                            "$ref": "#/definitions/dtos.MessageResponseType"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.AuthInput": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dtos.ChangePhotoURL": {
            "type": "object",
            "required": [
                "photoURL"
            ],
            "properties": {
                "photoURL": {
                    "type": "string"
                }
            }
        },
        "dtos.Contest": {
            "type": "object",
            "properties": {
                "countMatches": {
                    "type": "integer"
                },
                "rounds": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.Round"
                    }
                }
            }
        },
        "dtos.CreateTiktok": {
            "type": "object",
            "required": [
                "name",
                "url"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "dtos.CreateTournament": {
            "type": "object",
            "required": [
                "name",
                "photoURL",
                "tiktoks"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                },
                "size": {
                    "type": "integer",
                    "maximum": 64,
                    "minimum": 4
                },
                "tiktoks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CreateTiktok"
                    }
                }
            }
        },
        "dtos.EditTournament": {
            "type": "object",
            "required": [
                "name",
                "photoURL",
                "tiktoks"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                },
                "size": {
                    "type": "integer",
                    "maximum": 64,
                    "minimum": 4
                },
                "tiktoks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CreateTiktok"
                    }
                }
            }
        },
        "dtos.Match": {
            "type": "object",
            "properties": {
                "firstOption": {},
                "matchID": {
                    "type": "string"
                },
                "secondOption": {}
            }
        },
        "dtos.MessageResponseType": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "dtos.RegisterDetails": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dtos.Round": {
            "type": "object",
            "properties": {
                "matches": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.Match"
                    }
                },
                "round": {
                    "type": "integer"
                }
            }
        },
        "dtos.TournamentIds": {
            "type": "object",
            "required": [
                "tournamentIds"
            ],
            "properties": {
                "tournamentIds": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "dtos.TournamentWinner": {
            "type": "object",
            "required": [
                "tiktokURL"
            ],
            "properties": {
                "tiktokURL": {
                    "type": "string"
                }
            }
        },
        "dtos.TournamentsResponse": {
            "type": "object",
            "required": [
                "tournamentCount",
                "tournaments"
            ],
            "properties": {
                "tournamentCount": {
                    "type": "integer"
                },
                "tournaments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Tournament"
                    }
                }
            }
        },
        "dtos.WhoAmI": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Tiktok": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "tournament": {
                    "$ref": "#/definitions/models.Tournament"
                },
                "tournamentID": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "wins": {
                    "type": "integer"
                }
            }
        },
        "models.Tournament": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "timesPlayed": {
                    "type": "integer"
                },
                "user": {
                    "$ref": "#/definitions/models.User"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
