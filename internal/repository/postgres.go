package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jonboulle/clockwork"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/filters"
	"github.com/tempcke/rpm/internal/lib/log"
	"github.com/tempcke/rpm/usecase"
)

// Postgres repository should NOT be used in production
type Postgres struct {
	db    *sql.DB
	clock clockwork.Clock
}

// NewPostgresRepo constructs a Postgres repository
func NewPostgresRepo(db *sql.DB) Postgres {
	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to db: " + err.Error())
	}

	return Postgres{
		db:    db,
		clock: clockwork.NewRealClock(),
	}
}

func (r Postgres) NewProperty(street, city, state, zip string) entity.Property {
	return entity.NewProperty(street, city, state, zip)
}
func (r Postgres) StoreProperty(ctx context.Context, property entity.Property) error {
	const query = `
		INSERT INTO properties (
			id, street, city, state, zip, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)

		ON CONFLICT (id) DO UPDATE SET
			street=$2, city=$3, state=$4, zip=$5`

	qArgs := []any{
		property.ID,
		property.Street,
		property.City,
		property.StateCode,
		property.Zip,
		property.CreatedAt,
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecContext(ctx, qArgs...)

	if err != nil {
		return err
	}
	return nil
}
func (r Postgres) GetProperty(ctx context.Context, id string) (entity.Property, error) {
	const query = `
		SELECT id, street, city, state, zip, created_at
		FROM properties WHERE id = $1;`

	var (
		p        = entity.Property{}
		scanArgs = []any{
			&p.ID, &p.Street, &p.City,
			&p.StateCode, &p.Zip, &p.CreatedAt,
		}
	)
	err := r.db.QueryRowContext(ctx, query, id).Scan(scanArgs...)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return p, internal.ErrEntityNotFound
		}
		return p, err
	}
	p.CreatedAt = p.CreatedAt.Local()

	return p, nil
}
func (r Postgres) PropertyList(ctx context.Context, f usecase.PropertyFilter) ([]entity.Property, error) {
	var (
		search   = ""
		propList = make([]entity.Property, 0)
	)

	if s := f.Search; s != "" {
		search = "%" + removeChars(s, ".", ",") + "%"
	}

	const query = `
		SELECT id, street, city, state, zip, created_at
		FROM properties p
		WHERE $1 = '' OR TRANSLATE(CONCAT(p.street, ' ', p.city, ' ', p.state, ' ', p.zip), '.', '') ILIKE $1
	`

	rows, err := r.db.QueryContext(ctx, query, search)

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

func (r Postgres) StoreTenant(ctx context.Context, tenant entity.Tenant) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := r.storeTenant(ctx, tx, tenant); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (r Postgres) GetTenant(ctx context.Context, id entity.ID) (*entity.Tenant, error) {
	const query = `
		SELECT id, full_name, dl_num, dl_state, dob
		FROM tenants WHERE id=$1;`
	var (
		tenant   = entity.Tenant{}
		scanArgs = []any{&tenant.ID, &tenant.FullName, &tenant.DLNum, &tenant.DLState, &tenant.DateOfBirth}
	)
	if err := r.db.QueryRowContext(ctx, query, id).Scan(scanArgs...); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, internal.ErrEntityNotFound
		}
		return nil, err
	}
	phones, err := r.getTenantPhones(ctx, id)
	if err != nil {
		return nil, err
	}
	tenant.Phones = phones
	return &tenant, nil
}
func (r Postgres) ListTenants(ctx context.Context, filter ...filters.TenantFilter) ([]entity.Tenant, error) {
	const query = `SELECT id, full_name, dl_num, dl_state, dob FROM tenants;`
	var tenants []entity.Tenant
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var (
			tenant   entity.Tenant
			scanArgs = []any{
				&tenant.ID, &tenant.FullName,
				&tenant.DLNum, &tenant.DLState,
				&tenant.DateOfBirth,
			}
		)
		if err := rows.Scan(scanArgs...); err != nil {
			return tenants, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}
func (r Postgres) storeTenant(ctx context.Context, tx *sql.Tx, tenant entity.Tenant) error {
	const query = `
			INSERT INTO tenants (id, full_name, dl_num, dl_state, dob, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET full_name=$2, dl_num=$3, dl_state=$4, dob=$5, updated_at=$6;`
	qArgs := []any{
		tenant.ID,
		tenant.FullName,
		tenant.DLNum,
		tenant.DLState,
		tenant.DateOfBirth,
		r.clock.Now(),
	}
	if _, err := tx.ExecContext(ctx, query, qArgs...); err != nil {
		return err
	}
	return r.storeTenantPhones(ctx, tx, tenant)
}
func (r Postgres) storeTenantPhones(ctx context.Context, tx *sql.Tx, tenant entity.Tenant) error {
	const (
		delPhonesQuery = `DELETE FROM tenant_phones WHERE tenant_id=$1;`
		insPhonesQuery = `INSERT INTO tenant_phones (tenant_id, phone_num, note) VALUES($1, $2, $3);`
	)

	var delPhonesArgs = []any{tenant.ID}
	if _, err := tx.Exec(delPhonesQuery, delPhonesArgs...); err != nil {
		return err
	}

	if len(tenant.Phones) > 0 {
		stmt, err := tx.PrepareContext(ctx, insPhonesQuery)
		if err != nil {
			return err
		}
		for _, phone := range tenant.Phones {
			qArgs := []any{tenant.GetID(), phone.Number, phone.Note}
			if _, err := stmt.Exec(qArgs...); err != nil {
				return err
			}
		}
	}
	return nil
}
func (r Postgres) getTenantPhones(ctx context.Context, tenantID string) ([]entity.Phone, error) {
	const query = `SELECT phone_num, note FROM tenant_phones p WHERE p.tenant_id=$1`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	var phones []entity.Phone
	for rows.Next() {
		var (
			phone    entity.Phone
			scanArgs = []any{&phone.Number, &phone.Note}
		)
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}
	return phones, nil
}

func removeChars(s string, chars ...string) string {
	for _, char := range chars {
		s = strings.ReplaceAll(s, char, "")
	}
	return s
}
