// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/MarconiProtocol/go-methereum-lite/common"
	"github.com/MarconiProtocol/go-methereum-lite/common/hexutil"
	"github.com/MarconiProtocol/go-methereum-lite/core/types"
)

// Tests that ethash works correctly in test mode.
func TestTestMode(t *testing.T) {
	header := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(100)}

	ethash := NewTester(nil, false)
	defer ethash.Close()

	results := make(chan *types.Block)
	err := ethash.Seal(nil, types.NewBlockWithHeader(header), results, nil)
	if err != nil {
		t.Fatalf("failed to seal block: %v", err)
	}
	select {
	case block := <-results:
		header.Nonce = types.EncodeNonce(block.Nonce())
		header.MixDigest = block.MixDigest()
		if err := ethash.VerifySeal(nil, header); err != nil {
			t.Fatalf("unexpected verification error: %v", err)
		}
	case <-time.NewTimer(time.Second).C:
		t.Error("sealing result timeout")
	}
}

// Tests that block sealing and seal verification works correctly in
// quick test mode.
func TestQuickTestMode(t *testing.T) {
	head := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(100)}

	ethash := NewQuickTest(nil, false)
	defer ethash.Close()

	// Override the randomness used for nonce generation to have a
	// constant seed, that way this test is deterministic.
	ethash.rand = rand.New(rand.NewSource(0))
	results := make(chan *types.Block)
	err := ethash.Seal(nil, types.NewBlockWithHeader(head), results, nil)
	if err != nil {
		t.Fatalf("failed to seal block: %v", err)
	}
	var block *types.Block = nil
	select {
	case block = <-results:
		head.Nonce = types.EncodeNonce(block.Nonce())
		head.MixDigest = block.MixDigest()
	case <-time.NewTimer(time.Second).C:
		t.Fatalf("sealing result timeout")
	}

	if err := ethash.VerifySeal(nil, head); err != nil {
		t.Fatalf("unexpected verification error: %v", err)
	}
	// If we change the nonce, verification should now fail.
	head.Nonce = types.EncodeNonce(block.Nonce() + 1)
	if err := ethash.VerifySeal(nil, head); err != errInvalidMixDigest {
		t.Fatalf("expected invalid mix digest but got: %v", err)
	}
	// As a sanity check, changing the nonce back should succeed again.
	head.Nonce = types.EncodeNonce(block.Nonce())
	if err := ethash.VerifySeal(nil, head); err != nil {
		t.Fatalf("unexpected verification error: %v", err)
	}
}

// Tests that block sealing and seal verification works correctly in
// cryptonight mode.
func TestCryptonightMode(t *testing.T) {
  // Note if the difficulty is too high, since cryptonight gets a much
  // lower hash rate than other algorithms, on slower machines the
  // test may time out without finding a valid nonce. If that happens,
  // you can decrease the difficulty or increase the NewTimer's
  // timeout below.
	head := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(100)}

	ethash := NewCryptonight(nil, false)
	defer ethash.Close()

	// Override the randomness used for nonce generation to have a
	// constant seed, that way this test is deterministic.
	ethash.rand = rand.New(rand.NewSource(0))
	results := make(chan *types.Block)
	err := ethash.Seal(nil, types.NewBlockWithHeader(head), results, nil)
	if err != nil {
		t.Fatalf("failed to seal block: %v", err)
	}
	var block *types.Block = nil
	select {
	case block = <-results:
		head.Nonce = types.EncodeNonce(block.Nonce())
		head.MixDigest = block.MixDigest()
	case <-time.NewTimer(time.Second * 30).C:
		t.Fatalf("sealing result timeout")
	}

	if err := ethash.VerifySeal(nil, head); err != nil {
		t.Fatalf("unexpected verification error: %v", err)
	}
	// If we change the nonce, verification should now fail.
	head.Nonce = types.EncodeNonce(block.Nonce() + 1)
	if err := ethash.VerifySeal(nil, head); err != errInvalidMixDigest {
		t.Fatalf("expected invalid mix digest but got: %v", err)
	}
	// As a sanity check, changing the nonce back should succeed again.
	head.Nonce = types.EncodeNonce(block.Nonce())
	if err := ethash.VerifySeal(nil, head); err != nil {
		t.Fatalf("unexpected verification error: %v", err)
	}
}

