package pkg

import "log"

func CheckAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckAndLogFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckAndLog(err error) {
	if err != nil {
		log.Println(err)
	}
}
