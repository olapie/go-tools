{{define "value"}}

{{range .Values}}
type {{.Name}} struct {
    {{range .Fields}}
        {{.Name}} {{.Type}} {{if .Tag}} `{{.Tag}}` {{end}}
    {{end}}
}

{{end}}

{{end}}
