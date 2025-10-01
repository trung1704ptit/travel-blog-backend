package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
)

type CategoryService interface {
	Fetch(ctx context.Context, page, limit int) ([]domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (domain.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error)
	Update(ctx context.Context, cat *domain.Category) error
	Store(context.Context, *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error)
	GetRootCategories(ctx context.Context) ([]domain.Category, error)
	GetCategoryTree(ctx context.Context) ([]domain.Category, error)
	GetCategoryWithChildren(ctx context.Context, slug string) (domain.Category, error)
}

type CategoryHandler struct {
	Category CategoryService
}

const defaultCategoryLimit = 100
const defaultCategoryPage = 1

func NewCategoryHandler(e *echo.Echo, svc CategoryService) {
	handler := &CategoryHandler{
		Category: svc,
	}

	e.GET("/categories", handler.FetchCategory)
	e.GET("/categories/tree", handler.GetCategoryTree)
	e.GET("/categories/roots", handler.GetRootCategories)
	e.POST("/categories", handler.Store)
	e.PATCH("/categories/:id", handler.Update)
	e.GET("/categories/:slug", handler.GetBySlug)
	e.GET("/categories/:id", handler.GetByID)
	e.GET("/categories/:slug/children", handler.GetChildren)
	e.DELETE("/categories/:id", handler.Delete)
}

func (cat *CategoryHandler) FetchCategory(c echo.Context) error {
	// Parse page parameter
	pageS := c.QueryParam("page")
	page, err := strconv.Atoi(pageS)
	if err != nil || page < 1 {
		page = defaultCategoryPage
	}

	// Parse limit parameter
	limitS := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitS)
	if err != nil || limit < 1 {
		limit = defaultCategoryLimit
	}

	ctx := c.Request().Context()

	listCat, err := cat.Category.Fetch(ctx, page, limit)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, listCat)
}

// GetBySlug will get category by given slug
func (cat *CategoryHandler) GetBySlug(c echo.Context) error {
	slug := c.Param("slug")

	if slug == "" {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Slug is required"})
	}

	ctx := c.Request().Context()

	category, err := cat.Category.GetCategoryWithChildren(ctx, slug)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, category)
}

// GetByID will get category by given ID
func (cat *CategoryHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid UUID format"})
	}

	ctx := c.Request().Context()

	category, err := cat.Category.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, category)
}

func isCategoryRequestValid(m *domain.Category) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the category by given request body
func (cat *CategoryHandler) Store(c echo.Context) (err error) {
	var category domain.Category
	err = c.Bind(&category)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isCategoryRequestValid(&category); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Generate UUID if not provided
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	ctx := c.Request().Context()
	err = cat.Category.Store(ctx, &category)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, category)
}

// Update will update the category by given request body (PATCH - partial update)
func (cat *CategoryHandler) Update(c echo.Context) (err error) {
	idStr := c.Param("id")
	categoryID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid UUID format"})
	}

	// Get existing category
	ctx := c.Request().Context()
	existingCategory, err := cat.Category.GetByID(ctx, categoryID)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	// Parse partial update data
	updateData := make(map[string]interface{})
	err = c.Bind(&updateData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// Apply partial updates to existing category
	updatedCategory := existingCategory
	if name, ok := updateData["name"].(string); ok {
		updatedCategory.Name = name
	}
	if slug, ok := updateData["slug"].(string); ok {
		updatedCategory.Slug = slug
	}
	if description, ok := updateData["description"].(string); ok {
		updatedCategory.Description = description
	}
	if image, ok := updateData["image"].(string); ok {
		updatedCategory.Image = image
	}
	if parentID, ok := updateData["parent_id"].(string); ok {
		if parsedParentID, err := uuid.Parse(parentID); err == nil {
			updatedCategory.ParentID = &parsedParentID
		}
	}

	// Validate updated category
	var ok bool
	if ok, err = isCategoryRequestValid(&updatedCategory); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = cat.Category.Update(ctx, &updatedCategory)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, updatedCategory)
}

// Delete will delete category by given slug
func (cat *CategoryHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid UUID format"})
	}
	ctx := c.Request().Context()

	// First get the category to get its ID
	category, err := cat.Category.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	// Then delete by ID
	err = cat.Category.Delete(ctx, category.ID)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// GetCategoryTree retrieves the complete category tree
func (cat *CategoryHandler) GetCategoryTree(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := cat.Category.GetCategoryTree(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, categories)
}

// GetRootCategories retrieves all root categories
func (cat *CategoryHandler) GetRootCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := cat.Category.GetRootCategories(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, categories)
}

// GetChildren retrieves children of a category
func (cat *CategoryHandler) GetChildren(c echo.Context) error {
	slug := c.Param("slug")

	if slug == "" {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Slug is required"})
	}

	ctx := c.Request().Context()

	// First get the category to get its ID
	category, err := cat.Category.GetBySlug(ctx, slug)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	// Then get its children
	children, err := cat.Category.GetChildren(ctx, category.ID)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, children)
}
