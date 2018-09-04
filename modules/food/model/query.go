package model

import (
	"fmt"
	"strings"
)

func (m *Manager) CreateRestaurant(r *Restaurant) error {
	q := fmt.Sprintf(`INSERT INTO %s (name,tax,send,locations,packing,created_at) VALUES (?,?,?,?,?,?)`, RestaurantTableName)
	res, err := m.GetConn().Exec(q, r.Name, r.Tax, r.Send, r.Locations, r.Packing, r.CreatedAt)
	if err != nil {
		return err
	}
	num, _ := res.LastInsertId()
	r.ID = num
	return nil
}

func (m *Manager) CreateFood(r *Food) error {
	q := fmt.Sprintf(`INSERT INTO %s (title,price,resturant_id) VALUES (?,?,?)`, FoodTableName)
	res, err := m.GetConn().Exec(q, r.Title, r.Price, r.RestaurantID)
	if err != nil {
		return err
	}
	num, _ := res.LastInsertId()
	r.ID = num
	return nil
}

func (m *Manager) CreateOrder(r *Order) error {
	q := fmt.Sprintf(`INSERT INTO %s (restaurant_id,food_ids,price,created_at,updated_at) VALUES (?,?,?,?,?)`, OrderTableName)
	res, err := m.GetConn().Exec(q, r.RestaurantID,r.FoodIDs, r.Price, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return err
	}
	num, _ := res.LastInsertId()
	r.ID = num
	return nil
}

func (m *Manager) UpdateRestaurant(r *Restaurant) error {
	q := fmt.Sprintf(`UPDATE %s SET name=?,tax=?,send=?,locations=?,packing=? WHERE id=?`, RestaurantTableName)
	_, err := m.GetConn().Exec(q, r.Name, r.Tax, r.Send, r.Locations, r.Packing, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateFood(r *Food) error {
	q := fmt.Sprintf(`UPDATE %s SET title=?,price=? WHERE id=?`, FoodTableName)
	_, err := m.GetConn().Exec(q, r.Title, r.Price, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) FindRestaurant(id int64) (*Restaurant, error) {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id=?`, RestaurantTableName)
	var r Restaurant
	err := m.GetConn().SelectOne(&r, q, id)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (m *Manager) FindFood(id int64) (*Food, error) {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id=?`, FoodTableName)
	var r Food
	err := m.GetConn().SelectOne(&r, q, id)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (m *Manager) ListRestaurant() []*Restaurant {
	q := fmt.Sprintf(`SELECT * FROM %s`, RestaurantTableName)
	var r []*Restaurant
	_, err := m.GetConn().Select(&r, q)
	if err != nil {
		panic(err)
	}
	return r
}

func (m *Manager) DeleteRestaurant(id int64) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id=?`, RestaurantTableName)
	_, err := m.GetConn().Exec(q, id)
	if err != nil {
		panic(err)
	}
	return nil
}

func (m *Manager) FindFoods(rID int64, fIDs []int64) []*Food {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE resturant_id=? AND id in (%s)`,
		FoodTableName,
		func() string {
			return strings.TrimRight(strings.Repeat("?,", len(fIDs)), ",")
		}(),
	)

	var params = make([]interface{}, 0)
	params = append(params, rID)
	for i := range fIDs {
		params = append(params, fIDs[i])
	}

	var r []*Food
	_, err := m.GetConn().Select(&r, q, params...)
	if err != nil {
		panic(err)
	}
	return r
}
