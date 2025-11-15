package constructor

import (
    "github.com/kimvlry/avito-internship-assignment/api"
)

func ErrorResponse(code api.ErrorResponseErrorCode, message string) struct {
    Code    api.ErrorResponseErrorCode `json:"code"`
    Message string                     `json:"message"`
} {
    return struct {
        Code    api.ErrorResponseErrorCode `json:"code"`
        Message string                     `json:"message"`
    }{
        Code:    code,
        Message: message,
    }
}
