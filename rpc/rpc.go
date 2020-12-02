package rpc

var rpcs = map[uint8]map[string]func(string) ([]byte, error){} // rpcType -> {rpcMethodName -> actual func var}

func RegisterRPC(rpcType uint8, rpcMethods map[string]func(string) ([]byte, error)) {
	rpcs[rpcType] = rpcMethods
}

func GetRPCMethod(rpcType uint8, rpcMethodName string) func(string) ([]byte, error) {
	rpcMethods, ok := rpcs[rpcType]
	if !ok {
		return nil
	}
	actualMethod, ok := rpcMethods[rpcMethodName]
	if !ok {
		return nil
	}
	return actualMethod
}