// This test checks that cache lru logic doesn't crash under load.
// It reproduces https://github.com/ethereum/go-ethereum/issues/14943
func TestCacheFileEvict(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "ethash-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	e := New(Config{CachesInMem: 3, CachesOnDisk: 10, CacheDir: tmpdir, PowMode: ModeTest}, nil, false)
	defer e.Close()

	workers := 8
	epochs := 100
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go verifyTest(&wg, e, i, epochs)
	}
	wg.Wait()
}

func verifyTest(wg *sync.WaitGroup, e *Ethash, workerIndex, epochs int) {
	defer wg.Done()

	const wiggle = 4 * epochLength
	r := rand.New(rand.NewSource(int64(workerIndex)))
	for epoch := 0; epoch < epochs; epoch++ {
		block := int64(epoch)*epochLength - wiggle/2 + r.Int63n(wiggle)
		if block < 0 {
			block = 0
		}
		header := &types.Header{Number: big.NewInt(block), Difficulty: big.NewInt(100)}
		e.VerifySeal(nil, header)
	}
}

func TestRemoteSealer(t *testing.T) {
	ethash := NewTester(nil, false)
	defer ethash.Close()

	api := &API{ethash}
	if _, err := api.GetWork(); err != errNoMiningWork {
		t.Error("expect to return an error indicate there is no mining work")
	}
	header := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(100)}
	block := types.NewBlockWithHeader(header)
	sealhash := ethash.SealHash(header)

	// Push new work.
	results := make(chan *types.Block)
	ethash.Seal(nil, block, results, nil)

	var (
		work [4]string
		err  error
	)
	if work, err = api.GetWork(); err != nil || work[0] != sealhash.Hex() {
		t.Error("expect to return a mining work has same hash")
	}

	if res := api.SubmitWork(types.BlockNonce{}, sealhash, common.Hash{}); res {
		t.Error("expect to return false when submit a fake solution")
	}
	// Push new block with same block number to replace the original one.
	header = &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(1000)}
	block = types.NewBlockWithHeader(header)
	sealhash = ethash.SealHash(header)
	ethash.Seal(nil, block, results, nil)

	if work, err = api.GetWork(); err != nil || work[0] != sealhash.Hex() {
		t.Error("expect to return the latest pushed work")
	}
}

func TestHashRate(t *testing.T) {
	var (
		hashrate = []hexutil.Uint64{100, 200, 300}
		expect   uint64
		ids      = []common.Hash{common.HexToHash("a"), common.HexToHash("b"), common.HexToHash("c")}
	)
	ethash := NewTester(nil, false)
	defer ethash.Close()

	if tot := ethash.Hashrate(); tot != 0 {
		t.Error("expect the result should be zero")
	}

	api := &API{ethash}
	for i := 0; i < len(hashrate); i += 1 {
		if res := api.SubmitHashRate(hashrate[i], ids[i]); !res {
			t.Error("remote miner submit hashrate failed")
		}
		expect += uint64(hashrate[i])
	}
	if tot := ethash.Hashrate(); tot != float64(expect) {
		t.Error("expect total hashrate should be same")
	}
}

func TestClosedRemoteSealer(t *testing.T) {
	ethash := NewTester(nil, false)
	time.Sleep(1 * time.Second) // ensure exit channel is listening
	ethash.Close()

	api := &API{ethash}
	if _, err := api.GetWork(); err != errEthashStopped {
		t.Error("expect to return an error to indicate ethash is stopped")
	}

	if res := api.SubmitHashRate(hexutil.Uint64(100), common.HexToHash("a")); res {
		t.Error("expect to return false when submit hashrate to a stopped ethash")
	}
}
