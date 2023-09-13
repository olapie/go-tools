package domain

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go.olapie.com/naming"
	"go.olapie.com/utils"
	"go/format"
	"gocode/domain/templates"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"
)

type Model struct {
	Imports    []string     `xml:"import"`
	Entities   []*Entity    `xml:"entity"`
	Aliases    []*Alias     `xml:"alias"`
	Values     []*Value     `xml:"value"`
	Repos      []*Interface `xml:"repo"`
	Services   []*Interface `xml:"service"`
	JSONNaming string       `xml:"jsonNaming,attr"`
	BSONNaming string       `xml:"bsonNaming,attr"`

	ShortImports []string
	LongImports  []string
}

type Entity struct {
	Interface
	ImplName     string
	ImplReceiver string

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
	ValueField
	Readonly  bool   `xml:"readonly,attr"`
	BSON      string `xml:"bson,attr"`
	SetIfNil  bool   `xml:"setIfNil"`
	SetIfZero bool   `xml:"setIfZero"`
}

type Value struct {
	Name   string       `xml:"name,attr"`
	Fields []ValueField `xml:"field"`
}

type ValueField struct {
	Name    string `xml:",innerxml"`
	Type    string `xml:"type,attr"`
	Tag     string
	VarName string
}

type Alias struct {
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
	//
	//jsonData, _ := json.Marshal(m)
	//fmt.Println(string(jsonData))

	m.Imports = append(m.Imports, "context", "fmt", "slices")
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

	sort.Slice(m.Services, func(i, j int) bool {
		return m.Services[i].Name < m.Services[j].Name
	})

	sort.Slice(m.Repos, func(i, j int) bool {
		return m.Repos[i].Name < m.Repos[j].Name
	})

	for _, v := range m.Values {
		sort.Slice(v.Fields, func(i, j int) bool {
			return v.Fields[i].Name < v.Fields[j].Name
		})
	}

	for _, e := range m.Entities {
		sort.Slice(e.Fields, func(i, j int) bool {
			return e.Fields[i].Name < e.Fields[j].Name
		})
		e.Name = naming.ToPascal(e.Name)
		e.BuilderName = e.Name + "Builder"
		e.ModifierName = e.Name + "Modifier"

		camel := naming.ToCamel(e.Name)
		e.ImplReceiver = camel[:1]
		e.ImplName = camel + "Impl"
		e.BuilderImplName = camel + "BuilderImpl"
		e.ModifierImplName = camel + "ModifierImpl"
		e.FieldsName = e.Name + "Fields"
		for _, f := range e.Fields {
			f.Name = naming.ToPascal(f.Name)
			f.VarName = naming.ToCamel(f.Name)
			if slices.Contains(reservedNames, f.VarName) {
				f.VarName += "Val"
			}

			var tags []string
			if e.JSON {
				jsonName := naming.ToCamel(f.Name)
				if m.JSONNaming == "SnakeCase" {
					jsonName = naming.ToSnake(f.Name)
				}
				tags = append(tags, fmt.Sprintf(`json:"%s,omitempty"`, jsonName))
			}
			if e.BSON {
				if f.BSON != "" {
					tags = append(tags, fmt.Sprintf(`bson:"%s"`, f.BSON))
				} else {
					bsonName := naming.ToCamel(f.Name)
					if m.BSONNaming == "SnakeCase" {
						bsonName = naming.ToSnake(f.Name)
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
	tpl := loadTemplates()
	output := bytes.NewBuffer(nil)
	tplNames := []string{"import", "alias", "value", "entity", "repo", "service"}
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
