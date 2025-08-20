package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"embed"
)

type (
	// Swagger v3.0 结构体定义
	SwaggerSpec struct {
		OpenAPI    string              `json:"openapi"`
		Info       Info                `json:"info"`
		Paths      map[string]PathItem `json:"paths"`
		Components *Components         `json:"components,omitempty"`
	}

	Info struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	}

	PathItem struct {
		Get    *Operation `json:"get,omitempty"`
		Post   *Operation `json:"post,omitempty"`
		Put    *Operation `json:"put,omitempty"`
		Delete *Operation `json:"delete,omitempty"`
	}

	Operation struct {
		Tags        []string            `json:"tags,omitempty"`
		Summary     string              `json:"summary"`
		Description string              `json:"description,omitempty"`
		RequestBody *RequestBody        `json:"requestBody,omitempty"`
		Responses   map[string]Response `json:"responses"`
	}

	RequestBody struct {
		Required bool                 `json:"required"`
		Content  map[string]MediaType `json:"content"`
	}

	Response struct {
		Description string               `json:"description"`
		Content     map[string]MediaType `json:"content,omitempty"`
	}

	MediaType struct {
		Schema Schema `json:"schema"`
	}

	Schema struct {
		Type        string            `json:"type,omitempty"`
		Properties  map[string]Schema `json:"properties,omitempty"`
		Required    []string          `json:"required,omitempty"`
		Ref         string            `json:"$ref,omitempty"`
		Description string            `json:"description,omitempty"`
		Items       *Schema           `json:"items,omitempty"`
	}

	Components struct {
		Schemas map[string]Schema `json:"schemas"`
	}

	// SwaggerGenerator 是一个更通用的 Swagger 生成器
	SwaggerGenerator struct {
		spec *SwaggerSpec
	}
)

//go:embed ui
var UI embed.FS

// NewSwaggerGenerator 创建新的生成器
func NewSwaggerGenerator(title, version string) *SwaggerGenerator {
	return &SwaggerGenerator{
		spec: &SwaggerSpec{
			OpenAPI: "3.0.0",
			Info: Info{
				Title:   title,
				Version: version,
			},
			Paths: make(map[string]PathItem),
			Components: &Components{
				Schemas: make(map[string]Schema),
			},
		},
	}
}

func (sg *SwaggerGenerator)SetInfo(title, version string) {
	sg.spec.Info.Title = title
	sg.spec.Info.Version = version
}

// AddPath 添加一个新的路径
func (sg *SwaggerGenerator) AddPath(path, method string, inType, outType reflect.Type, summary, description string, tags ...string) {
	// 添加到 components schemas 并获取实际的 schema 名称
	inSchemaName := sg.addSchemaRecursively(inType, false)
	outSchema := sg.generateResponseSchema(outType)

	// 创建 operation
	operation := &Operation{
		Tags:        tags,
		Summary:     summary,
		Description: description,
		Responses: map[string]Response{
			"200": {
				Description: "Successful response",
				Content: map[string]MediaType{
					"application/json": {
						Schema: outSchema,
					},
				},
			},
		},
	}

	// 所有方法都添加 requestBody（如果 inType 不是空结构体）
	if inSchemaName != "" {
		operation.RequestBody = &RequestBody{
			Required: true,
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{
						Ref: fmt.Sprintf("#/components/schemas/%s", inSchemaName),
					},
				},
			},
		}
	}

	// 获取或创建 PathItem
	pathItem, exists := sg.spec.Paths[path]
	if !exists {
		pathItem = PathItem{}
	}

	// 根据方法设置 operation
	switch strings.ToUpper(method) {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	}

	sg.spec.Paths[path] = pathItem
}

// generateResponseSchema 生成响应的 schema，对于数组类型直接生成内联 schema
func (sg *SwaggerGenerator) generateResponseSchema(t reflect.Type) Schema {

	schemaType := "object"
	// 处理数组/切片类型
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
		schemaType = "array"
	}

	// 处理普通类型
	schemaName := sg.addSchemaRecursively(t, true)
	if schemaName == "" {
		return Schema{Type: schemaType}
	}

	return Schema{
		Type: schemaType,
		Ref: fmt.Sprintf("#/components/schemas/%s", schemaName),
	}
}

// addSchemaRecursively 递归地添加结构体及其嵌套结构体到 schemas
func (sg *SwaggerGenerator) addSchemaRecursively(t reflect.Type, needAnonymous bool) string {
	schemaName := t.Name()
	if schemaName == "" {
		if needAnonymous {
			// 对于匿名结构体，生成一个唯一的名称
			schemaName = sg.generateAnonymousSchemaName(t)
		} else {
			return ""
		}
	}

	// 如果已经添加过，就不重复添加
	if _, exists := sg.spec.Components.Schemas[schemaName]; exists {
		return schemaName
	}

	// 生成并添加当前结构体的 schema
	schema := sg.generateSchema(t)
	sg.spec.Components.Schemas[schemaName] = schema

	return schemaName
}

