/*
	Proof of Equality of Commitment

	The ModEq protocol gives a ZK proof of the fact that C_e (committed in RSA group) and c_e (committed in prime order group)
	commits to the same value modulo q
*/

package modEq

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	hash2prime "../Hashtoprime"
	generate "../Root"
)

//Prove function returns the proof of the fact that C-e and c_e commits to the same value
//q is taken as prime for group Zq
func Prove(crs, commitment, witness []big.Int, lambda int64) []big.Int {

	lambda_z := lambda - 3
	lambda_s := lambda - 2
	mu := 2*lambda - 2
	N, G, H, q, g, h := crs[0], crs[1], crs[2], crs[3], crs[4], crs[5]
	e, r, rq := witness[0], witness[2], witness[3]

	bound := new(big.Int).Exp(big.NewInt(2), big.NewInt(lambda_z+lambda_s+mu), nil)
	re := generate.Generate_random(*bound)

	N_ := new(big.Int).Div(&N, big.NewInt(4))
	bound1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(lambda_s+lambda_z), nil)
	bound1.Mul(bound1, N_)
	bound1.Sub(bound1, big.NewInt(1))
	rr := generate.Generate_random(*bound1)

	rrq, _ := rand.Int(rand.Reader, &q)

	alpha1 := generate.Generate_alpha(G, H, N, re, rr)
	alpha2 := generate.Generate_alpha(g, h, q, *new(big.Int).Mod(&re, &q), *rrq)

	x := []big.Int{alpha1, alpha2, commitment[0], commitment[1]}
	hash := sha256.New()
	for _, y := range x {
		hash.Write(y.Bytes())
	}
	hash_string := fmt.Sprintf("%x", hash.Sum(nil))
	c_hash, _ := new(big.Int).SetString(hash_string, 16)
	c := hash2prime.Fu(*c_hash)

	se := generate.Generate_s(re, c, e)
	sr := generate.Generate_s(rr, c, r)
	srq := generate.Generate_s(*rrq, c, rq)
	srq.Mod(&srq, &q)

	pi := []big.Int{alpha1, alpha2, se, sr, srq}
	return pi

}

// VerProof verifies the correctness of the proof and returns 1 if correct
// Else returns 0
func VerProof(crs, commitment, pi []big.Int) int {

	x := []big.Int{pi[0], pi[1], commitment[0], commitment[1]}
	hash := sha256.New()
	for _, y := range x {
		hash.Write(y.Bytes())
	}
	hash_string := fmt.Sprintf("%x", hash.Sum(nil))
	c_hash, _ := new(big.Int).SetString(hash_string, 16)
	c := hash2prime.Fu(*c_hash)

	alpha_1 := generate.Generate_alpha_ver(commitment[0], crs[1], crs[2], c, pi[2], pi[3], crs[0])
	alpha_2 := generate.Generate_alpha_ver(commitment[1], crs[4], crs[5], c, *new(big.Int).Mod(&pi[2], &crs[3]), pi[4], crs[3])
	fmt.Println(alpha_1, alpha_2)
	fmt.Println(pi[0], pi[1])

	w := generate.Generate_alpha_ver(commitment[1], crs[4], crs[5], c, *new(big.Int).Mod(&pi[2], &crs[3]), pi[4], crs[3])

	fmt.Println("xxx:", w)

	if pi[0].Cmp(&alpha_1) == 0 && pi[1].Cmp(&alpha_2) == 0 {
		return 1
	}

	return 0

}
