{{ define "model2" }}

type {{.Name}} struct {
{{range .ValueFields}} {{if eq .Type  "time.Time"}} {{.Name}} int64 `json:"{{.JsonName}},omitempty"`
{{else}} {{.Name}} {{.Type}} `json:"{{.JsonName}},omitempty"` {{end}}
{{end}}}

func New{{.Name}}() *{{.Name}} {
    return new({{.Name}})
}


{{end}}