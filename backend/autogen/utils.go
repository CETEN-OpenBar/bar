package autogen

import (
	"math"
	"time"
)

func OptionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (i *Item) RealPrice(r AccountPriceRole) uint64 {
	// TODO: modify this when modifying price roles
	var price uint64

	// Get price from the role of the user
	switch r {
	case AccountPriceCeten:
		price = i.Prices.Ceten
	case AccountPriceExterne:
		price = i.Prices.Externe
	case AccountPriceCoutant:
		price = i.Prices.Coutant
	case AccountPricePrivilegies:
		price = i.Prices.Privilegies
	case AccountPriceStaffBar:
		price = i.Prices.StaffBar
	case AccountPriceMenu:
		price = i.Prices.Menu
	default:
		price = i.Prices.Externe
	}

	if i.Promotion == nil {
		return price
	}
	if i.PromotionEndsAt != nil && uint64(time.Now().Unix()) > *i.PromotionEndsAt {
		return price
	}

	if price == 0 {
		price = i.Prices.Externe
	}

	return uint64(float64(price) * (1.0 - (float64(*i.Promotion) / 100.0)))
}

func (i *Item) RealPrices() ItemPrices {
	// TODO: modify this when modifying price roles
	prices := i.Prices

	if i.Promotion == nil {
		return prices
	}
	if i.PromotionEndsAt != nil && uint64(time.Now().Unix()) > *i.PromotionEndsAt {
		return prices
	}

	return ItemPrices{
		Ceten:       uint64(float64(i.Prices.Ceten) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Externe:     uint64(float64(i.Prices.Externe) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Coutant:     uint64(float64(i.Prices.Coutant) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Privilegies: uint64(float64(i.Prices.Privilegies) * (1.0 - (float64(*i.Promotion) / 100.0))),
		StaffBar:    uint64(float64(i.Prices.StaffBar) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Menu:        uint64(float64(i.Prices.Menu) * (1.0 - (float64(*i.Promotion) / 100.0))),
	}
}

func (a *Account) Name() string {
	if a.FirstName != "" && a.LastName != "" {
		return a.FirstName + " " + a.LastName
	}
	if a.EmailAddress != "" {
		return a.EmailAddress
	}
	return a.Id.String()
}

func (a *Account) HasPrivileges() bool {
	// TODO: modify this when modifying roles
	return a.Role == "admin" || a.Role == "superadmin" || a.Role == "member" || a.Role == "ghost"
}

func Pager(page *uint64, limit *uint64, count *uint64) (dbPage uint64, pageOut uint64, limitOut uint64, maxPage uint64) {
	// We make sure that page is filled, it's going to be 0 otherwise
	if page != nil {
		pageOut = *page
	}
	// We make sure that limit is filled, it's going to be 0 otherwise
	if limit != nil {
		limitOut = *limit
	}
	// We block the limit to 100
	if limitOut > 100 {
		limitOut = 100
	}

	// We calculate the max page (0 if count is nil)
	if count != nil {
		maxPageFloat := float64(*count) / float64(limitOut)
		maxPage = uint64(math.Ceil(maxPageFloat) - 1)
	}

	if pageOut > maxPage+1 {
		pageOut = maxPage
	}
	if pageOut > 0 {
		pageOut--
	}

	dbPage = pageOut
	pageOut++
	return
}
