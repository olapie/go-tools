package main

import (
	"fmt"
	"strings"

	"go.olapie.com/utils"
	"gopkg.in/yaml.v2"
)

type Field struct {
	Name string
	Type string
}

type Entity struct {
	Name   string
	Fields []*Field
}

type RepoModel struct {
	Name        string        `yaml:"name"`
	Table       string        `yaml:"table"`
	PrimaryKey  []string      `yaml:"primaryKey"`
	Columns     yaml.MapSlice `yaml:"columns"`
	JsonColumns []string      `yaml:"jsonColumns"`
}

func (r *RepoModel) IsKey(col string) bool {
	for _, v := range r.PrimaryKey {
		if v == col {
			return true
		}
	}
	return false
}

func (r *RepoModel) IsJSON(col string) bool {
	for _, v := range r.JsonColumns {
		if v == col {
			return true
		}
	}
	return false
}

func (r *RepoModel) IsArray(col string) bool {
	return strings.Index(col, "[]") == 0 && col != "[]byte"
}

func (r *RepoModel) Args() string {
	args := make([]string, len(r.Columns))
	for i, c := range r.Columns {
		name := c.Key.(string)
		value := c.Value.(string)
		args[i] = "v." + utils.ToClassName(name)
		if r.IsJSON(name) {
			args[i] = "xsql.JSON(" + args[i] + ")"
		} else if r.IsArray(value) {
			args[i] = "pq.Array(" + args[i] + ")"
		}
	}
	return strings.Join(args, ", ")
}

func (r *RepoModel) Placeholders() string {
	placeholders := make([]string, len(r.Columns))
	for i := range r.Columns {
		placeholders[i] = "$" + fmt.Sprint(i+1)
	}
	return strings.Join(placeholders, ", ")
}

func (r *RepoModel) UpdateColumns() string {
	updates := make([]string, 0, len(r.Columns))
	for i, c := range r.Columns {
		name := c.Key.(string)
		if !r.IsKey(name) {
			updates = append(updates, fmt.Sprintf("%s=$%d", name, i+1))
		}
	}

	return strings.Join(updates, ", ")
}

func (r *RepoModel) KeyConditions() string {
	keys := make([]string, 0, len(r.PrimaryKey))
	for i, c := range r.Columns {
		name := c.Key.(string)
		if r.IsKey(name) {
			keys = append(keys, fmt.Sprintf("%s=$%d", name, i+1))
		}
	}

	return strings.Join(keys, " AND ")
}

func (r *RepoModel) GetColType(col string) string {
	for _, c := range r.Columns {
		if c.Key == col {
			return c.Value.(string)
		}
	}
	return ""
}

func (r *RepoModel) KeyParams() string {
	params := make([]string, len(r.PrimaryKey))
	for i, k := range r.PrimaryKey {
		params[i] = fmt.Sprintf("%s %v", utils.ToCamel(k), r.GetColType(k))
	}
	return strings.Join(params, ", ")
}

func (r *RepoModel) BatchKeyParams() string {
	if len(r.PrimaryKey) != 1 {
		return ""
	}
	k := r.PrimaryKey[0]
	return fmt.Sprintf("%s []%v", utils.ToCamel(utils.Plural(k)), r.GetColType(k))
}

func (r *RepoModel) BatchKeyConditions() string {
	if len(r.PrimaryKey) != 1 {
		return ""
	}
	return fmt.Sprintf("%s=ANY($%d)", r.PrimaryKey[0], len(r.Columns)+1)
}

func (r *RepoModel) BatchKeyArgs() string {
	if len(r.PrimaryKey) != 1 {
		return ""
	}
	return strings.Split(r.BatchKeyParams(), " ")[0]
}

func (r *RepoModel) GetKeys() string {
	return strings.Join(r.PrimaryKey, ", ")
}

func (r *RepoModel) KeyArgs() string {
	args := make([]string, len(r.PrimaryKey))
	for i, k := range r.PrimaryKey {
		args[i] = utils.ToCamel(k)
	}
	return strings.Join(args, ", ")
}

func (r *RepoModel) ScanHolders() string {
	scanArgs := make([]string, len(r.Columns))
	for i, c := range r.Columns {
		name := c.Key.(string)
		value := c.Value.(string)
		scanArgs[i] = "&v." + utils.ToClassName(name)
		if r.IsJSON(name) {
			scanArgs[i] = "xsql.JSON(" + scanArgs[i] + ")"
		} else if r.IsArray(value) {
			scanArgs[i] = "pq.Array(" + scanArgs[i] + ")"
		}
	}

	return strings.Join(scanArgs, ", ")
}

func (r *RepoModel) GetColumns() string {
	cols := make([]string, len(r.Columns))
	for i := range r.Columns {
		cols[i] = r.Columns[i].Key.(string)
	}
	return strings.Join(cols, ", ")
}
