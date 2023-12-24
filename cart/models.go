package cart

import (
	"accounts"
	"db"
	"fmt"
	"products"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Cart struct {
	ID           int                `json:"id"`
	Customer     *accounts.Users    `json:"customer"`
	Item         *products.Products `json:"item"`
	Quantity     int                `json:"quantity"`
	Price        float64            `json:"price"`
	NetPrice     float64            `json:"net_price"`
	Currency     string             `json:"currency"`
	ModifiedDate time.Time          `json:"modified_date"`
	PubDate      time.Time          `json:"pubdate"`
}

func (c Cart) Save(customer string) (*Cart, error) {
	mydb := db.DBconnect()

	LastInsertID := 0

	var user = &accounts.Users{}

	err := mydb.QueryRow("SELECT id, firstname, lastname, email, is_active FROM users WHERE email=$1", customer).Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.IsActive,
	)

	if err != nil {
		return &Cart{}, err
	}

	myerr := mydb.QueryRow("INSERT INTO cart(customer, item, price, quantity, currency, modified_date) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		user.ID,
		c.Item,
		c.Price,
		c.Quantity,
		c.Currency,
		c.ModifiedDate,
	).Scan(&LastInsertID)

	if myerr != nil {
		return &Cart{}, err
	}

	c.ID = LastInsertID
	c.Customer = user
	return &c, nil
}

type CartRequest struct {
	ID           int       `json:"id"`
	Customer     int       `json:"customer"`
	Item         int       `json:"item"`
	Quantity     int       `json:"quantity"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	ModifiedDate time.Time `json:"modified_date"`
	PubDate      time.Time `json:"pubdate"`
}

func (crt CartRequest) Validate() error {
	return validation.ValidateStruct(&crt,
		validation.Field(&crt.Currency, validation.Required),
		validation.Field(&crt.Item, validation.Required),
		validation.Field(&crt.Quantity),
		validation.Field(&crt.Price, validation.Required),
	)
}

func (crt CartRequest) Save(customer int) (Cart, error) {
	// connect to the database
	mydb := db.DBconnect()
	LastInsertID := 0
	var user accounts.Users
	var product products.Products

	err := mydb.QueryRow("SELECT id, firstname, lastname, email, is_active FROM users WHERE id=$1", customer).Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.IsActive,
	)

	if err != nil {
		return Cart{}, err
	}

	err = mydb.QueryRow("SELECT id, name, slug, category, description, stock, quantity, price, discount, currency FROM products WHERE id=$1", crt.Item).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Category,
		&product.Description,
		&product.Stock,
		&product.Quantity,
		&product.Price,
		&product.Discount,
		&product.Currency,
	)

	if err != nil {
		return Cart{}, err
	}

	if validation.IsEmpty(crt.Quantity) {
		crt.Quantity += 1

		err = mydb.QueryRow("INSERT INTO cart(customer, item, quantity, price, net_price, currency) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
			&customer,
			&crt.Item,
			&crt.Quantity,
			&crt.Price,
			&crt.Price,
			&crt.Currency,
		).Scan(&LastInsertID)

		if err != nil {
			fmt.Printf("HERE LIES THE ERROR\n")
			return Cart{}, err
		}
		return Cart{
			ID:       LastInsertID,
			Customer: &user,
			Item:     &product,
			Price:    crt.Price,
			NetPrice: crt.Price,
			Currency: crt.Currency,
			Quantity: crt.Quantity,
			PubDate:  time.Now().Local(),
		}, nil
	} else {
		net_price := float64(crt.Quantity) * crt.Price
		err = mydb.QueryRow("INSERT INTO cart(customer, item, quantity, price, net_price, currency) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
			&customer,
			&crt.Item,
			&crt.Quantity,
			&crt.Price,
			&net_price,
			&crt.Currency,
		).Scan(&LastInsertID)

		if err != nil {
			return Cart{}, err
		}
		return Cart{
			ID:       LastInsertID,
			Customer: &user,
			Item:     &product,
			Price:    crt.Price,
			NetPrice: crt.Price,
			Currency: crt.Currency,
			Quantity: crt.Quantity,
			PubDate:  time.Now().Local(),
		}, nil
	}
}
