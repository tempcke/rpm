package mig

import (
	"database/sql"
	"strings"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

const DefaultTrackingTable = "mig_applied_migrations"
const DefaultDialect = "postgres"

type Runner struct {
	db      *sql.DB
	dialect string
	logger  logrus.FieldLogger
	flows   []*Flow
	migSet  migrate.MigrationSet
}

func NewRunner(db *sql.DB) Runner {
	r := Runner{db: db}
	return r.
		WithLogger(logrus.StandardLogger()).
		WithTrackingTable(DefaultTrackingTable).
		WithDialect(DefaultDialect)
}
func (r Runner) WithFlows(flows ...*Flow) Runner {
	if r.flows == nil {
		r.flows = make([]*Flow, 0, len(flows))
	}
	r.flows = append(r.flows, flows...)
	return r
}
func (r Runner) WithLogger(logger logrus.FieldLogger) Runner {
	r.logger = logger
	return r
}

// WithTrackingTable can be used to change the table name that migrations are recorded in
// if you qualify the table name with the schema it will also set the schema
// so tblName = "migrations" will set the table name leaving schemaName unchanged
// and tblName = "foo.migrations" will set both schema=foo, tbl=migrations
func (r Runner) WithTrackingTable(tblName string) Runner {
	var schema string
	if strings.Contains(tblName, ".") {
		toks := strings.Split(tblName, ".")
		schema, tblName = toks[0], toks[1]
		r.migSet.SchemaName = strings.Trim(schema, `"`)
	}
	r.migSet.TableName = strings.Trim(tblName, `"`)
	return r
}

// WithSchema sets the schemaName that the migration table is created in
// if that schema does not already exist it will try to create it for you
func (r Runner) WithSchema(schemaName string) Runner {
	r.migSet.SchemaName = schemaName
	return r
}

// WithDialect allows you to use something other than the default "postgres"
func (r Runner) WithDialect(dialect string) Runner {
	r.dialect = dialect
	return r
}

func (r *Runner) Up() error {
	migrationSource := r.migrations()

	if len(migrationSource.Migrations) == 0 {
		r.log().Warn("no migrations, was Up called more than once?")
		return nil
	}

	n, err := r.migSet.Exec(r.db, r.dialect, migrationSource, migrate.Up)
	if err != nil {
		return err
	}
	r.log().Infof("Applied %d migrations in %s schema!", n, r.migSet.SchemaName)
	return nil
}

func (r *Runner) migrations() *migrate.MemoryMigrationSource {
	var migrations = make([]*migrate.Migration, 0)
	for _, flow := range r.migFlows() {
		for _, step := range flow {
			migration := &migrate.Migration{
				Id:   step.ID,
				Up:   []string{step.Up},
				Down: []string{step.Down},
			}
			migrations = append(migrations, migration)
		}
	}
	return &migrate.MemoryMigrationSource{
		Migrations: migrations,
	}
}

// migFlows collects and returns the flows
// then it frees the memory by emptying the global vars
// this should prove how it works: https://go.dev/play/p/C07sQAZfpB8
// if you call this a 2nd time it will return an empty set
func (r *Runner) migFlows() []Flow {
	result := make([]Flow, len(r.flows))
	for i, flow := range r.flows {
		if flow != nil {
			result[i] = *flow
			*r.flows[i] = nil // free the ram
		}
	}
	r.flows = nil // free the ram
	return result
}

func (r Runner) log() logrus.FieldLogger {
	if r.logger == nil {
		r.logger = logrus.StandardLogger()
	}
	return r.logger.WithField("mig", "mig.Runner")
}
