package parallel

import (
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	var c atomic.Int64
	err := Parallel(Fn(func() error {
		time.Sleep(time.Millisecond)
		require.True(t, c.CAS(1, 2))
		return nil
	}), Fn(func() error {
		require.True(t, c.CAS(0, 1))
		return nil
	})).Do()
	require.NoError(t, err)
	require.Equal(t, int64(2), c.Load())
}

func TestOrdered(t *testing.T) {
	var c atomic.Int64
	err := Ordered(Fn(func() error {
		time.Sleep(time.Millisecond)
		require.True(t, c.CAS(0, 1))
		return nil
	}), Fn(func() error {
		require.True(t, c.CAS(1, 2))
		return nil
	})).Do()
	require.NoError(t, err)
	require.Equal(t, int64(2), c.Load())
}

func TestMix(t *testing.T) {
	var c atomic.Int64
	err := Parallel(
		Ordered(Fn(func() error {
			time.Sleep(time.Millisecond)
			require.True(t, c.CAS(2, 3))
			return nil
		}), Fn(func() error {
			require.True(t, c.CAS(3, 4))
			return nil
		})),
		Ordered(Fn(func() error {
			require.True(t, c.CAS(0, 1))
			return nil
		}), Fn(func() error {
			require.True(t, c.CAS(1, 2))
			return nil
		})),
	).Do()
	require.NoError(t, err)
	require.Equal(t, int64(4), c.Load())
}

func TestParallelErrs(t *testing.T) {
	err := Parallel(Fn(func() error {
		time.Sleep(time.Millisecond)
		return errors.New("a")
	}), Fn(func() error {
		return errors.New("b")
	})).Do()
	require.Equal(t, "b; a", err.Error())
}

func TestOrderedErrs(t *testing.T) {
	err := Ordered(
		Fn(func() error { return errors.New("a") }),
		Fn(func() error { return errors.New("b") }),
	).Do()
	require.Equal(t, "a", err.Error())
}

func TestMixErrs(t *testing.T) {
	err := Parallel(
		Ordered(
			Fn(func() error { time.Sleep(time.Millisecond); return errors.New("a") }),
			Fn(func() error { return errors.New("b") }),
		),
		Ordered(
			Fn(func() error { return errors.New("c") }),
			Fn(func() error { return errors.New("d") }),
		),
	).Do()
	require.Equal(t, "c; a", err.Error())
}
