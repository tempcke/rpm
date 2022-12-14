package migrate

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/internal/db/postgres/migrate/flows"
	"github.com/tempcke/rpm/pkg/mig"
)

var allFlows = []*mig.Flow{
	&flows.Flow001Properties,
}

func Up(db *sql.DB, log logrus.FieldLogger) error {
	log = log.WithField("operation", "migrate.Up")

	if len(allFlows) == 0 {
		log.Warn("no migrations, was Up called more than once?")
		return nil
	}

	m := mig.NewRunner(db).
		WithLogger(log).
		WithFlows(allFlows...)

	return m.Up()
}
