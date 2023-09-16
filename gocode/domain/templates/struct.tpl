{{define "struct"}}

{{range .Structs}}
type {{.Name}} struct {
    {{range .Embeddings}} {{.}}
    {{end}}
    {{range .Fields}}
        {{.Name}} {{.Type}} {{if .Tag}} `{{.Tag}}` {{end}}
    {{end}}
}

{{end}}

{{end}}
