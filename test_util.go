package gomu

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
