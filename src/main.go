package main

import (
	"fmt"
	"math/big"

	member "./SetMembership"
	setup "./Setup"
)

func main() {

	//SETUP------------------------------------------------------------------
	//PEDERSEN COMMITMENT----------------------------------------------------
	prime, g, h := setup.Pedersen_setup(512, 512)
	ck_pedersen := []big.Int{prime, g, h}
	u := big.NewInt(12345)

	cu, ru := setup.Pedersen_commit(ck_pedersen, prime, *u)

	ver := setup.Pedersen_ver(ck_pedersen, cu, *u, ru)
	if ver == 1 {
		fmt.Println("Pedersen:Commitment VERIFIED")
	}

	//SET COMMITMENT----------------------------------------------------------
	N, G := setup.Set_setup(512, 512)

	ck_set := []big.Int{N, G}
	set := []big.Int{*big.NewInt(12342), *big.NewInt(12343), *big.NewInt(12344), *big.NewInt(12345)}

	Acc, _ := setup.Set_commit(ck_set, set)

	ver1 := setup.Set_ver(ck_set, set, Acc)
	if ver1 == 1 {
		fmt.Println("Set commit:Commitment VERIFIED")
	}

	//SETMEMBERSHIP------------------------------------------------------------
	//KEYGEN-------------------------------------------------------------------
	ck_key := []big.Int{N, G, prime, g, h}
	crs := member.KeyGen(ck_key)

	//PROVE--------------------------------------------------------------------

	Ce, ce, pi_root, pi_mod, pi_hash := member.Prove(crs, set, Acc, cu, *u, ru)

	//VERIFICATION--------------------------------------------------------------

	BOOL := member.VerProof(crs, Acc, cu, Ce, ce, pi_root, pi_mod, pi_hash)

	if BOOL == 1 {
		fmt.Println("Verified")
	}

}
