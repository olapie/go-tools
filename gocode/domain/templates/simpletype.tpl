{{define "simpletype"}}
{{range .SimpleTypes}}
type {{.Name}} {{.Type}}
{{end}}
{{end}}
