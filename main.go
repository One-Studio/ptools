package main

import (
	"fmt"
	tool "github.com/One-Studio/ptools/pkg"
	"log"
)

func main() {

	tool.Test()

	t := ""
	t, err := tool.GetBinaryPath("wt")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(t)
	}

	_ = tool.ExecRealtime("ping baidu.com")
}
