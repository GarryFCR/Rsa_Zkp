package Setup

import (
	"math/big"
	"testing"

	hash2prime "../Hashtoprime"
)

//Test for the pedersen commitment
func TestPedersen(t *testing.T) {

	//checking if the recieved parameters are correct
	prime, g, h := Pedersen_setup(15, 15)
	if !prime.ProbablyPrime(10) || prime.BitLen() > 15 || g.BitLen() > 15 || h.BitLen() > 15 {
		t.Fatalf("Either %v is not prime or prime,g,h are bigger than expected ", prime)
	}

	ck_pedersen := []big.Int{prime, g, h}
	u := big.NewInt(12345)

	//checking if the commitment is correct
	cu, ru := Pedersen_commit(ck_pedersen, prime, *u)
	commitment := new(big.Int).Mul(g.Exp(&g, u, &prime), h.Exp(&h, &ru, &prime))
	commitment.Mod(commitment, &prime)

	if commitment.Cmp(&cu) != 0 {
		t.Fatalf("Commitment failed")

	}

	//checking if verification pass for a correct commitment
	ver := Pedersen_ver(ck_pedersen, cu, *u, ru)
	if ver == 0 {
		t.Fatalf("Pedersen verification failed")
	}

	//forcing  a failure
	commitment.Mul(commitment, big.NewInt(2))
	ver_fail := Pedersen_ver(ck_pedersen, *commitment, *u, ru)
	if ver_fail == 1 {
		t.Fatalf("Pedersen verification passed on an Invalid commitment")
	}

}

//Test for Set commitment
func TestSetcommit(t *testing.T) {

	//checking parameters generated
	N, G := Set_setup(15, 15)
	if G.BitLen() > N.BitLen() {
		t.Fatalf("Parameters are not as expected")
	}

	//checking commitment correctness

	ck_set := []big.Int{N, G}
	set := []big.Int{*big.NewInt(12342), *big.NewInt(12343), *big.NewInt(12344), *big.NewInt(12345)}
	for _, u := range set {
		prime := hash2prime.Hprime(u)
		G.Exp(&G, &prime, &N)
	}

	Acc, _ := Set_commit(ck_set, set)
	if Acc.Cmp(&G) != 0 {
		t.Fatalf("Accumulator not correct")
	}

	//checking if verification pass for a correct commitment
	ver := Set_ver(ck_set, set, Acc)
	if ver == 0 {
		t.Fatalf("Set commitment verification failed")
	}

	//checking if verification pass for a incorrect commitment
	ver_fail := Set_ver(ck_set, set, *Acc.Add(&Acc, big.NewInt(1)))
	if ver_fail == 1 {
		t.Fatalf("Set commitment verification passed for incorrect commitment")
	}

}
