package week_02

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"

	"github.com/latesun/go-training-camp/week_02/model"
)

type mocker struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func NewQueryMocker() (*mocker, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, errors.Wrap(err, "sqlmock.New")
	}
	mock.ExpectQuery("SELECT id, name, age FROM orders WHERE user_id = ?").
		WillReturnError(sql.ErrNoRows)

	return &mocker{db: db, mock: mock}, nil
}

func NewUpdateMocker() (*mocker, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, errors.Wrap(err, "sqlmock.New")
	}

	mock.ExpectExec("UPDATE orders").
		WithArgs("latesun", 2333).
		WillReturnError(sql.ErrNoRows)

	return &mocker{db: db, mock: mock}, nil
}

func (m *mocker) Close() {
	m.db.Close()
}

func (m *mocker) QueryByUserID(userID int) ([]model.Order, error) {
	query := "SELECT id, name, age FROM orders WHERE user_id = ?"
	res, err := (&model.Order{}).Query(m.db, query, userID)
	if err != nil {
		return nil, errors.WithMessage(err, "mocker could not query by user_id")
	}

	if err := m.mock.ExpectationsWereMet(); err != nil {
		return nil, errors.WithMessage(err, "query result was unfulfilled")
	}
	return res, nil
}

func (m *mocker) UpdateNameByUserID(name string, userID int) error {
	query := "UPDATE orders SET name = ? WHERE user_id = ?"
	err := (&model.Order{}).Update(m.db, query, name, userID)
	if err != nil {
		return errors.WithMessage(err, "mocker could not update name by user_id")
	}

	if err := m.mock.ExpectationsWereMet(); err != nil {
		return errors.WithMessage(err, "update result was unfulfilled")
	}
	return nil
}
