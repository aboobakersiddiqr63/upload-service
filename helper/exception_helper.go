package helper

func HandleException(err error, funcName string) {
	if err != nil {
		Log.Errorf("Error in func: %v, errMsg: %v", funcName, err)
	}
}

func DbExceptionHandler(err error, method string) {
	if err != nil {
		Log.Errorf("Error in func: %v, errMsg: %v", err, method)
	}
}
