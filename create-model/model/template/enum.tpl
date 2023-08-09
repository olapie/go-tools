{{ define "enum" }}
{{range .}}
{{$name := .Name}}
type {{$name}} {{.Type}}
{{if .Values}}
const (
{{range $key, $value := .Values}} {{$key}} {{$name}} = {{$value}}
{{end}})
{{end}}
{{end}}
{{end}}
