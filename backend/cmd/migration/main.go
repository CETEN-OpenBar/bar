package main

import (
	"bar/autogen"
	"bar/internal/config"
	"bar/internal/db"
	"bar/internal/db/mongo"
	"bar/internal/models"
	"bar/internal/storage"
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/sirupsen/logrus"
)

func main() {
	c := config.GetConfig()

	opts := db.NewDatabaseOptions(c.MongoConfig.ConnectionURI, c.MongoConfig.Database, time.Millisecond*time.Duration(c.MongoConfig.Timeout))
	mongoDB := mongo.NewMongoBackend(opts)

	err := mongoDB.Connect()
	if err != nil {
		logrus.Panic(err)
	}

	// We need to connect to MariaDB on localhost:3306
	// to execute the query that we will move to PMB
	mariaDB, err := sql.Open("mysql", "root:qwerty@tcp(localhost:3306)/pmb?parseTime=true")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	mariaDB.SetConnMaxLifetime(time.Minute * 3)
	mariaDB.SetMaxOpenConns(10)
	mariaDB.SetMaxIdleConns(10)

	// Set a seed for UUID
	r := rand.NewSource(3625)
	rdr := rand.New(r)
	uuid.SetRand(rdr)

	part1(mongoDB, mariaDB)
	part2(mongoDB, mariaDB)
	part3(mongoDB, mariaDB)
}

