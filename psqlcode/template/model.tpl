{{ define `model` }}package generate
import (
	"time"
)

{{range .Entities}}

type {{.Name}} struct {
{{range .Fields}}   {{toPascal .Name}} {{.Type}} `json:"{{toSnake .Name}}"`
{{end}}}

{{end}}

{{end}}