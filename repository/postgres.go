package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/tempcke/rpm/entity"
)

// Postgres repository should NOT be used in production
type Postgres struct {
	db *sql.DB
}

// NewPostgresRepo constructs an Postgres repository
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
func (r Postgres) StoreProperty(property entity.Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO properties
		(id, street, city, state, zip, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		property.ID,
		property.Street,
		property.City,
		property.StateCode,
		property.Zip,
		property.CreatedAt,
	)

	return err
}

// RetrieveProperty by id
func (r Postgres) RetrieveProperty(id string) (entity.Property, error) {

	p := entity.Property{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, street, city, state, zip, created_at
		FROM properties WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Street, &p.City,
		&p.StateCode, &p.Zip, &p.CreatedAt,
	)

	p.CreatedAt = p.CreatedAt.Local()

	return p, err
}

// PropertyList is used to list properties
func (r Postgres) PropertyList() ([]entity.Property, error) {
	propList := make([]entity.Property, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, street, city, state, zip, created_at
		FROM properties
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return propList, err
	}
	defer rows.Close()

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
func (r Postgres) DeleteProperty(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM properties WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}
