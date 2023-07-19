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
	UnMarkDeleteAccount(id string) error
	DeleteAccount(id string) error
	RestoreAccount(id string) error

	GetDeletedAccounts(page uint64, size uint64) ([]*models.Account, error)
	CountDeletedAccounts() (uint64, error)

	// CarouselText's CRUD
	CreateCarouselText(ct *models.CarouselText) error
	GetCarouselText(id string) (*models.CarouselText, error)
	UpdateCarouselText(ct *models.CarouselText) error
	MarkDeleteCarouselText(id, by string) error
	UnMarkDeleteCarouselText(id string) error
	DeleteCarouselText(id string) error
	RestoreCarouselText(id string) error

	GetDeletedCarouselTexts(page uint64, size uint64) ([]*models.CarouselText, error)
	CountDeletedCarouselTexts() (uint64, error)

	// CarouselImage's CRUD
	CreateCarouselImage(ci *models.CarouselImage) error
	GetCarouselImage(id string) (*models.CarouselImage, error)
	UpdateCarouselImage(ci *models.CarouselImage) error
	MarkDeleteCarouselImage(id, by string) error
	UnMarkDeleteCarouselImage(id string) error
	DeleteCarouselImage(id string) error
	RestoreCarouselImage(id string) error

	GetDeletedCarouselImages(page uint64, size uint64) ([]*models.CarouselImage, error)
	CountDeletedCarouselImages() (uint64, error)

	// Category's CRUD
	CreateCategory(c *models.Category) error
	GetCategory(id string) (*models.Category, error)
	UpdateCategory(c *models.Category) error
	MarkDeleteCategory(id, by string) error
	UnMarkDeleteCategory(id string) error
	DeleteCategory(id string) error
	RestoreCategory(id string) error

	GetDeletedCategories(page uint64, size uint64) ([]*models.Category, error)
	CountDeletedCategories() (uint64, error)

	// Item's CRUD
	CreateItem(i *models.Item) error
	GetItem(id string) (*models.Item, error)
	UpdateItem(i *models.Item) error
	MarkDeleteItem(id, by string) error
	UnMarkDeleteItem(id string) error
	DeleteItem(id string) error
	RestoreItem(id string) error

	GetDeletedItems(page uint64, size uint64) ([]*models.Item, error)
	CountDeletedItems() (uint64, error)

	// Refill's CRUD
	CreateRefill(r *models.Refill) error
	GetRefill(id string) (*models.Refill, error)
	UpdateRefill(r *models.Refill) error
	MarkDeleteRefill(id, by string) error
	UnMarkDeleteRefill(id string) error
	DeleteRefill(id string) error
	RestoreRefill(id string) error

	GetDeletedRefills(page uint64, size uint64) ([]*models.Refill, error)
	CountDeletedRefills() (uint64, error)

	// Transaction's CRUD
	CreateTransaction(t *models.Transaction) error
	GetTransaction(id string) (*models.Transaction, error)
	UpdateTransaction(t *models.Transaction) error
	MarkDeleteTransaction(id, by string) error
	UnMarkDeleteTransaction(id string) error
	DeleteTransaction(id string) error
	RestoreTransaction(id string) error

	GetDeletedTransactions(page uint64, size uint64) ([]*models.Transaction, error)
	CountDeletedTransactions() (uint64, error)

	GetTransactions(account string, page uint64, size uint64, state string) ([]*models.Transaction, error)
	CountTransactions(account string, state string) (uint64, error)

	GetAllTransactions(page uint64, size uint64, state string) ([]*models.Transaction, error)
	CountAllTransactions(state string) (uint64, error)

	// Other requests that are not CRUD but still needed
	GetAccountByGoogle(googleID string) (*models.Account, error)
	GetAccountByCard(card string) (*models.Account, error)
	GetAccounts(page uint64, size uint64) ([]*models.Account, error)
	CountAccounts() (uint64, error)
	GetRefills(account string, page uint64, size uint64, startAt, endAt uint64) ([]*models.Refill, error)
	CountRefills(account string, startAt, endAt uint64) (uint64, error)
	GetItems(categoryID string, page, size uint64, state string) ([]*models.Item, error)
	CountItems(categoryID string, state string) (uint64, error)

	GetAllRefills(page uint64, size uint64, startAt, endAt uint64) ([]*models.Refill, error)
	CountAllRefills(startAt, endAt uint64) (uint64, error)
	GetAllCategories() ([]*models.Category, error)
	GetAllCarouselImages() ([]*models.CarouselImage, error)
	GetAllCarouselTexts() ([]*models.CarouselText, error)
}
