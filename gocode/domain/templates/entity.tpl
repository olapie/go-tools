{{ define "entity" }}

{{range .Entities}}

{{$fieldsStructName := .FieldsName}}
{{$interfaceName := .Name}}
{{$implName := .ImplName}}
{{$builderName := .BuilderName}}
{{$builderImplName := .BuilderImplName}}
{{$receiver := .ImplReceiver}}
{{$modifierName := .ModifierName}}
{{$modifierImplName := .ModifierImplName}}

type {{$interfaceName}} interface {
{{range .Fields}} Get{{.Name}}() {{.Type}}
    {{- if not .Readonly }}
    Set{{.Name}}({{.VarName}} {{.Type}}) error {{end}}
{{end}}
{{range .Methods}}{{.}}
{{end}}

    Dirty() bool
    AppendValidator(validator any)
    Modifier() {{$modifierName}}

    // Unsafe returns underlying fields for efficient read only. DO NOT modify the fields
    Unsafe() *{{$fieldsStructName}}
}

type {{$builderName}} interface {
    {{range .Fields}} With{{.Name}}({{.VarName}} {{.Type}}) {{$builderName}}
    {{end}}
    Build() ({{$interfaceName}}, error)
    MustBuild() {{$interfaceName}}
}

type {{$modifierName}} interface {
    {{range .Fields}}

      {{- if not .Readonly }} Set{{.Name}}({{.VarName}} {{.Type}}) {{$modifierName}}
      Set{{.Name}}P({{.VarName}} *{{.Type}}) {{$modifierName}}
      {{end}}
    {{end}}
    Error() error
}

//{{$interfaceName}}_FieldsValidator validate all fields
type {{$interfaceName}}_FieldsValidator interface {
    ValidateFields(fields *{{$fieldsStructName}}) error
}

{{range .Fields}}
type {{$interfaceName}}_{{.Name}}Validator interface {
    Validate{{.Name}}({{.VarName}} {{.Type}}) ({{.Type}}, error)
}
{{end}}

type {{$fieldsStructName}} struct {
{{range .Fields}}{{.Name}} {{.Type}} {{if .Tag}} `{{.Tag}}` {{end}}
{{end}}}

type Unimplemented{{$interfaceName}} struct {

}
{{range .Methods}}
func (_ *Unimplemented{{$interfaceName}}){{.}} {
    panic("method is not implemented: {{.}}")
}
{{end}}

type {{$implName}} struct {
    Unimplemented{{$interfaceName}}

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

func ({{$receiver}} *{{$implName}}) Modifier() {{$modifierName}} {
    return &{{$modifierImplName}}  {
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
        var err error

{{else}}

func ({{$receiver}} *{{$implName}}) Set{{.Name}}({{.VarName}} {{.Type}}) error {
        var err error
        {{if .SetIfNil}}  if {{$receiver}}.fields.{{.Name}} != nil {
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
        if {{$validatorName}}, ok := validator.({{$interfaceName}}_{{.Name}}Validator); ok {
               if  {{.VarName}}, err = {{$validatorName}}.Validate{{.Name}}({{.VarName}}); err != nil {
                   return fmt.Errorf("invalid {{.VarName}}: %w", err)
               }
           }
       }

    if validator, ok := any({{$receiver}}).({{$interfaceName}}_{{.Name}}Validator); ok {
            if {{.VarName}}, err = validator.Validate{{.Name}}({{.VarName}}); err != nil {
                return fmt.Errorf("invalid {{.VarName}}: %w", err)
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

type {{$builderImplName}} struct {
    impl {{$implName}}
    err error
}

func New{{$builderName}}(validators ...any) {{$builderName}} {
    b := new({{$builderImplName}})
    b.impl.validators = validators
    return b
}

{{range .Fields}}

func (b *{{$builderImplName}}) With{{.Name}}({{.VarName}} {{.Type}}) {{$builderName}}  {
    if b.err == nil {
        {{if .Readonly}} b.err = b.impl.set{{.Name}}({{.VarName}})
        {{else}}  b.err = b.impl.Set{{.Name}}({{.VarName}})
        {{end}}
    }
    return b
}

{{end}}

func (b *{{$builderImplName}}) Build() ({{$interfaceName}}, error) {
    if b.err != nil {
        return nil, b.err
    }

       for _, validator := range b.impl.validators {
        if v, ok := validator.({{$interfaceName}}_FieldsValidator); ok {
               if err := v.ValidateFields(&b.impl.fields); err != nil {
                   b.err = err
                   return nil, err
               }
           }
       }
    return &b.impl, nil
}

func (b *{{$builderImplName}}) MustBuild() {{$interfaceName}} {
    v, err := b.Build()
    if err != nil {
        panic(err)
    }
    return v
}

type {{$modifierImplName}} struct {
    impl *{{$implName}}
    validatedFields bool
    err error
}

func (m *{{$modifierImplName}}) Error() error {
    if m.err != nil {
        return m.err
    }
    if !m.validatedFields {
        m.validatedFields = true
        for _, validator := range m.impl.validators {
            if v, ok := validator.({{$interfaceName}}_FieldsValidator); ok {
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
func (m *{{$modifierImplName}}) Set{{.Name}}({{.VarName}} {{.Type}}) {{$modifierName}} {
    if m.err != nil {
        return m
    }
    m.validatedFields = false
    m.err = m.impl.Set{{.Name}}({{.VarName}})
    return m
}

func (m *{{$modifierImplName}}) Set{{.Name}}P({{.VarName}} *{{.Type}}) {{$modifierName}} {
    if m.err != nil {
        return m
    }
    if {{.VarName}} == nil {
        return m
    }
    m.validatedFields = false
    m.err = m.impl.Set{{.Name}}(*{{.VarName}})
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

{{end}}