// generateAnonymousSchemaName 为匿名结构体生成唯一名称
func (sg *SwaggerGenerator) generateAnonymousSchemaName(t reflect.Type) string {
	// 基于结构体的字段生成一个唯一的名称
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		var fieldName string
		if jsonTag == "" {
			fieldName = field.Name
		} else {
			fieldName = strings.Split(jsonTag, ",")[0]
		}
		fieldNames = append(fieldNames, fieldName)
	}

	// 生成基础名称
	baseName := "AnonymousStruct"
	if len(fieldNames) > 0 {
		baseName = strings.Join(fieldNames, "") + "Struct"
	}else{
		return baseName  // 如果没有字段，则返回默认名称, 这个作为特殊
	}

	// 确保名称唯一
	counter := 1
	finalName := baseName
	for {
		if _, exists := sg.spec.Components.Schemas[finalName]; !exists {
			break
		}
		finalName = fmt.Sprintf("%s%d", baseName, counter)
		counter++
	}

	return finalName
}

func (sg *SwaggerGenerator) generateSchema(t reflect.Type) Schema {
	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
		Required:   []string{},
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 处理嵌入结构体（匿名字段）
		if field.Anonymous {
			// 递归处理嵌入结构体的字段
			embeddedSchema := sg.generateSchema(field.Type)
			// 将嵌入结构体的属性合并到当前 schema 中
			for propName, propSchema := range embeddedSchema.Properties {
				schema.Properties[propName] = propSchema
			}
			// 合并必需字段
			schema.Required = append(schema.Required, embeddedSchema.Required...)
			continue
		}

		jsonTag := field.Tag.Get("json")

		// 如果 json 标签是 "-"，跳过该字段
		if jsonTag == "-" {
			continue
		}

		var fieldName string
		if jsonTag == "" {
			// 没有 json 标签，使用原字段名
			fieldName = field.Name
		} else {
			// 解析 json tag，去掉 omitempty 等选项
			fieldName = strings.Split(jsonTag, ",")[0]
		}

		var fieldSchema Schema
		switch field.Type.Kind() {
		case reflect.String:
			fieldSchema = Schema{Type: "string"}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldSchema = Schema{Type: "integer"}
		case reflect.Float32, reflect.Float64:
			fieldSchema = Schema{Type: "number"}
		case reflect.Bool:
			fieldSchema = Schema{Type: "boolean"}
		case reflect.Slice, reflect.Array:
			// 处理数组/切片类型
			elemType := field.Type.Elem()
			switch elemType.Kind() {
			case reflect.String:
				fieldSchema = Schema{
					Type:  "array",
					Items: &Schema{Type: "string"},
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldSchema = Schema{
					Type:  "array",
					Items: &Schema{Type: "integer"},
				}
			case reflect.Float32, reflect.Float64:
				fieldSchema = Schema{
					Type:  "array",
					Items: &Schema{Type: "number"},
				}
			case reflect.Bool:
				fieldSchema = Schema{
					Type:  "array",
					Items: &Schema{Type: "boolean"},
				}
			case reflect.Struct:
				// 处理结构体数组
				elemSchemaName := sg.addSchemaRecursively(elemType, true)
				fieldSchema = Schema{
					Type: "array",
					Items: &Schema{
						Ref: fmt.Sprintf("#/components/schemas/%s", elemSchemaName),
					},
				}
			default:
				fieldSchema = Schema{
					Type:  "array",
					Items: &Schema{Type: "object"},
				}
			}
		case reflect.Struct:
			// 处理嵌套结构体，生成引用
			nestedSchemaName := sg.addSchemaRecursively(field.Type, true)
			fieldSchema = Schema{
				Ref: fmt.Sprintf("#/components/schemas/%s", nestedSchemaName),
			}
		default:
			fieldSchema = Schema{Type: "object"}
		}

		// 检查是否有 desc tag，如果有则添加描述
		if descTag := field.Tag.Get("desc"); descTag != "" {
			fieldSchema.Description = descTag
		}

		schema.Properties[fieldName] = fieldSchema

		if strings.Contains(field.Tag.Get("validate"), "required") {
			schema.Required = append(schema.Required, fieldName)
		}
	}

	return schema
}

// ToJSON 生成 JSON 格式的 Swagger 文档
func (sg *SwaggerGenerator) ToJSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(sg.spec, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
