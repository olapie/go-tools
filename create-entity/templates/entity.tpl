{{ define "entity" }}

{{$fieldsStructName := printf "%s%s" .UpperName "Fields"}}
{{$interfaceName := .UpperName}}
{{$implName := printf "%s%s"  .LowerName "Impl"}}
{{$builderName := printf "%s%s"  .LowerName "Builder"}}
{{$receiver := .Receiver}}
{{$modifierName := printf "%s%s"  .LowerName "Modifier"}}

type {{$interfaceName}} interface {
{{range .Fields}}
Get{{.Name}}() {{.Type}}
Set{{.Name}}({{.VarName}} {{.Type}}) error
{{end}}
{{range .Methods}}
{{.}}
{{end}}

AppendValidator(validator any)

Modifier() *{{$modifierName}}

// Unsafe returns underlying fields for efficient read only. DO NOT modify the fields
Unsafe() *{{$fieldsStructName}}
}

//{{$interfaceName}}FieldsValidator validate all fields
type {{$interfaceName}}FieldsValidator interface {
    ValidateFields(fields {{$fieldsStructName}}) error
}

{{range .Fields}}
type {{$interfaceName}}{{.Name}}Validator interface {
    Validate{{.Name}}({{$receiver}} {{$interfaceName}}, {{.VarName}} {{.Type}}) error
}
{{end}}

type {{$fieldsStructName}} struct {
{{range .Fields}}{{.Name}} {{.Type}} `json:"{{.SnakeName}},omitempty"`
{{end}}}

type {{$implName}} struct {
    fields {{$fieldsStructName}}
    validators []any
}

func ({{$receiver}} *{{$implName}}) AppendValidator(validator any) {
    if slices.Contains({{$receiver}}.validators, validator) {
        return
    }
    {{$receiver}}.validators = append({{$receiver}}.validators, validator)
}

func ({{$receiver}} *{{$implName}}) Modifier() *{{$modifierName}} {
    return &{{$modifierName}}  {
        impl: {{$receiver}},
    }
}

func ({{$receiver}} *{{$implName}}) Unsafe() *{{$fieldsStructName}} {
    return &{{$receiver}}.fields
}

{{range .Fields}}

func ({{$receiver}} *{{$implName}}) Get{{.Name}}() {{.Type}} {
    return {{$receiver}}.fields.{{.Name}}
}

func ({{$receiver}} *{{$implName}}) Set{{.Name}}({{.VarName}} {{.Type}}) error {
        {{- if .SetEmpty}}  if len({{$receiver}}.fields.{{.Name}}) != 0 {
                return errors.New("cannot overwrite field {{.Name}}")
            }  {{end}}
    {{- if .SetNX}}  var zero {{.Type}}
        if {{$receiver}}.fields.{{.Name}} != zero {
            return errors.New("cannot overwrite field {{.Name}}")
        }
        {{end}}
       {{- $validatorName := printf "%s%s" .VarName "Validator"}}
       for _, validator := range {{$receiver}}.validators {
        if {{$validatorName}}, ok := validator.({{$interfaceName}}{{.Name}}Validator); ok {
               if err := {{$validatorName}}.Validate{{.Name}}({{$receiver}}, {{.VarName}}); err != nil {
                   return err
               }
           }
       }

    if validator, ok := any({{$receiver}}).({{$interfaceName}}{{.Name}}Validator); ok {
            if err := validator.Validate{{.Name}}({{$receiver}}, {{.VarName}}); err != nil {
                return err
            }
        }
    {{$receiver}}.fields.{{.Name}} = {{.VarName}}
    return nil
}

{{end}}

type {{$builderName}} struct {
    impl {{$implName}}
    err error
}

func New{{$interfaceName}}Builder(validators ...any) *{{$builderName}} {
    b := new({{$builderName}})
    b.impl.validators = validators
    return b
}

func (b *{{$builderName}}) Build() ({{$interfaceName}}, error) {
    if b.err != nil {
        return nil, b.err
    }

       for _, validator := range b.impl.validators {
        if v, ok := validator.({{$interfaceName}}FieldsValidator); ok {
               if err := v.ValidateFields(b.impl.fields); err != nil {
                   b.err = err
                   return nil, err
               }
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

type {{$modifierName}} struct {
    impl *{{$implName}}
    err error
}

func (m *{{$modifierName}}) Error() error {
    return m.err
}

{{range .Fields}}
func (m *{{$modifierName}}) Set{{.Name}}({{.VarName}} {{.Type}}) *{{$modifierName}} {
    if m.err != nil {
        return m
    }
    m.err = m.impl.Set{{.Name}}({{.VarName}})
    return m
}
{{end}}

// Restore{{$interfaceName}} restores {{$interfaceName}} from storage e.g. database
// Bypass validation to improve performance
func Restore{{$interfaceName}}(fields *{{$fieldsStructName}}) {{$interfaceName}} {
    return &{{$implName}} {
        fields: *fields,
    }
}


{{end}}