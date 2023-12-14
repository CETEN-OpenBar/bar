package api

import (
	"bar/autogen"
	"bar/internal/models"
	"time"

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

	count, err := s.DBackend.CountAllRestocks(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetAllRestocks(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	restocks := make([]autogen.Restock, len(data))
	for i, r := range data {
		restocks[i] = r.Restock
	}

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
	usr, err := MustGetAdmin(c)
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
			DriverId:      body.DriverId,
			Id:            uuid.New(),
			TotalCostHt:   body.TotalCostHt,
			TotalCostTtc:  body.TotalCostTtc,
			Type:          body.Type,
			CreatedAt:     uint64(time.Now().Unix()),
			CreatedBy:     usr.Id,
			CreatedByName: usr.Name(),
		},
		CreatedAt: time.Now().Unix(),
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
					return nil, ErrorItemNotFound(c)
				}
				logrus.Error(err)
				return nil, Error500(c)
			}

			restockItem.ItemName = item.Name
			restockItem.ItemPictureUri = item.PictureUri

			item.AmountLeft += restockItem.AmountOfBundle * restockItem.AmountPerBundle
			item.LastTva = &restockItem.Tva
			item.Prices.Coutant = uint64((10000.0+float64(restockItem.Tva)) * float64(restockItem.BundleCostHt) / (10000.0 * float64(restockItem.AmountPerBundle)))

			err = s.DBackend.UpdateItem(c.Request().Context(), item)
			if err != nil {
				logrus.Error(err)
				return nil, Error500(c)
			}

			restock.Items = append(restock.Items, restockItem)
		}

		err = s.DBackend.CreateRestock(c.Request().Context(), &restock)
		if err != nil {
			logrus.Error(err)
			return nil, Error500(c)
		}

		return nil, nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	logrus.WithField("restock", restock.Id.String()).WithField("by", usr.Name()).Info("Created restock")
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

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		// Remove restock from all items
		restock, err := s.DBackend.GetRestock(c.Request().Context(), restockId.String())
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, ErrorRestockNotFound(c)
			}
			logrus.Error(err)
			return nil, Error500(c)
		}

		for _, item := range restock.Items {
			i, err := s.DBackend.GetItem(c.Request().Context(), item.ItemId.String())
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, ErrorItemNotFound(c)
				}
				logrus.Error(err)
				return nil, Error500(c)
			}

			i.AmountLeft -= item.AmountOfBundle * item.AmountPerBundle

			err = s.DBackend.UpdateItem(c.Request().Context(), i)
			if err != nil {
				logrus.Error(err)
				return nil, Error500(c)
			}
		}

		// Get restock from database
		err = s.DBackend.MarkDeleteRestock(c.Request().Context(), restockId.String(), account.Id.String())
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, ErrorRestockNotFound(c)
			}
			logrus.Error(err)
			return nil, Error500(c)
		}
		return nil, nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if c.Response().Committed {
		return nil
	}

	logrus.WithField("restock", restockId.String()).WithField("by", account.Name()).Info("Restock marked for deletion")
	autogen.DeleteRestock204Response{}.VisitDeleteRestockResponse(c.Response())
	return nil
}
