{{define "service"}}

{{range .Services}}
type {{.Name}} interface {
    {{range .Methods}}
        {{.}}
    {{end}}
}

{{end}}

{{end}}
