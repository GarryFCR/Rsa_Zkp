package HashEq

import (
	"crypto/rand"
	"math/big"
	"testing"

	hash2prime "../Hashtoprime"
	setup "../Setup"
)

func TestHashEq(t *testing.T) {

	u := big.NewInt(11111)
	e := hash2prime.Hprime(*u)
	fu := hash2prime.Fu(*u)
	j := new(big.Int).Sub(&e, &fu)

	prime, g, h := setup.Pedersen_setup(12, 12)
	ce, rq := setup.Pedersen_commit([]big.Int{prime, g, h}, prime, e)
	cu, ru := setup.Pedersen_commit([]big.Int{prime, g, h}, prime, *u)

	hash_commitment := []big.Int{ce, cu}
	hash_witness := []big.Int{e, *u, rq, ru, *j}

	pi_hash := Prove([]big.Int{prime, g, h}, hash_commitment, hash_witness, int64(512), int64(512), int64(512), int64(512))
	hash_bool := VerProof([]big.Int{prime, g, h}, hash_commitment, pi_hash)
	//Check if it fails for a correct proof
	if hash_bool == 0 {
		t.Fatalf("HashEq fails on proof generated")
	}

	//For a random j
	j_random, _ := rand.Int(rand.Reader, &prime)
	pi_hash1 := Prove([]big.Int{prime, g, h}, hash_commitment, []big.Int{e, *u, rq, ru, *j_random}, int64(512), int64(512), int64(512), int64(512))
	hash_bool1 := VerProof([]big.Int{prime, g, h}, hash_commitment, pi_hash1)
	if hash_bool1 == 1 {
		t.Fatalf("HashEq worked for a random j")
	}

}
