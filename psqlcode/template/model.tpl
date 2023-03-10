{{ define `model` }}package generate
import (
	"time"
)

{{range .Entities}}

type {{.Name}} struct {
{{range .Fields}}   {{toStructName .Name}} {{.Type}} `json:"{{toSnake .Name}}"`
{{end}}}

{{end}}

{{end}}