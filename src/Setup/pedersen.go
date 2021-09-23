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

// Pedersen_setup Generates prime and g and h
// prime for prime order group for committing the set member.
// g and h are random group element in that prime order group.

func Pedersen_setup(lambda, nu int64) (x, y, z *big.Int) {
	i := 0
	var prime_bound1, prime_bound2 *big.Int
	var prime *big.Int

	//upper := big.NewInt(int64(math.Pow(2, float64(nu))))
	//lower := big.NewInt(int64(math.Pow(2, float64(nu-1))))

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
	//fmt.Println(prime)
	g, _ := rand.Int(rand.Reader, prime)
	h, _ := rand.Int(rand.Reader, prime)
	//fmt.Println(prime, g, h)
	return prime, g, h

}

// Pedersen_commit returns the commitent c of set element u

func Pedersen_commit(ck []*big.Int, q, u *big.Int) (commitment, open *big.Int) {
	//var c *big.Int
	r, _ := rand.Int(rand.Reader, ck[0])
	//prime := ck[0]
	g, h := ck[1], ck[2]
	g.Exp(g, u, q)
	h.Exp(h, r, q)
	c := g.Mul(g, h)
	c.Mod(c, q)

	return c, r
}

// Pedersen_ver verifies the correctness of the commitment and returns 1 if its correct
// Else returns 0
func Pedersen_ver(ck []*big.Int, c, u, r *big.Int) int {

	prime, g, h := ck[0], ck[1], ck[2]
	g.Exp(g, u, prime)
	h.Exp(h, r, prime)
	c_ := g.Mul(g, h)
	c_.Mod(c_, prime)
	if c == c_ {
		return 1
	}
	return 0
}
