{{ define "model2" }}

type {{.Name}} struct {
{{range .ValueFields}} {{if eq .Type  "time.Time"}} {{.Name}} int64 `json:"{{.SnakeName}},omitempty"`
{{else}} {{.Name}} {{.Type}} `json:"{{.SnakeName}},omitempty"` {{end}}
{{end}}}

func New{{.Name}}() *{{.Name}} {
    return new({{.Name}})
}


{{end}}