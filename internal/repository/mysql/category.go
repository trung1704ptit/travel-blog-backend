package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
)

type CategoryRepository struct {
	Conn *sql.DB
}

// NewCategoryRepository will create an object that represent the category.Repository interface
func NewCategoryRepository(conn *sql.DB) *CategoryRepository {
	return &CategoryRepository{conn}
}

// GetByArticleID fetches all categories for a specific article
func (m *CategoryRepository) GetByArticleID(ctx context.Context, articleID uuid.UUID) ([]domain.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, c.description, c.image, c.parent_id, c.created_at, c.updated_at
		FROM category c
		INNER JOIN article_category ac ON c.id = ac.category_id
		WHERE ac.article_id = ?
		ORDER BY c.name
	`

	rows, err := m.Conn.QueryContext(ctx, query, articleID)
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

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		var parentID sql.NullString
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// Handle parent_id
		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				category.ParentID = &parentUUID
			}
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetByIDs fetches categories by their IDs
func (m *CategoryRepository) GetByIDs(ctx context.Context, categoryIDs []uuid.UUID) ([]domain.Category, error) {
	if len(categoryIDs) == 0 {
		return []domain.Category{}, nil
	}

	// Build the query with placeholders for IN clause
	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
		SELECT id, name, slug, description, image, parent_id, created_at, updated_at
		FROM category
		WHERE id IN (` + joinStrings(placeholders, ",") + `)
		ORDER BY name
	`

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

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		var parentID sql.NullString
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// Handle parent_id
		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				category.ParentID = &parentUUID
			}
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// Fetch retrieves categories with pagination
func (m *CategoryRepository) Fetch(ctx context.Context, cursor string, num int64) ([]domain.Category, string, error) {
	query := `SELECT id, name, slug, description, image, parent_id, created_at, updated_at
			  FROM category 
			  WHERE created_at > ? 
			  ORDER BY created_at DESC
			  LIMIT ?`

	// For simplicity, using created_at as cursor
	// In production, you might want to use a proper cursor implementation
	var decodedCursor time.Time
	if cursor != "" {
		// Decode cursor (simplified implementation)
		// You might want to use the same cursor logic as articles
		decodedCursor = time.Now().Add(-24 * time.Hour) // Default to 24 hours ago
	} else {
		decodedCursor = time.Time{} // Beginning of time
	}

	rows, err := m.Conn.QueryContext(ctx, query, decodedCursor, num)
	if err != nil {
		logrus.Error(err)
		return nil, "", err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		var parentID sql.NullString
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, "", err
		}

		// Handle parent_id
		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				category.ParentID = &parentUUID
			}
		}

		categories = append(categories, category)
	}

	// Generate next cursor
	var nextCursor string
	if len(categories) == int(num) {
		// Use the last category's created_at as next cursor
		nextCursor = categories[len(categories)-1].CreatedAt.Format(time.RFC3339)
	}

	return categories, nextCursor, nil
}

// GetBySlug retrieves a category by its slug
func (m *CategoryRepository) GetBySlug(ctx context.Context, slug string) (domain.Category, error) {
	query := `SELECT id, name, slug, description, image, parent_id, created_at, updated_at
			  FROM category 
			  WHERE slug = ?`

	row := m.Conn.QueryRowContext(ctx, query, slug)

	category := domain.Category{}
	var parentID sql.NullString
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&parentID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Category{}, domain.ErrNotFound
		}
		logrus.Error(err)
		return domain.Category{}, err
	}

	// Handle parent_id
	if parentID.Valid {
		parentUUID, err := uuid.Parse(parentID.String)
		if err == nil {
			category.ParentID = &parentUUID
		}
	}

	return category, nil
}

