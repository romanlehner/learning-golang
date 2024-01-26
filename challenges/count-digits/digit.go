package digit

func counter(number int) int {
	// devide the number by 10 until %10 is zero to identify all digits
	count := 0
	for i := 1; number%10 != 0; i++ {
		number /= 10
		count = i
	}

	return count
}