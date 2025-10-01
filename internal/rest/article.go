package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleService represent the article's usecases
//
//go:generate mockery --name ArticleService
type ArticleService interface {
	Fetch(ctx context.Context, page, limit int) ([]domain.ArticleResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.ArticleResponse, error)
	Update(ctx context.Context, ar *domain.Article) error
	UpdatePartial(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	GetBySlug(ctx context.Context, slug string) (domain.ArticleResponse, error)
	Store(context.Context, *domain.Article) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ArticleHandler  represent the httphandler for article
type ArticleHandler struct {
	Service ArticleService
}

const defaultLimit = 100
const defaultPage = 1

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(e *echo.Echo, svc ArticleService) {
	handler := &ArticleHandler{
		Service: svc,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/articles", handler.Store)
	e.PATCH("/articles/:id", handler.Update)
	e.GET("/articles/:id", handler.GetByID)
	e.GET("/articles/slug/:slug", handler.GetBySlug)
	e.DELETE("/articles/:id", handler.Delete)
}

// FetchArticle will fetch the article based on given params
func (a *ArticleHandler) FetchArticle(c echo.Context) error {
	// Parse page parameter
	pageS := c.QueryParam("page")
	page, err := strconv.Atoi(pageS)
	if err != nil || page < 1 {
		page = defaultPage
	}

	// Parse limit parameter
	limitS := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitS)
	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	ctx := c.Request().Context()

	listAr, err := a.Service.Fetch(ctx, page, limit)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get article by given id
func (a *ArticleHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid UUID format"})
	}

	ctx := c.Request().Context()

	art, err := a.Service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, art)
}

// GetByID will get article by given id
func (a *ArticleHandler) GetBySlug(c echo.Context) error {
	slug := c.Param("slug")

	if slug == "" {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Slug is required"})
	}

	ctx := c.Request().Context()

	art, err := a.Service.GetBySlug(ctx, slug)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the article by given request body
func (a *ArticleHandler) Store(c echo.Context) (err error) {
	var article domain.Article
	err = c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Generate UUID if not provided
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	ctx := c.Request().Context()
	err = a.Service.Store(ctx, &article)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, article)
}

// Update will update the article by given request body (PATCH - partial update)
func (a *ArticleHandler) Update(c echo.Context) (err error) {
	idStr := c.Param("id")
	articleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid UUID format"})
	}

	// Parse partial update data
	updateData := make(map[string]interface{})
	err = c.Bind(&updateData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// Convert JSON data to proper types
	processedUpdates := make(map[string]interface{})

	processedUpdates["updated_at"] = time.Now()

	// Handle string fields
	if title, ok := updateData["title"].(string); ok {
		processedUpdates["title"] = title
	}
	if slug, ok := updateData["slug"].(string); ok {
		processedUpdates["slug"] = slug
	}
	if content, ok := updateData["content"].(string); ok {
		processedUpdates["content"] = content
	}
	if thumbnail, ok := updateData["thumbnail"].(string); ok {
		processedUpdates["thumbnail"] = thumbnail
	}
	if image, ok := updateData["image"].(string); ok {
		processedUpdates["image"] = image
	}
	if shortDesc, ok := updateData["short_description"].(string); ok {
		processedUpdates["short_description"] = shortDesc
	}
	if metaDesc, ok := updateData["meta_description"].(string); ok {
		processedUpdates["meta_description"] = metaDesc
	}

	// Handle numeric fields
	if readingTime, ok := updateData["reading_time_minutes"].(float64); ok {
		processedUpdates["reading_time_minutes"] = int(readingTime)
	}
	if views, ok := updateData["views"].(float64); ok {
		processedUpdates["views"] = int(views)
	}
	if likes, ok := updateData["likes"].(float64); ok {
		processedUpdates["likes"] = int(likes)
	}
	if comments, ok := updateData["comments"].(float64); ok {
		processedUpdates["comments"] = int(comments)
	}

	// Handle boolean fields
	if published, ok := updateData["published"].(bool); ok {
		processedUpdates["published"] = published
	}

	// Handle time fields
	if publishedAt, ok := updateData["published_at"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, publishedAt); err == nil {
			processedUpdates["published_at"] = &parsedTime
		}
	}

	// Handle author
	if authorData, ok := updateData["author"].(map[string]interface{}); ok {
		if authorID, ok := authorData["id"].(string); ok {
			if parsedAuthorID, err := uuid.Parse(authorID); err == nil {
				processedUpdates["author_id"] = parsedAuthorID
			}
		}
	}

	// Handle keywords (JSONStringSlice)
	if keywordsData, ok := updateData["keywords"]; ok {
		if keywordsSlice, ok := keywordsData.([]interface{}); ok {
			var keywords []string
			for _, keyword := range keywordsSlice {
				if keywordStr, ok := keyword.(string); ok {
					keywords = append(keywords, keywordStr)
				}
			}
			processedUpdates["keywords"] = domain.JSONStringSlice(keywords)
		}
	}

	// Handle tags (JSONStringSlice)
	if tagsData, ok := updateData["tags"]; ok {
		if tagsSlice, ok := tagsData.([]interface{}); ok {
			var tags []string
			for _, tag := range tagsSlice {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, tagStr)
				}
			}
			processedUpdates["tags"] = domain.JSONStringSlice(tags)
		}
	}

	// Handle categories
	if categoriesData, ok := updateData["categories"]; ok {
		if categoriesSlice, ok := categoriesData.([]interface{}); ok {
			var categories []domain.Category
			for _, categoryData := range categoriesSlice {
				if categoryMap, ok := categoryData.(map[string]interface{}); ok {
					var category domain.Category
					if id, ok := categoryMap["id"].(string); ok {
						if parsedID, err := uuid.Parse(id); err == nil {
							category.ID = parsedID
						}
					}
					if name, ok := categoryMap["name"].(string); ok {
						category.Name = name
					}
					if slug, ok := categoryMap["slug"].(string); ok {
						category.Slug = slug
					}
					if description, ok := categoryMap["description"].(string); ok {
						category.Description = description
					}
					if image, ok := categoryMap["image"].(string); ok {
						category.Image = image
					}
					categories = append(categories, category)
				}
			}
			processedUpdates["categories"] = categories
		}
	}

	// Use the new UpdatePartial method
	ctx := c.Request().Context()
	err = a.Service.UpdatePartial(ctx, articleID, processedUpdates)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	// Return the updated article
	updatedArticle, err := a.Service.GetByID(ctx, articleID)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.JSON(http.StatusOK, updatedArticle.Article)
}

// Delete will delete article by given param
func (a *ArticleHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid UUID format"})
	}

	ctx := c.Request().Context()

	err = a.Service.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), getErrorResponse(err))
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func getErrorResponse(err error) interface{} {
	if err == nil {
		return nil
	}

	// Return raw error message for all errors
	return ResponseError{Message: err.Error()}
}
