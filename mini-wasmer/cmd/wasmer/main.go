package main

import (
	"fmt"
	"os"

	wasmer "github.com/sammyne/mastering-wasm/mini-wasmer"
	"github.com/sammyne/mastering-wasm/mini-wasmer/cmd/wasmer/tools"
	"github.com/sammyne/mastering-wasm/mini-wasmer/vm"
	flag "github.com/spf13/pflag"
)

var dump bool

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.PrintDefaults()
		fmt.Println("only 1 positional argument is allowed")
		os.Exit(-1)
	}

	module, err := wasmer.DecodeModuleFromFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	if dump {
		if err := tools.Dump(module); err != nil {
			panic(fmt.Sprintf("fail to dump: %v", err))
		}
		return
	}

	if err := vm.Run(module); err != nil {
		panic(fmt.Sprintf("fail to run: %v", err))
	}
}

func init() {
	flag.BoolVarP(&dump, "dump", "d", false, "")
}
