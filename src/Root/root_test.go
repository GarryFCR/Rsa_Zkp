package Root

import (
	"crypto/rand"
	"math/big"
	"testing"

	hash2prime "../Hashtoprime"
	setup "../Setup"
)

func TestProve(t *testing.T) {

	u := *big.NewInt(12345)
	e := hash2prime.Hprime(u)
	prime, g, h := setup.Pedersen_setup(12, 12)
	N, G := setup.Set_setup(12, 12)

	ck_set := []big.Int{N, G}
	set := []big.Int{*big.NewInt(12342), *big.NewInt(12343), *big.NewInt(12344), *big.NewInt(12345)}
	Acc, _ := setup.Set_commit(ck_set, set)

	var H *big.Int
	for {
		H, _ = rand.Int(rand.Reader, &N)
		if H.GCD(nil, nil, H, &N).Cmp(big.NewInt(1)) == 0 {
			break
		}
	}
	H_ := *new(big.Int).Exp(H, big.NewInt(2), &N)

	crs := []big.Int{N, G, H_, prime, g, h}

	Ce, r := setup.Pedersen_commit(crs[:3], crs[0], e)

	Primes := make([]big.Int, len(set))
	for i, u_dash := range set {
		Primes[i] = hash2prime.Hprime(u_dash)
		if u_dash.Cmp(&u) != 0 {

			G.Exp(&G, &Primes[i], &crs[0])
		}

	}

	W := G
	commit := []big.Int{Ce, Acc}
	root_witness := []big.Int{e, r, W}

	//Generation of proof of the root
	pi_root := Prove(crs[:3], commit, root_witness, int64(12), int64(12), int64(12))
	//verification of proof----------------------------------------------------------
	root_bool := VerProof(crs[:3], commit, pi_root, int64(12), int64(12), int64(12))

	if root_bool == 0 {
		t.Fatalf("The root verification failed")
	}

	//changing the proof---------------------------------------------------------------
	pi_root[0] = pi_root[1]
	root_bool_fail := VerProof(crs[:3], commit, pi_root, int64(12), int64(12), int64(12))

	if root_bool_fail == 1 {
		t.Fatalf("The root verification passed for an incorrect proof")
	}

	//Giving a non-member input---------------------------------------------------------
	x := big.NewInt(1111)
	Primes1 := make([]big.Int, len(set))
	G1 := crs[1]
	for i, u_dash := range set {
		Primes1[i] = hash2prime.Hprime(u_dash)
		if u_dash.Cmp(x) != 0 {

			G1.Exp(&G, &Primes1[i], &crs[0])
		}

	}

	W_ := G1
	root_witness_fail := []big.Int{e, r, W_}
	//Generation of proof of the root
	pi_root_fail := Prove(crs[:3], commit, root_witness_fail, int64(12), int64(12), int64(12))
	//verification of proof
	root_wronginput := VerProof(crs[:3], commit, pi_root_fail, int64(12), int64(12), int64(12))

	if root_wronginput == 1 {
		t.Fatalf("The root verification passed for a non-member")
	}

}
