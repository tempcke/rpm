package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

// Postgres repository should NOT be used in production
type Postgres struct {
	db *sql.DB
}

// NewPostgresRepo constructs a Postgres repository
func NewPostgresRepo(db *sql.DB) Postgres {
	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to db: " + err.Error())
	}

	return Postgres{
		db: db,
	}
}

// NewProperty constructs a property
func (r Postgres) NewProperty(street, city, state, zip string) entity.Property {
	return entity.NewProperty(street, city, state, zip)
}

// StoreProperty persists a property
func (r Postgres) StoreProperty(ctx context.Context, property entity.Property) error {
	query := `
		INSERT INTO properties (
			id, street, city, state, zip, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)

		ON CONFLICT (id) DO UPDATE SET
			street=$2, city=$3, state=$4, zip=$5
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx,
		property.ID,
		property.Street,
		property.City,
		property.StateCode,
		property.Zip,
		property.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

// RetrieveProperty by id
func (r Postgres) RetrieveProperty(ctx context.Context, id string) (entity.Property, error) {
	p := entity.Property{}

	query := `
		SELECT id, street, city, state, zip, created_at
		FROM properties WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Street, &p.City,
		&p.StateCode, &p.Zip, &p.CreatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return p, internal.ErrEntityNotFound
		}
		return p, err
	}
	p.CreatedAt = p.CreatedAt.Local()

	return p, nil
}

// PropertyList is used to list properties
func (r Postgres) PropertyList(ctx context.Context) ([]entity.Property, error) {
	propList := make([]entity.Property, 0)

	query := `
		SELECT id, street, city, state, zip, created_at
		FROM properties
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return propList, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		p := entity.Property{}

		err = rows.Scan(
			&p.ID, &p.Street, &p.City,
			&p.StateCode, &p.Zip, &p.CreatedAt,
		)
		if err != nil {
			return propList, err
		}

		p.CreatedAt = p.CreatedAt.Local()
		propList = append(propList, p)
	}

	return propList, nil
}

// DeleteProperty by id
func (r Postgres) DeleteProperty(ctx context.Context, id string) error {
	query := "DELETE FROM properties WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, id)
	return err
}
