package domain

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"go/format"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"

	"go.olapie.com/naming"
	"go.olapie.com/tools/gocode/domain/templates"
	"go.olapie.com/utils"
)

//go:embed testdata/domain.xml
var ExampleXML string

type Model struct {
	Imports     []string      `xml:"import"`
	Entities    []*Entity     `xml:"entity"`
	Aliases     []*SimpleType `xml:"alias"`
	SimpleTypes []*SimpleType `xml:"simpletype"`
	Structs     []*StructType `xml:"struct"`
	Interfaces  []*Interface  `xml:"interface"`
	JSONNaming  string        `xml:"jsonNaming,attr"`
	BSONNaming  string        `xml:"bsonNaming,attr"`
	Package     string        `json:"package,attr"`

	ShortImports []string
	LongImports  []string
}

type Entity struct {
	Interface
	ImplName        string
	ImplReceiver    string
	ValidatorPrefix string

	BuilderName     string
	BuilderImplName string

	ModifierName     string
	ModifierImplName string

	FieldsName string
	Fields     []*EntityField `xml:"field"`
	JSON       bool           `xml:"json,attr"`
	BSON       bool           `xml:"bson,attr"`
}

type EntityField struct {
	StructField
	Readonly  bool   `xml:"readonly,attr"`
	BSON      string `xml:"bson,attr"`
	SetIfNil  bool   `xml:"setIfNil"`
	SetIfZero bool   `xml:"setIfZero"`
	VarName   string
}

type StructType struct {
	Name       string         `xml:"name,attr"`
	JSON       bool           `xml:"json,attr"`
	Fields     []*StructField `xml:"field"`
	Embeddings []string       `xml:"embed"`
}

type StructField struct {
	Name string `xml:",innerxml"`
	Type string `xml:"type,attr"`
	Tag  string
}

type SimpleType struct {
	Name string `xml:",innerxml"`
	Type string `xml:"type,attr"`
}

type Interface struct {
	Name    string   `xml:"name,attr"`
	Methods []string `xml:"method"`
}

var reservedNames = []string{"type", "struct", "map", "slices", "maps", "string", "int", "os", "chan", "rune", "os",
	"io", "reflect", "http", "byte", "bytes", "time", "delete", "clear", "min", "max", "copy", "make", "new", "range",
	"switch", "select", "for", "main", "return", "error", "errors", "fmt", "panic", "defer"}

func loadTemplates() *template.Template {
	var t = template.New("")
	return template.Must(t.ParseFS(templates.FS, "*.tpl"))
}

func parseModel(xmlFilename string) *Model {
	data, err := os.ReadFile(xmlFilename)
	if err != nil {
		panic(err)
	}
	var m Model
	err = xml.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	m.Package = strings.TrimSpace(m.Package)
	if m.Package == "" {
		m.Package = "domain"
	}
	m.Imports = append(m.Imports, "fmt", "slices")
	m.Imports = utils.UniqueSlice(m.Imports)
	sort.Strings(m.Imports)
	for _, s := range m.Imports {
		if strings.Contains(s, ".") {
			m.LongImports = append(m.LongImports, s)
		} else {
			m.ShortImports = append(m.ShortImports, s)
		}
	}

	sort.Slice(m.Aliases, func(i, j int) bool {
		return m.Aliases[i].Name < m.Aliases[j].Name
	})

	sort.Slice(m.SimpleTypes, func(i, j int) bool {
		return m.SimpleTypes[i].Name < m.SimpleTypes[j].Name
	})

	sort.Slice(m.Interfaces, func(i, j int) bool {
		return m.Interfaces[i].Name < m.Interfaces[j].Name
	})

	for _, i := range m.Interfaces {
		sort.Strings(i.Methods)
		for j, method := range i.Methods {
			i.Methods[j] = strings.ReplaceAll(method, "\n", "")
		}
	}

	for _, v := range m.Structs {
		sort.Strings(v.Embeddings)
		sort.Slice(v.Fields, func(i, j int) bool {
			return v.Fields[i].Name < v.Fields[j].Name
		})

		for _, f := range v.Fields {
			originName := f.Name
			f.Name = naming.ToPascal(originName, naming.WithAcronym())
			var tags []string
			if v.JSON {
				jsonName := naming.ToCamel(originName)
				if m.JSONNaming == "SnakeCase" {
					jsonName = naming.ToSnake(originName)
				}
				tags = append(tags, fmt.Sprintf(`json:"%s,omitempty"`, jsonName))
			}
			if len(tags) > 0 {
				f.Tag = fmt.Sprintf(`%s`, strings.Join(tags, " "))
			}
		}
	}

	for _, e := range m.Entities {
		sort.Strings(e.Methods)
		for j, method := range e.Methods {
			e.Methods[j] = strings.ReplaceAll(method, "\n", "")
		}

		sort.Slice(e.Fields, func(i, j int) bool {
			return e.Fields[i].Name < e.Fields[j].Name
		})
		e.Name = naming.ToPascal(e.Name)
		name := strings.Replace(e.Name, "Entity", "", 1)
		e.ValidatorPrefix = name
		e.BuilderName = name + "Builder"
		e.ModifierName = name + "Modifier"

		camel := naming.ToCamel(name)
		e.ImplReceiver = camel[:1]
		e.ImplName = camel + "Impl"
		e.BuilderImplName = camel + "BuilderImpl"
		e.ModifierImplName = camel + "ModifierImpl"
		e.FieldsName = name + "Fields"
		for _, f := range e.Fields {
			originName := f.Name
			f.Name = naming.ToPascal(originName, naming.WithAcronym())
			f.VarName = naming.ToCamel(originName, naming.WithAcronym())
			if slices.Contains(reservedNames, f.VarName) {
				f.VarName += "Val"
			}

			var tags []string
			if e.JSON {
				jsonName := naming.ToCamel(originName)
				if m.JSONNaming == "SnakeCase" {
					jsonName = naming.ToSnake(originName)
				}
				tags = append(tags, fmt.Sprintf(`json:"%s,omitempty"`, jsonName))
			}
			if e.BSON {
				if f.BSON != "" {
					tags = append(tags, fmt.Sprintf(`bson:"%s"`, f.BSON))
				} else {
					bsonName := naming.ToCamel(originName)
					if m.BSONNaming == "SnakeCase" {
						bsonName = naming.ToSnake(originName)
					}
					tags = append(tags, fmt.Sprintf(`bson:"%s,omitempty"`, bsonName))
				}
			}

			if len(tags) > 0 {
				f.Tag = fmt.Sprintf(`%s`, strings.Join(tags, " "))
			}
		}
	}

	return &m
}

func Generate(xmlFilename, outputGoFilename string) {
	m := parseModel(xmlFilename)
	//
	//jsonData, _ := json.Marshal(m)
	//fmt.Println(string(jsonData))

	tpl := loadTemplates()
	output := bytes.NewBuffer(nil)
	tplNames := []string{"import", "alias", "simpletype", "struct", "interface", "entity"}
	for _, name := range tplNames {
		err := tpl.ExecuteTemplate(output, name, m)
		if err != nil {
			panic(fmt.Sprintf("render template %s: %v", name, err))
		}
	}

	data, err := format.Source(output.Bytes())
	if err != nil {
		fmt.Println(err)
		err = os.WriteFile("go-code-error.out", output.Bytes(), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(1)
	}

	err = os.WriteFile(outputGoFilename, data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
