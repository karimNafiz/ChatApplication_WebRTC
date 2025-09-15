package data

import (
	"database/sql"
	"time"
)

type Client struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}

type ClientModel struct {
	db *sql.DB
}

func (c ClientModel) Insert(client *Client) error {
	return nil
}

func (c ClientModel) Get(id string) (*Client, error) {
	return nil, nil
}

func (c ClientModel) Update(client *Client) error {
	return nil
}
func (c ClientModel) Delete(id string) error {
	return nil
}
