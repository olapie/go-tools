{{define "interface"}}

{{range .Interfaces}}
type {{.Name}} interface {
    {{range .Methods}}
        {{.}}
    {{end}}
}

{{end}}

{{end}}
