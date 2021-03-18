package internal

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
