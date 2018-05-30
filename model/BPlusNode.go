package main

import (
	"math/big"
	"fmt"
)

func main() {
   n := new(big.Int)
   n.SetString("/core/genaro",61)
   var s string
   s = n.String()
   fmt.Print(s)
}
