package main

import (
	"fmt"
	"math/big"

	member "github.com/GarryFCR/Rsazkp/src/SetMembership"
	setup "github.com/GarryFCR/Rsazkp/src/Setup"
)

func main() {

	//SETUP------------------------------------------------------------------
	//PEDERSEN COMMITMENT----------------------------------------------------
	prime, g, h := setup.Pedersen_setup(512, 512)
	ck_pedersen := []big.Int{prime, g, h}
	u := big.NewInt(12345)

	cu, ru := setup.Pedersen_commit(ck_pedersen, prime, *u)

	ver := setup.Pedersen_ver(ck_pedersen, cu, *u, ru)
	fmt.Println("--------------Commiting a set-member--------------")
	if ver == 1 {
		fmt.Println("Pedersen:Commitment VERIFIED")
	} else {
		fmt.Println("Pedersen verification failed")

	}

	//SET COMMITMENT----------------------------------------------------------
	N, G := setup.Set_setup(512, 512)

	ck_set := []big.Int{N, G}
	set := []big.Int{*big.NewInt(12342), *big.NewInt(12343), *big.NewInt(12344), *big.NewInt(12345)}

	Acc, _ := setup.Set_commit(ck_set, set)

	ver1 := setup.Set_ver(ck_set, set, Acc)
	fmt.Println()
	fmt.Println("-------------------Commiting the set--------------")
	if ver1 == 1 {
		fmt.Println("Set commit:Commitment VERIFIED")
	} else {
		fmt.Println("Set verification failed")

	}

	//SETMEMBERSHIP------------------------------------------------------------
	//KEYGEN-------------------------------------------------------------------
	ck_key := []big.Int{N, G, prime, g, h}
	fmt.Println()
	fmt.Println("Generating key for protocol...")
	crs := member.KeyGen(ck_key)

	//PROVE--------------------------------------------------------------------
	fmt.Println()
	fmt.Println("---------------------Proving----------------------")

	Ce, ce, pi_root, pi_mod, pi_hash := member.Prove(crs, set, Acc, cu, *u, ru)
	fmt.Println("Proof generated")
	//VERIFICATION--------------------------------------------------------------

	BOOL := member.VerProof(crs, Acc, cu, Ce, ce, pi_root, pi_mod, pi_hash)
	fmt.Println("---------------------Verification-----------------")
	if BOOL == 1 {
		fmt.Println("Verified")
	} else {
		fmt.Println("Verification failed")
	}

}
