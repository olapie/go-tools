{{define "repo"}}

{{range .Repos}}
type {{.Name}} interface {
    {{range .Methods}}
        {{.}}
    {{end}}
}

{{end}}

{{end}}
