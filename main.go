package main

import (
	"fmt"
	tool "github.com/One-Studio/ptools/pkg"
)

func main() {

	tool.Test()
	path, err := tool.GetBinaryPath("wt")
	fmt.Println(path)
	fmt.Println(err)

	path, err = tool.GetBinaryPath("temp/ffmpeg")
	fmt.Println(path)
	fmt.Println(err)

}
