package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/usecase"
)

type (
	propertyRepo = usecase.PropertyRepo
	tenantRepo   = usecase.TenantRepo
)

var ctx = context.Background()

func assertTimestampMatch(t *testing.T, t1, t2 time.Time) {
	t.Helper()

	// ensure matches within 1 second
	timeDiff := t1.Sub(t2).Seconds()
	assert.LessOrEqual(t, timeDiff, 1.0)

	t1Zone, _ := t1.Zone()
	t2Zone, _ := t2.Zone()

	// ensure time converted back to local time
	assert.Equal(t, t1Zone, t2Zone)
}
func assertEntityInSet[T entity.Entity](t testing.TB, id entity.ID, set ...T) {
	t.Helper()
	for _, got := range set {
		if got.GetID() == id {
			return
		}
	}
	t.Fatalf("entityID '%s' not found in set of length %d", id, len(set))
}
