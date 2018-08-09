package utils

// Check panic if  an error have occured
func Check(e error) {
	if e != nil {
		println("failed:", e)
		panic(e)
	}
}
