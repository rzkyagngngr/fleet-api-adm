package docs

import (
	"github.com/swaggo/swag"
)

func init() {
	swag.Register(swag.Name, &s{})
}

type s struct{}

func (s *s) ReadDoc() string {
	return `{
    "swagger": "2.0",
    "info": {
        "description": "API for Fleet Administration Backend.",
        "title": "Gin Boilerplate API",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {
        "/auth/register": {
            "post": {
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Register payload",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UserRegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Login payload",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.LoginRequest": {
            "type": "object",
            "required": [
                "employee_id",
                "password"
            ],
            "properties": {
                "employee_id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dto.UserRegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "employee_id",
                "full_name",
                "password"
            ],
            "properties": {
                "branch_code": {
                    "type": "integer"
                },
                "branch_name": {
                    "type": "string"
                },
                "company_code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "employee_id": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "job_title": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "personnel_area": {
                    "type": "string"
                },
                "personnel_sub_area": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                },
                "sub_unit_name": {
                    "type": "string"
                },
                "terminal_code": {
                    "type": "integer"
                }
            }
        },
        "utils.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}
`
}
