package sum

func Sum(numbers []int, ch chan int) {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	ch <- sum
}
