package api

import (
	"bar/autogen"
	"bar/internal/models"
	"encoding/csv"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /account)
func (s *Server) GetAccount(c echo.Context) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	accountID, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	account, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// Return account
	resp := autogen.GetAccount200JSONResponse{
		Account: &account.Account,
	}
	resp.VisitGetAccountResponse(c.Response())
	return nil
}

// (GET /accounts)
func (s *Server) GetAccounts(c echo.Context, params autogen.GetAccountsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// Set up parameters
	var page uint64 = 1
	var limit uint64 = 10
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	if page > 0 {
		page -= 1
	}

	if limit > 100 {
		limit = 100
	}

	// Calculate max page
	count, err := s.DBackend.CountAccounts()
	if err != nil {
		return Error500(c)
	}

	maxPage := uint64(count) / limit

	// Get accounts from database
	accounts, err := s.DBackend.GetAccounts(page, limit)
	if err != nil {
		return Error500(c)
	}

	var ac []autogen.Account
	for _, account := range accounts {
		ac = append(ac, account.Account)
	}

	autogen.GetAccounts200JSONResponse{
		Accounts: ac,
		Limit:    limit,
		Page:     page,
		MaxPage:  maxPage,
	}.VisitGetAccountsResponse(c.Response())
	return nil
}

// (POST /accounts)
func (s *Server) PostAccounts(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var req autogen.NewAccount
	err := c.Bind(&req)
	if err != nil {
		return Error400(c)
	}

	var cardId string
	if req.CardId != nil {
		cardId = *req.CardId
	}

	account := &models.Account{
		Account: autogen.Account{
			Balance:      req.Balance,
			CardId:       cardId,
			CardPin:      "9af15b336e6a9619928537df30b2e6a2376569fcf9d7e773eccede65606529a0",
			EmailAddress: req.EmailAddress,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Role:         req.Role,
			State:        autogen.AccountOk,
		},
	}

	err = s.CreateAccount(account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Error409(c)
		}
		return Error500(c)
	}

	logrus.Info("Account created: ", account.Account.Id, " by ", adminID)
	autogen.PostAccounts200JSONResponse(account.Account).VisitPostAccountsResponse(c.Response())
	return nil
}

// (DELETE /accounts/{account_id})
func (s *Server) MarkDeleteAccountId(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	err := s.DBackend.MarkDeleteAccount(accountId.String(), adminID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	logrus.Info("Account marked as deleted: ", accountId, " by ", adminID)
	autogen.DeleteAccount204Response{}.VisitDeleteAccountResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id})
func (s *Server) GetAccountId(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	account, err := s.DBackend.GetAccount(accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	logrus.Info("Account retrieved: ", accountId, " by ", adminID)
	autogen.GetAccountId200JSONResponse(account.Account).VisitGetAccountIdResponse(c.Response())
	return nil
}

// (PATCH /accounts/{account_id})
func (s *Server) PatchAccountId(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var req autogen.UpdateAccountAdmin
	err := c.Bind(&req)
	if err != nil {
		return Error400(c)
	}

	account, err := s.DBackend.GetAccount(accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	if req.Balance != nil {
		account.Account.Balance = *req.Balance
	}
	if req.CardId != nil {
		account.Account.CardId = *req.CardId
		account.Account.CardPin = "9af15b336e6a9619928537df30b2e6a2376569fcf9d7e773eccede65606529a0"
	}
	if req.EmailAddress != nil {
		account.Account.EmailAddress = *req.EmailAddress
	}
	if req.FirstName != nil {
		account.Account.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		account.Account.LastName = *req.LastName
	}
	if req.Role != nil {
		account.Account.Role = *req.Role
	}
	if req.State != nil {
		account.Account.State = *req.State
	}
	if req.Restrictions != nil {
		account.Account.Restrictions = *req.Restrictions
	}

	err = s.UpdateAccount(account)
	if err != nil {
		return Error500(c)
	}

	logrus.Info("Account updated: ", accountId, " by ", adminID)
	autogen.PatchAccountId200JSONResponse(account.Account).VisitPatchAccountIdResponse(c.Response())
	return nil
}

// (POST /import/accounts)
func (s *Server) ImportAccounts(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return Error400(c)
	}

	// Parse CSV file
	f, err := file.Open()
	if err != nil {
		return Error400(c)
	}

	// Read CSV file
	r := csv.NewReader(f)
	colNames, err := r.Read()
	if err != nil {
		return Error400(c)
	}

	// Create assignment map for columns
	// map[string]uint64{"email": 0, "first_name": 1, "last_name": 2, "role": 3, "balance": 4} meaning that the email is in the first column
	// Using reflection to get the field name of the struct
	var req autogen.NewAccount
	var assignments = make(map[string]int)

	v := reflect.ValueOf(req)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("json")
		for j, colName := range colNames {
			if tag == colName {
				assignments[tag] = j
			}
		}
	}

	records, err := r.ReadAll()
	if err != nil {
		return Error400(c)
	}

	var notProcessed []string

	for _, record := range records {
		// Check balance
		balance, err := strconv.ParseUint(record[assignments["balance"]], 10, 64)
		if err != nil {
			notProcessed = append(notProcessed, record[0])
			continue
		}

		account := &models.Account{
			Account: autogen.Account{
				Balance:      balance,
				CardId:       record[assignments["card_id"]],
				CardPin:      "9af15b336e6a9619928537df30b2e6a2376569fcf9d7e773eccede65606529a0",
				EmailAddress: record[assignments["email"]],
				FirstName:    record[assignments["first_name"]],
				LastName:     record[assignments["last_name"]],
				Role:         autogen.AccountRole(record[assignments["role"]]),
				State:        autogen.AccountOk,
			},
		}

		err = s.CreateAccount(account)
		if err != nil {
			notProcessed = append(notProcessed, record[0])
			continue
		}
	}

	logrus.Info("Accounts imported: ", len(records)-len(notProcessed), " by ", adminId)
	autogen.ImportAccounts200JSONResponse{
		NotAccepted: &notProcessed,
	}.VisitImportAccountsResponse(c.Response())
	return nil
}
