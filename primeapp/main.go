package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	intro()

	doneChan := make(chan bool) //create channel

	go readUserInput(os.Stdin, doneChan) //start go routine to read user input
	//block until doneChan gets the value
	<-doneChan

	close(doneChan)
	fmt.Println("Goodbye.")
}

func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)
	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChan <- true
			return
		}
		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	scanner.Scan()
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}
	numToCheck, err := strconv.Atoi(scanner.Text())

	if err != nil {
		return "Please enter a whole number!", false
	}

	_, msg := isPrime(numToCheck)
	return msg, false
}

func intro() {
	fmt.Println("Is it Prime?")
	fmt.Println("------------")
	fmt.Println("Enter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
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
