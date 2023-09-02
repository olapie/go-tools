{{ define "model" }}

{{$entityName := .Name}}
{{$valueName := .ValueName}}
{{$builderName := .BuilderName}}

type {{$valueName}} struct {
{{range .Fields}}{{.Name}} {{.Type}} `json:"{{.JsonName}},omitempty"`
{{end}}}

func New{{$valueName}}() *{{$valueName}} {
    return new({{$valueName}})
}

type {{$entityName}} struct {
    legal bool
    err error
    v {{$valueName}}
}

func (e *{{$entityName}}) ensureConstructed() {
    if e.legal {
        return
    }
    panic("{{$entityName}} is not created with builder")
}

func (e *{{$entityName}}) Error() error {
    return e.err
}

func (e *{{$entityName}}) Value() *{{$valueName}} {
    return To{{$valueName}}(e)
}

{{range .Fields}}

func (e *{{$entityName}}) Get{{.Name}}() {{.Type}} {
    e.ensureConstructed()
    return e.v.{{.Name}}
}

func (e *{{$entityName}}) Set{{.Name}}({{.VarName}} {{.Type}}) {
    if e.err != nil {
        return
    }
    f, ok := any(e).(interface {
        Validate{{.Name}}({{.Type}}) error
    })
    if ok {
        e.err = f.Validate{{.Name}}({{.VarName}})
        if e.err != nil {
            return
        }
    }
    e.v.{{.Name}} = {{.VarName}}
}

{{end}}

type {{$builderName}} struct {
    entity {{$entityName}}
}

func New{{$builderName}}() *{{$builderName}} {
    b := new({{$builderName}})
    b.entity.legal = true
    return b
}

func (b *{{$builderName}}) Build() (*{{$entityName}}, error) {
    if err := b.entity.Error(); err != nil {
        return nil, err
    }

    f, ok := any(&b.entity).(interface {
        Validate() error
    })

    if !ok {
        f, _ = any(b.entity).(interface {
            Validate() error
        })
    }

    if f != nil {
        if err := f.Validate(); err != nil {
            return nil, err
        }
    }

    return &b.entity, nil
}

func (b *{{$builderName}}) MustBuild() *{{$entityName}} {
    v, err := b.Build()
    if err != nil {
        panic(err)
    }
    return v
}

{{range .Fields}}

func (b *{{$builderName}}) With{{.Name}}({{.VarName}} {{.Type}}) *{{$builderName}}  {
    b.entity.Set{{.Name}}({{.VarName}})
    return b
}

{{end}}

func To{{$valueName}}(e *{{$entityName}}) *{{$valueName}} {
    return &{{$valueName}}{
{{range .ValueFields}} {{.Name}}: e.v.{{.Name}},
{{end}}
    }
}

func To{{$valueName}}List(a []*{{$entityName}}) []*{{$valueName}} {
    if a == nil {
        return nil
    }

    l := make([]*{{$valueName}}, len(a))
    for idx, e := range a {
        l[idx] = To{{$valueName}}(e)
    }
    return l
}

{{if not .Unexposed}}

// New{{$entityName}} creates {{$entityName}} from {{$valueName}}
func New{{$entityName}}(v *{{$valueName}}) *{{$entityName}} {
    b := New{{$builderName}}()
        b.entity.legal = true
    {{range .ValueFields}}b.With{{.Name}}(v.{{.Name}})
    {{end}}return b.MustBuild()
}

func New{{$entityName}}List(a []*{{$valueName}}) []*{{$entityName}} {
    if a == nil {
        return nil
    }

    l := make([]*{{$entityName}}, len(a))
    for idx, e := range a {
        l[idx] = New{{$entityName}}(e)
    }
    return l
}


{{end}}

{{end}}