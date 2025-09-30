package article

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ArticleRepository represent the article's repository contract
//
//go:generate mockery --name ArticleRepository
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Article, error)
	GetBySlug(ctx context.Context, title string) (domain.Article, error)
	Update(ctx context.Context, ar *domain.Article) error
	Store(ctx context.Context, a *domain.Article) error
	Delete(ctx context.Context, id uuid.UUID) error
	SlugExistsExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error)
}

// AuthorRepository represent the author's repository contract
//
//go:generate mockery --name AuthorRepository
type AuthorRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Author, error)
}

type Service struct {
	articleRepo ArticleRepository
	authorRepo  AuthorRepository
}

// NewService will create a new article service object
func NewService(a ArticleRepository, ar AuthorRepository) *Service {
	return &Service{
		articleRepo: a,
		authorRepo:  ar,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *Service) fillAuthorDetails(ctx context.Context, data []domain.Article) ([]domain.Article, error) {
	g, ctx := errgroup.WithContext(ctx)
	// Get the author's id
	mapAuthors := map[uuid.UUID]domain.Author{}

	for _, article := range data { //nolint
		mapAuthors[article.Author.ID] = domain.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		defer close(chanAuthor)
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}

	}()

	for author := range chanAuthor {
		if author.ID != uuid.Nil {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data { //nolint
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *Service) GetByID(ctx context.Context, id uuid.UUID) (res domain.Article, err error) {
	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}
	res.Author = resAuthor
	return
}

func (a *Service) Update(ctx context.Context, ar *domain.Article) (err error) {
	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *Service) GetBySlug(ctx context.Context, slug string) (res domain.Article, err error) {
	res, err = a.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *Service) Store(ctx context.Context, m *domain.Article) (err error) {
	existedArticle, _ := a.GetBySlug(ctx, m.Slug) // ignore if any error
	if existedArticle.ID != uuid.Nil {
		return domain.ErrConflict
	}

	if m.Slug == "" {
		m.Slug = generateSlug(m.Title)
	}

	m.Slug = a.ensureUniqueSlug(ctx, m.Slug, uuid.Nil)

	// Generate UUID if not set
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *Service) Delete(ctx context.Context, id uuid.UUID) (err error) {
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle.ID == uuid.Nil {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

func (a *Service) ensureUniqueSlug(ctx context.Context, baseSlug string, excludeID uuid.UUID) string {
	slug := baseSlug
	counter := 1

	for {
		exists, err := a.articleRepo.SlugExistsExcludingID(ctx, slug, excludeID)
		if err != nil || !exists {
			break
		}

		// Generate new slug with counter
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug
}
