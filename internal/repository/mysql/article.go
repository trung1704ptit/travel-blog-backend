package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
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
			&t.MetaDescription,
			&t.Keywords,
			&t.Tags,
			&t.ReadingTimeMinutes,
			&t.Views,
			&t.Likes,
			&t.Comments,
			&t.Published,
			&t.PublishedAt,
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

func (m *ArticleRepository) Fetch(ctx context.Context, page, limit int) (res []domain.Article, err error) {
	// Calculate offset for pagination
	offset := (page - 1) * limit

	query := `SELECT id, title, slug, content, thumbnail, image, short_description, meta_description, keywords, tags, reading_time_minutes, views, likes, comments, published, published_at, author_id, updated_at, created_at
  						FROM article ORDER BY created_at DESC LIMIT ? OFFSET ? `

	res, err = m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (m *ArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (res domain.Article, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, meta_description, keywords, tags, reading_time_minutes, views, likes, comments, published, published_at, author_id, updated_at, created_at
  						FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, fmt.Errorf("article with ID '%s' not found", id)
	}

	return
}

func (m *ArticleRepository) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, meta_description, keywords, tags, reading_time_minutes, views, likes, comments, published, published_at, author_id, updated_at, created_at
  						FROM article WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, fmt.Errorf("article with title '%s' not found", title)
	}
	return
}

func (m *ArticleRepository) GetBySlug(ctx context.Context, slug string) (res domain.Article, err error) {
	query := `SELECT id, title, slug, content, thumbnail, image, short_description, meta_description, keywords, tags, reading_time_minutes, views, likes, comments, published, published_at, author_id, updated_at, created_at
  						FROM article WHERE slug = ?`

	list, err := m.fetch(ctx, query, slug)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, fmt.Errorf("article with slug '%s' not found", slug)
	}
	return
}

func (m *ArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	// Start transaction
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT article SET id=?, title=?, slug=?, content=?, thumbnail=?, image=?, short_description=?, meta_description=?, keywords=?, tags=?, reading_time_minutes=?, views=?, likes=?, comments=?, published=?, published_at=?, author_id=?, updated_at=?, created_at=?`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	// Generate UUID if not set
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	_, err = stmt.ExecContext(ctx, a.ID, a.Title, a.Slug, a.Content, a.Thumbnail, a.Image, a.ShortDescription, a.MetaDescription, a.Keywords, a.Tags, a.ReadingTimeMinutes, a.Views, a.Likes, a.Comments, a.Published, a.PublishedAt, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return
	}

	// Link categories if provided, or assign default category
	categories := a.Categories
	if len(categories) == 0 {
		// Assign default "Uncategorized" category
		defaultCategoryID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		categories = []domain.Category{{ID: defaultCategoryID}}
	}

	categoryQuery := `INSERT INTO article_category (id, article_id, category_id, created_at) VALUES (?, ?, ?, ?)`
	categoryStmt, err := tx.PrepareContext(ctx, categoryQuery)
	if err != nil {
		return err
	}
	defer categoryStmt.Close()

	for _, category := range categories {
		categoryLinkID := uuid.New()
		_, err = categoryStmt.ExecContext(ctx, categoryLinkID, a.ID, category.ID, a.CreatedAt)
		if err != nil {
			return err
		}
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
	// Start transaction
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `UPDATE article set title=?, slug=?, content=?, thumbnail=?, image=?, short_description=?, meta_description=?, keywords=?, tags=?, reading_time_minutes=?, views=?, likes=?, comments=?, published=?, published_at=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Slug, ar.Content, ar.Thumbnail, ar.Image, ar.ShortDescription, ar.MetaDescription, ar.Keywords, ar.Tags, ar.ReadingTimeMinutes, ar.Views, ar.Likes, ar.Comments, ar.Published, ar.PublishedAt, ar.Author.ID, ar.UpdatedAt, ar.ID)
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

	// Update categories if provided
	if len(ar.Categories) > 0 {
		// First, lock the rows to prevent deadlock - select existing links
		lockQuery := `SELECT id FROM article_category WHERE article_id = ? ORDER BY id FOR UPDATE`
		_, err = tx.ExecContext(ctx, lockQuery, ar.ID)
		if err != nil {
			return err
		}

		// Delete existing category links
		deleteQuery := `DELETE FROM article_category WHERE article_id = ?`
		_, err = tx.ExecContext(ctx, deleteQuery, ar.ID)
		if err != nil {
			return err
		}

		// Insert new category links in sorted order to prevent deadlock
		categoryQuery := `INSERT INTO article_category (id, article_id, category_id, created_at) VALUES (?, ?, ?, ?)`
		categoryStmt, err := tx.PrepareContext(ctx, categoryQuery)
		if err != nil {
			return err
		}
		defer categoryStmt.Close()

		for _, category := range ar.Categories {
			categoryLinkID := uuid.New()
			_, err = categoryStmt.ExecContext(ctx, categoryLinkID, ar.ID, category.ID, ar.UpdatedAt)
			if err != nil {
				return err
			}
		}
	}

	return
}

func (m *ArticleRepository) SlugExistsExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM article WHERE slug = ? AND id != ?`
	var count int
	err := m.Conn.QueryRowContext(ctx, query, slug, excludeID).Scan(&count)
	return count > 0, err
}
