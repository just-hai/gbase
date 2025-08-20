package handler

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

var (
	validate = validator.New()
	trans    ut.Translator
)

func init() {
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	zhTrans.RegisterDefaultTranslations(validate, trans)
}

func newType(typ reflect.Type, c echo.Context) (reflect.Value, error) {
	requestType := typ
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	requestObj := reflect.New(requestType)
	pathAndQueryParams := c.QueryParams()
	for _, name := range c.ParamNames() {
		value := c.Param(name)
		for _, maybeName := range strings.Split(name, ",") {
			pathAndQueryParams[maybeName] = []string{value}
		}
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(requestObj.Interface(), pathAndQueryParams)
	if err != nil {
		return requestObj, err
	}
	for i := 0; i < requestType.NumField(); i++ {
		field := requestType.Field(i)

		if field.Name == "Body" || field.Anonymous {
			theType := field.Type
			var value interface{}
			if theType.Kind() == reflect.Ptr {
				value = reflect.New(field.Type.Elem()).Interface()
			} else {
				value = reflect.New(field.Type).Interface()
			}

			if field.Name == "Body" {
				buf, err := io.ReadAll(c.Request().Body)
				if err != nil {
					return requestObj, err
				}
				c.Request().Body = io.NopCloser(bytes.NewBuffer(buf))
				if err = c.Bind(value); err != nil {
					return requestObj, err
				}
				c.Request().Body = io.NopCloser(bytes.NewBuffer(buf)) // for next handler
			} else {
				err = decoder.Decode(value, pathAndQueryParams)
			}
			if err != nil {
				return requestObj, err
			}

			targetField := requestObj.Elem().FieldByName(field.Name)
			if targetField.CanSet() {
				if theType.Kind() == reflect.Ptr {
					targetField.Set(reflect.ValueOf(value))
				} else {
					targetField.Set(reflect.ValueOf(value).Elem())
				}
			}
		}
	}

	if err = validate.Struct(requestObj.Interface()); err != nil {
		return requestObj, translateValidationErr(err)
	}
	return requestObj, nil
}

func translateValidationErr(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, v := range errs.Translate(trans) {
			return errors.New(v)
		}
	}
	return err
}
