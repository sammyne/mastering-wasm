package native

import (
	"strings"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

func parseNameAndSig(nameAndSig string) (string, types.FuncType) {
	idxOfLPar := strings.IndexByte(nameAndSig, '(')
	name := nameAndSig[:idxOfLPar]
	sig := nameAndSig[idxOfLPar:]
	return name, parseSig(sig)
}

func parseSig(sig string) types.FuncType {
	paramsAndResults := strings.SplitN(sig, "->", 2)
	return types.FuncType{
		ParamTypes:  parseValTypes(paramsAndResults[0]),
		ResultTypes: parseValTypes(paramsAndResults[1]),
	}
}

func parseValTypes(list string) []types.ValueType {
	list = strings.TrimSpace(list)
	list = list[1 : len(list)-1] // remove ()

	var valTypes []types.ValueType
	for _, t := range strings.Split(list, ",") {
		switch strings.TrimSpace(t) {
		case "i32":
			valTypes = append(valTypes, types.ValueTypeI32)
		case "i64":
			valTypes = append(valTypes, types.ValueTypeI64)
		case "f32":
			valTypes = append(valTypes, types.ValueTypeF32)
		case "f64":
			valTypes = append(valTypes, types.ValueTypeF64)
		}
	}
	return valTypes
}
