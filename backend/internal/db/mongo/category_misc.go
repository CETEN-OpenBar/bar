package mongo

import "bar/internal/models"

func (b *Backend) GetAllCategories() ([]*models.Category, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var categories []*models.Category
	cursor, err := b.db.Collection(CategoriesCollection).Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
