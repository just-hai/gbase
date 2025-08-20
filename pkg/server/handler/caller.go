package handler

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

func CallHandler(handler interface{}, inParams map[reflect.Type]reflect.Value, c echo.Context) ([]reflect.Value, error) {
	handlerRef := reflect.ValueOf(handler)
	var params []reflect.Value
	for i := 0; i < handlerRef.Type().NumIn(); i++ {
		paramType := handlerRef.Type().In(i)
		v, ok := inParams[paramType]
		var err error
		if !ok {
			v, err = newType(paramType, c)
			if err != nil {
				return nil, err
			}
			inParams[paramType] = v
		}
		params = append(params, v)
	}
	values := handlerRef.Call(params)

	var notErrors []reflect.Value
	var err error
	for _, value := range values {
		if value.Interface() != nil {
			if !isErrorType(value) {
				inParams[value.Type()] = value
				notErrors = append(notErrors, value)
			} else {
				err = value.Interface().(error)
			}
		}
	}
	return notErrors, err
}

func isErrorType(v reflect.Value) bool {
	return v.MethodByName("Error").IsValid()
}
