/*
	SetCommitment is a (non-hiding) commitment scheme for set of strings.
	In this set Commitment scheme we map the set members to set of prime using Hash2prime function
	and then the Set Commitment is build as an RSA Accumulator to the set of those primes derived from
	set of strings using Has2prime
*/


package Setup

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"

	hash2prime "../Hashtoprime"
)

// Set_setup generates N := pq where p and q are lambda(security parameter) bit prime for RSA group
// It also generates random G which is in Quadratic Residue Modulo N 
func Set_setup(lambda, nu int) (n, f *big.Int) {

	pk, _ := rsa.GenerateKey(rand.Reader, lambda)
	var F *big.Int
	N := pk.PublicKey.N

	for {
		F, _ = rand.Int(rand.Reader, N)
		if F != pk.Primes[0] && F != pk.Primes[1] && F != big.NewInt(1) {
			break
		}
	}

	G := F.Exp(F, big.NewInt(2), pk.PublicKey.N)

	return N, G
}

//Set_commit generates the RSA Accumulator for the set of primes derived from set of strings using Hash2prime

func Set_commit(ck, U []*big.Int) (c, o *big.Int) {

	P := make([]*big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = hash2prime.Hprime(u)
		G.Exp(G, P[i], N)
	}

	fmt.Println(P)
	return G, nil

}

// Set_ver verifies the correctness of the set commitment (RSA Accumulator) and returns 1 if correct
// Else returns 0
func Set_ver(ck, U []*big.Int, Acc *big.Int) int {

	P := make([]*big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = hash2prime.Hprime(u)
		G.Exp(G, P[i], N)
	}

	if Acc == G {
		return 1
	}
	return 0

}
