package api

import (
	"bar/autogen"
	"bar/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /restocks)
func (s *Server) GetRestocks(c echo.Context, params autogen.GetRestocksParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var page uint64
	if params.Page != nil {
		page = *params.Page
	}
	if page > 0 {
		page--
	}

	var limit uint64 = 10
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > 100 {
		limit = 100
	}

	count, err := s.DBackend.CountAllRestocks(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetAllRestocks(c.Request().Context(), page, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	restocks := make([]autogen.Restock, len(data))
	for i, r := range data {
		restocks[i] = r.Restock
	}

	page++
	autogen.GetRestocks200JSONResponse{
		Restocks: restocks,
		Limit:    limit,
		Page:     page,
		MaxPage:  maxPage,
	}.VisitGetRestocksResponse(c.Response())
	return nil
}

// (POST /restocks)
func (s *Server) CreateRestock(c echo.Context) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var body autogen.CreateRestockJSONRequestBody
	if err := c.Bind(&body); err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	restock := models.Restock{
		Restock: autogen.Restock{
			DriverId:     body.DriverId,
			Id:           uuid.New(),
			TotalCostHt:  body.TotalCostHt,
			TotalCostTtc: body.TotalCostTtc,
			Type:         body.Type,
		},
	}

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		for _, item := range body.Items {
			restockItem := autogen.RestockItem{
				AmountOfBundle:  item.AmountOfBundle,
				AmountPerBundle: item.AmountPerBundle,
				BundleCostHt:    item.BundleCostHt,
				ItemId:          item.ItemId,
				Tva:             item.Tva,
			}

			// Get item to check if it exists
			item, err := s.DBackend.GetItem(c.Request().Context(), item.ItemId.String())
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return ErrorItemNotFound(c)
				}
				logrus.Error(err)
				return Error500(c)
			}

			restockItem.ItemName = item.Name
			restockItem.ItemPictureUri = item.PictureUri

			item.AmountLeft += restockItem.AmountOfBundle * restockItem.AmountPerBundle

			restock.Items = append(restock.Items, restockItem)
		}

		err = s.DBackend.CreateRestock(c.Request().Context(), &restock)
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if c.Response().Committed {
		return nil
	}

	autogen.CreateRestock201JSONResponse(restock.Restock).VisitCreateRestockResponse(c.Response())
	return nil
}

// (DELETE /restocks/{restock_id})
func (s *Server) DeleteRestock(c echo.Context, restockId autogen.UUID) error {
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	// Get restock from database
	err = s.DBackend.MarkDeleteRestock(c.Request().Context(), restockId.String(), account.Id.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRestockNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Restock %s marked as deleted by %s", restockId.String(), account.Id)
	autogen.DeleteRestock204Response{}.VisitDeleteRestockResponse(c.Response())
	return nil
}
