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

func Setup_Set(lambda int, neu int) ( n, f *big.Int)  {
  pk, _ := rsa.GenerateKey(rand.Reader , lambda)
  var F *big.Int

  for  {
    F, _ = rand.Int(rand.Reader, pk.PublicKey.N)
    if F!=pk.Primes[0] && F!=pk.Primes[1] && F!= big.NewInt(1) {
        break
      }
}

G := F.Exp(F, F, pk.PublicKey.N)
N := pk.PublicKey.N

return N, G
}


func fu(x *big.Int) *big.Int {
  temp1 :=big.NewInt(0)
  temp1.Add(x,big.NewInt(2))
  temp1.Mul(temp1,big.NewInt(2) )
  temp2 :=big.NewInt(0)
  temp2.Add(x,big.NewInt(1))
  temp2.Mul(temp2, temp2)

//2*(x+2)*
  fx := math.Log2(temp2)

  return fx

}



//Function to map the set element to prime

func Hprime(u *big.Int) *big.Int {


  for j:= add(fu(u),1) ; j< fu(add(u,1) ; j++ {

    x := fu(u)+j
    if x.ProbablyPrime(10) == true {
      break
    }
    return x
  }


}

//Set Commitment function

func Commit(N, G, U[] *big.Int) (c, o *big.Int)  {

  var P [] *big.Int
  var Acc *big.Int
  for i :=0; i< len(U); i++ {

    u := U[i]
    P[i] = Hprime(u)
    Acc *= G.Exp(G, P[i], pk.PublicKey.N)
    Acc.Exp(Acc, 1, pk.PublicKey.N)
   }
   return Acc, nil

}
