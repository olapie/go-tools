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
{{ if not .Readonly }} Set{{.Name}}({{.VarName}} {{.Type}}) error {{end}}
{{end}}
{{range .Methods}}
{{.}}
{{end}}
Dirty() bool

AppendValidator(validator any)

Modifier() {{$interfaceName}}Modifier

// Unsafe returns underlying fields for efficient read only. DO NOT modify the fields
Unsafe() *{{$fieldsStructName}}
}

type {{$interfaceName}}Builder interface {
    {{range .Fields}} With{{.Name}}({{.VarName}} {{.Type}}) {{$interfaceName}}Builder
    {{end}}
    Build() ({{$interfaceName}}, error)
    MustBuild() {{$interfaceName}}
}

type {{$interfaceName}}Modifier interface {
    {{range .Fields}}

      {{- if not .Readonly }} Set{{.Name}}({{.VarName}} {{.Type}}) {{$interfaceName}}Modifier {{end}}
    {{end}}
    Error() error
}

//{{$interfaceName}}FieldsValidator validate all fields
type {{$interfaceName}}FieldsValidator interface {
    ValidateFields(fields *{{$fieldsStructName}}) error
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
    dirty bool
}

func ({{$receiver}} *{{$implName}}) AppendValidator(validator any) {
    if slices.Contains({{$receiver}}.validators, validator) {
        return
    }
    {{$receiver}}.validators = append({{$receiver}}.validators, validator)
}

func ({{$receiver}} *{{$implName}}) Modifier() {{$interfaceName}}Modifier {
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

{{if .Readonly}}

func ({{$receiver}} *{{$implName}}) set{{.Name}}({{.VarName}} {{.Type}}) error {

{{else}}

func ({{$receiver}} *{{$implName}}) Set{{.Name}}({{.VarName}} {{.Type}}) error {
        {{- if .SetIfNil}}  if {{$receiver}}.fields.{{.Name}} != nil {
                return errors.New("cannot overwrite field {{.Name}}")
            }
        {{else if .SetIfZero}}  var zero {{.Type}}
            if {{$receiver}}.fields.{{.Name}} != zero {
                return errors.New("cannot overwrite field {{.Name}}")
            }
        {{end}}
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
    {{$receiver}}.dirty = true
    return nil
}

{{end}}

func ({{$receiver}} *{{$implName}}) Dirty() bool {
    return {{$receiver}}.dirty
}

type {{$builderName}} struct {
    impl {{$implName}}
    err error
}

func New{{$interfaceName}}Builder(validators ...any) {{$interfaceName}}Builder {
    b := new({{$builderName}})
    b.impl.validators = validators
    return b
}

{{range .Fields}}

func (b *{{$builderName}}) With{{.Name}}({{.VarName}} {{.Type}}) {{$interfaceName}}Builder  {
    if b.err == nil {
        {{if .Readonly}} b.err = b.impl.set{{.Name}}({{.VarName}})
        {{else}}  b.err = b.impl.Set{{.Name}}({{.VarName}})
        {{end}}
    }
    return b
}

{{end}}

func (b *{{$builderName}}) Build() ({{$interfaceName}}, error) {
    if b.err != nil {
        return nil, b.err
    }

       for _, validator := range b.impl.validators {
        if v, ok := validator.({{$interfaceName}}FieldsValidator); ok {
               if err := v.ValidateFields(&b.impl.fields); err != nil {
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

type {{$modifierName}} struct {
    impl *{{$implName}}
    validatedFields bool
    err error
}

func (m *{{$modifierName}}) Error() error {
    if m.err != nil {
        return m.err
    }
    if !m.validatedFields {
        m.validatedFields = true
        for _, validator := range m.impl.validators {
            if v, ok := validator.({{$interfaceName}}FieldsValidator); ok {
                   if err := v.ValidateFields(&m.impl.fields); err != nil {
                       m.err = err
                       break
                   }
               }
           }
    }
    return m.err
}

{{range .Fields}}
{{if not .Readonly}}
func (m *{{$modifierName}}) Set{{.Name}}({{.VarName}} {{.Type}}) {{$interfaceName}}Modifier {
    if m.err != nil {
        return m
    }
    m.validatedFields = false
    m.err = m.impl.Set{{.Name}}({{.VarName}})
    return m
}
{{end}}
{{end}}

// Restore{{$interfaceName}} restores {{$interfaceName}} from storage e.g. database
// Bypass validation to improve performance
func Restore{{$interfaceName}}(fields *{{$fieldsStructName}}) {{$interfaceName}} {
    return &{{$implName}} {
        fields: *fields,
    }
}


{{end}}