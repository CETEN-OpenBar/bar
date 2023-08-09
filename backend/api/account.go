package api

import (
	"bar/autogen"
	"bar/internal/models"
	"encoding/csv"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /account)
func (s *Server) GetAccount(c echo.Context) error {
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	account := c.Get("userAccount").(*models.Account)

	// Return account
	resp := autogen.GetAccount200JSONResponse{
		Account: &account.Account,
	}
	resp.VisitGetAccountResponse(c.Response())
	return nil
}

// (PATCH /account)
func (s *Server) PatchAccount(c echo.Context) error {
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	account := c.Get("userAccount").(*models.Account)

	var param autogen.PatchAccountJSONBody
	err := c.Bind(&param)
	if err != nil {
		return Error400(c)
	}

	// sha256 both pins
	if !account.VerifyPin(param.OldCardPin) {
		return Error400(c)
	}

	account.SetPin(param.NewCardPin)

	err = s.UpdateAccount(c.Request().Context(), account)
	if err != nil {
		return Error500(c)
	}

	autogen.PatchAccount200JSONResponse{
		Account: &account.Account,
	}.VisitPatchAccountResponse(c.Response())
	return nil
}

// (GET /accounts)
func (s *Server) GetAccounts(c echo.Context, params autogen.GetAccountsParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
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

	var search string
	if params.Search != nil {
		search = *params.Search
	}

	// Calculate max page
	count, err := s.DBackend.CountAccounts(c.Request().Context(), search)
	if err != nil {
		return Error500(c)
	}

	maxPage := uint64(count) / limit

	// Get accounts from database
	accounts, err := s.DBackend.GetAccounts(c.Request().Context(), page, limit, search)
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
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	var req autogen.NewAccount
	err := c.Bind(&req)
	if err != nil {
		return Error400(c)
	}

	var cardId string
	if req.CardId != nil {
		cardId = *req.CardId
	}

	var priceRole = autogen.AccountPriceNormal
	if req.PriceRole != nil {
		priceRole = *req.PriceRole
	}

	account := &models.Account{
		Account: autogen.Account{
			Balance:      req.Balance,
			CardId:       cardId,
			EmailAddress: req.EmailAddress,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Role:         req.Role,
			PriceRole:    priceRole,
			State:        autogen.AccountOK,
		},
	}
	account.SetPin("1234")

	err = s.CreateAccount(c.Request().Context(), account)
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
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	err := s.DBackend.MarkDeleteAccount(c.Request().Context(), accountId.String(), adminID)
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
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId.String())
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
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	var req autogen.UpdateAccountAdmin
	err := c.Bind(&req)
	if err != nil {
		return Error400(c)
	}

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId.String())
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
		account.SetPin("1234")
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
		r := *req.Role

		// Can only set the roles to something below the current role
		switch account.Role {
		case autogen.AccountSuperAdmin:
			account.Account.Role = r
		case autogen.AccountAdmin:
			if r == autogen.AccountSuperAdmin {
				return Error400(c)
			}
			account.Account.Role = r
		case autogen.AccountGhost:
			// Can't set to anything
			return Error400(c)
		case autogen.AccountMember:
			if r == autogen.AccountSuperAdmin || r == autogen.AccountAdmin {
				return Error400(c)
			}
			account.Account.Role = r
		}
	}
	if req.PriceRole != nil {
		r := *req.PriceRole

		account.Account.PriceRole = r
	}
	if req.State != nil {
		account.Account.State = *req.State
	}
	if req.Restrictions != nil {
		account.Account.Restrictions = *req.Restrictions
	}

	err = s.UpdateAccount(c.Request().Context(), account)
	if err != nil {
		return Error500(c)
	}

	logrus.Info("Account updated: ", accountId, " by ", adminID)
	autogen.PatchAccountId200JSONResponse(account.Account).VisitPatchAccountIdResponse(c.Response())
	return nil
}

// (POST /import/accounts)
func (s *Server) ImportAccounts(c echo.Context) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

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
		tag = strings.Split(tag, ",")[0]
		for j, colName := range colNames {
			colName = strings.ToLower(colName)
			colName = strings.ReplaceAll(colName, " ", "_")
			colName = strings.ReplaceAll(colName, "-", "_")
			colName = strings.TrimSpace(colName)
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
		balance, err := strconv.ParseInt(record[assignments["balance"]], 10, 64)
		if err != nil {
			notProcessed = append(notProcessed, record[0])
			continue
		}

		account := &models.Account{
			Account: autogen.Account{
				Balance:      balance,
				Id:           uuid.New(),
				CardId:       record[assignments["card_id"]],
				EmailAddress: record[assignments["email_address"]],
				FirstName:    record[assignments["first_name"]],
				LastName:     record[assignments["last_name"]],
				Role:         autogen.AccountRole(record[assignments["role"]]),
				PriceRole:    autogen.AccountPriceRole(record[assignments["price_role"]]),
				State:        autogen.AccountNotOnBoarded,
			},
		}
		account.SetPin("1234")

		err = s.CreateAccount(c.Request().Context(), account)
		if err != nil {
			logrus.Error(err)
			notProcessed = append(notProcessed, record[0])
			continue
		}
	}

	logrus.Info("Accounts imported: ", len(records)-len(notProcessed), " by ", adminID)
	autogen.ImportAccounts200JSONResponse{
		NotAccepted: &notProcessed,
	}.VisitImportAccountsResponse(c.Response())
	return nil
}

// (GET /account/admin)
func (s *Server) GetAccountAdmin(ctx echo.Context) error {
	logged := ctx.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(ctx)
	}

	canRestore := ctx.Get("adminAccountRole").(autogen.AccountRole) == autogen.AccountSuperAdmin

	// Return account
	resp := autogen.GetAccountAdmin200JSONResponse{
		IsAllowed:  true,
		CanRestore: canRestore,
	}
	resp.VisitGetAccountAdminResponse(ctx.Response())
	return nil
}
