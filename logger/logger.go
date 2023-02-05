package logger

// private variables
var verbosetotal = false

func PrepareLogger(Verbosetotal bool) {
	verbosetotal = Verbosetotal
	Log("Logger prepared !")
}

func Log(message string) {
	println(message)
}

func Debug(message string) {
	if verbosetotal {
		println(message)
	}
}
