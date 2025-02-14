// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sql

import (
	"context"
	"database/sql"
	"fmt"

	"go.temporal.io/api/serviceerror"

	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/persistence"
	"go.temporal.io/server/common/persistence/sql/sqlplugin"
)

type sqlShardStore struct {
	SqlStore
	currentClusterName string
}

// newShardPersistence creates an instance of ShardManager
func newShardPersistence(
	db sqlplugin.DB,
	currentClusterName string,
	logger log.Logger,
) (persistence.ShardStore, error) {
	return &sqlShardStore{
		SqlStore:           NewSqlStore(db, logger),
		currentClusterName: currentClusterName,
	}, nil
}

func (m *sqlShardStore) GetClusterName() string {
	return m.currentClusterName
}

func (m *sqlShardStore) CreateShard(
	request *persistence.InternalCreateShardRequest,
) error {
	ctx, cancel := newExecutionContext()
	defer cancel()
	if _, err := m.GetShard(&persistence.InternalGetShardRequest{
		ShardID: request.ShardID,
	}); err == nil {
		return &persistence.ShardAlreadyExistError{
			Msg: fmt.Sprintf("CreateShard operaiton failed. Shard with ID %v already exists.", request.ShardID),
		}
	}

	row := &sqlplugin.ShardsRow{
		ShardID:      request.ShardID,
		RangeID:      request.RangeID,
		Data:         request.ShardInfo.Data,
		DataEncoding: request.ShardInfo.EncodingType.String(),
	}

	if _, err := m.Db.InsertIntoShards(ctx, row); err != nil {
		return serviceerror.NewInternal(fmt.Sprintf("CreateShard operation failed. Failed to insert into shards table. Error: %v", err))
	}

	return nil
}

func (m *sqlShardStore) GetShard(
	request *persistence.InternalGetShardRequest,
) (*persistence.InternalGetShardResponse, error) {
	ctx, cancel := newExecutionContext()
	defer cancel()
	row, err := m.Db.SelectFromShards(ctx, sqlplugin.ShardsFilter{
		ShardID: request.ShardID,
	})
	switch err {
	case nil:
		return &persistence.InternalGetShardResponse{
			ShardInfo: persistence.NewDataBlob(row.Data, row.DataEncoding),
		}, nil
	case sql.ErrNoRows:
		return nil, serviceerror.NewNotFound(fmt.Sprintf("GetShard operation failed. Shard with ID %v not found. Error: %v", request.ShardID, err))
	default:
		return nil, serviceerror.NewInternal(fmt.Sprintf("GetShard operation failed. Failed to get record. ShardId: %v. Error: %v", request.ShardID, err))
	}
}

func (m *sqlShardStore) UpdateShard(
	request *persistence.InternalUpdateShardRequest,
) error {
	ctx, cancel := newExecutionContext()
	defer cancel()
	return m.txExecute(ctx, "UpdateShard", func(tx sqlplugin.Tx) error {
		if err := lockShard(ctx,
			tx,
			request.ShardID,
			request.PreviousRangeID,
		); err != nil {
			return err
		}
		result, err := tx.UpdateShards(ctx, &sqlplugin.ShardsRow{
			ShardID:      request.ShardID,
			RangeID:      request.RangeID,
			Data:         request.ShardInfo.Data,
			DataEncoding: request.ShardInfo.EncodingType.String(),
		})
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("rowsAffected returned error for shardID %v: %v", request.ShardID, err)
		}
		if rowsAffected != 1 {
			return fmt.Errorf("rowsAffected returned %v shards instead of one", rowsAffected)
		}
		return nil
	})
}

// initiated by the owning shard
func lockShard(
	ctx context.Context,
	tx sqlplugin.Tx,
	shardID int32,
	oldRangeID int64,
) error {

	rangeID, err := tx.WriteLockShards(ctx, sqlplugin.ShardsFilter{
		ShardID: shardID,
	})
	switch err {
	case nil:
		if rangeID != oldRangeID {
			return &persistence.ShardOwnershipLostError{
				ShardID: shardID,
				Msg:     fmt.Sprintf("Failed to update shard. Previous range ID: %v; new range ID: %v", oldRangeID, rangeID),
			}
		}
		return nil
	case sql.ErrNoRows:
		return serviceerror.NewInternal(fmt.Sprintf("Failed to lock shard with ID %v that does not exist.", shardID))
	default:
		return serviceerror.NewInternal(fmt.Sprintf("Failed to lock shard with ID: %v. Error: %v", shardID, err))
	}
}

// initiated by the owning shard
func readLockShard(
	ctx context.Context,
	tx sqlplugin.Tx,
	shardID int32,
	oldRangeID int64,
) error {
	rangeID, err := tx.ReadLockShards(ctx, sqlplugin.ShardsFilter{
		ShardID: shardID,
	})
	switch err {
	case nil:
		if rangeID != oldRangeID {
			return &persistence.ShardOwnershipLostError{
				ShardID: shardID,
				Msg:     fmt.Sprintf("Failed to lock shard. Previous range ID: %v; new range ID: %v", oldRangeID, rangeID),
			}
		}
		return nil
	case sql.ErrNoRows:
		return serviceerror.NewInternal(fmt.Sprintf("Failed to lock shard with ID %v that does not exist.", shardID))
	default:
		return serviceerror.NewInternal(fmt.Sprintf("Failed to lock shard with ID: %v. Error: %v", shardID, err))
	}
}
