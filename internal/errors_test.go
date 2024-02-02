package internal_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal"
)

func TestErrors(t *testing.T) {
	t.Run("Is", func(t *testing.T) {
		var (
			err1 = errors.New(uuid.NewString())
			err2 = errors.New(uuid.NewString())
			err3 = errors.New(uuid.NewString())
			err4 = errors.New(uuid.NewString())
		)
		errs1 := internal.NewErrors(err1, err2)
		errs2 := errs1.Append(err3, err4)

		assert.True(t, errors.Is(errs1, err1))
		assert.True(t, errors.Is(errs1, err2))
		assert.False(t, errors.Is(errs1, err3))
		assert.False(t, errors.Is(errs1, err4))

		assert.True(t, errors.Is(errs2, err1))
		assert.True(t, errors.Is(errs2, err2))
		assert.True(t, errors.Is(errs2, err3))
		assert.True(t, errors.Is(errs2, err4))
	})
	t.Run("format", func(t *testing.T) {
		var (
			err1 = errors.New("first")
			err2 = errors.New("second")
			err3 = errors.New("third")
			err4 = errors.New("fourth")
		)
		errs1 := internal.NewErrors(err1, err2).Append(err3).Append(err4)
		assert.Equal(t, "fourth: third: second: first", errs1.Error())

		errs2 := internal.NewErrors(err1, internal.NewErrors(err2, err3), err4)
		assert.Equal(t, "fourth: third: second: first", errs2.Error())
	})
	t.Run("error or nil", func(t *testing.T) {
		require.Nil(t, internal.Errors{}.ErrorOrNil())
		require.NotNil(t, internal.Errors{internal.ErrInternal}.ErrorOrNil())
	})
}
func TestKnownErrors(t *testing.T) {
	assert.True(t, internal.IsKnownErr(internal.ErrInternal))
	assert.False(t, internal.IsKnownErr(errors.New(uuid.NewString())))
	wrappedErr := internal.MakeErr(internal.ErrBadRequest, uuid.NewString())
	require.True(t, internal.IsKnownErr(wrappedErr))
}
