package injecter

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Handle interface {
	// mysql> show variables like '%version_%';
	// +-------------------------+------------------------------+
	// | Variable_name           | Value                        |
	// +-------------------------+------------------------------+
	// | slave_type_conversions  |                              |
	// | version_comment         | MySQL Community Server - GPL |
	// | version_compile_machine | x86_64                       |
	// | version_compile_os      | Win64                        |
	// | version_compile_zlib    | 1.2.11                       |
	// +-------------------------+------------------------------+
	// result, _ := Query("show variables like '%version_%'")
	// result[2]["Variable_name"] = "x86_64"
	Query(query string, args ...interface{}) ([]map[string]string, error)
}

type handle struct {
	db *sql.DB
}

func connect(address, username, password string) (*handle, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/mysql", username, password, address)
	return connectWithDSN(dsn)
}

func connectWithDSN(dsn string) (*handle, error) {
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = d.Ping()
	if err != nil {
		return nil, err
	}
	d.SetConnMaxLifetime(time.Minute)
	d.SetMaxOpenConns(1)
	d.SetMaxIdleConns(1)
	return &handle{db: d}, nil
}

func (h *handle) Query(query string, args ...interface{}) ([]map[string]string, error) {
	rows, err := h.db.Query(query, args...)
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
	return list, nil
}
