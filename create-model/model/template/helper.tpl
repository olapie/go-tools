{{ define "helper" }}

{{$name := .Name}}
{{$entityName := .Name}}
{{$valueName := .ValueName}}
{{$builderName := .BuilderName}}
{{$receiver := receiver .Name}}

/*


type {{$valueName}} struct {
{{range .Fields}}   {{.Name}} {{.Type}} `json:"{{.JsonName}},omitempty"`
{{end}}
}


func from{{$entityName}}({{$receiver}} *{{$entityName}}) *{{$valueName}} {
    return &{{$valueName}}{
{{range .Fields}} {{.Name}}: {{$receiver}}.Get{{.Name}}(),
{{end}}
    }
}

func from{{$entityName}}List(a []*{{$entityName}}) []*{{$valueName}} {
    if a == nil {
        return nil
    }

    l := make([]*{{$valueName}}, len(a))
    for idx, {{$receiver}} := range a {
        l[idx] = from{{$entityName}}({{$receiver}})
    }
    return l
}

func to{{$entityName}}({{$receiver}} *{{$valueName}}) *{{$entityName}} {
    b := New{{$entityName}}Builder()
{{range .Fields}}b.With{{.Name}}({{$receiver}}.{{.Name}})
{{end}}
    return b.MustBuild()
}


func to{{$entityName}}List(a []*{{$valueName}}) []*{{$entityName}} {
    if a == nil {
        return nil
    }

    l := make([]*{{$entityName}}, len(a))
    for idx, v := range a {
        l[idx] = to{{$entityName}}(v)
    }
    return l
}



*/

{{end}}