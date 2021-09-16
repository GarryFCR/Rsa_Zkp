package Hashtoprime

import (
	"math"
	"math/big"
)

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
