// Copyright (c) 2018, Decred Labs
// Copyright (c) 2017, Jonathan Chappelow
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dcrlabs/hdaddy"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
)

var xkey = flag.String("key", "", "extended key")
var index = flag.Int("index", 0, "The starting index of nodes for which to derive addresses")
var count = flag.Int("count", 1, "The number of addresses to derive")
var branch = flag.Int("branch", -1, "(Optional) Derive addresses on a child branch (e.g. external/0 or internal/1)")
var testnet = flag.Bool("testnet", false, "key is for testnet")
var simnet = flag.Bool("simnet", false, "key is for simnet")

func runMain() error {
	defer func() {
		_ = os.Stdout.Sync()
	}()
	flag.Parse()

	chainParams := &chaincfg.MainNetParams
	if *testnet {
		chainParams = &chaincfg.TestNet2Params
		if *simnet {
			return fmt.Errorf("Cannot specify both simnet and testnet")
		}
	} else if *simnet {
		chainParams = &chaincfg.SimNetParams
	}

	startUint := uint32(*index)
	countUint := uint32(*count)
	branchUint := uint32(*branch)

	var addrs []dcrutil.Address
	var err error
	// branch -1 indicates to derive keys (and addresses) for immediate children
	// of the input extended key. A non-negative branch will derive keys one
	// level deeper, on the specified branch.
	if *branch == -1 {
		addrs, _, err = hdaddy.AddressRangeExtendedKeyStr(*xkey, startUint, countUint, chainParams)
		if err != nil {
			return err
		}
	} else {
		addrs, _, err = hdaddy.AddressRangeFromAccountKeyString(*xkey, startUint, countUint, branchUint, chainParams)
		if err != nil {
			return err
		}
	}

	for i := range addrs {
		fmt.Printf("%s\n", addrs[i].String())
	}

	return nil
}

func main() {
	if err := runMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
