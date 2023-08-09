package model_test

import (
	"testing"

	"create-model/model"
)

func TestGenerate(t *testing.T) {
	model.Generate("testdata/model.json")
}
