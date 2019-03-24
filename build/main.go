
package main

import (
	"fmt"
	"os"
)

var funcNames = []string{"totalAmount","bankName","balances","withdraw","owner","deposit"}

func Usage() {
	fmt.Printf("%s 1 -- deploy\n", os.Args[0])
	num := 2
	for _, v := range funcNames {
		fmt.Printf("%s %d -- %s\n", os.Args[0], num, v)
		num++
	}
}


func main() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(0)
	}
	if os.Args[1] == "1" {
		CallDeploy()
	} else if os.Args[1] == "2" {
		CallTotalAmount()
	}  else if os.Args[1] == "3" {
		CallBankName()
	}  else if os.Args[1] == "4" {
		CallBalances()
	}  else if os.Args[1] == "5" {
		CallWithdraw()
	}  else if os.Args[1] == "6" {
		CallOwner()
	}  else if os.Args[1] == "7" {
		CallDeposit()
	} 
}

