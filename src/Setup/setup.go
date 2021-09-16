package Setup

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math"
	"math/big"
)

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

	G := F.Exp(F, F, pk.PublicKey.N)

	return N, G
}

func fu(x *big.Int) *big.Int {

	temp1 := big.NewInt(2)
	temp1.Add(x, temp1)
	two := big.NewInt(2)
	temp1.Mul(temp1, two)

	temp2 := big.NewInt(1)
	temp2.Add(x, temp2)

	bit := temp2.BitLen()
	one := big.NewFloat(1)

	divisor := new(big.Float).SetMantExp(one, bit-1)

	f := new(big.Float).SetInt(temp2)

	z := new(big.Float).Quo(f, divisor)
	w, _ := z.Float64()

	y := math.Pow((math.Log2(w) + float64(bit-1)), 2)

	temp1.Mul(temp1, big.NewInt(int64(y)))

	return temp1

}

//Function to map the set element to prime

func Hprime(u *big.Int) *big.Int {

	Huj := fu(u)
	j := fu(u)

	for {

		temp := Huj
		temp.Add(temp, j)
		if temp.ProbablyPrime(10) {

			return temp

		}

		j.Add(j, big.NewInt(1))
	}
	//return big.NewInt(-1)

}

//Set Commitment function

func Set_commit(ck, U []*big.Int) (c, o *big.Int) {

	P := make([]*big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = Hprime(u)
		G.Exp(G, P[i], N)
	}

	fmt.Println(P)
	return G, nil

}

func Set_ver(ck, U []*big.Int, Acc *big.Int) int {

	P := make([]*big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = Hprime(u)
		G.Exp(G, P[i], N)
	}

	if Acc == G {
		return 1
	}
	return 0

}
