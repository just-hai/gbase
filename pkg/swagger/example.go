package swagger

import (
	"fmt"
	"reflect"
)

type UserGet struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
	UserGet2
	Inte        interface{} `json:"inte"`
}

type UserGet2 struct {
	ID          int         `json:"id"`
	Address     string      `json:"address" desc:"地址"`
	Address2    string      `json:"address2" desc:"地址2"`
	HiddenField string      `json:"-"`           // 这个字段不会显示
	NoJsonTag   string      `desc:"没有json标签的字段"` // 这个字段会使用原字段名
}

// 定义更多的结构体示例
type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

type DeleteResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func RunAdvancedExample() string {
	// 高级示例 - 使用生成器
	generator := NewSwaggerGenerator("User Management API", "2.0.0")

	// 添加创建用户的接口 (POST)
	generator.AddPath(
		"/api/users",
		"POST",
		reflect.TypeOf(UserCreateRequest{}),
		reflect.TypeOf(UserResponse{}),
		"Create a new user",
		"Create a new user with the provided information",
		"Users",
	)


	// 添加获取用户列表的接口 (GET)
	generator.AddPath(
		"/api/users",
		"GET",
		reflect.TypeOf(UserGet{}), // GET 请求通常没有 body
		reflect.TypeOf(UserListResponse{}),
		"Get all users",
		"Retrieve a list of all users",
		"Users",
	)

	// 添加更新用户的接口 (PUT)
	generator.AddPath(
		"/api/users/{id}",
		"PUT",
		reflect.TypeOf(UserCreateRequest{}),
		reflect.TypeOf(UserResponse{}),
		"Update user",
		"Update an existing user with new information",
		"Users",
	)

	// 添加删除用户的接口 (DELETE)
	generator.AddPath(
		"/api/users/{id}",
		"DELETE",
		reflect.TypeOf(struct{}{}), // DELETE 请求通常没有 body
		reflect.TypeOf(DeleteResponse{}),
		"Delete user",
		"Delete an existing user by ID",
		"Users",
	)

	// 生成完整的 Swagger 文档
	advancedJSON, err := generator.ToJSON()
	if err != nil {
		fmt.Printf("Error generating advanced swagger: %v\n", err)
		return ""
	}

	return advancedJSON
}
