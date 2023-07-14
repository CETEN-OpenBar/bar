package mongo

import "bar/internal/models"

func (b *Backend) GetAllCarouselImages() ([]*models.CarouselImage, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var carouselImages []*models.CarouselImage
	cursor, err := b.db.Collection(CarouselImagesCollection).Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &carouselImages); err != nil {
		return nil, err
	}

	return carouselImages, nil
}
