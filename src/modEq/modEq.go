package modEq

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	generate "../Root"
)

//q is taken as prime for group Zq
func Prove(crs, commitment, witness []*big.Int, lambda int64) []*big.Int {

	lambda_z := lambda - 3
	lambda_s := lambda - 2
	mu := 2*lambda - 2

	bound := new(big.Int).Exp(big.NewInt(2), big.NewInt(lambda_z+lambda_s+mu), nil)
	re := generate.Generate_random(bound)

	N_ := new(big.Int).Div(crs[0], big.NewInt(4))
	bound1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(lambda_s+lambda_z), nil)
	bound1.Mul(bound1, N_)
	bound1.Sub(bound1, big.NewInt(1))
	rr := generate.Generate_random(bound1)

	rrq, _ := rand.Int(rand.Reader, crs[3])

	alpha1 := generate.Generate_alpha(crs[1], crs[2], crs[0], re, rr)
	alpha2 := generate.Generate_alpha(crs[4], crs[5], crs[3], re.Mod(re, crs[3]), rrq)

	//c := big.NewInt(12345)
	x := []*big.Int{alpha1, alpha2, commitment[0], commitment[1]}
	hash := sha256.New()
	for _, y := range x {
		hash.Write(y.Bytes())
	}
	hash_string := fmt.Sprintf("%x", hash.Sum(nil))
	c, _ := new(big.Int).SetString(hash_string, 16)

	se := generate.Generate_s(re, c, witness[0])
	sr := generate.Generate_s(rr, c, witness[2])
	srq := generate.Generate_s(rrq, c, witness[3])
	srq.Mod(srq, crs[3])

	pi := []*big.Int{alpha1, alpha2, se, sr, srq}
	return pi

}

func VerProof(crs, commitment, pi []*big.Int) int {

	x := []*big.Int{pi[0], pi[1], commitment[0], commitment[1]}
	hash := sha256.New()
	for _, y := range x {
		hash.Write(y.Bytes())
	}
	hash_string := fmt.Sprintf("%x", hash.Sum(nil))
	c, _ := new(big.Int).SetString(hash_string, 16)

	alpha_1 := generate.Generate_alpha_ver(commitment[0], crs[1], crs[2], c, pi[2], pi[3], crs[0])
	alpha_2 := generate.Generate_alpha_ver(commitment[1], crs[4], crs[5], c, pi[2].Mod(pi[2], crs[0]), pi[4], crs[0])

	if pi[0] == alpha_1 && pi[1] == alpha_2 {
		return 1
	}

	return 0

}
