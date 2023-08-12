package flows

import "github.com/tempcke/rpm/internal/lib/mig"

var Flow002Tenants = mig.Flow{
	{
		ID: mig.MakeID(idPrefix, 2, 1),
		Up: `
			CREATE TABLE IF NOT EXISTS tenants (
				id          VARCHAR(36)  PRIMARY KEY,
				full_name   VARCHAR(128) NOT NULL,
				dl_num      VARCHAR(32)  NOT NULL,
				dl_state    VARCHAR(32)  NOT NULL,
				dob         date,

				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
			);`,
	},
}
