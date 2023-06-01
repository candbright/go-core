package rest

import (
	"github.com/candbright/go-log/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

var ErrFunc = func(err error) {
	log.Error(err)
}

func GET(context *gin.Context, handler func() (interface{}, error)) {
	var (
		err error
	)
	res, err := handler()
	if err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Error(context, err)
		return
	}
	Ok(context, res)
}

func POST[T any](context *gin.Context, handler func(receive T) (interface{}, error)) {
	var (
		err error
	)
	receive := new(T)
	if err = context.ShouldBindJSON(receive); err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Error(context, NewHttpError(err, CodeBindJsonFailed, http.StatusBadRequest))
		return
	}
	if preCheck, ok := reflect.ValueOf(receive).Interface().(Check); ok {
		success, message := preCheck.Check()
		if !success {
			checkErr := errors.New(message)
			if ErrFunc != nil {
				ErrFunc(checkErr)
			}
			Error(context, NewHttpError(checkErr, CodePreCheckFailed, http.StatusBadRequest))
			return
		}
	}
	res, err := handler(*receive)
	if err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Error(context, err)
		return
	}
	Ok(context, res)
}

func PUT[T any](context *gin.Context, handler func(receive T) (interface{}, error)) {
	POST(context, handler)
}

func DELETE(context *gin.Context, handler func() (interface{}, error)) {
	GET(context, handler)
}
