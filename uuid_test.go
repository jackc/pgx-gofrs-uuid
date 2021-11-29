package uuid_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype/testutil"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/stretchr/testify/require"
)

func TestGetter(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustCloseContext(t, conn)

	pgxuuid.Register(conn.ConnInfo())

	original, err := uuid.NewV4()
	require.NoError(t, err)

	rows, err := conn.Query(context.Background(), `select $1::uuid`, original)
	require.NoError(t, err)

	for rows.Next() {
		values, err := rows.Values()
		require.NoError(t, err)

		require.Len(t, values, 1)
		v0, ok := values[0].(uuid.UUID)
		require.True(t, ok)
		require.Equal(t, original, v0)
	}

	require.NoError(t, rows.Err())

	rows, err = conn.Query(context.Background(), `select $1::uuid`, nil)
	require.NoError(t, err)

	for rows.Next() {
		values, err := rows.Values()
		require.NoError(t, err)

		require.Len(t, values, 1)
		require.Equal(t, nil, values[0])
	}

	require.NoError(t, rows.Err())
}
