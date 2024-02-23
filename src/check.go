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
	u := big.NewInt(10000011)

	cu, ru := setup.Pedersen_commit(ck_pedersen, prime, *u)

	//SET COMMITMENT----------------------------------------------------------
	N, G := setup.Set_setup(512, 512)

	ck_set := []big.Int{N, G}
	set := []big.Int{*big.NewInt(10000001), *big.NewInt(10000002),
		*big.NewInt(10000003), *big.NewInt(10000000), *big.NewInt(10000011)}

	Acc, _ := setup.Set_commit(ck_set, set)

	//SETMEMBERSHIP------------------------------------------------------------
	//KEYGEN-------------------------------------------------------------------
	ck_key := []big.Int{N, G, prime, g, h}
	crs := member.KeyGen(ck_key)

	//PROVE--------------------------------------------------------------------
	fmt.Println()

	//VERIFICATION--------------------------------------------------------------
	var correct int
	var wrong int
	var iter = 10
	fmt.Println("Computing")
	for i := 0; i < iter; i++ {
		Ce, ce, pi_root, pi_mod, pi_hash := member.Prove(crs, set, Acc, cu, *u, ru)
		BOOL := member.VerProof(crs, Acc, cu, Ce, ce, pi_root, pi_mod, pi_hash)
		if BOOL == 1 {

			correct++
		} else {

			wrong++
		}
	}

	// Convert to float64 for accurate percentage calculation
	fmt.Println("Done")
	correctPercentage := (float64(correct) / float64(iter)) * 100
	wrongPercentage := (float64(wrong) / float64(iter)) * 100

	fmt.Printf("Probability of correct: %.2f%%\n", correctPercentage)
	fmt.Printf("Probability of wrong: %.2f%%\n", wrongPercentage)

}
