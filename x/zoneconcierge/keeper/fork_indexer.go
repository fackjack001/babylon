package keeper

import (
	"bytes"

	sdkerrors "cosmossdk.io/errors"
	"github.com/babylonchain/babylon/x/zoneconcierge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetForks returns a list of forked headers at a given height
func (k Keeper) GetForks(ctx sdk.Context, chainID string, height uint64) *types.Forks {
	store := k.forkStore(ctx, chainID)
	heightBytes := sdk.Uint64ToBigEndian(height)
	// if no fork at the moment, create an empty struct
	if !store.Has(heightBytes) {
		return &types.Forks{
			Headers: []*types.IndexedHeader{},
		}
	}
	forksBytes := store.Get(heightBytes)
	var forks types.Forks
	k.cdc.MustUnmarshal(forksBytes, &forks)
	return &forks
}

// insertForkHeader inserts a forked header to the list of forked headers at the same height
func (k Keeper) insertForkHeader(ctx sdk.Context, chainID string, header *types.IndexedHeader) error {
	if header == nil {
		return sdkerrors.Wrapf(types.ErrInvalidHeader, "header is nil")
	}
	store := k.forkStore(ctx, chainID)
	forks := k.GetForks(ctx, chainID, header.Height) // if no fork at the height, forks will be an empty struct rather than nil
	// if the header is already in forks, discard this header and return directly
	for _, h := range forks.Headers {
		if bytes.Equal(h.Hash, header.Hash) {
			return nil
		}
	}
	forks.Headers = append(forks.Headers, header)
	forksBytes := k.cdc.MustMarshal(forks)
	store.Set(sdk.Uint64ToBigEndian(header.Height), forksBytes)
	return nil
}

// forkStore stores the forks for each CZ
// prefix: ForkKey || chainID
// key: height that this fork starts from
// value: a list of IndexedHeader, representing each header in the fork
func (k Keeper) forkStore(ctx sdk.Context, chainID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	forkStore := prefix.NewStore(store, types.ForkKey)
	chainIDBytes := []byte(chainID)
	return prefix.NewStore(forkStore, chainIDBytes)
}
