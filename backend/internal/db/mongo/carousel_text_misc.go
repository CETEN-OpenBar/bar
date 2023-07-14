package mongo

import "bar/internal/models"

func (b *Backend) GetAllCarouselTexts() ([]*models.CarouselText, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var carouselTexts []*models.CarouselText
	cursor, err := b.db.Collection(CarouselTextsCollection).Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &carouselTexts); err != nil {
		return nil, err
	}

	return carouselTexts, nil
}
