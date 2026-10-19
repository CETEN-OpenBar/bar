// Package seed fills a local development database with sample data.
package seed

import (
	"bar/autogen"
	"bar/internal/db"
	"bar/internal/models"
	"bar/internal/storage"
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// namespace makes seeded IDs stable, so running the seed twice finds the
// documents it created the first time instead of duplicating them.
var namespace = uuid.MustParse("6f1c2d9e-3b4a-4f60-9a1e-5d2c7b8e0f13")

// savePicture is replaced in tests.
var savePicture = storage.SaveFile

type sampleItem struct {
	name  string
	price uint64 // cents
}

type sampleCategory struct {
	name  string
	color color.RGBA
	items []sampleItem
}

var categories = []sampleCategory{
	{"Boissons", color.RGBA{0x3b, 0x82, 0xf6, 0xff}, []sampleItem{
		{"Eau 50cl", 50},
		{"Coca-Cola 33cl", 100},
		{"Ice Tea 33cl", 100},
		{"Café", 50},
	}},
	{"À manger", color.RGBA{0xf9, 0x73, 0x16, 0xff}, []sampleItem{
		{"Croque-monsieur", 200},
		{"Panini", 250},
	}},
	{"Snacks", color.RGBA{0x22, 0xc5, 0x5e, 0xff}, []sampleItem{
		{"Chips", 80},
		{"Barre chocolatée", 70},
	}},
}

type sampleAccount struct {
	card      string
	firstName string
	role      autogen.AccountRole
	balance   int64 // cents
}

var accounts = []sampleAccount{
	{"2", "Membre", autogen.AccountMember, 0},
	{"3", "Étudiant", autogen.AccountStudent, 2000},
	{"4", "Admin", autogen.AccountAdmin, 0},
}

// IsLocalURI reports whether every host in a mongodb:// URI is this machine.
func IsLocalURI(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme != "mongodb" || u.Host == "" {
		return false
	}
	for _, h := range strings.Split(u.Host, ",") {
		host := h
		if hp, _, err := net.SplitHostPort(h); err == nil {
			host = hp
		}
		host = strings.Trim(host, "[]")
		if host != "localhost" && host != "127.0.0.1" && host != "::1" {
			return false
		}
	}
	return true
}

// Run creates the sample categories, items and accounts that don't exist yet
// and returns how many documents it created.
func Run(ctx context.Context, b db.DBackend) (int, error) {
	created := 0

	for pos, sc := range categories {
		catID := uuid.NewSHA1(namespace, []byte("category:"+sc.name))
		ok, err := missing(b.GetCategory(ctx, catID.String()))
		if err != nil {
			return created, err
		}
		if ok {
			if err := savePicture("categories/"+catID.String(), solidPNG(sc.color)); err != nil {
				return created, err
			}
			err := b.CreateCategory(ctx, &models.Category{Category: autogen.Category{
				Id:         catID,
				Name:       sc.name,
				Position:   uint64(pos),
				PictureUri: "/categories/" + catID.String() + "/picture",
			}})
			if err != nil {
				return created, err
			}
			created++
		}

		for _, si := range sc.items {
			itemID := uuid.NewSHA1(namespace, []byte("item:"+si.name))
			ok, err := missing(b.GetItem(ctx, itemID.String()))
			if err != nil {
				return created, err
			}
			if !ok {
				continue
			}
			if err := savePicture("items/"+itemID.String(), solidPNG(sc.color)); err != nil {
				return created, err
			}
			err = b.CreateItem(ctx, &models.Item{Item: autogen.Item{
				Id:            itemID,
				CategoryId:    catID,
				Name:          si.name,
				State:         autogen.ItemBuyable,
				AmountLeft:    50,
				OptimalAmount: 20,
				Prices: autogen.ItemPrices{
					Ceten:       si.price,
					Coutant:     si.price,
					Externe:     si.price + 20,
					Menu:        si.price,
					Privilegies: si.price,
					StaffBar:    si.price,
				},
				PictureUri: "/categories/" + catID.String() + "/items/" + itemID.String() + "/picture",
			}})
			if err != nil {
				return created, err
			}
			created++
		}
	}

	for _, sa := range accounts {
		ok, err := missing(b.GetAccountByCard(ctx, sa.card))
		if err != nil {
			return created, err
		}
		if !ok {
			continue
		}
		acc := &models.Account{Account: autogen.Account{
			Id:     uuid.NewSHA1(namespace, []byte("account:"+sa.card)),
			CardId: autogen.OptionalString(sa.card),
			// email_address has a unique index, empty strings included
			EmailAddress: "carte" + sa.card + "@bar.local",
			FirstName:    sa.firstName,
			LastName:     "Test",
			Role:         sa.role,
			State:        autogen.AccountOK,
			PriceRole:    autogen.AccountPriceCeten,
			Balance:      sa.balance,
		}}
		acc.SetPin("1234")
		if err := b.CreateAccount(ctx, acc); err != nil {
			return created, err
		}
		created++
	}

	return created, nil
}

// missing turns a Get* result into "not found yet" (true) or an error.
func missing[T any](_ T, err error) (bool, error) {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return true, nil
	}
	return false, err
}

func solidPNG(c color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = c.R, c.G, c.B, c.A
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
