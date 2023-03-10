{{ define `repotest` }}package generate

import (
	"context"
    "testing"

	"github.com/lib/pq"
)

func setupTest{{.Name}}(t *testing.T) *{{.Name}} {
    db := setupTestDB(t)
    _, err := db.Exec(`TRUNCATE TABLE {{.Table}}`)
    if err != nil {
        t.Error(err)
    }
    return New{{.Name}}(db)
}

func newTest{{.Entity.Name}}() *{{.Entity.Name}} {
    return new({{.Entity.Name}})
}

func TestInsert{{.Entity.Name}}(t *testing.T) {
    ctx := context.TODO()
    v := newTest{{.Entity.Name}}()
    r := setupTest{{.Name}}(t)
    err := r.Insert(ctx, v)
    if err != nil {
        t.Error(err)
    }
}

func TestUpdate{{.Entity.Name}}(t *testing.T) {
    ctx := context.TODO()
    v := newTest{{.Entity.Name}}()
    r := setupTest{{.Name}}(t)
    err := r.Insert(ctx, v)
    if err != nil {
        t.Error(err)
    }
    // TODO:
    // update columns
    err = r.Update(ctx, v)
    if err != nil {
        t.Error(err)
    }
}

func TestSave{{.Entity.Name}}(t *testing.T) {
    ctx := context.TODO()
    v := newTest{{.Entity.Name}}()
    r := setupTest{{.Name}}(t)
    err := r.Save(ctx, v)
    if err != nil {
        t.Error(err)
    }
    // TODO:
    // update columns
    err = r.Save(ctx, v)
    if err != nil {
        t.Error(err)
    }

//    got, err := r.Get(ctx, {{.KeyArgs}})
//    if err != nil {
//        t.Error(err)
//    }
//    t.Log(got)
}

func TestBatchGet{{.Entity.Name}}(t *testing.T) {
    // TODO:
}

func TestList{{.Entity.Name}}(t *testing.T) {
    // TODO:
}

{{end}}