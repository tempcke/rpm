package migrate

import (
	"database/sql"
	"log/slog"

	"github.com/tempcke/rpm/internal/db/postgres/migrate/flows"
	"github.com/tempcke/rpm/internal/lib/mig"
)

var allFlows = []*mig.Flow{
	&flows.Flow001Properties,
	&flows.Flow002Tenants,
}

func Up(db *sql.DB, log *slog.Logger) error {
	log = log.With("operation", "migrate.Up")

	if len(allFlows) == 0 {
		log.Warn("no migrations, was Up called more than once?")
		return nil
	}

	m := mig.NewRunner(db).
		WithLogger(log).
		WithFlows(allFlows...)

	if err := m.Up(); err != nil {
		return err
	}
	allFlows = nil
	return nil
}
