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
	MarkDeleteAccount(id, by string) error
	DeleteAccount(id string) error
	RestoreAccount(id string) error

	// CarouselText's CRUD
	CreateCarouselText(ct *models.CarouselText) error
	GetCarouselText(id string) (*models.CarouselText, error)
	UpdateCarouselText(ct *models.CarouselText) error
	MarkDeleteCarouselText(id, by string) error
	DeleteCarouselText(id string) error
	RestoreCarouselText(id string) error

	// CarouselImage's CRUD
	CreateCarouselImage(ci *models.CarouselImage) error
	GetCarouselImage(id string) (*models.CarouselImage, error)
	UpdateCarouselImage(ci *models.CarouselImage) error
	MarkDeleteCarouselImage(id, by string) error
	DeleteCarouselImage(id string) error
	RestoreCarouselImage(id string) error

	// Category's CRUD
	CreateCategory(c *models.Category) error
	GetCategory(id string) (*models.Category, error)
	UpdateCategory(c *models.Category) error
	MarkDeleteCategory(id, by string) error
	DeleteCategory(id string) error
	RestoreCategory(id string) error

	// Item's CRUD
	CreateItem(i *models.Item) error
	GetItem(id string) (*models.Item, error)
	UpdateItem(i *models.Item) error
	MarkDeleteItem(id, by string) error
	DeleteItem(id string) error
	RestoreItem(id string) error

	// Refill's CRUD
	CreateRefill(r *models.Refill) error
	GetRefill(id string) (*models.Refill, error)
	UpdateRefill(r *models.Refill) error
	MarkDeleteRefill(id, by string) error
	DeleteRefill(id string) error
	RestoreRefill(id string) error

	// Transaction's CRUD
	CreateTransaction(t *models.Transaction) error
	GetTransaction(id string) (*models.Transaction, error)
	UpdateTransaction(t *models.Transaction) error
	MarkDeleteTransaction(id, by string) error
	DeleteTransaction(id string) error
	RestoreTransaction(id string) error

	// Other requests that are not CRUD but still needed
	GetAccountByCard(card string) (*models.Account, error)
	GetAccounts(page int, size int) ([]*models.Account, error)
	CountAccounts() (int64, error)
	GetTransactions(account string, page int, size int, state string) ([]*models.Transaction, error)
	GetRefills(account string, page int, size int, startAt, endAt int64) ([]*models.Refill, error)
	GetAllRefills(page int, size int, startAt, endAt int64) ([]*models.Refill, error)
	GetAllCategories() ([]*models.Category, error)
	GetAllCarouselImages() ([]*models.CarouselImage, error)
	GetAllCarouselTexts() ([]*models.CarouselText, error)
}
