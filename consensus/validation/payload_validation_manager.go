package validation

var validations = map[uint8]map[string]func([]byte) error{} // rpcType -> {rpcMethodName -> actual func var}

func RegisterValidation(typ uint8, methods map[string]func([]byte) error) {
	validations[typ] = methods
}

func GetValidationMethod(typ uint8, methodName string) func([]byte) error {
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
