/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"create-entity/templates"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go.olapie.com/utils"
	"go/format"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{filename}",
	Short: "Generate models",
	Long:  `Generate models from json file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			return
		}
		Generate(args[0])
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gocode.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type JSONModel struct {
	Imports  []string                     `json:"imports"`
	Entities map[string]map[string]string `json:"entities"`
}

func Generate(fileName string) {

	var reservedNames = []string{"type", "struct", "map", "slices", "maps", "string", "int", "os"}

	var globalTemplate = template.New("")

	globalTemplate = template.Must(globalTemplate.ParseFS(templates.FS, "*.tpl"))

	var jsonModel JSONModel
	utils.MustNil(json.Unmarshal(utils.MustGet(os.ReadFile(fileName)), &jsonModel))

	var entities []*Entity
	for name, m := range jsonModel.Entities {
		e := &Entity{
			UpperName: utils.ToClassName(name),
			LowerName: utils.ToCamel(name),
		}
		e.Receiver = e.LowerName[0:1]
		for field, attr := range m {
			if strings.HasPrefix(field, "$method") {
				e.Methods = append(e.Methods, attr)
				continue
			}
			f := &Field{
				Name:      utils.ToClassName(field),
				Type:      strings.Split(attr, ",")[0],
				SetNX:     strings.Contains(attr, "setnx"),
				SetEmpty:  strings.Contains(attr, "setempty"),
				VarName:   utils.ToCamel(field),
				SnakeName: utils.ToSnake(field),
			}
			if slices.Contains(reservedNames, f.VarName) {
				f.VarName = f.VarName + "Val"
			}
			e.Fields = append(e.Fields, f)
		}
		sort.Slice(e.Fields, func(i, j int) bool {
			return e.Fields[i].Name < e.Fields[j].Name
		})
		entities = append(entities, e)
	}

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].LowerName < entities[j].LowerName
	})

	var b bytes.Buffer
	for _, e := range entities {
		err := globalTemplate.ExecuteTemplate(&b, "entity", e)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generate entity %s\n", e.UpperName)
	}

	output := "package entity\n"
	output += "import(\n"

	var shortImports []string
	var longImports []string
	for _, s := range jsonModel.Imports {
		s = fmt.Sprintf("\"%s\"", s)
		if strings.Contains(s, ".") {
			longImports = append(longImports, s)
		} else {
			shortImports = append(shortImports, s)
		}
	}
	sort.Strings(shortImports)
	sort.Strings(longImports)
	output += strings.Join(shortImports, "\n")
	output += "\n\n"
	output += strings.Join(longImports, "\n")
	output += ")\n"
	output += b.String()
	for {
		replaced := strings.ReplaceAll(output, "\n\n\n", "\n\n")
		if output == replaced {
			break
		}
		output = replaced
	}

	data, err := format.Source([]byte(output))
	if err != nil {
		log.Fatalln(err)
	}

	os.Mkdir("gen", 0755)
	err = os.WriteFile("gen/entity.gen.go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	SetNX    bool   `json:"set_nx"`
	SetEmpty bool   `json:"set_empty"`

	VarName   string `json:"var_name"`
	SnakeName string `json:"snake_name"`
}

type Entity struct {
	UpperName string
	LowerName string
	Receiver  string
	Fields    []*Field `json:"fields"`
	Methods   []string
}
