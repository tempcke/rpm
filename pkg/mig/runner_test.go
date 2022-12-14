//go:build withDocker
// +build withDocker

package mig_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal/test"
	"github.com/tempcke/rpm/pkg/mig"
)

var (
	schema   = test.RandString(6) // will work when set to "public" also
	idPrefix = schema
)

// in addition to being a test, this is really an example of how to use it
func TestRunner(t *testing.T) {
	var db = test.DB(t)

	// construct the runner
	runner := mig.NewRunner(db).
		WithSchema(schema).
		WithFlows( // we use pointers so that the vars can be cleared when done
			&Flow001Companies,
			&Flow002Addresses)

	// run the migrations
	require.NoError(t, runner.Up())

	// assert the migrations ran...
	assertTableExistsInSchema(t, db, schema, mig.DefaultTrackingTable)
	assertTableExistsInSchema(t, db, schema, "companies")
	assertTableExistsInSchema(t, db, schema, "addresses")

	// the goal is to not leave a bunch of text variables in RAM the entire time the application is running
	// so here we assert that they have been emptied
	assert.Empty(t, Flow001Companies)
	assert.Empty(t, Flow002Addresses)

	// cleanup
	if schema != "" && schema != "public" {
		_, err := db.Exec(`DROP SCHEMA IF EXISTS ` + schema + ` CASCADE`)
		require.NoError(t, err)
	}
}

func assertTableExistsInSchema(t testing.TB, db *sql.DB, schema, table string) {
	t.Helper()

	if schema == "" {
		schema = "public"
	}

	query := fmt.Sprintf(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = '%s'
		AND    table_name   = '%s'
	);`, schema, table)

	var exists bool
	require.NoError(t, db.QueryRow(query).Scan(&exists))
	require.True(t, exists)
}

// test flows
// for your project you would likely have a flows sub package
// and then one file per flow

// flows/001_companies.go
var Flow001Companies = mig.Flow{
	{
		ID: mig.MakeID(idPrefix, 1, 1),
		Up: `
			CREATE TABLE IF NOT EXISTS ` + schema + `.companies (
				org_id         text                     NOT NULL,
				name           text                     NOT NULL,
				created_at     timestamp with time zone NOT NULL DEFAULT now(),
				updated_at     timestamp with time zone NOT NULL DEFAULT now(),
		
				PRIMARY KEY (org_id) NOT DEFERRABLE
			);
		`,
		Down: `DROP TABLE IF EXISTS ` + schema + `.companies;`,
	},
	{
		ID: mig.MakeID(idPrefix, 1, 2),
		Up: `
			ALTER TABLE IF EXISTS ` + schema + `.companies
				ADD version  integer NOT NULL DEFAULT 1,
				ADD status   text    NOT NULL DEFAULT 'ACTIVE';`,
		Down: `
			ALTER TABLE IF EXISTS ` + schema + `.companies
				DROP version,
				DROP status;`,
	},
}

// flows/002_addresses.go
var Flow002Addresses = mig.Flow{
	{
		ID: mig.MakeID(idPrefix, 2, 1),
		Up: `
			CREATE TABLE IF NOT EXISTS ` + schema + `.addresses (
				org_id           text                  NOT NULL,
				id               text                  NOT NULL,
				name             text                  NOT NULL,
				street_address   text                  NOT NULL,
				city             text                  NOT NULL,
				country_code     text                  NOT NULL,
				region           text                  NOT NULL,
				postal_code      text                  NOT NULL,
			
				PRIMARY KEY (id) NOT DEFERRABLE,
				CONSTRAINT fk_org_id FOREIGN KEY(org_id) REFERENCES ` + schema + `.companies(org_id) ON DELETE CASCADE
			);
			
			CREATE INDEX IF NOT EXISTS addresses__org ON ` + schema + `.addresses (org_id);
		`,
		Down: `
			DROP TABLE IF EXISTS ` + schema + `.addresses;
		`,
	},
}
