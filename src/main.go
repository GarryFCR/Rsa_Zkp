package main

import (
	"fmt"
	"math/big"

	HashEq "./HashEq"
	hash2prime "./Hashtoprime"
	root "./Root"
	setm "./SetMembership"
	setup "./Setup"
	mod "./modEq"
)

func main() {

	//PEDERSEN COMMITMENT-----------------------------------------
	prime, g, h := setup.Pedersen_setup(12, 12)
	ck := []big.Int{prime, g, h}
	u := big.NewInt(12345)

	c, o := setup.Pedersen_commit(ck, prime, *u)

	ver := setup.Pedersen_ver(ck, c, *u, o)
	if ver == 1 {
		fmt.Println("Pedersen:Commitment VERIFIED")
	}

	//SET COMMITMENT--------------------------------------------
	N, G, p, q := setup.Set_setup(12, 12)

	ck1 := []big.Int{N, G}
	set := []big.Int{*big.NewInt(12342), *big.NewInt(12343), *big.NewInt(12344), *big.NewInt(12345)}

	com, _ := setup.Set_commit(ck1, set)

	ver1 := setup.Set_ver(ck1, set, com)
	if ver1 == 1 {
		fmt.Println("Set commit:Commitment VERIFIED")
	}

	//Setmembership(keygen)--------------------------------------------------------------------
	ck_key := []big.Int{N, G, prime, g, h}
	crs := setm.KeyGen(ck_key)

	//Root-----------------------------------------------------------------------------------
	e := hash2prime.Hprime(*u)
	Ce, r := setup.Pedersen_commit(crs[:3], N, e)

	phi := new(big.Int).Mul(new(big.Int).Sub(&q, big.NewInt(1)), new(big.Int).Sub(&p, big.NewInt(1)))
	inverse := new(big.Int).ModInverse(&e, phi)
	W := new(big.Int).Exp(&com, inverse, &N)

	commit := []big.Int{Ce, com}
	root_witness := []big.Int{e, r, *W}

	pi_root := root.Prove(crs[0:3], commit, root_witness, int64(12), int64(12), int64(12))

	ver2 := root.VerProof(crs[0:3], commit, pi_root, int64(12), int64(12), int64(12))
	if ver2 == 1 {
		fmt.Println("Root :Root VERIFIED")
	}

	//ModEq-----------------------------------------------------------------------------------
	ce, rq := setup.Pedersen_commit(crs[3:], prime, e)
	commit_mod := []big.Int{Ce, ce}
	mod_witness := []big.Int{e, e, r, rq}
	pi_mod := mod.Prove(crs, commit_mod, mod_witness, int64(64))
	//fmt.Println(ce)

	ver3 := mod.VerProof(crs, commit_mod, pi_mod)
	if ver3 == 1 {
		fmt.Println("ModEq :Modeq VERIFIED")
	}

	//HashEq-----------------------------------------------------------------------------------
	fu := hash2prime.Fu(*u)
	j := new(big.Int).Sub(&e, &fu)

	commit_hash := []big.Int{ce, c}
	hash_witness := []big.Int{e, *u, rq, o, *j}

	pi_hash := HashEq.Prove(crs[3:], commit_hash, hash_witness, int64(12), int64(12), int64(12), int64(12))

	ver4 := HashEq.VerProof(crs[3:], commit_hash, pi_hash)
	if ver4 == 1 {
		fmt.Println("HashEq :hash VERIFIED")
	}
}
