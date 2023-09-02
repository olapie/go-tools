package model

import (
	"embed"
	"strings"
	"text/template"

	"go.olapie.com/naming"
	"go.olapie.com/utils"
)

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`

	VarName  string `json:"var_name"`
	JsonName string `json:"snake_name"`
}

type Entity struct {
	Name        string   `json:"name"`
	Unexposed   []string `json:"unexposed"`
	Fields      []*Field `json:"fields"`
	ValueFields []*Field `json:"value_fields"`

	ValueName   string `json:"value_name"`
	BuilderName string `json:"builder_name"`
}

func (e *Entity) Exported(name string) bool {
	for _, s := range e.Unexposed {
		if s == name {
			return false
		}
	}
	return true
}

type Struct struct {
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
}

type Enum struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Values map[string]string `json:"values"`
}

type Model struct {
	Enums    []*Enum   `json:"enums"`
	Entities []*Entity `json:"entities"`
	Structs  []*Struct `json:"structs"`
}

func (m *Model) ContainsType(name string) bool {
	for _, e := range m.Enums {
		if e.Name == name {
			return true
		}
	}

	for _, e := range m.Entities {
		if e.Name == name {
			return true
		}
	}

	for _, e := range m.Structs {
		if e.Name == name {
			return true
		}
	}

	return false
}

//go:embed template
var tplFS embed.FS
var globalTemplate = template.New("")

func init() {
	globalTemplate = globalTemplate.Funcs(template.FuncMap{
		"toStructName": utils.ToClassName,
		"toCamel":      naming.ToCamel,
		"toSnake":      naming.ToSnake,
		"toEntityName": func(s string) string {
			return utils.ToClassName(s) + "Entity"
		},
		"toValueName": func(s string) string {
			return utils.ToClassName(s)
		},
		"toBuilderName": func(s string) string {
			return naming.ToCamel(s) + "EntityBuilder"
		},
		"toModifierName": func(s string) string {
			return naming.ToCamel(s) + "EntityModifier"
		},
		"first": func(s string) string {
			return s[:1]
		},
		"receiver": func(s string) string {
			return strings.ToLower(s[:1])
		},
	})
	globalTemplate = template.Must(globalTemplate.ParseFS(tplFS, "template/*.tpl"))
}
