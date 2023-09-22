package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

const (
	CodeUnknown        = -1
	CodeSuccess        = 0
	CodeBindJsonFailed = 1
	CodePreCheckFailed = 2
)

type Result struct {
	Err        error       `json:"-"`
	Code       int         `json:"code"`
	HttpStatus int         `json:"-"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
}

func NewResult(code int, data interface{}, err error, status int) Result {
	return Result{
		err,
		code,
		status,
		data,
		"",
	}
}

func Success(c *gin.Context, data interface{}) {
	if data != nil {
		Response(c, NewResult(CodeSuccess, data, nil, http.StatusOK))
	} else {
		Response(c, NewResult(CodeSuccess, nil, nil, http.StatusNoContent))
	}
}

func Failed(c *gin.Context, err error) {
	if err == nil {
		Response(c, NewResult(CodeUnknown, nil, nil, http.StatusInternalServerError))
	} else {
		if resultErr, ok := err.(HttpError); ok {
			Response(c, NewResult(resultErr.Code, nil, resultErr.Err, resultErr.HttpStatus))
		} else {
			Response(c, NewResult(CodeUnknown, nil, err, http.StatusInternalServerError))
		}
	}
}

func Response(c *gin.Context, result Result) {
	if result.Data == nil {
		if result.Err == nil {
			c.AbortWithStatus(result.HttpStatus)
		} else {
			if result.Code == CodeSuccess {
				result.Code = CodeUnknown
			}
			result.Message = errors.Cause(result.Err).Error()
			c.AbortWithStatusJSON(result.HttpStatus, result)
		}
	} else {
		if result.Err != nil {
			result.Message = errors.Cause(result.Err).Error()
		} else {
			c.AbortWithStatusJSON(result.HttpStatus, result)
		}
	}
}
