package uuid

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
)

type UUID uuid.UUID

func (d *UUID) DecodeUUID(u *pgtype.UUID) error {
	if !u.Valid {
		return fmt.Errorf("cannot decode uuid NULL into %T", d)
	}

	*d = UUID(u.Bytes)
	return nil
}

type NullUUID uuid.NullUUID

func (d *NullUUID) DecodeUUID(u *pgtype.UUID) error {
	if u.Valid {
		*d = NullUUID{UUID: u.Bytes, Valid: true}
	} else {
		*d = NullUUID{}
	}
	return nil
}

func UUIDDecoderWrapper(value interface{}) pgtype.UUIDDecoder {
	switch value := value.(type) {
	case *uuid.UUID:
		return (*UUID)(value)
	case *uuid.NullUUID:
		return (*NullUUID)(value)
	default:
		return nil
	}
}

func Getter(a pgtype.UUID) interface{} {
	if !a.Valid {
		return nil
	}

	var b UUID
	err := b.DecodeUUID(&a)
	if err != nil {
		panic(err) // Can't happen
	}

	return uuid.UUID(b)
}

// Register registers the github.com/gofrs/uuid integration with a pgtype.ConnInfo.
func Register(ci *pgtype.ConnInfo) {
	ci.PreferAssignToOverSQLScannerForType(&uuid.UUID{})
	ci.PreferAssignToOverSQLScannerForType(&uuid.NullUUID{})
	ci.RegisterDataType(pgtype.DataType{
		Value: &pgtype.UUID{
			UUIDDecoderWrapper: UUIDDecoderWrapper,
			Getter:             Getter,
		},
		Name: "uuid",
		OID:  pgtype.UUIDOID,
	})
}
