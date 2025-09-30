package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/repository"
)

type ArticleRepository struct {
	Conn *sql.DB
}

// NewArticleRepository will create an object that represent the article.Repository interface
func NewArticleRepository(conn *sql.DB) *ArticleRepository {
	return &ArticleRepository{conn}
}

func (m *ArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Article, 0)
	for rows.Next() {
		t := domain.Article{}
		var authorID uuid.UUID
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Slug,
			&t.Content,
			&t.Thumbnail,
			&t.Image,
			&t.ShortDescription,
			&t.Keywords,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = domain.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *ArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, keywords, author_id, updated_at, created_at
  						FROM article WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}
func (m *ArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (res domain.Article, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, keywords, author_id, updated_at, created_at
  						FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *ArticleRepository) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, keywords, author_id, updated_at, created_at
  						FROM article WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *ArticleRepository) GetBySlug(ctx context.Context, slug string) (res domain.Article, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, keywords, author_id, updated_at, created_at
  						FROM article WHERE slug = ?`

	list, err := m.fetch(ctx, query, slug)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *ArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	query := `INSERT article SET id=?, title=?, slug=?, content=?, thumbnail=?, image=?, short_description=?, keywords=?, author_id=?, updated_at=?, created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	// Generate UUID if not set
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	_, err = stmt.ExecContext(ctx, a.ID, a.Title, a.Slug, a.Content, a.Thumbnail, a.Image, a.ShortDescription, a.Keywords, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return
	}
	return
}

func (m *ArticleRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
func (m *ArticleRepository) Update(ctx context.Context, ar *domain.Article) (err error) {
	query := `UPDATE article set title=?, slug=?, content=?, thumbnail=?, image=?, short_description=?, keywords=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Slug, ar.Content, ar.Thumbnail, ar.Image, ar.ShortDescription, ar.Keywords, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (m *ArticleRepository) SlugExistsExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM article WHERE slug = ? AND id != ?`
	var count int
	err := m.Conn.QueryRowContext(ctx, query, slug, excludeID).Scan(&count)
	return count > 0, err
}
