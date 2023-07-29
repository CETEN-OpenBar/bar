package autogen

import "time"

func (i *Item) RealPrice() uint64 {
	if i.Promotion == nil {
		return i.Price
	}
	if i.PromotionEndsAt != nil && uint64(time.Now().Unix()) > *i.PromotionEndsAt {
		return i.Price
	}
	return uint64(float64(i.Price) * (1.0 - (float64(*i.Promotion) / 100.0)))
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
