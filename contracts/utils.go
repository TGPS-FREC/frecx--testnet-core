// Copyright (c) 2018 FRECNET
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package contracts

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/FRECNET/accounts"
	"github.com/FRECNET/accounts/abi/bind"
	"github.com/FRECNET/common"
	"github.com/FRECNET/common/hexutil"
	"github.com/FRECNET/consensus"

	"github.com/FRECNET/consensus/S2PoS"
	"github.com/FRECNET/consensus/S2PoS/utils"
	"github.com/FRECNET/contracts/blocksigner/contract"
	randomizeContract "github.com/FRECNET/contracts/randomize/contract"
	"github.com/FRECNET/core"
	"github.com/FRECNET/core/state"
	stateDatabase "github.com/FRECNET/core/state"
	"github.com/FRECNET/core/types"
	"github.com/FRECNET/ethdb"
	"github.com/FRECNET/log"
	"github.com/FRECNET/params"
)

const (
	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

type rewardLog struct {
	Sign   uint64   `json:"sign"`
	Reward *big.Int `json:"reward"`
}

var TxSignMu sync.RWMutex

// Send tx sign for block number to smart contract blockSigner.
func CreateTransactionSign(chainConfig *params.ChainConfig, pool *core.TxPool, manager *accounts.Manager, block *types.Block, chainDb ethdb.Database, eb common.Address) error {
	TxSignMu.Lock()
	defer TxSignMu.Unlock()
	if chainConfig.S2PoS != nil {
		// Find active account.
		account := accounts.Account{}
		var wallet accounts.Wallet
		etherbaseAccount := accounts.Account{
			Address: eb,
			URL:     accounts.URL{},
		}
		if wallets := manager.Wallets(); len(wallets) > 0 {
			if w, err := manager.Find(etherbaseAccount); err == nil && w != nil {
				wallet = w
				account = etherbaseAccount
			} else {
				wallet = wallets[0]
				if accts := wallets[0].Accounts(); len(accts) > 0 {
					account = accts[0]
				}
			}
		}

		// Create and send tx to smart contract for sign validate block.
		nonce := pool.State().GetNonce(account.Address)
		for retryCount := 1; ; retryCount++ {
			tx := CreateTxSign(block.Number(), block.Hash(), nonce, common.HexToAddress(common.BlockSigners))
			txSigned, err := wallet.SignTx(account, tx, chainConfig.ChainId)
			if err != nil {
				log.Error("Fail to create tx sign", "error", err)
				return err
			}
			// Add tx signed to local tx pool.
			err = pool.AddLocal(txSigned)
			if err == nil {
				break
			}
			newNonce := pool.State().GetNonce(account.Address)
			if retryCount >= 10 {
				log.Error("Fail to add signed tx to local pool", "number", block.NumberU64(), "hash", block.Hash().Hex(), "from", account.Address, "nonce", nonce, "error", err)
				return err
			}
			log.Warn("CreateTransactionSign", "number", block.NumberU64(), "from", account.Address, "nonce", nonce, "newNonce", newNonce, "retry", retryCount, "error", err)
			nonce = newNonce
		}

		// Create secret tx.
		blockNumber := block.Number().Uint64()
		checkNumber := blockNumber % chainConfig.S2PoS.Epoch
		// Generate random private key and save into chaindb.
		randomizeKeyName := []byte("randomizeKey")
		exist, _ := chainDb.Has(randomizeKeyName)

		// Set secret for randomize.
		if !exist && checkNumber > 0 && common.EpocBlockSecret <= checkNumber && common.EpocBlockOpening > checkNumber {
			// Only process when private key empty in state db.
			// Save randomize key into state db.
			randomizeKeyValue := RandStringByte(32)
			nonce := pool.State().GetNonce(account.Address)
			tx, err := BuildTxSecretRandomize(nonce, common.HexToAddress(common.RandomizeSMC), chainConfig.S2PoS.Epoch, randomizeKeyValue)
			if err != nil {
				log.Error("Fail to get tx opening for randomize", "error", err)
				return err
			}
			txSigned, err := wallet.SignTx(account, tx, chainConfig.ChainId)
			if err != nil {
				log.Error("Fail to create tx secret", "error", err)
				return err
			}
			// Add tx signed to local tx pool.
			err = pool.AddLocal(txSigned)
			if err != nil {
				log.Error("Fail to add tx secret to local pool.", "error", err, "number", block.NumberU64(), "hash", block.Hash().Hex(), "from", account.Address, "nonce", nonce)
				return err
			}

			// Put randomize key into chainDb.
			chainDb.Put(randomizeKeyName, randomizeKeyValue)
		}

		// Set opening for randomize.
		if exist && checkNumber > 0 && common.EpocBlockOpening <= checkNumber && common.EpocBlockRandomize >= checkNumber {
			randomizeKeyValue, err := chainDb.Get(randomizeKeyName)
			if err != nil {
				log.Error("Fail to get randomize key from state db.", "error", err)
				return err
			}

			nonce := pool.State().GetNonce(account.Address)
			tx, err := BuildTxOpeningRandomize(nonce, common.HexToAddress(common.RandomizeSMC), randomizeKeyValue)
			if err != nil {
				log.Error("Fail to get tx opening for randomize", "error", err)
				return err
			}
			txSigned, err := wallet.SignTx(account, tx, chainConfig.ChainId)
			if err != nil {
				log.Error("Fail to create tx opening", "error", err)
				return err
			}
			// Add tx to pool.
			err = pool.AddLocal(txSigned)
			if err != nil {
				log.Error("Fail to add tx opening to local pool.", "error", err, "number", block.NumberU64(), "hash", block.Hash().Hex(), "from", account.Address, "nonce", nonce)
				return err
			}

			// Clear randomize key in state db.
			chainDb.Delete(randomizeKeyName)
		}
	}

	return nil
}

// Create tx sign.
func CreateTxSign(blockNumber *big.Int, blockHash common.Hash, nonce uint64, blockSigner common.Address) *types.Transaction {
	data := common.Hex2Bytes(common.HexSignMethod)
	inputData := append(data, common.LeftPadBytes(blockNumber.Bytes(), 32)...)
	inputData = append(inputData, common.LeftPadBytes(blockHash.Bytes(), 32)...)
	tx := types.NewTransaction(nonce, blockSigner, big.NewInt(0), 200000, big.NewInt(0), inputData)

	return tx
}

// Send secret key into randomize smartcontract.
func BuildTxSecretRandomize(nonce uint64, randomizeAddr common.Address, epocNumber uint64, randomizeKey []byte) (*types.Transaction, error) {
	data := common.Hex2Bytes(common.HexSetSecret)
	rand.Seed(time.Now().UnixNano())
	secretNumb := rand.Intn(int(epocNumber))

	// Append randomize suffix in -1, 0, 1.
	secrets := []int64{int64(secretNumb)}
	sizeOfArray := int64(32)

	// Build extra data for tx with first position is size of array byte and second position are length of array byte.
	arrSizeOfSecrets := common.LeftPadBytes(new(big.Int).SetInt64(sizeOfArray).Bytes(), 32)
	arrLengthOfSecrets := common.LeftPadBytes(new(big.Int).SetInt64(int64(len(secrets))).Bytes(), 32)
	inputData := append(data, arrSizeOfSecrets...)
	inputData = append(inputData, arrLengthOfSecrets...)
	for _, secret := range secrets {
		encryptSecret := Encrypt(randomizeKey, new(big.Int).SetInt64(secret).String())
		inputData = append(inputData, common.LeftPadBytes([]byte(encryptSecret), int(sizeOfArray))...)
	}
	tx := types.NewTransaction(nonce, randomizeAddr, big.NewInt(0), 200000, big.NewInt(0), inputData)

	return tx, nil
}

// Send opening to randomize SMC.
func BuildTxOpeningRandomize(nonce uint64, randomizeAddr common.Address, randomizeKey []byte) (*types.Transaction, error) {
	data := common.Hex2Bytes(common.HexSetOpening)
	inputData := append(data, randomizeKey...)
	tx := types.NewTransaction(nonce, randomizeAddr, big.NewInt(0), 200000, big.NewInt(0), inputData)

	return tx, nil
}

// Get signers signed for blockNumber from blockSigner contract.
func GetSignersFromContract(state *stateDatabase.StateDB, block *types.Block) ([]common.Address, error) {
	return stateDatabase.GetSigners(state, block), nil
}

// Get signers signed for blockNumber from blockSigner contract.
func GetSignersByExecutingEVM(addrBlockSigner common.Address, client bind.ContractBackend, blockHash common.Hash) ([]common.Address, error) {
	blockSigner, err := contract.NewBlockSigner(addrBlockSigner, client)
	if err != nil {
		log.Error("Fail get instance of blockSigner", "error", err)
		return nil, err
	}
	opts := new(bind.CallOpts)
	addrs, err := blockSigner.GetSigners(opts, blockHash)
	if err != nil {
		log.Error("Fail get block signers", "error", err)
		return nil, err
	}
	return addrs, nil
}

// Get random from randomize contract.
func GetRandomizeFromContract(client bind.ContractBackend, addrMasternode common.Address) (int64, error) {
	randomize, err := randomizeContract.NewFRERandomize(common.HexToAddress(common.RandomizeSMC), client)
	if err != nil {
		log.Error("Fail to get instance of randomize", "error", err)
	}
	opts := new(bind.CallOpts)
	secrets, err := randomize.GetSecret(opts, addrMasternode)
	if err != nil {
		log.Error("Fail get secrets from randomize", "error", err)
	}
	opening, err := randomize.GetOpening(opts, addrMasternode)
	if err != nil {
		log.Error("Fail get opening from randomize", "error", err)
	}

	return DecryptRandomizeFromSecretsAndOpening(secrets, opening)
}

// Generate m2 listing from randomize array.
func GenM2FromRandomize(randomizes []int64, lenSigners int64) ([]int64, error) {
	fmt.Println("GenM2FromRandomize", lenSigners)
	blockValidator := NewSlice(int64(0), lenSigners, 1)
	randIndexs := make([]int64, lenSigners)
	total := int64(0)
	var temp int64 = 0
	for _, j := range randomizes {
		total += j
	}
	rand.Seed(total)
	for i := len(blockValidator) - 1; i >= 0; i-- {
		blockLength := len(blockValidator) - 1
		if blockLength <= 1 {
			blockLength = 1
		}
		randomIndex := int64(rand.Intn(blockLength))
		temp = blockValidator[randomIndex]
		blockValidator[randomIndex] = blockValidator[i]
		blockValidator[i] = temp
		blockValidator = append(blockValidator[:i], blockValidator[i+1:]...)
		randIndexs[i] = temp
	}

	return randIndexs, nil
}

// Get validators from m2 array integer.
func BuildValidatorFromM2(listM2 []int64) []byte {
	var validatorBytes []byte
	for _, numberM2 := range listM2 {
		// Convert number to byte.
		m2Byte := common.LeftPadBytes([]byte(fmt.Sprintf("%d", numberM2)), utils.M2ByteLength)
		validatorBytes = append(validatorBytes, m2Byte...)
	}

	return validatorBytes
}

// Decode validator hex string.
func DecodeValidatorsHexData(validatorsStr string) ([]int64, error) {
	validatorsByte, err := hexutil.Decode(validatorsStr)
	if err != nil {
		return nil, err
	}

	return utils.ExtractValidatorsFromBytes(validatorsByte), nil
}

// Decrypt randomize from secrets and opening.
func DecryptRandomizeFromSecretsAndOpening(secrets [][32]byte, opening [32]byte) (int64, error) {
	var random int64
	if len(secrets) > 0 {
		for _, secret := range secrets {
			trimSecret := bytes.TrimLeft(secret[:], "\x00")
			decryptSecret := Decrypt(opening[:], string(trimSecret))
			if isInt(decryptSecret) {
				intNumber, err := strconv.Atoi(decryptSecret)
				if err != nil {
					log.Error("Can not convert string to integer", "error", err)
					return -1, err
				}
				random = int64(intNumber)
			}
		}
	}

	return random, nil
}

// Calculate reward for reward checkpoint.
func GetRewardForCheckpoint(c *S2PoS.S2PoS, chain consensus.ChainReader, header *types.Header, rCheckpoint uint64, totalSigner *uint64) (map[common.Address]*rewardLog, error) {
	// Not reward for singer of genesis block and only calculate reward at checkpoint block.
	number := header.Number.Uint64()
	prevCheckpoint := number - (rCheckpoint * 2)
	startBlockNumber := prevCheckpoint + 1
	endBlockNumber := startBlockNumber + rCheckpoint - 1
	signers := make(map[common.Address]*rewardLog)
	mapBlkHash := map[uint64]common.Hash{}

	data := make(map[common.Hash][]common.Address)
	for i := prevCheckpoint + (rCheckpoint * 2) - 1; i >= startBlockNumber; i-- {
		header = chain.GetHeader(header.ParentHash, i)

		mapBlkHash[i] = header.Hash()
		signData, ok := c.GetCachedSigningTxs(header.Hash())

		if !ok {
			log.Debug("Failed get from cached", "hash", header.Hash().String(), "number", i)
			block := chain.GetBlock(header.Hash(), i)
			txs := block.Transactions()
			if !chain.Config().IsTIPSigning(header.Number) {
				receipts := core.GetBlockReceipts(c.GetDb(), header.Hash(), i)
				signData = c.CacheNoneTIPSigningTxs(header, txs, receipts)
			} else {
				signData = c.CacheSigningTxs(header.Hash(), txs)
			}
		}
		txs := signData.([]*types.Transaction)
		for _, tx := range txs {
			blkHash := common.BytesToHash(tx.Data()[len(tx.Data())-32:])
			from := *tx.From()
			data[blkHash] = append(data[blkHash], from)
		}
	}
	header = chain.GetHeader(header.ParentHash, prevCheckpoint)
	masternodes := utils.GetMasternodesFromCheckpointHeader(header)

	epoch := chain.Config().S2PoS.Epoch
	for i := startBlockNumber; i <= endBlockNumber; i++ {

		if i%epoch < common.MergeSignRange || i%common.MergeSignRange == 0 || !chain.Config().IsTIP2019(big.NewInt(int64(i))) {
			addrs := data[mapBlkHash[i]]
			fmt.Println("utils.go:::addrs::GetRewardForCheckpoint", addrs)
			// Filter duplicate address.
			if len(addrs) > 0 {
				addrSigners := make(map[common.Address]bool)
				for _, masternode := range masternodes {
					for _, addr := range addrs {
						if addr == masternode {
							if _, ok := addrSigners[addr]; !ok {
								addrSigners[addr] = true
							}
							break
						}
					}
				}

				for addr := range addrSigners {
					_, exist := signers[addr]
					if exist {
						signers[addr].Sign++
					} else {
						signers[addr] = &rewardLog{1, new(big.Int)}
					}
					*totalSigner++
				}
			}
		}
	}

	log.Info("Calculate reward at checkpoint", "startBlock", startBlockNumber, "endBlock", endBlockNumber)

	return signers, nil
}

// // Calculate reward for signers.
// func CalculateRewardForSigner(chainReward *big.Int, signers map[common.Address]*rewardLog, totalSigner uint64, state *state.StateDB) (map[common.Address]*big.Int, error) {
// 	resultSigners := make(map[common.Address]*big.Int)
// 	fmt.Println("signers:::::", signers)

// 	totalStaking := new(big.Int)
// 	stakes := make(map[common.Address]*big.Int)

// 	for candidate, _ := range signers {
// 		candidatestake := stateDatabase.GetCandidateCap(state, candidate)
// 		stakes[candidate] = candidatestake
// 		fmt.Println("candidate:::::", candidate)
// 		fmt.Println("candidatestake:::::", candidatestake)
// 		totalStaking.Add(totalStaking, candidatestake)
// 	}
// 	fmt.Println("totalStaking:::::", totalStaking)

// 	// Add reward for signers.
// 	if totalSigner > 0 {
// 		for signer, rLog := range signers {
// 			stake := stakes[signer]

// 			if stake == nil {
// 				stake = big.NewInt(0)
// 			}

// 			// Calculate the stake ratio
// 			stakeRatio := new(big.Int).Mul(stake, big.NewInt(1e18))
// 			stakeRatio.Div(stakeRatio, totalStaking)

// 			// Calculate the block ratio
// 			blockRatio := new(big.Int).Mul(new(big.Int).SetUint64(rLog.Sign), big.NewInt(1e18))
// 			blockRatio.Div(blockRatio, new(big.Int).SetUint64(totalSigner))

// 			// Calculate the reward ratio
// 			rewardRatio := new(big.Int).Mul(stakeRatio, blockRatio)
// 			rewardRatio.Div(rewardRatio, big.NewInt(1e18))

// 			// Calculate the reward as (reward_ratio * chainReward)
// 			rewardValue := new(big.Int).Mul(rewardRatio, chainReward)
// 			rewardValue.Div(rewardValue, big.NewInt(1e18))

// 			rLog.Reward = rewardValue
// 			resultSigners[signer] = rewardValue

// 			fmt.Println("signer::::", signer)
// 			fmt.Println("rLog.Reward::::", rLog.Reward)
// 			fmt.Println("calcReward::::", rewardValue)
// 			fmt.Println("rewardRatio::::", rewardRatio)
// 			fmt.Println("stakeRatio::::", stakeRatio)
// 			fmt.Println("blockRatio::::", blockRatio)

// 		}
// 	}
// 	jsonSigners, err := json.Marshal(signers)
// 	if err != nil {
// 		log.Error("Fail to parse json signers", "error", err)
// 		return nil, err
// 	}
// 	log.Info("Signers data", "signers", string(jsonSigners), "totalSigner", totalSigner, "totalReward", chainReward)

// 	return resultSigners, nil
// }

// Calculate reward for signers.
func CalculateRewardForSigner(chainReward *big.Int, signers map[common.Address]*rewardLog, totalSigner uint64, state *state.StateDB, signReward *big.Int) (map[common.Address]*big.Int, error) {
	resultSigners := make(map[common.Address]*big.Int)
	// Add reward for signers.
	candicates := stateDatabase.GetCandidates(state)
	fmt.Println("candidates:::", candicates)
	targetValue := new(big.Int)
	targetValue.SetString("1000000000000000000000000", 10)
	fmt.Println("signers::::::::", signers)
	fmt.Println("chainReward::::::::", chainReward)

	if totalSigner > 0 {
		for signer, rLog := range signers {
			candidatestake := stateDatabase.GetCandidateCap(state, signer)
			calcReward := new(big.Int)

			if candidatestake.Cmp(targetValue) == 0 {
				// Add Reward for Master nodes per epoch
				fmt.Println("stakeCap is equal to 1000000000000000000000000")
				calcReward.Div(chainReward, new(big.Int).SetUint64(5))

			} else {
				// Add reward for signers per epoch.
				fmt.Println("stakeCap is not equal to 1000000000000000000000000")
				stakeInEth := new(big.Int).Div(candidatestake, new(big.Int).SetInt64(1000000000000000000)) // Convert from wei to ether
				stakeMultiplier := new(big.Int).Div(stakeInEth, new(big.Int).SetInt64(1000))               // Divide by 1000 ether to get multiplier
				calcReward = new(big.Int).Mul(signReward, stakeMultiplier)                                 // 0.046 * 1e18 (to account for wei)
			}
			fmt.Println("signReward::inClaculateRewardSigner::::", signReward)
			fmt.Println("calcReward:::inRewardSigners::::", calcReward)
			rLog.Reward = calcReward
			resultSigners[signer] = calcReward
		}
	}
	jsonSigners, err := json.Marshal(signers)
	if err != nil {
		log.Error("Fail to parse json signers", "error", err)
		return nil, err
	}
	log.Info("Signers data", "signers", string(jsonSigners), "totalSigner", totalSigner, "totalReward", chainReward)

	return resultSigners, nil
}

// Get candidate owner by address.
func GetCandidatesOwnerBySigner(state *state.StateDB, signerAddr common.Address) common.Address {
	owner := stateDatabase.GetCandidateOwner(state, signerAddr)
	fmt.Println("GetCandidatesOwnerBySigner", owner)
	return owner
}

func CalculateRewardForHolders(foundationWalletAddr common.Address, state *state.StateDB, signer common.Address, calcReward *big.Int, blockNumber uint64) (error, map[common.Address]*big.Int) {
	rewards, err := GetRewardBalancesRate(foundationWalletAddr, state, signer, calcReward, blockNumber)
	if err != nil {
		return err, nil
	}
	return nil, rewards
}

func GetRewardBalancesRate(foundationWalletAddr common.Address, state *state.StateDB, masterAddr common.Address, totalReward *big.Int, blockNumber uint64) (map[common.Address]*big.Int, error) {
	fmt.Println("masterAddr at GetRewardBalancesRate", masterAddr)
	owner := GetCandidatesOwnerBySigner(state, masterAddr)

	stakeCap := stateDatabase.GetCandidateCap(state, masterAddr)
	// fmt.Printf("stakeCap:::%t", stakeCap)
	// fmt.Println("stakeCap:::", stakeCap)

	// Assuming you want to compare stakeCap with a large number
	targetValue := new(big.Int)
	targetValue.SetString("1000000000000000000000000", 10) // The big number you're comparing to

	fmt.Println("owner:::", owner)
	balances := make(map[common.Address]*big.Int)

	// var rewardPercent uint64
	if stakeCap.Cmp(targetValue) == 0 {
		fmt.Println("stakeCap is equal to 1000000000000000000000000")
		foundationReward := new(big.Int).Mul(totalReward, new(big.Int).SetInt64(common.RewardFoundationPercent))
		foundationReward = new(big.Int).Div(foundationReward, new(big.Int).SetInt64(100))
		fmt.Println("foundationReward:::", foundationReward)
		balances[foundationWalletAddr] = foundationReward
		masterReward := new(big.Int).Sub(totalReward, foundationReward)
		balances[masterAddr] = masterReward
		fmt.Println("Reward Distributed", masterReward)
	} else {
		fmt.Println("stakeCap is not equal to 1000000000000000000000000")
		fmt.Println("Reward Distributed", totalReward)
		balances[masterAddr] = totalReward
	}

	// // Calculate the reward
	// reward := new(big.Int).Mul(totalReward, big.NewInt(int64(rewardPercent)))
	// reward.Div(reward, big.NewInt(100))

	// fmt.Println("Reward:::", reward)

	// Update the balance
	fmt.Println("totalReward::::", totalReward)
	// balances[masterAddr] = totalReward

	// Get voters for masternode.
	voters := stateDatabase.GetVoters(state, masterAddr)

	fmt.Println("voters for masternode:::", voters)

	//this logic wont make an impact, as VoterPercent is zero always.
	if len(voters) > 0 {
		totalVoterReward := new(big.Int).Mul(totalReward, new(big.Int).SetUint64(common.RewardVoterPercent))
		totalVoterReward = new(big.Int).Div(totalVoterReward, new(big.Int).SetUint64(100))
		totalCap := new(big.Int)
		fmt.Println("voters rewards:::", totalVoterReward)

		// Get voters capacities.
		voterCaps := make(map[common.Address]*big.Int)
		for _, voteAddr := range voters {
			if _, ok := voterCaps[voteAddr]; ok && common.TIP2019Block.Uint64() <= blockNumber {
				continue
			}
			voterCap := stateDatabase.GetVoterCap(state, masterAddr, voteAddr)
			fmt.Println("voterCap:::", voterCap)
			totalCap.Add(totalCap, voterCap)
			voterCaps[voteAddr] = voterCap
		}
		if totalCap.Cmp(new(big.Int).SetInt64(0)) > 0 {
			for addr, voteCap := range voterCaps {
				// Only valid voter has cap > 0.
				if voteCap.Cmp(new(big.Int).SetInt64(0)) > 0 {
					rcap := new(big.Int).Mul(totalVoterReward, voteCap)
					rcap = new(big.Int).Div(rcap, totalCap)
					if balances[addr] != nil {
						balances[addr].Add(balances[addr], rcap)
					} else {
						balances[addr] = rcap
					}
				}
			}
		}
	}

	jsonHolders, err := json.Marshal(balances)
	if err != nil {
		log.Error("Fail to parse json holders", "error", err)
		return nil, err
	}
	log.Trace("Holders reward", "holders", string(jsonHolders), "masternode", masterAddr.String())

	return balances, nil
}

// Dynamic generate array sequence of numbers.
func NewSlice(start int64, end int64, step int64) []int64 {
	s := make([]int64, end-start)
	for i := range s {
		s[i] = start
		start += step
	}

	return s
}

// Shuffle array.
func Shuffle(slice []int64) []int64 {
	newSlice := make([]int64, len(slice))
	copy(newSlice, slice)

	for i := 0; i < len(slice)-1; i++ {
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(newSlice))
		x := newSlice[i]
		newSlice[i] = newSlice[randIndex]
		newSlice[randIndex] = x
	}

	return newSlice
}

// encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("Fail to encrypt", "err", err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(cryptoRand.Reader, iv); err != nil {
		log.Error("Fail to encrypt iv", "err", err)
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("Fail to decrypt", "err", err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Error("ciphertext too short")
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

// Generate random string.
func RandStringByte(n int) []byte {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	b := make([]byte, n)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

// Helper function check string is numeric.
func isInt(strNumber string) bool {
	if _, err := strconv.Atoi(strNumber); err == nil {
		return true
	} else {
		return false
	}
}
