package errs

var stackSkip = 4
var isTrace = true

func SetTrace(trace bool) {
	isTrace = trace
}

func SetStackSkip(skip int) {
	stackSkip = skip
}
