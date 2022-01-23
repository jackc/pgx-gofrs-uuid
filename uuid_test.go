package uuid_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5/pgtype/testutil"
	"github.com/stretchr/testify/require"
)

func TestCodecDecodeValue(t *testing.T) {
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

func TestArray(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustCloseContext(t, conn)

	pgxuuid.Register(conn.ConnInfo())

	inputSlice := []uuid.UUID{}

	for i := 0; i < 10; i++ {
		u, err := uuid.NewV4()
		require.NoError(t, err)
		inputSlice = append(inputSlice, u)
	}

	var outputSlice []uuid.UUID
	err := conn.QueryRow(context.Background(), `select $1::uuid[]`, inputSlice).Scan(&outputSlice)
	require.NoError(t, err)
	require.Equal(t, inputSlice, outputSlice)
}
