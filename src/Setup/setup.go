package Setup

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func Pedersen_setup(lambda, nu int64) (x, y, z *big.Int) {
	i := 0
	var prime_bound1, prime_bound2 big.Int
	var prime *big.Int

	upper := big.NewInt(int64(math.Pow(2, float64(nu))))
	lower := big.NewInt(int64(math.Pow(2, float64(nu-1))))

	for i == 0 {
		prime, _ = rand.Prime(rand.Reader, int(lambda))
		fmt.Println(prime)
		prime_bound1.Sub(upper, prime)
		prime_bound2.Sub(prime, lower)
		if prime_bound1.Sign() == 1 && prime_bound2.Sign() == 1 {
			break
		}
	}
	//fmt.Println(prime)
	g, _ := rand.Int(rand.Reader, prime)
	h, _ := rand.Int(rand.Reader, prime)
	//fmt.Println(prime, g, h)
	return prime, g, h

}

func Pedersen_commit(ck []*big.Int, u *big.Int) (commitment, open *big.Int) {
	//var c *big.Int
	r, _ := rand.Int(rand.Reader, ck[0])
	prime := ck[0]
	g, h := ck[1], ck[2]
	g.Exp(g, u, prime)
	h.Exp(h, r, prime)
	c := g.Mul(g, h)
	c.Mod(c, prime)

	return c, r
}

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
