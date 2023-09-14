package domain

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	//m := parseModel("testdata/domain.xml")
	//b, _ := json.Marshal(m)
	//t.Log(string(b))
	//
	Generate("testdata/domain.xml", "testdata/domain.gen.go")
}
