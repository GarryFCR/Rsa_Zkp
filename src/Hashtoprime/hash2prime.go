/*
	RSA Accumulator has a limitation that it deals with only set of primes
	To overcome this limitation we map arbitrary values to primes in collision-resistant manner

	We are using deterministic map H(u,j) = f(u)+j for Hash to ptime where f(u)= 2*(u+2)*log2(u+1)^2
	the prime is ensured using Cramer's conjecture
*/

package Hashtoprime

import (
	"math"
	"math/big"
)

func Fu(x big.Int) big.Int {

	u := x
	temp1 := new(big.Int).Add(&u, big.NewInt(2))
	temp1.Mul(temp1, big.NewInt(2))

	temp2 := new(big.Int).Add(&u, big.NewInt(1))

	bit := temp2.BitLen()
	one := big.NewFloat(1)

	divisor := new(big.Float).SetMantExp(one, bit-1)

	f := new(big.Float).SetInt(temp2)

	z := new(big.Float).Quo(f, divisor)
	w, _ := z.Float64()

	y := math.Pow((math.Log2(w) + float64(bit-1)), 2)

	temp1.Mul(temp1, big.NewInt(int64(y)))

	return *temp1

}

//Hprime returns the prime which is mapped to set element u in collision resistant manner

func Hprime(u big.Int) big.Int {

	Huj := Fu(u)
	j := Fu(u)
	var temp big.Int

	for {

		temp = Huj
		prime := new(big.Int).Add(&temp, &j)

		//temp.Add(&temp, &j)
		if prime.ProbablyPrime(10) {

			return *prime
		}

		j.Add(&j, big.NewInt(1))
	}
	//return big.NewInt(-1)

}
