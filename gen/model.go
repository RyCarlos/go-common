package gen

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"reflect"
	"strings"
)

type Model struct {
	Comment   string // 模型注释
	Name      string // 模型名称
	Fields    []*Field
	HasDelete bool
}

type Field struct {
	Name    string
	Type    string
	Comment string
	NotNull bool
}

var ignoreFields = []string{"Id", "CreateTime", "UpdateTime", "SoftDelete"}

func (f *Field) IsShow() bool {
	return !slice.Contain(ignoreFields, f.Name)
}

func (f *Field) GetNameLowerCamelCase() string {
	return strutil.CamelCase(f.Name)
}

// GetNameUpperCamelCase 小驼峰 foobarTest
func (m *Model) GetNameUpperCamelCase() string {
	return strutil.CamelCase(m.Name)
}

// GetNameLowerCamelCase 大驼峰 FoobarTest
func (m *Model) GetNameLowerCamelCase() string {
	return strutil.CamelCase(m.Name)
}

// GetNameSnakeCase foobar_test
func (m *Model) GetNameSnakeCase() string {
	return strutil.SnakeCase(m.Name)
}

func NewModel(m interface{}) (instance *Model) {
	instance = &Model{
		Name:   "",
		Fields: nil,
	}
	// 类型反射
	t := reflect.TypeOf(m)
	elem := t.Elem()
	// 循环模型属性
	fields := make([]*Field, 0)
	for i := 0; i < elem.NumField(); i++ {
		refField := elem.Field(i)
		field := &Field{
			Name: refField.Name,
			Type: refField.Type.String(),
		}
		if refField.Name == "AutoId" {
			field.Name = "Id"
		}
		if refField.Name == "SoftDelete" {
			instance.HasDelete = true
		}
		tag := refField.Tag.Get("gorm")
		if tag != "" {
			parseTagToField(tag, field)
		}
		fields = append(fields, field)
	}
	instance.Name = elem.Name()
	instance.Fields = fields
	// 值反射
	v := reflect.ValueOf(m)
	method := v.MethodByName("TableComment")
	if method.IsValid() && len(method.Call(nil)) > 0 {
		instance.Comment = method.Call(nil)[0].String()
	}
	return
}

func parseTagToField(tag string, field *Field) {
	values := strings.Split(tag, ";")
	for _, value := range values {
		result := strings.Split(value, ":")
		switch result[0] {
		case "comment":
			field.Comment = strings.TrimSpace(result[1])
		}
	}
}
