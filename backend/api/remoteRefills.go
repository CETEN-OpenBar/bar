package api

import (
	"bar/autogen"
	"bar/autogen/helloasso"
	"bar/internal/config"
	"bar/internal/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	if remote_refill.State != autogen.RemoteRefillStarted {
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

	// Update the refill status and the account balance
	refill := &models.Refill{
		Refill: autogen.Refill{
			AccountId:    user.Id,
			AccountName:  user.Name(),
			Amount:       int64(remote_refill.Amount),
			Type:         autogen.RefillHelloAsso,
			Id:           uuid.New(),
			IssuedAt:     uint64(time.Now().Unix()),
			IssuedBy:     user.Id,
			IssuedByName: user.Name(),
			State:        autogen.RefillStateValid,
		},
	}

	remote_refill.State = autogen.RemoteRefillProcessed
	remote_refill.RefillId = &refill.Id
	remote_refill.OrderId = order.Id

	user.Balance += int64(remote_refill.Amount)

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		
		err := s.DBackend.CreateRefill(ctx, refill)
		if err != nil {
			return nil, errors.New("failed to create refill")
		}

		err = s.DBackend.UpdateAccount(ctx, user)
		if err != nil {
			return nil, errors.New("failed to update account")
		}

		err = s.DBackend.UpdateRemoteRefill(ctx, remote_refill)
		if err != nil {
			return nil, errors.New("failed to update remote refill")
		}
		return nil, nil
	})

	if err != nil {
		logrus.Error(err)
		autogen.SelfValidateRemoteRefill500JSONResponse{
			ErrorCode: autogen.ErrInternalServerError,
			Message: autogen.MsgInternalServerError,
		}.VisitSelfValidateRemoteRefillResponse(c.Response())
		return err
	}

	if c.Response().Committed {
		return nil
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

	// TODO
	logrus.Panic("Not implemented")

	return nil
}