package storage

import (
	"context"
)

type OrderPair struct {
	OrderUID string
	Data     string
}

func (s *Storage) Create(orderUID, data string) error {
	_, err := s.Db.Exec(context.TODO(), `
		insert into orders (order_uid, data)
		values ($1, $2)`,
		orderUID, data,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get(uid string) (string, error) {
	var data string
	err := s.Db.QueryRow(context.TODO(), `
		select data
		from orders
		where order_uid = $1
	`, uid).Scan(
		&data,
	)

	if err != nil {
		return "", err
	}
	return data, nil

}

func (s *Storage) GetAll() []OrderPair {
	rows, err := s.Db.Query(context.TODO(), `
		select order_uid, data
		from orders
	`)
	if err != nil {
		return []OrderPair{}
	}

	result := make([]OrderPair, 0)

	for rows.Next() {
		var orderStorage OrderPair
		err = rows.Scan(
			&orderStorage.OrderUID,
			&orderStorage.Data,
		)
		if err != nil {
			return []OrderPair{}
		}
		result = append(result, orderStorage)
	}

	return result
}
