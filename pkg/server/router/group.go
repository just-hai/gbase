package router

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/server/handler"
	"github.com/just-hai/gbase/pkg/swagger"
	"github.com/labstack/echo/v4"
)

type HandlerFunc any

type RouterGroup struct {
	EchoGroup *echo.Group
	Tag       string
}

func (r *RouterGroup) Group(path, tag string) *RouterGroup {
	return &RouterGroup{
		EchoGroup: r.EchoGroup.Group(path),
		Tag:       tag,
	}
}

func (r *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	route := r.EchoGroup.GET(relativePath, Handles(handlers...))
	buildSwagger(route.Path, http.MethodGet, r.Tag, handlers...)
}

func (r *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	route := r.EchoGroup.PUT(relativePath, Handles(handlers...))
	buildSwagger(route.Path, http.MethodPut, r.Tag, handlers)
}

func (r *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	route := r.EchoGroup.POST(relativePath, Handles(handlers...))
	buildSwagger(route.Path, http.MethodPost, r.Tag, handlers)
}

func (r *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	route := r.EchoGroup.DELETE(relativePath, Handles(handlers...))
	buildSwagger(route.Path, http.MethodDelete, r.Tag, handlers)
}

func Handles(handlers ...HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		inParams := make(map[reflect.Type]reflect.Value)
		inParams[reflect.TypeOf((*echo.Context)(nil)).Elem()] = reflect.ValueOf(ctx)

		var outs []reflect.Value
		var err error
		for _, h := range handlers {
			if reflect.TypeOf(h).Kind() != reflect.Func {
				continue
			}

			outs, err = handler.CallHandler(h, inParams, ctx)
			if err != nil {
				return handler.RespJson(&ctx, nil, err)
			}
			if len(outs) > 1 {
				return handler.RespJson(&ctx, nil, fmt.Errorf("too many return values"))
			}
		}
		return handler.RespJson(&ctx, outs, nil)
	}
}

var swaggerGenerator *swagger.SwaggerGenerator

// GetSwaggerJSON 返回 Swagger JSON 文档
func GetSwaggerJSON() (string, error) {
	if swaggerGenerator == nil {
		return "", nil
	}
	return swaggerGenerator.ToJSON()
}

func buildSwagger(path, method, tag string, handlers ...HandlerFunc) {

	if !config.Conf.Swagger {
		return
	}

	if swaggerGenerator == nil {
		swaggerGenerator = swagger.NewSwaggerGenerator(config.Conf.AppName, config.Conf.AppVersion)
	}

	var summary string
	var handle reflect.Value
	for _, h := range handlers {
		if reflect.TypeOf(h).Kind() == reflect.String {
			summary = h.(string)
		} else if reflect.TypeOf(h).Kind() == reflect.Func {
			handle = reflect.ValueOf(h)
		}
	}
	if !handle.IsValid() {
		return
	}

	var inType reflect.Type
	for i := range handle.Type().NumIn() {
		inType = handle.Type().In(i)
		if inType.Kind() == reflect.Ptr {
			inType = inType.Elem()
		}
	}

	var outType reflect.Type
	if handle.Type().NumOut() > 0 {
		if handle.Type().Out(0) == reflect.TypeOf((*error)(nil)).Elem() {
			outType = reflect.TypeOf(struct{}{})
		} else {
			outType = handle.Type().Out(0)
			if outType.Kind() == reflect.Ptr {
				outType = outType.Elem()
			}
		}
	}

	swaggerGenerator.AddPath(
		path,
		method,
		inType,
		outType,
		summary,
		"",
		tag,
	)
}