var IDmap = make(map[string]uuid.UUID)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func part1(mongoDB db.DBackend, mariaDB *sql.DB) {
	// SELECT client.clientId, firstname, lastname, points, balance, status, A.groupId, A.name as promotion FROM client
	// JOIN (SELECT * FROM client_group JOIN client_group_membership ON client_group.groupId = client_group_membership.clientGroupId) as A
	// 	ON A.clientId = client.clientId;

	rq := `SELECT client.clientId, firstname, lastname, points, balance, status, A.groupId, A.name as promotion FROM client
	JOIN (SELECT * FROM client_group JOIN client_group_membership ON client_group.groupId = client_group_membership.clientGroupId) as A
		ON A.clientId = client.clientId;`

	rows, err := mariaDB.Query(rq)
	if err != nil {
		panic(err)
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	for rows.Next() {
		// Put in a map
		// clientId 	firstname 	lastname 	points 	balance 	status  	groupId 	name
		var clientId string
		var firstname string
		var lastname string
		var points int
		var balance string
		var status string
		var groupId string
		var promotion string

		err = rows.Scan(&clientId, &firstname, &lastname, &points, &balance, &status, &groupId, &promotion)
		if err != nil {
			panic(err)
		}

		uuid, _ := uuid.NewUUID()

		IDmap[clientId] = uuid

		bal, _ := strconv.Atoi(strings.ReplaceAll(balance, ".", ""))
		pts := int64(points)

		// -- client : 0
		// -- membre ceten : 1
		// -- vip : 2
		// -- membre bar : 3
		// -- membre privilégié : 4
		// -- membre du bureau : 5

		var priceRole autogen.AccountPriceRole
		var role autogen.AccountRole

		switch status {
		case "0":
			priceRole = autogen.AccountPriceInterne
			role = autogen.AccountStudent
		case "1":
			priceRole = autogen.AccountPriceCeten
			role = autogen.AccountStudent
		case "2":
			priceRole = autogen.AccountPriceVIP
			role = autogen.AccountStudent
		case "3":
			priceRole = autogen.AccountPriceStaff
			role = autogen.AccountStudent
		case "4":
			priceRole = autogen.AccountPriceMembrePrivilegie
			role = autogen.AccountStudent
		case "5":
			priceRole = autogen.AccountPriceMembreBureau
			role = autogen.AccountMember
		}

		acc := &models.Account{
			Account: autogen.Account{
				Balance:   int64(bal),
				Id:        uuid,
				Points:    pts,
				FirstName: firstname,
				LastName:  lastname,
				State:     autogen.AccountOK,
				Role:      role,
				PriceRole: priceRole,
			},
		}

		// Rândom RANDOM should become random.random@tn.net
		firstname, _, _ = transform.String(t, strings.ToLower(firstname))
		lastname, _, _ = transform.String(t, strings.ToLower(lastname))
		firstname = strings.ReplaceAll(firstname, " ", "-")
		lastname = strings.ReplaceAll(lastname, " ", "-")

		// Insert into mongo
		email := fmt.Sprintf("%s.%s@telecomnancy.net", firstname, lastname)
		h := sha256.Sum256([]byte(email))
		acc.GooglePicture = autogen.OptionalString(fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=mp&f=y", h))

		acc.EmailAddress = email

		acc.SetPin("1234")

		err = mongoDB.CreateAccount(context.Background(), acc)
		if err != nil {
			logrus.Error(err)
		}

		fmt.Println(clientId)
	}
}

var Categories = make(map[string]uuid.UUID)

func part2(mongoDB db.DBackend, mariaDB *sql.DB) {
	// SELECT product.name, price as exte, cetenPrice as ceten, privilegePrice as vip, barMemberPrice staff, privilegePrice as privilegie, cost as bureau, product.image, stock, type, product_type.name, product_type.image from product join product_type on product.type = product_type.productTypeId;

	rq := `SELECT product.hidden, product.criticalStock, product.name, price as exte, cetenPrice as ceten, privilegePrice as vip, barMemberPrice staff, privilegePrice as privilegie, cost as bureau, product.image, stock, type, product_type.name as categorie, product_type.image as image_categorie, product_type.boolean as categorie_hidden from product join product_type on product.type = product_type.productTypeId;`

	rows, err := mariaDB.Query(rq)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		// name 	exte 	ceten 	vip 	staff 	privilegie 	bureau 	image 	stock 	type 	name 	image
		var hidden bool
		var criticalStock int
		var productName string
		var priceExte string
		var priceCeten string
		var priceVip string
		var priceStaff string
		var pricePrivilegie string
		var priceBureau string
		var image []byte
		var stock int
		var type_ string
		var categorie string
		var imageCategorie []byte
		var categorieHidden bool

		if strings.Contains(productName, "OLD") {
			continue
		}

		var err = rows.Scan(&hidden, &criticalStock, &productName, &priceExte, &priceCeten, &priceVip, &priceStaff, &pricePrivilegie, &priceBureau, &image, &stock, &type_, &categorie, &imageCategorie, &categorieHidden)
		if err != nil {
			panic(err)
		}

		// Find categorie
		if _, ok := Categories[categorie]; !ok {
			uid := uuid.New()
			storage.SaveFile("categories/"+uid.String(), imageCategorie)
			cat := &models.Category{
				Category: autogen.Category{
					Name:       categorie,
					Id:         uid,
					Position:   0,
					PictureUri: "/categories/" + uid.String() + "/picture",
				},
			}
			if categorieHidden {
				system := uuid.MustParse("00000000-0000-0000-0000-000000000000")
				now := uint64(time.Now().Unix())
				cat.DeletedAt = &now
				cat.DeletedBy = &system
			}
			mongoDB.CreateCategory(context.Background(), cat)
			Categories[categorie] = uid
		}

		categoryId := Categories[categorie]

		cetenPrice, _ := strconv.Atoi(strings.ReplaceAll(priceCeten, ".", ""))
		vipPrice, _ := strconv.Atoi(strings.ReplaceAll(priceVip, ".", ""))
		staffPrice, _ := strconv.Atoi(strings.ReplaceAll(priceStaff, ".", ""))
		privilegiePrice, _ := strconv.Atoi(strings.ReplaceAll(pricePrivilegie, ".", ""))
		bureauPrice, _ := strconv.Atoi(strings.ReplaceAll(priceBureau, ".", ""))
		extePrice, _ := strconv.Atoi(strings.ReplaceAll(priceExte, ".", ""))

		prices := autogen.ItemPrices{
			Ceten:            uint64(cetenPrice),
			Vip:              uint64(vipPrice),
			Staff:            uint64(staffPrice),
			MembrePrivilegie: uint64(privilegiePrice),
			MembreBureau:     uint64(bureauPrice),
			Exte:             uint64(extePrice),
			Interne:          uint64(extePrice),
		}

		uid := uuid.New()
		storage.SaveFile("items/"+uid.String(), image)
		item := &models.Item{
			Item: autogen.Item{
				Id:            uid,
				CategoryId:    categoryId,
				Name:          productName,
				Prices:        prices,
				PictureUri:    "/categories/" + categoryId.String() + "/items/" + uid.String() + "/picture",
				State:         autogen.ItemBuyable,
				AmountLeft:    uint64(stock),
				OptimalAmount: uint64(criticalStock),
			},
		}

		if hidden {
			item.State = autogen.ItemNotBuyable
		}

		// Save item to database
		err = mongoDB.CreateItem(context.Background(), item)
		if err != nil {
			logrus.Error(err)
		}

	}
}

func part3(mongoDB db.DBackend, mariaDB *sql.DB) {
	rq := `SELECT client_transaction.cancelled, client_transaction.clientId, client_transaction.operator, client_transaction.date,client_transaction.method, client_transaction.amount from client_transaction; `

	rows, err := mariaDB.Query(rq)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		// clientId 	operator 	date 	method 	amount
		var canceled bool
		var clientId string
		var operator string
		var date time.Time
		var method string
		var amount string

		err = rows.Scan(&canceled, &clientId, &operator, &date, &method, &amount)
		if err != nil {
			panic(err)
		}

		var type_ autogen.RefillType

		switch method {
		case "0":
			type_ = autogen.RefillCash
		case "1":
			type_ = autogen.RefillCheck
		case "2":
			type_ = autogen.RefillTransfer
		case "3":
			type_ = autogen.RefillCard
		case "4":
			type_ = autogen.RefillOther
		}

		// Find account
		accountId, nok := IDmap[clientId]
		operatorId, nok2 := IDmap[operator]

		if !nok || !nok2 {
			fmt.Println("Skipped refill")
			continue
		}

		acc, _ := mongoDB.GetAccount(context.Background(), accountId.String())
		op, _ := mongoDB.GetAccount(context.Background(), operatorId.String())

		howMuch, _ := strconv.Atoi(strings.ReplaceAll(amount, ".", ""))

		refill := &models.Refill{
			Refill: autogen.Refill{
				AccountId:    accountId,
				AccountName:  acc.Name(),
				IssuedAt:     uint64(date.Unix()),
				IssuedBy:     operatorId,
				Type:         type_,
				IssuedByName: op.Name(),
				State:        autogen.Valid,
				Amount:       int64(howMuch),
				Id:           uuid.New(),
			},
		}

		fmt.Println(refill)

		if canceled {
			refill.Refill.State = autogen.Canceled
		}

		err = mongoDB.CreateRefill(context.Background(), refill)
		if err != nil {
			logrus.Error(err)
		}
	}
}
