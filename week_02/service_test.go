package week_02

import (
	"testing"
)

func TestQueryAndUpdate(t *testing.T) {
	t.Run("query_by_user_id", func(t *testing.T) {
		mocker, err := NewQueryMocker()
		if err != nil {
			t.Fatalf("Init db failed: %+v\n", err)
		}
		defer mocker.Close()

		orders, err := mocker.QueryByUserID(2333)
		if err != nil {
			t.Fatalf("\nQuery failed: %+v\n", err)
		}

		t.Log("result of query by user_id:", orders)
	})

	t.Run("update_name_by_user_id", func(t *testing.T) {
		mocker, err := NewUpdateMocker()
		if err != nil {
			t.Fatalf("Init db failed: %+v\n", err)
		}
		defer mocker.Close()

		name, userID := "latesun", 2333
		err = mocker.UpdateNameByUserID(name, userID)
		if err != nil {
			t.Fatalf("\nUpdate failed: %+v\n", err)
		}
		t.Log("Update succeeded.")
	})
}
