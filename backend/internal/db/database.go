package db

import (
	"bar/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	CreateAccount(ctx context.Context, acc *models.Account) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	UpdateAccount(ctx context.Context, acc *models.Account) error
	MarkDeleteAccount(ctx context.Context, id, by string) error
	UnMarkDeleteAccount(ctx context.Context, id string) error
	DeleteAccount(ctx context.Context, id string) error
	RestoreAccount(ctx context.Context, id string) error

	GetDeletedAccounts(ctx context.Context, page uint64, size uint64, query string) ([]*models.Account, error)
	CountDeletedAccounts(ctx context.Context, query string) (uint64, error)
	ListenForChanges(ctx context.Context, coll string, fn func(*mongo.ChangeStream)) error

	// CarouselText's CRUD
	CreateCarouselText(ctx context.Context, ct *models.CarouselText) error
	GetCarouselText(ctx context.Context, id string) (*models.CarouselText, error)
	UpdateCarouselText(ctx context.Context, ct *models.CarouselText) error
	MarkDeleteCarouselText(ctx context.Context, id, by string) error
	UnMarkDeleteCarouselText(ctx context.Context, id string) error
	DeleteCarouselText(ctx context.Context, id string) error
	RestoreCarouselText(ctx context.Context, id string) error

	GetDeletedCarouselTexts(ctx context.Context, page uint64, size uint64) ([]*models.CarouselText, error)
	CountDeletedCarouselTexts(ctx context.Context) (uint64, error)

	// CarouselImage's CRUD
	CreateCarouselImage(ctx context.Context, ci *models.CarouselImage) error
	GetCarouselImage(ctx context.Context, id string) (*models.CarouselImage, error)
	UpdateCarouselImage(ctx context.Context, ci *models.CarouselImage) error
	MarkDeleteCarouselImage(ctx context.Context, id, by string) error
	UnMarkDeleteCarouselImage(ctx context.Context, id string) error
	DeleteCarouselImage(ctx context.Context, id string) error
	RestoreCarouselImage(ctx context.Context, id string) error

	GetDeletedCarouselImages(ctx context.Context, page uint64, size uint64) ([]*models.CarouselImage, error)
	CountDeletedCarouselImages(ctx context.Context) (uint64, error)

	// Category's CRUD
	CreateCategory(ctx context.Context, c *models.Category) error
	GetCategory(ctx context.Context, id string) (*models.Category, error)
	UpdateCategory(ctx context.Context, c *models.Category) error
	MarkDeleteCategory(ctx context.Context, id, by string) error
	UnMarkDeleteCategory(ctx context.Context, id string) error
	DeleteCategory(ctx context.Context, id string) error
	RestoreCategory(ctx context.Context, id string) error

	GetDeletedCategories(ctx context.Context, page uint64, size uint64) ([]*models.Category, error)
	CountDeletedCategories(ctx context.Context) (uint64, error)

	// Item's CRUD
	CreateItem(ctx context.Context, i *models.Item) error
	GetItem(ctx context.Context, id string) (*models.Item, error)
	UpdateItem(ctx context.Context, i *models.Item) error
	MarkDeleteItem(ctx context.Context, id, by string) error
	UnMarkDeleteItem(ctx context.Context, id string) error
	DeleteItem(ctx context.Context, id string) error
	RestoreItem(ctx context.Context, id string) error

	GetDeletedItems(ctx context.Context, page uint64, size uint64) ([]*models.Item, error)
	CountDeletedItems(ctx context.Context) (uint64, error)

	// Refill's CRUD
	CreateRefill(ctx context.Context, r *models.Refill) error
	GetRefill(ctx context.Context, id string) (*models.Refill, error)
	UpdateRefill(ctx context.Context, r *models.Refill) error
	MarkDeleteRefill(ctx context.Context, id, by string) error
	UnMarkDeleteRefill(ctx context.Context, id string) error
	DeleteRefill(ctx context.Context, id string) error

	GetDeletedRefills(ctx context.Context, page uint64, size uint64) ([]*models.Refill, error)
	CountDeletedRefills(ctx context.Context) (uint64, error)

	// Transaction's CRUD
	CreateTransaction(ctx context.Context, t *models.Transaction) error
	GetTransaction(ctx context.Context, id string) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, t *models.Transaction) error
	MarkDeleteTransaction(ctx context.Context, id, by string) error
	UnMarkDeleteTransaction(ctx context.Context, id string) error
	DeleteTransaction(ctx context.Context, id string) error
	RestoreTransaction(ctx context.Context, id string) error

	GetDeletedTransactions(ctx context.Context, page uint64, size uint64) ([]*models.Transaction, error)
	CountDeletedTransactions(ctx context.Context) (uint64, error)

	GetTransactions(ctx context.Context, account string, page uint64, size uint64, state string) ([]*models.Transaction, error)
	CountTransactions(ctx context.Context, account string, state string) (uint64, error)

	GetAllTransactions(ctx context.Context, page uint64, size uint64, state string, name string) ([]*models.Transaction, error)
	CountAllTransactions(ctx context.Context, state string, name string) (uint64, error)

	// Restock's CRUD
	CreateRestock(ctx context.Context, t *models.Restock) error
	GetRestock(ctx context.Context, id string) (*models.Restock, error)
	UpdateRestock(ctx context.Context, t *models.Restock) error
	MarkDeleteRestock(ctx context.Context, id, by string) error
	UnMarkDeleteRestock(ctx context.Context, id string) error
	DeleteRestock(ctx context.Context, id string) error
	RestoreRestock(ctx context.Context, id string) error

	GetDeletedRestocks(ctx context.Context, page uint64, size uint64) ([]*models.Restock, error)
	CountDeletedRestocks(ctx context.Context) (uint64, error)

	GetRestocks(ctx context.Context, account string, page uint64, size uint64) ([]*models.Restock, error)
	CountRestocks(ctx context.Context, account string) (uint64, error)

	GetAllRestocks(ctx context.Context, page uint64, size uint64) ([]*models.Restock, error)
	CountAllRestocks(ctx context.Context) (uint64, error)

	// CashMovement's CRUD
	CreateCashMovement(ctx context.Context, t *models.CashMovement) error
	GetCashMovement(ctx context.Context, id string) (*models.CashMovement, error)
	UpdateCashMovement(ctx context.Context, t *models.CashMovement) error
	MarkDeleteCashMovement(ctx context.Context, id, by string) error
	UnMarkDeleteCashMovement(ctx context.Context, id string) error
	DeleteCashMovement(ctx context.Context, id string) error
	RestoreCashMovement(ctx context.Context, id string) error

	GetLatestCashMovement(ctx context.Context) (*models.CashMovement, error)

	GetDeletedCashMovements(ctx context.Context, page uint64, size uint64) ([]*models.CashMovement, error)
	CountDeletedCashMovements(ctx context.Context) (uint64, error)

	GetAllCashMovements(ctx context.Context, page uint64, size uint64, search string) ([]*models.CashMovement, error)
	CountAllCashMovements(ctx context.Context, search string) (uint64, error)

	// Other requests that are not CRUD but still needed
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	GetAccountByGoogle(ctx context.Context, googleID string) (*models.Account, error)
	GetAccountByCard(ctx context.Context, card string) (*models.Account, error)
	GetAccounts(ctx context.Context, page uint64, size uint64, query string) ([]*models.Account, error)
	CountAccounts(ctx context.Context, query string) (uint64, error)
	GetRefills(ctx context.Context, account string, page uint64, size uint64, startAt, endAt uint64) ([]*models.Refill, error)
	CountRefills(ctx context.Context, account string, startAt, endAt uint64) (uint64, error)
	GetItems(ctx context.Context, categoryID string, page, size uint64, state string, name string) ([]*models.Item, error)
	CountItems(ctx context.Context, categoryID string, state string, name string) (uint64, error)

	GetAllRefills(ctx context.Context, page uint64, size uint64, startAt, endAt uint64) ([]*models.Refill, error)
	CountAllRefills(ctx context.Context, startAt, endAt uint64) (uint64, error)
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetAllCarouselImages(ctx context.Context) ([]*models.CarouselImage, error)
	GetAllCarouselTexts(ctx context.Context) ([]*models.CarouselText, error)

	// Mongo transactions
	WithTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
}
