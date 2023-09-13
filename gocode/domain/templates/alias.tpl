{{define "alias"}}
{{range .Aliases}}
type {{.Name}} = {{.Type}}
{{end}}
{{end}}
