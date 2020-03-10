package burgers

import (
	"context"
	"crud/pkg/crud/models"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BurgersSvc struct {
	pool *pgxpool.Pool // dependency
}

func NewBurgersSvc(pool *pgxpool.Pool) *BurgersSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil")) // <- be accurate
	}
	return &BurgersSvc{pool: pool}
}

// func BurgersList(...)
// func BurgersListWithContext(...)

func (service *BurgersSvc) BurgersList(ctx context.Context) (list []models.Burger, err error) {
	list = make([]models.Burger, 0) // TODO: for REST API
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		return nil, err // TODO: wrap to specific error
	}
	defer conn.Release()
	rows, err := conn.Query(ctx, "SELECT id, name, price, fileName FROM burgers WHERE removed = FALSE")
	if err != nil {
		return nil, err // TODO: wrap to specific error
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Burger{}
		err := rows.Scan(&item.Id, &item.Name, &item.Price, &item.FileName)
		if err != nil {
			return nil, err // TODO: wrap to specific error
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (service *BurgersSvc) Save(ctx context.Context, model models.Burger) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return errors.New("can't execute pool: ")
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), "INSERT INTO burgers(name, price, fileName) VALUES ($1, $2, $3);", model.Name, model.Price, model.FileName)
	if err != nil {
		return errors.New("can't save burger: ")
	}
	return nil
}

func (service *BurgersSvc) RemoveById(ctx context.Context, id int) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return errors.New("can't execute pool: ")
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), "UPDATE burgers SET removed = true where id = $1;", id)
	if err != nil {
		return errors.New("can't remove burger: ")
	}
	return nil
}
