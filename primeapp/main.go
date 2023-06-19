package main

import "fmt"

func main() {
	n := 2
	_, msg := isPrime(n)
	fmt.Println(msg)
}

func isPrime(n int) (bool, string) {

	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not a prime number", n)
	}

	if n < 0 {
		return false, "negative numbers are not prime number"
	}

	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not a prime number because divisible by %d", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number", n)
}
