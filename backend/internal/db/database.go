package db

import (
	"bar/internal/models"
	"time"
)

type DatabaseOptions struct {
	ConnString   string
	Database     string
	QueryTimeout time.Duration
}

func NewDatabaseOptions(connString, database string, timeout time.Duration) *DatabaseOptions {
	return &DatabaseOptions{
		Database:     database,
		ConnString:   connString,
		QueryTimeout: timeout,
	}
}

type DBackend interface {
	Connect() error
	Disconnect() error

	// Account's CRUD
	CreateAccount(acc *models.Account) error
	GetAccount(id string) (*models.Account, error)
	UpdateAccount(acc *models.Account) error
	DeleteAccount(id, by string) error

	// CarouselText's CRUD
	CreateCarouselText(ct *models.CarouselText) error
	GetCarouselText(id string) (*models.CarouselText, error)
	UpdateCarouselText(ct *models.CarouselText) error
	DeleteCarouselText(id, by string) error

	// CarouselImage's CRUD
	CreateCarouselImage(ci *models.CarouselImage) error
	GetCarouselImage(id string) (*models.CarouselImage, error)
	UpdateCarouselImage(ci *models.CarouselImage) error
	DeleteCarouselImage(id, by string) error

	// Category's CRUD
	CreateCategory(c *models.Category) error
	GetCategory(id string) (*models.Category, error)
	UpdateCategory(c *models.Category) error
	DeleteCategory(id, by string) error

	// Item's CRUD
	CreateItem(i *models.Item) error
	GetItem(id string) (*models.Item, error)
	UpdateItem(i *models.Item) error
	DeleteItem(id, by string) error

	// Refill's CRUD
	CreateRefill(r *models.Refill) error
	GetRefill(id string) (*models.Refill, error)
	UpdateRefill(r *models.Refill) error
	DeleteRefill(id, by string) error

	// Transaction's CRUD
	CreateTransaction(t *models.Transaction) error
	GetTransaction(id string) (*models.Transaction, error)
	UpdateTransaction(t *models.Transaction) error
	DeleteTransaction(id, by string) error

	// Other requests that are not CRUD but still needed
	GetAccounts(page int, size int) ([]*models.Account, error)
	CountAccounts() (int64, error)
	GetTransactions(account string, page int, size int, state string) ([]*models.Transaction, error)
	GetRefills(account string, page int, size int, startAt, endAt int64) ([]*models.Refill, error)
	GetAllRefills(page int, size int, startAt, endAt int64) ([]*models.Refill, error)
	GetAllCategories() ([]*models.Category, error)
	GetAllCarouselImages() ([]*models.CarouselImage, error)
	GetAllCarouselTexts() ([]*models.CarouselText, error)
}
