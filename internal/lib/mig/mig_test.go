package mig_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/internal/lib/mig"
)

// prefix should be unique per module
const prefix = "46855208-d306-4da2-bd18-30f7bc"

func TestMakeID(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 50, 123}
	t.Run("same index should give same id every time", func(t *testing.T) {
		for _, flowNum := range nums {
			for _, stepNum := range nums {
				id := mig.MakeID(prefix, flowNum, stepNum)
				for i := 0; i < 10; i++ {
					assert.Equal(t, id, mig.MakeID(prefix, flowNum, stepNum))
				}
			}
		}
	})

	t.Run("different flowNum must result in different id", func(t *testing.T) {
		// make sure no index numbers are reused
		var stepNum = 0
		seen := make(map[string]int)
		for i := 0; i <= 100; i++ {
			flowNum := rand.Intn(mig.FlowNumLimit)
			id := mig.MakeID(prefix, flowNum, stepNum)
			prevFlowNum, ok := seen[id]

			// assert we have not seen this id before or assert the last time we saw it, it was created with the same index
			assert.True(t, !ok || prevFlowNum == flowNum)
			seen[id] = flowNum
		}
		// this is here just to prove we arnt hitting just a couple numbers over and over again
		// it is commented out because in very rare cases it could lead to a test failure
		// assert.True(t, len(seen) > 90)
	})

	t.Run("different stepNum must result in different id", func(t *testing.T) {
		// make sure no index numbers are reused
		var flowNum = 1
		seen := make(map[string]int)
		for i := 0; i <= 100; i++ {
			stepNum := rand.Intn(mig.StepNumLimit)
			id := mig.MakeID(prefix, flowNum, stepNum)
			prevStepNum, ok := seen[id]

			// assert we have not seen this id before or assert the last time we saw it, it was created with the same index
			assert.True(t, !ok || prevStepNum == stepNum)
			seen[id] = stepNum
		}
		// this is here just to prove we arnt hitting just a couple numbers over and over again
		// it is commented out because in very rare cases it could lead to a test failure
		// assert.True(t, len(seen) > 90)
	})
}
