package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

var ErrFunc func(err error)

func SimpleReq(context *gin.Context, handler func() (interface{}, error)) {
	var (
		err error
	)
	res, err := handler()
	if err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Failed(context, err)
		return
	}
	Success(context, res)
}

func GenericReq[T any](context *gin.Context, handler func(receive T) (interface{}, error)) {
	var (
		err error
	)
	receive := new(T)
	if err = context.ShouldBindJSON(receive); err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Failed(context, NewHttpError(err, CodeBindJsonFailed, http.StatusBadRequest))
		return
	}
	if preCheck, ok := reflect.ValueOf(receive).Interface().(Check); ok {
		success, message := preCheck.Check()
		if !success {
			checkErr := errors.New(message)
			if ErrFunc != nil {
				ErrFunc(checkErr)
			}
			Failed(context, NewHttpError(checkErr, CodePreCheckFailed, http.StatusBadRequest))
			return
		}
	}
	res, err := handler(*receive)
	if err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Failed(context, err)
		return
	}
	Success(context, res)
}
