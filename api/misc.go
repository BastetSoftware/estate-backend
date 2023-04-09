package api

func UnknownFPlug(_ []byte) (interface{}, error) {
	return Response{Code: ENoFun}, nil
}

func HandleFPing(_ []byte) (interface{}, error) {
	return Response{Code: 0}, nil
}