// GetByID retrieves a category by its ID
func (m *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	query := `SELECT id, name, slug, description, image, parent_id, created_at, updated_at
			  FROM category 
			  WHERE id = ?`

	row := m.Conn.QueryRowContext(ctx, query, id)

	category := domain.Category{}
	var parentID sql.NullString
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.Image,
		&parentID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Category{}, domain.ErrNotFound
		}
		logrus.Error(err)
		return domain.Category{}, err
	}

	// Handle parent_id
	if parentID.Valid {
		parentUUID, err := uuid.Parse(parentID.String)
		if err == nil {
			category.ParentID = &parentUUID
		}
	}

	return category, nil
}

// Store creates a new category
func (m *CategoryRepository) Store(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO category (id, name, slug, description, image, parent_id, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// Generate UUID if not set
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	_, err := m.Conn.ExecContext(ctx, query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.Image,
		category.ParentID,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

// Update modifies an existing category
func (m *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `UPDATE category 
			  SET name = ?, slug = ?, description = ?, image = ?, parent_id = ?, updated_at = ?
			  WHERE id = ?`

	category.UpdatedAt = time.Now()

	result, err := m.Conn.ExecContext(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Image,
		category.ParentID,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete removes a category
func (m *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM category WHERE id = ?`

	result, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// SlugExistsExcludingID checks if a slug exists for a different category
func (m *CategoryRepository) SlugExistsExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM category WHERE slug = ? AND id != ?`
	var count int
	err := m.Conn.QueryRowContext(ctx, query, slug, excludeID).Scan(&count)
	return count > 0, err
}

// GetChildren retrieves all children of a category
func (m *CategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error) {
	query := `SELECT id, name, slug, description, image, parent_id, created_at, updated_at
			  FROM category 
			  WHERE parent_id = ?
			  ORDER BY name`

	rows, err := m.Conn.QueryContext(ctx, query, parentID)
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

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		var parentID sql.NullString
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// Handle parent_id
		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				category.ParentID = &parentUUID
			}
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetRootCategories retrieves all root categories (no parent)
func (m *CategoryRepository) GetRootCategories(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT id, name, slug, description, image, parent_id, created_at, updated_at
			  FROM category 
			  WHERE parent_id IS NULL
			  ORDER BY name`

	rows, err := m.Conn.QueryContext(ctx, query)
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

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		var parentID sql.NullString
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// Handle parent_id
		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				category.ParentID = &parentUUID
			}
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryTree retrieves the complete category tree
func (m *CategoryRepository) GetCategoryTree(ctx context.Context) ([]domain.Category, error) {
	// First get all categories
	query := `SELECT id, name, slug, description, image, parent_id, created_at, updated_at
			  FROM category 
			  ORDER BY name`

	rows, err := m.Conn.QueryContext(ctx, query)
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

	var allCategories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		var parentID sql.NullString
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// Handle parent_id
		if parentID.Valid {
			parentUUID, err := uuid.Parse(parentID.String)
			if err == nil {
				category.ParentID = &parentUUID
			}
		}

		allCategories = append(allCategories, category)
	}

	// Build the tree structure
	return m.buildCategoryTree(allCategories), nil
}

// buildCategoryTree builds a hierarchical tree from flat category list
func (m *CategoryRepository) buildCategoryTree(categories []domain.Category) []domain.Category {
	// Create a map for quick lookup
	categoryMap := make(map[uuid.UUID]*domain.Category)
	var rootCategories []*domain.Category

	// First pass: create map and identify root categories
	for i := range categories {
		category := &categories[i]
		categoryMap[category.ID] = category

		if category.ParentID == nil {
			rootCategories = append(rootCategories, category)
		}
	}

	// Second pass: build parent-child relationships
	for i := range categories {
		category := &categories[i]
		if category.ParentID != nil {
			if parent, exists := categoryMap[*category.ParentID]; exists {
				parent.Children = append(parent.Children, *category)
			}
		}
	}

	// Convert pointers back to values for return
	var result []domain.Category
	for _, root := range rootCategories {
		result = append(result, *root)
	}

	return result
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
