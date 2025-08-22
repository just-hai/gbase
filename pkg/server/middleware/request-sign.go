package middleware

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/just-hai/gbase/pkg/types"
	"github.com/just-hai/gbase/pkg/utils"
	"github.com/labstack/echo/v4"
)

func CheckSign() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timestamp := c.Request().Header.Get("X-Timestamp")
			sign := c.Request().Header.Get("X-Sign")
			fmt.Println(utils.SHA256(timestamp+getAllParams(c), ""))
			if timestamp == "" || sign == "" || !validateTimestamp(timestamp) ||
				sign != utils.SHA256(timestamp+getAllParams(c), "") {
				return types.SignError
			}
			return next(c)
		}
	}
}

func getAllParams(c echo.Context) string {
	paramMap := make(map[string]string)
	// 获取 URL 查询参数
	queryParams := c.QueryParams()
	for key, values := range queryParams {
		if len(values) > 0 {
			paramMap[key] = values[0] // 只取第一个值
		}
	}
	// 获取表单参数 (application/x-www-form-urlencoded)
	if c.Request().Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		formParams, err := c.FormParams()
		if err == nil {
			for key, values := range formParams {
				if len(values) > 0 {
					// 如果查询参数中没有这个key，或者表单参数优先级更高
					if _, exists := paramMap[key]; !exists {
						paramMap[key] = values[0] // 只取第一个值
					}
				}
			}
		}
	}

	// 按参数名排序
	var keys []string
	for key := range paramMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	// 构建排序后的参数列表
	var params string
	for _, key := range keys {
		params += "&" + key + "=" + paramMap[key]
	}
	params = strings.TrimPrefix(params, "&")
	return params
}

func validateTimestamp(timestamp string) bool {
	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	// 5分钟过期
	if time.Now().Add(-5 * time.Minute).Before(time.UnixMilli(t)) {
		return false
	}
	v := t / 1000
	v2 := t % 1000000 / 1000
	if v2 == 0 {
		v2 = 1
	}
	v3 := v % v2
	return v*1000+v3 == t
}
