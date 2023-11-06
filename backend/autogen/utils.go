package autogen

import "time"

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
	case AccountPriceExte:
		price = i.Prices.Exte
	case AccountPriceInterne:
		price = i.Prices.Interne
	case AccountPriceMembreBureau:
		price = i.Prices.MembreBureau
	case AccountPriceMembrePrivilegie:
		price = i.Prices.MembrePrivilegie
	case AccountPriceStaff:
		price = i.Prices.Staff
	case AccountPriceVIP:
		price = i.Prices.Vip
	default:
		price = i.Prices.Exte
	}

	if i.Promotion == nil {
		return price
	}
	if i.PromotionEndsAt != nil && uint64(time.Now().Unix()) > *i.PromotionEndsAt {
		return price
	}

	if price == 0 {
		price = i.Prices.Exte
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
		Ceten:            uint64(float64(i.Prices.Ceten) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Exte:             uint64(float64(i.Prices.Exte) * (1.0 - (float64(*i.Promotion) / 100.0))),
		MembreBureau:     uint64(float64(i.Prices.MembreBureau) * (1.0 - (float64(*i.Promotion) / 100.0))),
		MembrePrivilegie: uint64(float64(i.Prices.MembrePrivilegie) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Staff:            uint64(float64(i.Prices.Staff) * (1.0 - (float64(*i.Promotion) / 100.0))),
		Vip:              uint64(float64(i.Prices.Vip) * (1.0 - (float64(*i.Promotion) / 100.0))),
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
