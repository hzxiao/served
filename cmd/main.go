package main

import (
	"fmt"
	"github.com/hzxiao/served"
	"os"
)

func main()  {
	err := served.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
