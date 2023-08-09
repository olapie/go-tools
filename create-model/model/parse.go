package model

import (
	"encoding/json"
	"go.olapie.com/utils"
	"log"
	"os"
	"sort"
	"strings"
)

//
//func ParseYAML(filename string) *Model {
//	type YAMLEnum struct {
//		Type   string        `yaml:"type"`
//		Values yaml.MapSlice `yaml:"values"`
//	}
//
//	type YAMLModel struct {
//		Enums    map[string]*YAMLEnum     `yaml:"enums"`
//		Entities map[string]yaml.MapSlice `yaml:"entities"`
//		Structs  map[string]yaml.MapSlice `yaml:"structs"`
//	}
//
//	data, err := os.ReadFile(filename)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	var yamlModel YAMLModel
//	err = yaml.Unmarshal(data, &yamlModel)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	var model Model
//	for name, enum := range yamlModel.Enums {
//		e := new(Enum)
//		e.Name = name
//		e.Type = enum.Type
//		e.Values = make([]*EnumValue, len(enum.Values))
//		for i, ev := range enum.Values {
//			enumName, ok := ev.Key.(string)
//			if !ok {
//				log.Panicf("%v is not string", ev.Key)
//			}
//			if e.Type == "string" {
//				if str, ok := ev.Value.(string); !ok || str == "" {
//					log.Panicf("%v is not string", ev.Value)
//				} else {
//					if str[0] != '"' {
//						ev.Value = fmt.Sprintf(`"%s"`, str)
//					}
//				}
//			}
//			e.Values[i] = &EnumValue{Name: enumName, Value: ev.Value}
//		}
//		model.Enums = append(model.Enums, e)
//	}
//
//	for name, entity := range yamlModel.Entities {
//		e := new(Entity)
//		e.Name = name
//		e.Fields = make([]*Field, len(entity))
//		for i, ev := range entity {
//			e.Fields[i] = &Field{
//				Name: ev.Key.(string),
//				Type: ev.Value.(string),
//			}
//		}
//		model.Entities = append(model.Entities, e)
//	}
//
//	for name, sv := range yamlModel.Structs {
//		e := new(Struct)
//		e.Name = name
//		e.Fields = make([]*Field, len(sv))
//		for i, ev := range sv {
//			e.Fields[i] = &Field{
//				Name: ev.Key.(string),
//				Type: ev.Value.(string),
//			}
//		}
//		model.Structs = append(model.Structs, e)
//	}
//
//	for _, e := range model.Entities {
//		for _, f := range e.Fields {
//			if model.ContainsType(f.Type) {
//				f.Type = utils.ToClassName(f.Type)
//			}
//		}
//	}
//
//	for _, e := range model.Structs {
//		for _, f := range e.Fields {
//			if model.ContainsType(f.Type) {
//				f.Type = utils.ToClassName(f.Type)
//			}
//		}
//	}
//
//	return &model
//}

func ParseJSON(filename string) *Model {
	var m *Model

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Fatalln(err)
	}

	sort.Slice(m.Entities, func(i, j int) bool {
		return m.Entities[i].Name < m.Entities[j].Name
	})

	sort.Slice(m.Enums, func(i, j int) bool {
		return m.Enums[i].Name < m.Enums[j].Name
	})

	sort.Slice(m.Structs, func(i, j int) bool {
		return m.Structs[i].Name < m.Structs[j].Name
	})

	for _, e := range m.Entities {

		sort.Slice(e.Fields, func(i, j int) bool {
			return e.Fields[i].Name < e.Fields[j].Name
		})

		e.ValueName = utils.ToClassName(e.Name)
		e.Name = e.ValueName + "Entity"
		e.BuilderName = e.Name + "Builder"
		for _, f := range e.Fields {
			if f.Type == "" {
				f.Type = "string"
			}
			name := f.Name
			f.Name = utils.ToClassName(f.Name)
			if f.SnakeName == "" {
				f.SnakeName = utils.ToSnake(f.Name)
			}
			if f.VarName == "" {
				f.VarName = utils.ToCamel(f.SnakeName)
			}

			if e.Exported(name) {
				f2 := *f
				prefix := "model."
				if strings.HasPrefix(f2.Type, prefix) {
					f2.Type = f2.Type[len(prefix):]
				}
				prefix = "*model."
				if strings.HasPrefix(f2.Type, prefix) {
					f2.Type = "*" + f2.Type[len(prefix):]
				}
				e.ValueFields = append(e.ValueFields, &f2)
			}
		}
	}

	for _, e := range m.Structs {
		e.Name = utils.ToClassName(e.Name)
		for _, f := range e.Fields {
			if f.Type == "" {
				f.Type = "string"
			}
			f.Name = utils.ToClassName(f.Name)
			if f.SnakeName == "" {
				f.SnakeName = utils.ToSnake(f.Name)
			}
			if f.VarName == "" {
				f.VarName = utils.ToCamel(f.SnakeName)
			}
		}
	}

	for _, e := range m.Enums {
		e.Name = utils.ToClassName(e.Name)
		vm := make(map[string]string)
		for k, v := range e.Values {
			k = e.Name + utils.ToClassName(k)
			vm[k] = v
		}
		e.Values = vm
	}

	return m
}
