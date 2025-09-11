package api

import (
	"bar/autogen"
	"bar/autogen/helloasso"
	"bar/internal/config"
	"bar/internal/models"
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Try to validate started refills that were not validated yet.
// The payment might not have been validated yet when the user returns on the
// callback page, or he might lose the connection and never trigger the callback.
func (s *Server) ProcessStartedRefills(c context.Context) error {

	// Get all started online refills
	refills, err := s.DBackend.GetAllRemoteRefillsWithState(c, autogen.RemoteRefillStarted)
	if err != nil {
		return err
	}

	for _, refill := range refills {

		expired := false

		// Check if the payment has been made
		valid, order, err := s.validateHelloAssoCheckout(c, *refill.CheckoutIntentId, refill.Amount)

		// Cannot get checkout information from HelloAsso
		// It might be expired if it was created more than 45 minutes ago
		if err != nil {
			logrus.WithField("refill_id", refill.Id).Warn("Could not get checkout information from HelloAsso")

			if uint64(time.Now().Unix()) > refill.CreatedAt + 45*60 {
				expired = true
			}
		}

		// If the checkout cannot be validated after 45 minutes,
		// it can be considered abandonned (https://dev.helloasso.com/docs/validation-de-vos-paiements)
		if !valid && uint64(time.Now().Unix()) > refill.CreatedAt + 45*60 {
			expired = true
		}

		if valid {

			// Process the payment
			_, err = s.processRefillPayment(c, refill, *order.Id)

			if err != nil {
				logrus.WithField("refill_id", refill.Id).Warn("Could not process payment")
				continue
			}

		} else if expired {
			// Abandon the refill
			_, err := s.DBackend.UpdateRemoteRefillStateAtomic(c, refill, autogen.RemoteRefillAbandoned)

			if err != nil {
				logrus.WithField("refill_id", refill.Id).Warn("Could not update refill state to expired")
				continue
			}
			// No need to check if the state was sucessfully updated, if it was not it means the refill
			// has been validated in the meantime, although very unlikely
		}

	}

	return nil
}


// Atomically process a validated refill and update the user's balance
func (s *Server) processRefillPayment(c context.Context, r *models.RemoteRefill, orderId int32) (*models.Refill, error) {
	
	// Create refill transaction
	refill := &models.Refill{
		Refill: autogen.Refill{
			AccountId:    r.AccountId,
			AccountName:  r.AccountName,
			Amount:       int64(r.Amount),
			Type:         autogen.RefillHelloAsso,
			Id:           uuid.New(),
			IssuedAt:     uint64(time.Now().Unix()),
			IssuedBy:     r.AccountId,
			IssuedByName: r.AccountName,
			State:        autogen.RefillStateValid,
		},
	}

	_, err := s.DBackend.WithTransaction(c, func(ctx mongo.SessionContext) (interface{}, error) {
		
		// Update the remote refill atomically
		updated, err := s.DBackend.UpdateRemoteRefillStateAtomic(ctx, r, autogen.RemoteRefillProcessed)

		if err != nil {
			return nil, errors.New("failed to update refill state")
		} 

		// The refill was already validated
		if !updated {
			return nil, errors.New("refill state changed during processing")
		}

		// If the state was updated, we can proceed to the transaction

		// Update the rest of the refill
		r.RefillId = &refill.Id
		r.OrderId = &orderId
		err = s.DBackend.UpdateRemoteRefill(ctx, r)
		if err != nil {
			return nil, errors.New("failed to update remote refill")
		}

		// Create the transaction
		err = s.DBackend.CreateRefill(ctx, refill)
		if err != nil {
			return nil, errors.New("failed to create refill")
		}

		// Update user balance
		user, err := s.DBackend.GetAccount(ctx, r.AccountId.String())
		if err != nil {
			return nil, errors.New("could not get user account")
		}

		user.Balance += int64(r.Amount)

		err = s.DBackend.UpdateAccount(ctx, user)
		if err != nil {
			return nil, errors.New("failed to update account")
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return refill, nil
}

// Returns true if an HelloAsso checkout has been payed.
// The error returned can be displayed to the user directly
func (s *Server) validateHelloAssoCheckout(ctx context.Context, checkoutIntentId int32, expectedAmount int32) (bool, *helloasso.HelloAssoApiV5ModelsStatisticsOrderDetail, error){
	
	// Get the checkout from HelloAsso
	checkout, err := s.HelloAssoClient.GetOrganizationsOrganizationSlugCheckoutIntentsCheckoutIntentIdWithResponse(
		ctx,
		config.GetConfig().HelloAssoConfig.Slug,
		checkoutIntentId,
		nil,
	)

	if err != nil {
		logrus.WithField("checkout_intent_id", checkoutIntentId).Error("Error getting checkout information", err)
		return false, nil, errors.New(string(autogen.MsgInternalServerError))
	}

	if checkout.StatusCode() != 200 || checkout.JSON200 == nil {
		logrus.WithField("checkout_intent_id", checkoutIntentId).Error("Error getting checkout information", checkout.Status(), checkout.Body)
		return false, nil, errors.New(string(autogen.MsgInternalServerError))
	}

	// Payment not validated by helloasso yet, reject the refill
	if checkout.JSON200.Order == nil {
		return false, nil, nil
	}

	// Sanity check
	if *checkout.JSON200.Order.Amount.Total != expectedAmount {
		logrus.WithField("checkout_id", checkoutIntentId).Errorf("Order amount mismatch : expected %d, got %d", expectedAmount, *checkout.JSON200.Order.Amount.Total)
		return false, checkout.JSON200.Order, errors.New("internal server error : refill amount does not match order amount")
	}

	return true, checkout.JSON200.Order, nil
}

// (POST /account/remote-refills/start)
func (s *Server) StartRemoteRefill(c echo.Context, params autogen.StartRemoteRefillParams) error {
	// Get account from cookie
	user, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	if params.Amount < 50 {
		return autogen.StartRemoteRefill400Response{}.VisitStartRemoteRefillResponse(c.Response())
	}

	conf := config.GetConfig()

	// Start a checkout for this user

	var metadata interface{} = map[string] interface{}{
		"user_id": user.Id,
		"amount": params.Amount,
	};

	resp, err := s.HelloAssoClient.PostOrganizationsOrganizationSlugCheckoutIntentsWithResponse(
			c.Request().Context(),
			conf.HelloAssoConfig.Slug,
			helloasso.HelloAssoApiV5ModelsCartsInitCheckoutBody{
				BackUrl:           conf.ApiConfig.FrontendBasePath + "/client/index/refill",
				ErrorUrl:          conf.ApiConfig.FrontendBasePath + "/client/index/refill/callback",
				ReturnUrl:         conf.ApiConfig.FrontendBasePath + "/client/index/refill/callback",
				InitialAmount:     int32(params.Amount),
				TotalAmount:       int32(params.Amount),
				ContainsDonation:  false,
				ItemName:          "Rechargement Bar",
				Metadata:          &metadata,
			},
		)
	
	if err != nil {
		logrus.Error(err)
		return autogen.StartRemoteRefill500JSONResponse{}.VisitStartRemoteRefillResponse(c.Response())
	}

	if resp.StatusCode() != 200 || resp.JSON200 == nil {
		logrus.Error("Invalid response code ", resp.StatusCode())
		logrus.Debug(string(resp.Body))
		return autogen.StartRemoteRefill500JSONResponse{}.VisitStartRemoteRefillResponse(c.Response())
	}


	// Save the checkout intent

	refill := &models.RemoteRefill{
		RemoteRefill: autogen.RemoteRefill{
			Id:           uuid.New(),
			State: 		  autogen.RemoteRefillStarted,
			AccountId:    user.Id,
			AccountName:  user.Name(),
			Amount:       params.Amount,
			CheckoutIntentId: resp.JSON200.Id,
		},
	}

	err = s.DBackend.CreateRemoteRefill(c.Request().Context(), refill)
	if err != nil {
		logrus.Error(err)
		return autogen.StartRemoteRefill500JSONResponse{}.VisitStartRemoteRefillResponse(c.Response())
	}

	// Return the redirection url

	return autogen.StartRemoteRefill200JSONResponse{
		RedirectUrl: *resp.JSON200.RedirectUrl,
	}.VisitStartRemoteRefillResponse(c.Response())
}

// (POST /account/remote-refills/validate)
func (s *Server) SelfValidateRemoteRefill(c echo.Context, params autogen.SelfValidateRemoteRefillParams) error {
	
	// Get account from cookie
	user, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	// Get the remote refill to validate
	remote_refill, err := s.DBackend.FindRemoteRefillForAccount(c.Request().Context(), user.Id.String(), params.CheckoutIntentId)
	if err != nil {
		logrus.WithField("account_id", user.Id).Error("Checkout id not found : ", params.CheckoutIntentId)
		autogen.SelfValidateRemoteRefill404Response{}.VisitSelfValidateRemoteRefillResponse(c.Response())
		return nil
	}

	// Prevent validating a refill twice
	// Abandonned refills can still be manually validated
	if remote_refill.State != autogen.RemoteRefillStarted && remote_refill.State != autogen.RemoteRefillAbandoned {
		logrus.Errorf("Trying to validate checkout %d more than once", params.CheckoutIntentId)
		autogen.SelfValidateRemoteRefill409Response{}.VisitSelfValidateRemoteRefillResponse(c.Response())
		return nil
	}

	// Check if the payment has been made
	valid, order, err := s.validateHelloAssoCheckout(c.Request().Context(), params.CheckoutIntentId, remote_refill.Amount)
	
	if err != nil {
		autogen.SelfValidateRemoteRefill500JSONResponse{
			ErrorCode: autogen.ErrInternalServerError,
			Message: autogen.Messages(err.Error()),
		}.VisitSelfValidateRemoteRefillResponse(c.Response())
		return err
	}


	// Payment not validated by helloasso yet, reject the refill
	if !valid {
		autogen.SelfValidateRemoteRefill402Response{}.VisitSelfValidateRemoteRefillResponse(c.Response())
		logrus.WithField("account_id", user.Id).Errorf("Checkout %d not yet processed", params.CheckoutIntentId)
		return nil
	}
	// Process the payment
	refill, err := s.processRefillPayment(c.Request().Context(), remote_refill, *order.Id)

	if err != nil {
		logrus.Error(err)
		autogen.SelfValidateRemoteRefill500JSONResponse{
			ErrorCode: autogen.ErrInternalServerError,
			Message: autogen.MsgInternalServerError,
		}.VisitSelfValidateRemoteRefillResponse(c.Response())
		return err
	}

	autogen.SelfValidateRemoteRefill200JSONResponse(refill.Refill).VisitSelfValidateRemoteRefillResponse(c.Response())
	return nil
}

// (GET /remote-refills)
func (s *Server) GetRemoteRefills(c echo.Context, params autogen.GetRemoteRefillsParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var startsAt uint64 = 0
	if params.StartDate != nil {
		t, err := time.Parse("2006-01-02", *params.StartDate) // 2006-01-02 is the reference time in Go
		if err == nil {
			startsAt = uint64(t.Unix())
		}
	}
	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		t, err := time.Parse("2006-01-02", *params.EndDate) // Putting the same date doesn't activate the date filter
		if err == nil {
			endsAt = uint64(t.Unix())
		}
	}

	count, err := s.DBackend.CountAllRemoteRefills(c.Request().Context(), params.AccountName, params.State, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetAllRemoteRefills(c.Request().Context(), dbpage, limit, params.AccountName, params.State, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.RemoteRefill

	for _, refill := range data {
		refills = append(refills, refill.RemoteRefill)
	}

	autogen.GetRemoteRefills200JSONResponse{
		RemoteRefills: refills,
		Limit: limit,
		Page: page,
		MaxPage: maxPage,
	}.VisitGetRemoteRefillsResponse(c.Response())
	return nil
}


// (GET /account/remote-refills/pending)
func (s *Server) GetPendingRemoteRefills(c echo.Context) error {
	// Get account from cookie
	user, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	data, err := s.DBackend.GetAllPendingRemoteRefillsForAccount(c.Request().Context(), user.Id.String())
	if err != nil {
		return Error500(c);
	}

	var refills []autogen.RemoteRefill = make([]autogen.RemoteRefill, 0);

	for _, refill := range data {
		refills = append(refills, refill.RemoteRefill)
	}

	autogen.GetPendingRemoteRefills200JSONResponse{
		RemoteRefills: refills,
	}.VisitGetPendingRemoteRefillsResponse(c.Response());
	return nil;
}

// (POST /remote-refills/validate)
func (s *Server) VerifyRemoteRefill(c echo.Context, params autogen.VerifyRemoteRefillParams) error {
	
	user, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	logrus.WithField("account_id", user.Account.Id).WithField("remote_refill_id", params.Id.String()).Debug("Manually verifying remote refill")

	// Get the remote refill to validate
	remote_refill, err := s.DBackend.GetRemoteRefill(c.Request().Context(), params.Id.String())
	if err != nil {
		// No refill found
		autogen.VerifyRemoteRefill404Response{}.VisitVerifyRemoteRefillResponse(c.Response())
		return nil
	}

	// Prevent validating a refill twice
	// Abandonned refills can still be manually validated
	if remote_refill.State != autogen.RemoteRefillStarted && remote_refill.State != autogen.RemoteRefillAbandoned {
		autogen.VerifyRemoteRefill409Response{}.VisitVerifyRemoteRefillResponse(c.Response())
		return nil
	}

	// Check if the payment has been made
	valid, order, err := s.validateHelloAssoCheckout(c.Request().Context(), *remote_refill.CheckoutIntentId, remote_refill.Amount)
	
	if err != nil {
		autogen.VerifyRemoteRefill500JSONResponse{
			ErrorCode: autogen.ErrInternalServerError,
			Message: autogen.Messages(err.Error()),
		}.VisitVerifyRemoteRefillResponse(c.Response())
		return err
	}

	// Payment not validated by helloasso yet, reject the refill
	if !valid {
		autogen.VerifyRemoteRefill402Response{}.VisitVerifyRemoteRefillResponse(c.Response())
		return nil
	}
	// Process the payment
	refill, err := s.processRefillPayment(c.Request().Context(), remote_refill, *order.Id)

	if err != nil {
		logrus.Error(err)
		autogen.VerifyRemoteRefill500JSONResponse{
			ErrorCode: autogen.ErrInternalServerError,
			Message: autogen.MsgInternalServerError,
		}.VisitVerifyRemoteRefillResponse(c.Response())
		return err
	}

	autogen.VerifyRemoteRefill200JSONResponse(refill.Refill).VisitVerifyRemoteRefillResponse(c.Response())
	return nil
}