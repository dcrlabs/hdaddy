// Copyright (c) 2018, Decred Labs
// Copyright (c) 2017, Jonathan Chappelow
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdaddy

import (
	"fmt"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/wallet/udb"
)

const (
	// ExternalBranch indicates the derivation path for addresses intended for receiving payments
	ExternalBranch uint32 = udb.ExternalBranch
	// InternalBranch indicates the derivation path for change addresses
	InternalBranch uint32 = udb.InternalBranch
)

// AddressRangeExtendedKey computes a range of addresses derived from an
// extended key according to BIP0032.
func AddressRangeExtendedKey(xkey *hdkeychain.ExtendedKey, start, count uint32,
	params *chaincfg.Params) ([]dcrutil.Address, uint32, error) {

	addresses := make([]dcrutil.Address, 0, count)

	// Convert child keys on this chain to addresses, keeping track of the index on the chain
	var i, numInvalid uint32
	for i < count {
		// Tiny chance that consecutive indexes are not possible
		child, err := xkey.Child(start + i + numInvalid)
		if err == hdkeychain.ErrInvalidChild {
			numInvalid++
			continue
		}
		if err != nil {
			return nil, 0, err
		}

		addy, err := child.Address(params)
		if err != nil {
			return nil, 0, err
		}

		addresses = append(addresses, addy)
		i++
	}

	end := start + count + numInvalid

	return addresses, end, nil
}

// AddressRangeExtendedKeyStr is the same as AddressRangeExtendedKey, but
// accepts a string representation of the extended key.
func AddressRangeExtendedKeyStr(xkey string, start, count uint32,
	params *chaincfg.Params) ([]dcrutil.Address, uint32, error) {

	xkeyhd, err := hdkeychain.NewKeyFromString(xkey)
	if err != nil {
		return nil, 0, err
	}
	if !xkeyhd.IsForNet(params) {
		return nil, 0, fmt.Errorf("extended key is for wrong network")
	}

	return AddressRangeExtendedKey(xkeyhd, start, count, params)
}

// AddressAtBranchAndIndex computes the address derived from an extended
// key on given branch and index.
func AddressAtBranchAndIndex(xkey *hdkeychain.ExtendedKey, index uint32,
	branch uint32, params *chaincfg.Params) (dcrutil.Address, error) {
	// Get the extended key for the branch origin node
	branchKey, err := xkey.Child(branch)
	if err != nil {
		return nil, err
	}

	// Get the extended key of the node with the index on the branch
	key, err := branchKey.Child(index)
	if err != nil {
		return nil, err
	}

	// Convert to standard P2PKH address
	return key.Address(params)
}

// AddressRangeFromAccountKey treats the input key as the non-hardened extended
// key for a wallet account, deriving addresses for the specified branch (e.g.
// internal or external) and index range [start, start+count).
func AddressRangeFromAccountKey(xkey *hdkeychain.ExtendedKey, start, count, branch uint32,
	params *chaincfg.Params) ([]dcrutil.Address, uint32, error) {
	// Get the extended key for the branch origin node
	xkeyBranch, err := xkey.Child(branch)
	if err != nil {
		return nil, 0, err
	}

	return AddressRangeExtendedKey(xkeyBranch, start, count, params)
}

// AddressRangeFromAccountKeyString is the same as AddressRangeFromAccountKey
// except that the extended key is a string.
func AddressRangeFromAccountKeyString(xkey string, start, count, branch uint32,
	params *chaincfg.Params) ([]dcrutil.Address, uint32, error) {

	xkeyhd, err := hdkeychain.NewKeyFromString(xkey)
	if err != nil {
		return nil, 0, err
	}
	if !xkeyhd.IsForNet(params) {
		return nil, 0, fmt.Errorf("extended key is for wrong network")
	}

	return AddressRangeFromAccountKey(xkeyhd, start, count, branch, params)
}
