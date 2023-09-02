{{ define "struct" }}
{{$name := toStructName .Name}}

type {{$name}} struct {
{{range .Fields}}   {{.Name}} {{.Type}} `json:"{{.JsonName}},omitempty"`
{{end}}

}


{{end}}