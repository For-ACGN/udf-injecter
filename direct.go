package injecter

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type direct struct {
	*sql.DB
}

func Connect(address, username, password string) (Handle, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/mysql", username, password, address)
	return ConnectWithDSN(dsn)
}

func ConnectWithDSN(dsn string) (Handle, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return &direct{DB: db}, nil
}

func (d *direct) Query(query string) ([]map[string]string, error) {
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnsLen := len(columns)
	cache := make([]interface{}, columnsLen)
	for i := 0; i < columnsLen; i++ {
		var a sql.RawBytes
		cache[i] = &a
	}
	var list []map[string]string
	for rows.Next() {
		row := make(map[string]string)
		err = rows.Scan(cache...)
		if err != nil {
			return nil, err
		}
		for i := 0; i < columnsLen; i++ {
			row[columns[i]] = string(*cache[i].(*sql.RawBytes))
		}
		list = append(list, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("no result")
	}
	return list, nil
}

func (d *direct) Exec(query string) error {
	_, err := d.DB.Exec(query)
	return err
}

func (d *direct) Close() {
	_ = d.DB.Close()
}
