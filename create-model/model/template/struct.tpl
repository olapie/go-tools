{{ define "struct" }}
{{$name := toPascal .Name}}

type {{$name}} struct {
{{range .Fields}}   {{.Name}} {{.Type}} `json:"{{.JsonName}},omitempty"`
{{end}}

}


{{end}}