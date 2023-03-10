{{ define `repo` }}package generate

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type {{.Name}} struct {
    db *sql.DB 
}

func New{{.Name}}(db *sql.DB) *{{.Name}} {
    r := new({{.Name}})
    r.db = db 
    return r
}

func (r *{{.Name}}) Insert(ctx context.Context, v *{{.Entity.Name}}) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO {{.Table}}({{.Columns}}) VALUES({{.Placeholders}})`, {{.Args}})
	return err 
}

func (r *{{.Name}}) InsertTx(ctx context.Context, tx *sql.Tx, v *{{.Entity.Name}}) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO {{.Table}}({{.Columns}}) VALUES({{.Placeholders}})`, {{.Args}})
	return err
}

func (r *{{.Name}}) Update(ctx context.Context, v *{{.Entity.Name}}) error {
	_, err := r.db.ExecContext(ctx, `UPDATE {{.Table}} SET {{.UpdateColumns}} WHERE {{.KeyConditions}}`, {{.Args}})
	return err
}

func (r *{{.Name}}) UpdateTx(ctx context.Context, tx *sql.Tx, v *{{.Entity.Name}}) error {
	_, err := tx.ExecContext(ctx, `UPDATE {{.Table}} SET {{.UpdateColumns}} WHERE {{.KeyConditions}}`, {{.Args}})
	return err
}

func (r *{{.Name}}) Save(ctx context.Context, v *{{.Entity.Name}}) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO {{.Table}}({{.Columns}}) VALUES({{.Placeholders}}) ON CONFLICT({{.Keys}}) DO UPDATE SET {{.UpdateColumns}}`, {{.Args}})
	return err
}

func (r *{{.Name}}) SaveTx(ctx context.Context, tx *sql.Tx, v *{{.Entity.Name}}) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO {{.Table}}({{.Columns}}) VALUES({{.Placeholders}}) ON CONFLICT({{.Keys}}) DO UPDATE SET {{.UpdateColumns}}`, {{.Args}})
	return err
}

func (r *{{.Name}}) Delete(ctx context.Context, {{.KeyParams}}) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
	return err
}

func (r *{{.Name}}) DeleteTx(ctx context.Context, tx *sql.Tx, {{.KeyParams}}) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
	return err
}

func (r *{{.Name}}) Get(ctx context.Context, {{.KeyParams}}) (v *{{.Entity.Name}}, err error) {
    v = new({{.Entity.Name}})
	row := r.db.QueryRowContext(ctx, `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
	err = row.Scan({{.ScanHolders}})
	if err != nil {
	    return nil, err
	}
	return v, nil
}

func (r *{{.Name}}) GetTx(ctx context.Context, tx *sql.Tx, {{.KeyParams}}) (v *{{.Entity.Name}}, err error) {
    v = new({{.Entity.Name}})
	row := tx.QueryRowContext(ctx, `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.KeyConditions}}`, {{.KeyArgs}})
	err = row.Scan({{.ScanHolders}})
	if err != nil {
	    return nil, err
	}
	return v, nil
}

{{if eq .NumKeys 1}}
func (r *{{.Name}}) BatchGet(ctx context.Context, {{.BatchKeyParams}}) (list []*{{.Entity.Name}}, err error) {
	rows, err :=  r.db.QueryContext(ctx, `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.BatchKeyConditions}}`, pq.Array({{.BatchKeyArgs}}))
	if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        v := new({{.Entity.Name}})
        err = rows.Scan({{.ScanHolders}})
        if err != nil {
            return nil, err
        }
        list = append(list, v)
    }
	return list, nil
}

func (r *{{.Name}}) BatchGetTx(ctx context.Context, tx *sql.Tx, {{.BatchKeyParams}}) (list []*{{.Entity.Name}}, err error) {
	rows, err :=  tx.QueryContext(ctx, `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.BatchKeyConditions}}`, pq.Array({{.BatchKeyArgs}}))
	if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        v := new({{.Entity.Name}})
        err = rows.Scan({{.ScanHolders}})
        if err != nil {
            return nil, err
        }
        list = append(list, v)
    }
	return list, nil
}

func (r *{{.Name}}) BatchDelete(ctx context.Context, {{.BatchKeyParams}}) error {
	_, err :=  r.db.ExecContext(ctx, `DELETE FROM {{.Table}} WHERE {{.BatchKeyConditions}}`, pq.Array({{.BatchKeyArgs}}))
	return err
}

func (r *{{.Name}}) BatchDeleteTx(ctx context.Context, tx *sql.Tx, {{.BatchKeyParams}}) error {
	_, err :=  tx.ExecContext(ctx, `DELETE FROM {{.Table}} WHERE {{.BatchKeyConditions}}`, pq.Array({{.BatchKeyArgs}}))
	return err
}

{{end}}

func (r *{{.Name}}) List(ctx context.Context) (list []*{{.Entity.Name}}, err error) {
	rows, err := r.db.QueryContext(ctx, `SELECT {{.Columns}} FROM {{.Table}}`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        v := new({{.Entity.Name}})
        err = rows.Scan({{.ScanHolders}})
        if err != nil {
            return nil, err
        }
        list = append(list, v)
    }
	return list, nil
}

func (r *{{.Name}}) ListTx(ctx context.Context, tx *sql.Tx) (list []*{{.Entity.Name}}, err error) {
	rows, err := tx.QueryContext(ctx, `SELECT {{.Columns}} FROM {{.Table}}`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        v := new({{.Entity.Name}})
        err = rows.Scan({{.ScanHolders}})
        if err != nil {
            return nil, err
        }
        list = append(list, v)
    }
	return list, nil
}

{{end}}