package utils

import (
	"github.com/FRECNET/common/hexutil"
	"github.com/FRECNET/core/types"
)

// S2PoS solid-state-proof-of-stake protocol constants.
var (
	//epochchange 4200
	EpochLength = uint64(200) // Default number of blocks after which to checkpoint and reset the pending votes

	ExtraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	ExtraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	NonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new signer
	NonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a signer.

	UncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
)

const (
	InmemorySnapshots      = 128 // Number of recent vote snapshots to keep in memory
	BlockSignersCacheLimit = 9000
	M2ByteLength           = 4
)
