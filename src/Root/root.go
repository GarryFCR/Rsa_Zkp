/*
The Root Protocol is the NIZK proof of a committed root of public RSA group element Acc
It takes an integer commitment to an element e and proves knowledge of an e-th root of Acc ie, W=Acc^(1/e) in zero knowledge
*/
package Root

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	//"unsafe"

	hash2prime "github.com/GarryFCR/Rsazkp/src/Hashtoprime"
)

func Generate_random(bound big.Int, lambdas ...int64) big.Int {
	var lambda int64
	if len(lambdas) > 0 {
		lambda = lambdas[0]
	} else {
		lambda = -1 // Indicates no lambda provided
	}

	upper := bound
	lower := new(big.Int).Neg(&upper)

	max := new(big.Int).Sub(&upper, lower)
	max.Add(max, big.NewInt(1))

	for {
		random, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("random number generation failed") // Handle error appropriately
		}
		random.Add(random, lower)

		// Check if within lambda bounds, if lambda is provided
		if lambda == -1 || isWithinLambda(random, lambda) {
			return *random
		}
	}
}

func Generate_s(a, b, c big.Int) big.Int {

	s := new(big.Int).Mul(&b, &c)
	s.Sub(&a, s)

	return *s

}

func Generate_alpha(G, H, N, x, y big.Int) big.Int {

	if x.Sign() == -1 {
		G = *new(big.Int).ModInverse(&G, &N)
		x.Neg(&x)
	}
	if y.Sign() == -1 {
		H = *new(big.Int).ModInverse(&H, &N)
		y.Neg(&y)
	}

	alpha := big.NewInt(0).Exp(&G, &x, &N)
	temp := big.NewInt(0).Exp(&H, &y, &N)
	alpha.Mul(alpha, temp)
	alpha.Mod(alpha, &N)

	return *alpha

}

func Generate_alpha_ver(G, H, C, x, y, z, N big.Int) big.Int {

	if x.Sign() == -1 {
		G = *new(big.Int).ModInverse(&G, &N)
		x.Neg(&x)
	}
	if y.Sign() == -1 {
		H = *new(big.Int).ModInverse(&H, &N)
		y.Neg(&y)
	}
	if z.Sign() == -1 {
		C = *new(big.Int).ModInverse(&C, &N)
		z.Neg(&z)
	}

	alpha := new(big.Int).Exp(&G, &x, &N)
	temp := new(big.Int).Exp(&H, &y, &N)
	temp1 := new(big.Int).Exp(&C, &z, &N)
	alpha.Mul(alpha, temp)
	alpha.Mod(alpha, &N)
	alpha.Mul(alpha, temp1)
	alpha.Mod(alpha, &N)

	return *alpha
}

