package flows

import "github.com/tempcke/rpm/internal/lib/mig"

const idPrefix = "rpm"

var Flow001Properties = mig.Flow{
	{
		ID: mig.MakeID(idPrefix, 1, 1),
		Up: `
			CREATE TABLE IF NOT EXISTS properties (
				id         VARCHAR(36)  PRIMARY KEY,
				street     VARCHAR(255) NOT NULL,
				city       VARCHAR(32)  NOT NULL,
				state      VARCHAR(32)  NOT NULL,
				zip        VARCHAR(10)  NOT NULL,

				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
			)
		`,
	},
}
