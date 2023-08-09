package model

import (
	"bytes"
	"go/format"
	"os"
	"strings"

	"log"
)

func Generate(fileName string) {
	model := ParseJSON(fileName)
	var b bytes.Buffer
	err := globalTemplate.ExecuteTemplate(&b, "enum", model.Enums)
	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range model.Entities {
		err := globalTemplate.ExecuteTemplate(&b, "model", e)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generate entity %s\n", e.Name)
	}

	output := "package model\n"
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
	err = os.WriteFile("gen/model.gen.go", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	b.Reset()
	for _, e := range model.Structs {
		err := globalTemplate.ExecuteTemplate(&b, "struct", e)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generate struct %s\n", e.Name)
	}

	for _, e := range model.Entities {
		err := globalTemplate.ExecuteTemplate(&b, "model", e)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generate model %s\n", e.Name)
	}

	//output = "package model\n"
	//output += b.String()
	//for {
	//	replaced := strings.ReplaceAll(output, "\n\n\n", "\n\n")
	//	if output == replaced {
	//		break
	//	}
	//	output = replaced
	//}
	//
	//data, err = format.Source([]byte(output))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//os.Mkdir("gen", 0755)
	//err = os.WriteFile("gen/model.gen.go", data, 0644)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	b.Reset()
	for _, e := range model.Entities {
		err := globalTemplate.ExecuteTemplate(&b, "helper", e)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generate model %s\n", e.Name)
	}

	output = b.String()
	for {
		replaced := strings.ReplaceAll(output, "\n\n\n", "\n\n")
		if output == replaced {
			break
		}
		output = replaced
	}

	data, err = format.Source([]byte(output))
	if err != nil {
		log.Fatalln(err)
	}

	os.Mkdir("gen", 0755)
	err = os.WriteFile("gen/helper.txt", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}

}
