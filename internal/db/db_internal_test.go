package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClose_NilGuards(t *testing.T) {
	var d *Db
	require.Error(t, d.Close())

	d2 := &Db{DB: nil}
	require.Error(t, d2.Close())
}

func TestPing_NilGuards(t *testing.T) {
	var d *Db
	require.Error(t, d.Ping(context.Background()))

	d2 := &Db{DB: nil}
	require.Error(t, d2.Ping(context.Background()))
}
