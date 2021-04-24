package model

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type Order struct {
	ID     int
	Name   string
	Age    int
	UserID int
	Remark string
}

func (o *Order) Query(db *sql.DB, query string, args ...interface{}) ([]Order, error) {
	orders := make([]Order, 0, 10)
	rows, err := db.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orders, nil
		}
		return nil, errors.Wrap(err, fmt.Sprintf("exec: '%s' %+v failed", query, args))
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err := rows.Scan(&o.ID, &o.Name, &o.Age)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan:")
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (o *Order) Update(db *sql.DB, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("exec: '%s' %v failed", query, args))
	}

	return nil
}
