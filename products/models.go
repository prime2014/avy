package products

import (
	"db"
	"time"
)

type Categories struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Products struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Category    *Categories `json:"category"`
	Description string      `json:"description"`
	Stock       int         `json:"stock"`
	Quantity    float64     `json:"quantity"`
	Price       float64     `json:"price"`
	Discount    float64     `json:"discount"`
	Currency    string      `json:"currency"`
	PubDate     time.Time   `json:"pubdate"`
}

func (p Products) Save() Products {
	mydb := db.DBconnect()

	LastInsertID := 0

	err := mydb.QueryRow("INSERT INTO products(name, slug, description, stock, price, currency, discount) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		p.Name,
		p.Slug,
		p.Description,
		p.Stock,
		p.Price,
		p.Currency,
		p.Discount,
	).Scan(
		&LastInsertID,
	)

	if err != nil {
		panic(err)
	}

	p.ID = LastInsertID
	p.PubDate = time.Now().Local()
	return p
}
