package mig

import (
	"context"
	"fmt"
	"log/slog"
)

// limits
const (
	FlowNumLimit = 4095 // FFF
	StepNumLimit = 255  // FF
)

// MakeID creates a migration ID from a given flow and step num
// flowNum and stepNum are converted to hex and appended to the prefix
// the intention is for each file to be numbered by its flowNum 001, 002 etc.
// but allow it to add related migrations in the future with stepNums
// this way create table and future alter table statements can be together in the same file
//
// if you want your ID's to be valid uuid's then simply remove the last 6 chars
// from a baseID to create your prefix, such as "46855208-d306-4da2-bd18-30f7bc"
// as the flowNum and stepNum will always result in 6 chars, 3 for each number
func MakeID(prefix string, flowNum, stepNum int) string {
	checkLimits(flowNum, stepNum)
	return fmt.Sprintf("%s%03x%03x", prefix, flowNum, stepNum)
}

// checkLimits makes sure that the nums are within range
// if they exceed the limits then the id may be to long and cause mig.MustID to panic
func checkLimits(flowNum, stepNum int) {
	if flowNum > FlowNumLimit || stepNum > StepNumLimit {
		ctx := context.Background()
		msg := "flow or step num exceeds limit"
		slog.Default().LogAttrs(ctx, slog.LevelError, msg,
			slog.String("func", "mig.MakeID"),
			slog.Int("flowNum", flowNum),
			slog.Int("stepNum", stepNum),
			slog.Int("FlowNumLimit", FlowNumLimit),
			slog.Int("StepNumLimit", StepNumLimit),
		)
		// this will never panic, because tests ensure it is set correctly
		panic(msg)
	}
}
