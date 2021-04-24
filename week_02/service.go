package week_02

import "github.com/latesun/go-training-camp/week_02/model"

type Order interface {
	QueryByUserID(userID int) ([]model.Order, error)
	UpdateNameByUserID(name string, userID int) error
}
