package domain

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	//m := parseModel("testdata/model.xml")
	//b, _ := json.Marshal(m)
	//t.Log(string(b))
	//
	Generate("testdata/model.xml", "test.gen.go")
}
