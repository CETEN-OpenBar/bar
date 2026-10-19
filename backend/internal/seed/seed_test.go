package seed

import (
	"bar/internal/db"
	"bar/internal/models"
	"context"
	"errors"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

// fakeBackend keeps seeded documents in memory. Only the methods Run uses are
// implemented; calling anything else panics through the nil embedded interface.
type fakeBackend struct {
	db.DBackend
	categories map[string]*models.Category
	items      map[string]*models.Item
	accounts   map[string]*models.Account
}

func newFakeBackend() *fakeBackend {
	return &fakeBackend{
		categories: map[string]*models.Category{},
		items:      map[string]*models.Item{},
		accounts:   map[string]*models.Account{},
	}
}

func (f *fakeBackend) GetCategory(_ context.Context, id string) (*models.Category, error) {
	if c, ok := f.categories[id]; ok {
		return c, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (f *fakeBackend) CreateCategory(_ context.Context, c *models.Category) error {
	f.categories[c.Id.String()] = c
	return nil
}

func (f *fakeBackend) GetItem(_ context.Context, id string) (*models.Item, error) {
	if i, ok := f.items[id]; ok {
		return i, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (f *fakeBackend) CreateItem(_ context.Context, i *models.Item) error {
	f.items[i.Id.String()] = i
	return nil
}

func (f *fakeBackend) GetAccountByCard(_ context.Context, card string) (*models.Account, error) {
	if a, ok := f.accounts[card]; ok {
		return a, nil
	}
	return nil, mongo.ErrNoDocuments
}

// CreateAccount mirrors the unique index on accounts.email_address, which
// also applies to empty strings.
func (f *fakeBackend) CreateAccount(_ context.Context, a *models.Account) error {
	for _, existing := range f.accounts {
		if existing.EmailAddress == a.EmailAddress {
			return errors.New("E11000 duplicate key error: email_address " + a.EmailAddress)
		}
	}
	f.accounts[*a.CardId] = a
	return nil
}

func stubPictures(t *testing.T) map[string][]byte {
	t.Helper()
	saved := map[string][]byte{}
	orig := savePicture
	savePicture = func(name string, data []byte) error {
		saved[name] = data
		return nil
	}
	t.Cleanup(func() { savePicture = orig })
	return saved
}

func TestIsLocalURI(t *testing.T) {
	tests := []struct {
		uri  string
		want bool
	}{
		{"mongodb://localhost:27017", true},
		{"mongodb://127.0.0.1:27017", true},
		{"mongodb://[::1]:27017", true},
		{"mongodb://127.0.0.1", true},
		{"mongodb+srv://cluster0.x.mongodb.net", false},
		{"mongodb://db.example.org:27017", false},
		{"mongodb://localhost:27017,db.example.org:27017", false},
		{"not a uri", false},
	}

	for _, tt := range tests {
		if got := IsLocalURI(tt.uri); got != tt.want {
			t.Errorf("IsLocalURI(%q) = %v, want %v", tt.uri, got, tt.want)
		}
	}
}

func TestRunIsIdempotent(t *testing.T) {
	stubPictures(t)
	b := newFakeBackend()

	first, err := Run(context.Background(), b)
	if err != nil {
		t.Fatalf("first Run: %v", err)
	}
	if first == 0 {
		t.Fatal("first Run created nothing")
	}
	counts := [3]int{len(b.categories), len(b.items), len(b.accounts)}

	second, err := Run(context.Background(), b)
	if err != nil {
		t.Fatalf("second Run: %v", err)
	}
	if second != 0 {
		t.Errorf("second Run created %d documents, want 0", second)
	}
	if got := [3]int{len(b.categories), len(b.items), len(b.accounts)}; got != counts {
		t.Errorf("document counts changed from %v to %v", counts, got)
	}
}

func TestRunSavesPicturesWhereTheAPIReadsThem(t *testing.T) {
	saved := stubPictures(t)
	b := newFakeBackend()

	if _, err := Run(context.Background(), b); err != nil {
		t.Fatalf("Run: %v", err)
	}

	for id, c := range b.categories {
		if _, ok := saved["categories/"+id]; !ok {
			t.Errorf("category %q has no picture at categories/%s", c.Name, id)
		}
		if want := "/categories/" + id + "/picture"; c.PictureUri != want {
			t.Errorf("category %q PictureUri = %q, want %q", c.Name, c.PictureUri, want)
		}
	}
	for id, i := range b.items {
		if _, ok := saved["items/"+id]; !ok {
			t.Errorf("item %q has no picture at items/%s", i.Name, id)
		}
		if !strings.HasSuffix(i.PictureUri, "/items/"+id+"/picture") {
			t.Errorf("item %q PictureUri = %q", i.Name, i.PictureUri)
		}
	}
}
