package src

import (
	"fmt"
	"os"
)

func throw(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}
