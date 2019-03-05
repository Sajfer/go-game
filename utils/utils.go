package utils

// Check panic if  an error have occured
func Check(e error, message string) {
	if e != nil {
		println("%s: %s",message, e)
		panic(e)
	}
}
