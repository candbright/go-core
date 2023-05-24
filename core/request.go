package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var ErrFunc func(err error)

func GET(context *gin.Context, handler func() (interface{}, error)) {
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
		Failed(context, NewResultErr(CodeBindJsonFailed, err, http.StatusBadRequest))
		return
	}
	//if preCheck, ok := reflect.ValueOf(receive).Interface().(IPreCheck); ok {
	//	checkErr := preCheck.PreCheck()
	//	if checkErr != nil && checkErr.Error() != "" {
	//		if ErrFunc != nil {
	//			ErrFunc(err)
	//		}
	//		Failed(context, NewResultErr(CodePreCheckFailed, checkErr, http.StatusBadRequest))
	//		return
	//	}
	//}
	res, err := handler(*receive)
	if err != nil {
		if ErrFunc != nil {
			ErrFunc(err)
		}
		Failed(context, err)
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
