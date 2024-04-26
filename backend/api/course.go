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

	data, err := s.DBackend.GetItems(c.Request().Context(), "", 0, 0, "", "", search)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	for _, item := range data {
		var amount_needed uint64 = item.OptimalAmount - item.AmountLeft
		if amount_needed > 0 && item.AmountPerBundle != nil {
			course = append(course, autogen.CourseItem{
				AmountToBuy: amount_needed / *item.AmountPerBundle + (amount_needed % *item.AmountPerBundle) * 2 / *item.AmountPerBundle,
				Item:        item.Item,
			})
		}
	}

	autogen.GetCourse200JSONResponse{
		Items: course,
	}.VisitGetCourseResponse(c.Response())
	return nil
}
