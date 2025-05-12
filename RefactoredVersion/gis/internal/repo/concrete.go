package repo

import (
	"context"
	"errors"
	"github.com/SwanHtetAungPhyo/gis/internal/model/model"
	"github.com/SwanHtetAungPhyo/gis/internal/model/req"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GigRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewGigRepository(db *gorm.DB, log *logrus.Logger) *GigRepository {
	return &GigRepository{db: db, log: log}
}

// CreateGig handles the complete gig creation flow
func (r *GigRepository) CreateGig(ctx context.Context, req *req.CreateGigRequest) (*model.Gig, error) {
	var gig *model.Gig

	err := r.db.Transaction(func(tx *gorm.DB) error {
		category, err := r.getCategory(tx, req.CategoryID)
		if err != nil {
			return err
		}

		// 2. Process tags (create if not exists)
		tags, err := r.processTags(tx, req.Tags)
		if err != nil {
			return err
		}

		// 3. Prepare gig model
		gig = &model.Gig{
			Title:       req.Title,
			Description: req.Description,
			IsActive:    true,
			CategoryID:  category.ID,
			SellerID:    uuid.MustParse(req.SellerID),
			Tags:        tags,
		}

		if err := r.addPackages(tx, gig, req.Packages); err != nil {
			return err
		}

		if err := r.addImages(tx, gig, req.Images); err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Create(gig).Error; err != nil {
			r.log.WithError(err).Error("Failed to create gig")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return gig, nil
}

func (r *GigRepository) getCategory(tx *gorm.DB, categoryName string) (*model.Category, error) {
	var category model.Category
	if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
		r.log.WithField("category", categoryName).Error("Category not found")
		return nil, errors.New("invalid category")
	}
	return &category, nil
}

func (r *GigRepository) processTags(tx *gorm.DB, tagNames []string) ([]model.GigTag, error) {
	var tags []model.GigTag

	for _, name := range tagNames {
		var tag model.GigTag
		err := tx.Where("label = ?", name).First(&tag).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = model.GigTag{Label: name}
			if err := tx.Create(&tag).Error; err != nil {
				r.log.WithField("tag", name).Error("Failed to create tag")
				continue
			}
		} else if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *GigRepository) addPackages(tx *gorm.DB, gig *model.Gig, packages []req.GigPackageRequest) error {
	for _, pkg := range packages {
		gigPackage := model.GigPackage{
			Title:        pkg.Title,
			Description:  pkg.Description,
			Price:        pkg.Price,
			DeliveryTime: pkg.DeliveryDays,
		}

		// Add package features
		for _, feature := range pkg.Features {
			gigPackage.Features = append(gigPackage.Features, model.GigPackageFeature{
				Title:       feature.Title,
				Description: &feature.Description,
				Included:    feature.Included,
			})
		}

		gig.Packages = append(gig.Packages, gigPackage)
	}
	return nil
}

func (r *GigRepository) addImages(tx *gorm.DB, gig *model.Gig, images []req.GigImageRequest) error {
	if len(images) == 0 {
		return nil
	}

	for _, img := range images {
		gig.Images = append(gig.Images, model.GigImage{
			URL:       img.URL,
			Alt:       &img.AltText,
			IsPrimary: img.IsPrimary,
		})
	}
	return nil
}

// UpdateGig updates gig metadata and returns the updated gig
func (r *GigRepository) UpdateGig(
	ctx context.Context,
	request *req.UpdateGigRequest,
) (*model.Gig, error) {
	if request.GigId == uuid.Nil {
		return nil, errors.New("gig ID cannot be empty")
	}

	var existingGig model.Gig
	if err := r.db.WithContext(ctx).
		Where("id = ? AND sellerId = ?", request.GigId, request.SellerId).
		First(&existingGig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gig not found or unauthorized")
		}
		r.log.WithError(err).Error("failed to verify gig ownership")
		return nil, err
	}

	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Description != "" {
		updates["description"] = request.Description
	}
	if request.IsActive != nil {
		updates["is_active"] = *request.IsActive
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Gig{}).
		Where("id = ?", request.GigId).
		Updates(updates).Error; err != nil {
		r.log.WithFields(logrus.Fields{
			"gigId": request.GigId,
			"error": err,
		}).Error("failed to update gig")
		return nil, err
	}

	var updatedGig model.Gig
	if err := r.db.WithContext(ctx).
		Preload("Packages").
		Preload("Images").
		First(&updatedGig, "id = ?", request.GigId).Error; err != nil {
		r.log.WithError(err).Error("failed to fetch updated gig")
		return nil, err
	}

	return &updatedGig, nil
}

// PartialUpdate Generic partial update helper
func (r *GigRepository) PartialUpdate(
	ctx context.Context,
	gigId uuid.UUID,
	updates map[string]interface{},
) (*model.Gig, error) {
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	result := r.db.WithContext(ctx).
		Model(&model.Gig{}).
		Where("id = ?", gigId).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("gig not found")
	}

	var gig model.Gig
	if err := r.db.First(&gig, "id = ?", gigId).Error; err != nil {
		return nil, err
	}
	return &gig, nil
}

func (r *GigRepository) CountTheGig() (*int64, error) {
	var total int64
	if err := r.db.
		WithContext(context.TODO()).
		Model(&model.Gig{}).
		Count(&total).Error; err != nil {
		r.log.WithError(err).Error("failed to count the gig")
		return nil, errors.New("failed to count the gig")
	}
	return &total, nil
}

func (r *GigRepository) GetByOffset(total *int64, page, perPage int, offset int) ([]*model.Gig, error) {

	var gigs []*model.Gig
	if err := r.db.WithContext(context.TODO()).
		Model(&model.Gig{}).
		Limit(perPage).
		Offset(offset).
		Find(&gigs).Error; err != nil {
		r.log.WithError(err).Error("failed to get gigs", "offset", offset, "total", total)
		return nil, err
	}
	return gigs, nil
}

func (r *GigRepository) AddPackageToGig(id uuid.UUID, request *req.GigPackageRequest) (*model.GigPackage, error) {
	pkgToCreate := &model.GigPackage{
		Title:        request.Title,
		Description:  request.Description,
		Price:        request.Price,
		DeliveryTime: request.DeliveryDays,
		GigID:        id,
	}
	pkgToCreate.Features = make([]model.GigPackageFeature, len(request.Features))
	for _, feature := range request.Features {
		pkgToCreate.Features = append(pkgToCreate.Features, model.GigPackageFeature{
			Title:       feature.Title,
			Description: &feature.Description,
			Included:    feature.Included,
		})
	}

	if err := r.db.
		WithContext(context.TODO()).
		Model(&model.GigPackage{}).
		Create(pkgToCreate).Error; err != nil {
		r.log.WithError(err).Error("failed to create gig package")
		return nil, err
	}

	var gigPkg model.GigPackage
	if err := r.db.
		WithContext(context.TODO()).
		Model(&model.GigPackage{}).
		Preload("Features").First(&gigPkg, "id = ?", pkgToCreate.ID).Error; err != nil {
		r.log.WithError(err).Error("failed to fetch gig package")
	}
	return &gigPkg, nil
}
