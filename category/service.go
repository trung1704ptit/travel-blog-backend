package category

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/bxcodec/go-clean-arch/domain"
)

// CategoryRepository represent the category's repository contract
//
//go:generate mockery --name CategoryRepository
type CategoryRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Category, string, error)
	GetBySlug(ctx context.Context, slug string) (domain.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Store(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	SlugExistsExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error)
	GetByArticleID(ctx context.Context, articleID uuid.UUID) ([]domain.Category, error)
	GetByIDs(ctx context.Context, categoryIDs []uuid.UUID) ([]domain.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error)
	GetRootCategories(ctx context.Context) ([]domain.Category, error)
	GetCategoryTree(ctx context.Context) ([]domain.Category, error)
}

type Service struct {
	categoryRepo CategoryRepository
}

// NewService will create a new category service object
func NewService(cr CategoryRepository) *Service {
	return &Service{
		categoryRepo: cr,
	}
}

func (c *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Category, nextCursor string, err error) {
	res, nextCursor, err = c.categoryRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	return
}

func (c *Service) GetBySlug(ctx context.Context, slug string) (res domain.Category, err error) {
	res, err = c.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return
	}
	return
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (res domain.Category, err error) {
	res, err = c.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (c *Service) Update(ctx context.Context, category *domain.Category) (err error) {
	category.UpdatedAt = time.Now()
	return c.categoryRepo.Update(ctx, category)
}

func (c *Service) Store(ctx context.Context, category *domain.Category) (err error) {
	// Check if category with same slug already exists
	existedCategory, _ := c.GetBySlug(ctx, category.Slug) // ignore if any error
	if existedCategory.ID != uuid.Nil {
		return domain.ErrConflict
	}

	if category.Slug == "" {
		category.Slug = generateSlug(category.Name)
	}

	category.Slug = c.ensureUniqueSlug(ctx, category.Slug, uuid.Nil)

	// Generate UUID if not set
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	err = c.categoryRepo.Store(ctx, category)
	return
}

func (c *Service) Delete(ctx context.Context, id uuid.UUID) (err error) {
	// Check if category exists by trying to delete it
	// The repository will return ErrNotFound if it doesn't exist
	return c.categoryRepo.Delete(ctx, id)
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

func (c *Service) ensureUniqueSlug(ctx context.Context, baseSlug string, excludeID uuid.UUID) string {
	slug := baseSlug
	counter := 1

	for {
		exists, err := c.categoryRepo.SlugExistsExcludingID(ctx, slug, excludeID)
		if err != nil || !exists {
			break
		}

		// Generate new slug with counter
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug
}

// GetChildren retrieves all children of a category
func (c *Service) GetChildren(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error) {
	return c.categoryRepo.GetChildren(ctx, parentID)
}

// GetRootCategories retrieves all root categories
func (c *Service) GetRootCategories(ctx context.Context) ([]domain.Category, error) {
	return c.categoryRepo.GetRootCategories(ctx)
}

// GetCategoryTree retrieves the complete category tree
func (c *Service) GetCategoryTree(ctx context.Context) ([]domain.Category, error) {
	return c.categoryRepo.GetCategoryTree(ctx)
}

// GetCategoryWithChildren retrieves a category with its children
func (c *Service) GetCategoryWithChildren(ctx context.Context, slug string) (domain.Category, error) {
	category, err := c.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return domain.Category{}, err
	}

	// Get children
	children, err := c.categoryRepo.GetChildren(ctx, category.ID)
	if err != nil {
		return domain.Category{}, err
	}

	category.Children = children
	return category, nil
}
