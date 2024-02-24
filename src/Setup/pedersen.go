/*
	Pederson commitment is used to commit to set elements in some large prime order group.
	we have used Z/pZ group for our implementation, where p is prime of nu bits.
*/

package Setup

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func Pedersen_setup(lambda, nu int64) (x, y, z big.Int) {

	i := 0
	var prime_bound1, prime_bound2 *big.Int
	var prime *big.Int

	upper := big.NewInt(2)
	upper.Exp(upper, big.NewInt(nu), nil)
	lower := big.NewInt(2)
	lower.Exp(lower, big.NewInt(nu-1), nil)

	for i == 0 {
		prime, _ = rand.Prime(rand.Reader, int(lambda))

		prime_bound1 = new(big.Int).Sub(upper, prime)
		prime_bound2 = new(big.Int).Sub(prime, lower)
		if prime_bound1.Sign() == 1 && prime_bound2.Sign() == 1 {
			break
		}
		fmt.Println(prime, prime_bound1.Sign(), prime_bound2.Sign())
	}

	g, _ := rand.Int(rand.Reader, prime)
	h, _ := rand.Int(rand.Reader, prime)
	//fmt.Println(prime, g, h)
	return *prime, *g, *h

}

func Pedersen_commit(ck []big.Int, q, u big.Int) (commitment, open big.Int) {

	r, _ := rand.Int(rand.Reader, &q)
	var g, h = ck[1], ck[2]
	g.Exp(&g, &u, &q)
	h.Exp(&h, r, &q)

	c := new(big.Int).Mul(&g, &h)
	c.Mod(c, &q)

	return *c, *r
}

func Pedersen_ver(ck []big.Int, c, u, r big.Int) int {

	prime, g, h := ck[0], ck[1], ck[2]
	g.Exp(&g, &u, &prime)
	h.Exp(&h, &r, &prime)
	c_ := new(big.Int).Mul(&g, &h)
	c_.Mod(c_, &prime)

	if c_.Cmp(&c) == 0 {
		return 1
	}
	return 0
}
