package api

import (
	"bar/autogen"
	"bar/autogen/helloasso"
	"bar/internal/config"
	"bar/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// (GET /remote-refills/start)
func (s *Server) StartRemoteRefill(c echo.Context, params autogen.StartRemoteRefillParams) error {
	// Get account from cookie
	user, err := MustGetUser(c)
	if err != nil {
		return nil
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
				BackUrl:           conf.ApiConfig.FrontendBasePath + "/client/index",
				ErrorUrl:          conf.ApiConfig.FrontendBasePath + "/remote-refill/error",
				ReturnUrl:         conf.ApiConfig.FrontendBasePath + "/remote-refill/success",
				InitialAmount:     int32(params.Amount),
				TotalAmount:       int32(params.Amount),
				ContainsDonation:  false,
				ItemName:          "Rechargement Bar",
				Metadata:          &metadata,
			},
		)
	
	if err != nil {
		return autogen.StartRemoteRefill500JSONResponse{}.VisitStartRemoteRefillResponse(c.Response())
	}

	if resp.StatusCode() != 200 || resp.JSON200 == nil {
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
