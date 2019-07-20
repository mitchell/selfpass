package migration

import (
	"fmt"
	"os"
	"runtime"
)

func Check(err error) {
	if err != nil {
		_, _, line, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("%v: %s\n", line, err)
			os.Exit(1)
		}

		fmt.Println(err)
	}
}
