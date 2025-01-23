package api

import (
	"bar/autogen"
	"bar/internal/webhook"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// PATCH /restocks/{restock_id}
func (s *Server) UpdateRestock(c echo.Context, restockId autogen.UUID) error {
	usr, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var body autogen.UpdateRestockJSONRequestBody
	if err := c.Bind(&body); err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	restock, err := s.DBackend.GetRestock(c.Request().Context(), restockId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	restock.Items = []autogen.RestockItem{}
	restock.TotalCostHt = body.TotalCostHt
	restock.TotalCostTtc = body.TotalCostTtc
	restock.Type = body.Type
	restock.CreatedAt = uint64(time.Now().Unix())

	oldState := restock.State
	if restock.State != autogen.RestockFinished {
		restock.State = body.State // Cannot go from finished to anything else
	}

	for _, newRestockItem := range body.Items {
		restockItem := autogen.RestockItem{
			AmountOfBundle:  newRestockItem.AmountOfBundle,
			AmountPerBundle: newRestockItem.AmountPerBundle,
			BundleCostHt:    newRestockItem.BundleCostHt,
			BundleCostTtc:   newRestockItem.BundleCostTtc,
			ItemId:          newRestockItem.ItemId,
			Tva:             newRestockItem.Tva,
			ItemName:        newRestockItem.ItemName,
		}

		item, err := s.DBackend.GetItem(c.Request().Context(), restockItem.ItemId.String())
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
		restockItem.ItemPictureUri = item.PictureUri

		restock.Items = append(restock.Items, restockItem)

		if oldState != autogen.RestockFinished && body.State == autogen.RestockFinished {

			category, err := s.DBackend.GetCategory(c.Request().Context(), item.CategoryId.String())
			if err != nil {
				logrus.Error(err)
				return Error500(c)
			}
			item = UpdateItem(item, category, restockItem)
			err = s.DBackend.UpdateItem(c.Request().Context(), item)
			if err != nil {
				logrus.Error(err)
				return Error500(c)
			}
			logrus.WithField("restock", restock.Id.String()).WithField("by", usr.Name()).Info("Items updated")
		}
	}
	err = s.DBackend.UpdateRestock(c.Request().Context(), restock)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}
	if oldState != autogen.RestockFinished && body.State == autogen.RestockFinished && restock.Type == autogen.RestockViennoiserie {
		err := webhook.SendDiscordWebhook(*restock)
		if err != nil {
			logrus.Errorf("Error sending webhook: %v\n", err)
		}
	}

	if oldState == body.State {
		return nil
	}

	logrus.WithField("restock", restock.Id.String()).WithField("by", usr.Name()).Info("Patch restock")
	return nil
}
