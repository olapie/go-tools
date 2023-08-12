{{ define "entity" }}

{{$fieldsStructName := printf "%s%s" .UpperName "Fields"}}
{{$interfaceName := .UpperName}}
{{$implName := printf "%s%s"  .LowerName "Impl"}}
{{$builderName := printf "%s%s"  .LowerName "Builder"}}

type {{$interfaceName}} interface {
{{range .Fields}}
{{.Name}}() {{.Type}}
Set{{.Name}}({{.VarName}} {{.Type}}) error
{{end}}
}

//{{$interfaceName}}FieldsValidator validate all fields
type {{$interfaceName}}FieldsValidator interface {
    ValidateFields(fields {{$fieldsStructName}}) error
}

{{range .Fields}}
type {{$interfaceName}}{{.Name}}Validator interface {
    Validate{{.Name}}({{.VarName}} {{.Type}}) error
}
{{end}}

type {{$fieldsStructName}} struct {
{{range .Fields}}{{.Name}} {{.Type}} `json:"{{.SnakeName}},omitempty"`
{{end}}}

type {{$implName}} struct {
    fields {{$fieldsStructName}}
    validator any
}

{{range .Fields}}

func (i *{{$implName}}) {{.Name}}() {{.Type}} {
    return i.fields.{{.Name}}
}

func (i *{{$implName}}) Set{{.Name}}({{.VarName}} {{.Type}}) error {
    if {{.SetNX}} {
        var zero {{.Type}}
        if i.fields.{{.Name}} != zero {
            return errors.New("cannot overwrite field {{.Name}}")
        }
    }

    v, ok := i.validator.({{$interfaceName}}{{.Name}}Validator)
    if ok {
        if err := v.Validate{{.Name}}({{.VarName}}); err != nil {
            return err
        }
    }
    i.fields.{{.Name}} = {{.VarName}}
    return nil
}

{{end}}

type {{$builderName}} struct {
    impl {{$implName}}
    err error
}

func New{{$interfaceName}}Builder(validator any) *{{$builderName}} {
    b := new({{$builderName}})
    b.impl.validator = validator
    return b
}

func (b *{{$builderName}}) Build() ({{$interfaceName}}, error) {
    if b.err != nil {
        return nil, b.err
    }

    if v, ok := b.impl.validator.({{$interfaceName}}FieldsValidator); ok {
        if err := v.ValidateFields(b.impl.fields); err != nil {
            return nil, err
        }
    }
    return &b.impl, nil
}

func (b *{{$builderName}}) MustBuild() {{$interfaceName}} {
    v, err := b.Build()
    if err != nil {
        panic(err)
    }
    return v
}

{{range .Fields}}

func (b *{{$builderName}}) With{{.Name}}({{.VarName}} {{.Type}}) *{{$builderName}}  {
    if b.err == nil {
        b.err = b.impl.Set{{.Name}}({{.VarName}})
    }
    return b
}

{{end}}


{{end}}