package main

import (
	"fmt"
	"github.com/PaluMacil/merkle/merkle"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalln("please run this command with filenames separated by commas")
	}
	filenames := strings.Split(args[1], ",")
	layer, err := merkle.From(filenames...)
	if err != nil {
		log.Fatalf("starting merkle tree: %s\n", err)
	}
	root := layer.Root()
	fmt.Printf("the root hash for this input is %x", root)
}
