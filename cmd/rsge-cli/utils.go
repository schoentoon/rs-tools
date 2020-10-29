package main

import (
	"fmt"
	"io"
	"math"
)

func prettyPrintPrice(out io.Writer, price int) {
	if price == math.MaxInt32 {
		fmt.Fprintf(out, "Max cash stack, actual price likely higher\n")
	} else if price >= 1000000000 { // 1b
		fmt.Fprintf(out, "%d.%03d.%03d.%03d (%d.%02db)\n",
			price/1000000000,
			(price/1000000)%1000,
			(price/1000)%1000,
			price%1000,
			price/1000000000,
			(price/10000000)%100)
	} else if price >= 1000000 { // 1m
		fmt.Fprintf(out, "%d.%03d.%03d (%d.%02dm)\n",
			price/1000000,
			(price/1000)%1000,
			price%1000,
			price/1000000,
			(price%1000000)/10000)
	} else if price >= 1000 { // 1k
		fmt.Fprintf(out, "%d.%03d (%d.%1dk)\n",
			price/1000,
			price%1000,
			price/1000,
			(price%1000)/100)
	} else {
		fmt.Fprintf(out, "%d\n", price)
	}
}
