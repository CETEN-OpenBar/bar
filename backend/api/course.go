package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// (GET /course)
func (s *Server) GetCourse(c echo.Context, params autogen.GetCourseParams) error {
	// Get admin account from cookie
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}
	search := ""
	if params.Fournisseur != nil {
		search = *params.Fournisseur
	}
	var course []autogen.CourseItem

	data, err := s.DBackend.GetItems(c.Request().Context(), "", 0, 1000, "", "", search, false)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	for _, item := range data {
		var amount_needed = item.OptimalAmount - item.AmountLeft
		if item.OptimalAmount > item.AmountLeft && item.AmountPerBundle != nil {
			amountToBuy := amount_needed / *item.AmountPerBundle + (amount_needed%*item.AmountPerBundle)*2 / *item.AmountPerBundle
			if amountToBuy > 0 {
				course = append(course, autogen.CourseItem{
					AmountToBuy: amountToBuy,
					Item:        item.Item,
				})
			}
		}
	}

	autogen.GetCourse200JSONResponse{
		Items: course,
	}.VisitGetCourseResponse(c.Response())
	return nil
}
