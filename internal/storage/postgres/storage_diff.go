package postgres

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// StorageDiff -
type StorageDiff struct {
	*postgres.Table[*storage.StorageDiff]
}

// NewStorageDiff -
func NewStorageDiff(db *database.PgGo) *StorageDiff {
	return &StorageDiff{
		Table: postgres.NewTable[*storage.StorageDiff](db),
	}
}

// GetOnBlock -
func (sd *StorageDiff) GetOnBlock(ctx context.Context, height, contractId uint64, key []byte) (diff storage.StorageDiff, err error) {
	query := sd.DB().ModelContext(ctx, &diff).
		Where("contract_id = ?", contractId).
		Where("key = ?", key)

	if height > 0 {
		query = query.Where("height >= ?", height)
	}

	err = query.Order("id desc").
		Limit(1).
		Select(&diff)
	return
}

// InsertByCopy -
func (sd *StorageDiff) InsertByCopy(diffs []storage.StorageDiff) (io.Reader, string, error) {
	if len(diffs) == 0 {
		return nil, "", nil
	}
	builder := new(strings.Builder)

	for i := range diffs {
		if err := writeUint64(builder, diffs[i].Height); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, diffs[i].ContractID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeBytes(builder, diffs[i].Key); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeBytes(builder, diffs[i].Value); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte('\n'); err != nil {
			return nil, "", err
		}
	}

	query := fmt.Sprintf(`COPY %s (
		height, contract_id, key, value
	) FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.StorageDiff{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (sd *StorageDiff) Filter(ctx context.Context, fltr storage.StorageDiffFilter, opts ...storage.FilterOption) ([]storage.StorageDiff, error) {
	q := sd.DB().ModelContext(ctx, (*storage.StorageDiff)(nil))
	q = integerFilter(q, "id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = addressFilter(q, "hash", fltr.Contract, "Contract")
	q = equalityFilter(q, "key", fltr.Key)
	q = optionsFilter(q, opts...)

	var result []storage.StorageDiff
	err := q.Select(&result)
	return result, err
}
