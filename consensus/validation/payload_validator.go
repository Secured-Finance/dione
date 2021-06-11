package validation

import "github.com/Secured-Finance/dione/types"

var validations = map[uint8]map[string]func(*types.DioneTask) error{} // rpcType -> {rpcMethodName -> actual func var}

func RegisterValidation(typ uint8, methods map[string]func(*types.DioneTask) error) {
	validations[typ] = methods
}

func GetValidationMethod(typ uint8, methodName string) func(*types.DioneTask) error {
	rpcMethods, ok := validations[typ]
	if !ok {
		return nil
	}
	actualMethod, ok := rpcMethods[methodName]
	if !ok {
		return nil
	}
	return actualMethod
}