func Prove(crs, commitment, witness []big.Int, lambda_s, lambda_z, mu int64) []big.Int {

	N, G, H := crs[0], crs[1], crs[2]
	e, r, W := witness[0], witness[1], witness[2]

	bound := *new(big.Int).Div(&N, big.NewInt(4))

	r2 := Generate_random(*new(big.Int).Sub(&bound, big.NewInt(1)))

	r3 := Generate_random(*new(big.Int).Sub(&bound, big.NewInt(1)))

	Cw := Generate_alpha(W, H, N, *big.NewInt(1), r2)
	Cr := Generate_alpha(G, H, N, r2, r3)

	boundLambdaMu := new(big.Int).Exp(big.NewInt(2), big.NewInt(lambda_s+lambda_z+mu), nil)
	boundLambda := new(big.Int).Exp(big.NewInt(2), big.NewInt(lambda_s+lambda_z), nil)
	boundLambda.Mul(boundLambda, &bound).Sub(boundLambda, big.NewInt(1))
	boundLambdaMuSpecial := boundLambdaMu.Mul(boundLambdaMu, &bound)

	// Generate random numbers within lambda bounds
	re := Generate_random(*boundLambdaMu, lambda_s+lambda_z+mu)
	rr := Generate_random(*boundLambda, lambda_s+lambda_z)
	rr2 := Generate_random(*boundLambda, lambda_s+lambda_z)
	rr3 := Generate_random(*boundLambda, lambda_s+lambda_z)
	r_beta := Generate_random(*boundLambdaMuSpecial, lambda_s+lambda_z+mu)
	r_delta := Generate_random(*boundLambdaMuSpecial, (lambda_s + lambda_z + mu))

	alpha1 := Generate_alpha(G, H, N, re, rr)
	alpha2 := Generate_alpha(G, H, N, rr2, rr3)
	alpha3 := Generate_alpha(Cw, *new(big.Int).ModInverse(&H, &N), N, re, r_beta)
	alpha4 := Generate_alpha_ver(*new(big.Int).ModInverse(&G, &N), *new(big.Int).ModInverse(&H, &N), Cr, r_beta, r_delta, re, N)
	//fmt.Println("xxx:", alpha1, alpha2, alpha3, alpha4)

	list := []big.Int{alpha1, alpha2, alpha3, alpha4, commitment[0], commitment[1]}
	h := sha256.New()
	for _, y := range list {
		h.Write(y.Bytes())
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	c_hash := new(big.Int)
	c_hash.SetString(hash, 16)
	c := hash2prime.Fu(*c_hash)

	se := Generate_s(re, c, e)
	sr := Generate_s(rr, c, r)
	sr2 := Generate_s(rr2, c, r2)
	sr3 := Generate_s(rr3, c, r3)
	s_beta := Generate_s(r_beta, *new(big.Int).Mul(&c, &e), r2)
	s_delta := Generate_s(r_delta, *new(big.Int).Mul(&c, &e), r3)

	pi := []big.Int{Cw, Cr, alpha1, alpha2, alpha3, alpha4, se, sr, sr2, sr3, s_beta, s_delta}
	//fmt.Println("pi:", pi[2:6])
	//fmt.Printf("Size:%d\n", unsafe.Sizeof(pi))

	return pi

}

func VerProof(crs, commitment, pi []big.Int, lambda, lambda_s, mu int64) int {

	N, G, H := crs[0], crs[1], crs[2]
	Ce, Acc := commitment[0], commitment[1]
	alpha1, alpha2, alpha3, alpha4 := pi[2], pi[3], pi[4], pi[5]

	list := []big.Int{alpha1, alpha2, alpha3, alpha4, Ce, Acc}
	h := sha256.New()
	for _, y := range list {
		h.Write(y.Bytes())
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	c_hash := new(big.Int)
	c_hash.SetString(hash, 16)
	c := hash2prime.Fu(*c_hash)

	alpha_1 := Generate_alpha_ver(Ce, G, H, c, pi[6], pi[7], N)
	alpha_2 := Generate_alpha_ver(pi[1], G, H, c, pi[8], pi[9], N)
	alpha_3 := Generate_alpha_ver(Acc, pi[0], *new(big.Int).ModInverse(&H, &N), c, pi[6], pi[10], N)
	alpha_4 := Generate_alpha_ver(pi[1], *new(big.Int).ModInverse(&G, &N), *new(big.Int).ModInverse(&H, &N), pi[6], pi[10], pi[11], N)
	//fmt.Println(alpha_1, alpha_2, alpha_3, alpha_4)

	upper := big.NewInt(2)
	upper.Exp(upper, big.NewInt(lambda+lambda_s+mu+1), nil)
	lower := new(big.Int).Neg(upper)

	var se_bool bool
	upper_bound := new(big.Int).Sub(upper, &pi[6])
	lower_bound := new(big.Int).Sub(&pi[6], lower)

	if upper_bound.Sign() == 1 {
		if pi[6].Sign() == -1 && lower_bound.Sign() == -1 {
			se_bool = true

		}
		if pi[6].Sign() == 1 && lower_bound.Sign() == 1 {
			se_bool = true

		}
	}

	if alpha1.Cmp(&alpha_1) == 0 && alpha2.Cmp(&alpha_2) == 0 && alpha3.Cmp(&alpha_3) == 0 && alpha4.Cmp(&alpha_4) == 0 && se_bool {
		return 1
	} else {
		fmt.Println("Root verification failed")

	}

	//fmt.Println(alpha1.Cmp(&alpha_1), alpha2.Cmp(&alpha_2), alpha3.Cmp(&alpha_3), alpha4.Cmp(&alpha_4), se_bool)

	return 0

}

func isWithinLambda(value *big.Int, lambda int64) bool {
	bound := big.NewInt(1)         // 2^0 = 1
	bound.Lsh(bound, uint(lambda)) // 2^lambda
	return value.Cmp(bound) == -1  // check if value < 2^lambda
}
