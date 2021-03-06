package main

import (
	"fmt"
	"os"

	wasmer "github.com/sammyne/mastering-wasm/wavm"
	"github.com/sammyne/mastering-wasm/wavm/cmd/wavm/tools"

	flag "github.com/spf13/pflag"
)

var (
	dump  bool
	check bool
)

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
			panicf("fail to dump: %v", err)
		}
		return
	} else if check {
		if err := tools.Check(module); err != nil {
			panicf("invalid module: %v", err)
		}
		fmt.Println("module is good :)")
		return
	}

	if err := tools.InstantiateAndExecMainFunc(module); err != nil {
		panicf("instantiate and exec main func: %v", err)
	}

	//if err := vm.Run(module); err != nil {
	//	panic(fmt.Sprintf("fail to run: %v", err))
	//}
}

func init() {
	flag.BoolVarP(&dump, "dump", "d", false, "")
	flag.BoolVarP(&check, "check", "c", false, "check wasm file")
}

func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
