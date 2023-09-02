package main

import (
	"embed"
	"text/template"

	"go.olapie.com/naming"
)

//go:embed template
var tplFS embed.FS
var globalTemplate = template.New("")

func init() {
	globalTemplate = globalTemplate.Funcs(template.FuncMap{
		"toPascal": naming.ToPascal,
		"toCamel":  naming.ToCamel,
		"toSnake":  naming.ToSnake,
		"toEntityName": func(s string) string {
			return naming.ToPascal(s) + "Entity"
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
	})
	globalTemplate = template.Must(globalTemplate.ParseFS(tplFS, "template/*.tpl"))
}
