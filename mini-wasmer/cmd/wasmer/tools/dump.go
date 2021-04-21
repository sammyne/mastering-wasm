package tools

import (
	"fmt"

	wasmer "github.com/sammyne/mastering-wasm/mini-wasmer"
	"github.com/sammyne/mastering-wasm/mini-wasmer/tools"
	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

type ImportCounts struct {
	Func   int
	Table  int
	Memory int
	Global int
}

func Dump(m *wasmer.Module) error {
	fmt.Printf("Version: 0x%02x\n", m.Version)

	dumpTypes(m.Types)

	importCounts, err := dumpImports(m.Imports)
	if err != nil {
		return fmt.Errorf("bad imports: %w", err)
	}

	dumpFunctions(m.Functions, importCounts.Func)
	dumpTables(m.Tables, importCounts.Table)
	dumpMemories(m.Memories, importCounts.Memory)
	dumpGlobals(m.Globals, importCounts.Global)
	dumpExports(m.Exports)
	dumpStart(m.Start)
	dumpElements(m.Elements)
	dumpCodes(m.Codes, m.Types, importCounts.Func)
	dumpData(m.Data)
	dumpCustoms(m.Customs)

	return nil
}

func dumpCodes(codes []types.Code, types_ []types.FuncType, offset int) {
	fmt.Printf("Code[%d]:\n", len(codes))
	for i, code := range codes {
		fmt.Printf("  func[%d]: locals=[", offset+i) // TODO
		if len(code.Locals) > 0 {
			for i, locals := range code.Locals {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s x %d", types.StringifyValueType(locals.Type), locals.N)
			}
		}
		fmt.Println("]")
		dumpExpr("    ", types_, code.Expr)
	}
}

func dumpCustoms(customs []types.Custom) {
	fmt.Printf("Custom[%d]:\n", len(customs))
	for i, v := range customs {
		fmt.Printf("  custom[%d]: name=%s\n", i, v.Name) // TODO
	}
}

func dumpData(data []types.Data) {
	fmt.Printf("Data[%d]:\n", len(data))
	for i, v := range data {
		fmt.Printf("  data[%d]: mem=%d\n", i, v.MemoryIdx) // TODO
	}
}

func dumpElements(elements []types.Element) {
	fmt.Printf("Element[%d]:\n", len(elements))
	for i, elem := range elements {
		fmt.Printf("  elem[%d]: table=%d\n", i, elem.TableIdx) // TODO
	}
}

func dumpExports(exports []types.Export) {
	fmt.Printf("Export[%d]:\n", len(exports))
	for _, v := range exports {
		switch v.Description.Tag {
		case types.PortTagFunc:
			fmt.Printf("  func[%d]: name=%s\n", int(v.Description.Idx), v.Name)
		case types.PortTagTable:
			fmt.Printf("  table[%d]: name=%s\n", int(v.Description.Idx), v.Name)
		case types.PortTagMemory:
			fmt.Printf("  memory[%d]: name=%s\n", int(v.Description.Idx), v.Name)
		case types.PortTagGlobal:
			fmt.Printf("  global[%d]: name=%s\n", int(v.Description.Idx), v.Name)
		}
	}
}

func dumpExpr(indent string, types_ []types.FuncType, expr types.Expr) {
	for _, v := range expr {
		switch v.Opcode {
		case types.OpcodeBlock, types.OpcodeLoop:
			block := v.Args.(*types.Block)
			blockType := tools.ParseBlockSig(block.BlockType, types_)
			fmt.Printf("%s%s %s\n", indent, v.GetOpname(), blockType)
			dumpExpr(indent+"  ", types_, block.Instructions)
			fmt.Printf("%send\n", indent)
		case types.OpcodeIf:
			blockIf := v.Args.(*types.BlockIf)
			blockType := tools.ParseBlockSig(blockIf.BlockType, types_)
			fmt.Printf("%sif %s\n", indent, blockType)
			dumpExpr(indent+"  ", types_, blockIf.Instructions1)
			fmt.Printf("%selse\n", indent)
			dumpExpr(indent+"  ", types_, blockIf.Instructions2)
			fmt.Printf("%send\n", indent)
		default:
			if v.Args != nil {
				fmt.Printf("%s%s %v\n", indent, v.GetOpname(), v.Args)
			}
		}
	}
}

func dumpFunctions(funcs []types.TypeIdx, offset int) {
	fmt.Printf("Function[%d]:\n", len(funcs))
	for i, v := range funcs {
		fmt.Printf("  func[%d]: sig=%d\n", offset+i, v)
	}
}

func dumpGlobals(globals []types.Global, offset int) {
	fmt.Printf("Global[%d]:\n", len(globals))
	for i, g := range globals {
		fmt.Printf("  global[%d]: %s\n", offset+i, g.Type)
	}
}

func dumpImports(imports []types.Import) (*ImportCounts, error) {
	fmt.Printf("Import[%d]:\n", len(imports))

	var out ImportCounts
	for _, v := range imports {
		switch v.Description.Tag {
		case types.PortTagFunc:
			fmt.Printf("  func[%d]: %s.%s, sig=%d\n", out.Func, v.Module, v.Name, v.Description.Func)
			out.Func++
		case types.PortTagTable:
			fmt.Printf("  table[%d]: %s.%s, %s\n",
				out.Table, v.Module, v.Name, v.Description.Table.Limits)
			out.Table++
		case types.PortTagMemory:
			fmt.Printf("  memory[%d]: %s.%s, %s\n", out.Memory, v.Module, v.Name, v.Description.Memory)
			out.Memory++
		case types.PortTagGlobal:
			fmt.Printf("  global[%d]: %s.%s, %s\n", out.Global, v.Module, v.Name, v.Description.Global)
			out.Global++
		default:
			return nil, fmt.Errorf("unknown tag: %02x", v.Description.Tag)
		}
	}

	return &out, nil
}

func dumpMemories(memories []types.Memory, offset int) {
	fmt.Printf("Memory[%d]:\n", len(memories))
	for i, limits := range memories {
		fmt.Printf("  memory[%d]: %s\n", offset+i, limits)
	}
}

func dumpStart(start *uint32) {
	fmt.Printf("Start:\n")
	if start != nil {
		fmt.Printf("  func=%d\n", start)
	}
}

func dumpTables(tables []types.Table, offset int) {
	fmt.Printf("Table[%d]:\n", len(tables))
	for i, t := range tables {
		fmt.Printf("  table[%d]: %s\n", offset+i, t.Limits)
	}
}

func dumpTypes(types []types.FuncType) {
	fmt.Printf("Type[%d]:\n", len(types))
	for i, ft := range types {
		fmt.Printf("  type[%d]: %s\n", i, ft)
	}
}
