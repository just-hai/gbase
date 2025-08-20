package handler

import (
	"net/http"
	"reflect"

	"github.com/just-hai/gbase/pkg/types"
	"github.com/labstack/echo/v4"
)

type RespObj struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,string"`
}

func RespJson(ctx *echo.Context, outs []reflect.Value, err error) error {
	var rest RespObj
	if len(outs) > 0 {
		rest.Data = outs[0].Interface()
	}
	if err != nil {
		if bizError, ok := err.(*types.BizError); ok {
			rest.Errno = bizError.Code()
			rest.Msg = bizError.Error()
		} else {
			rest.Errno = http.StatusInternalServerError
			rest.Msg = err.Error()
		}
	}
	return (*ctx).JSON(http.StatusOK, rest)
}

var HTTPErrorHandler = func(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}
	ctx.JSON(500, RespObj{Errno: 500, Msg: err.Error()})
}
