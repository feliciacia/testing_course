package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number"},
		{"not prime", 8, false, "8 is not a prime number because divisible by 2"},
		{"zero", 0, false, "0 is not a prime number"},
		{"one", 1, false, "1 is not a prime number"},
		{"negative number", -11, false, "negative numbers are not prime number"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)

		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}
		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_prompt(t *testing.T) {
	oldout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	prompt()
	_ = w.Close()             // close writer
	os.Stdout = oldout        //reset
	out, _ := io.ReadAll(r)   //read the output of prompt func from read pipe
	if string(out) != "-> " { //perform test //convert to string func cuz ReadAll return as slice of byte
		t.Errorf("incorrect prompt: expected -> but got %s", string(out))
	}
}

func Test_intro(t *testing.T) {
	oldout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	intro()
	_ = w.Close()                                               // close writer
	os.Stdout = oldout                                          //reset
	out, _ := io.ReadAll(r)                                     //read the output of prompt func from read pipe
	if !strings.Contains(string(out), "Enter a whole number") { //strings.Contains() contains substr of string means true but sensitive with upper/lower case
		t.Errorf("intro text not correct; got %s", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: "Please enter a whole number!"},
		{name: "zero", input: "0", expected: "0 is not a prime number"},
		{name: "one", input: "1", expected: "1 is not a prime number"},
		{name: "two", input: "2", expected: "2 is a prime number"},
		{name: "three", input: "3", expected: "3 is a prime number"},
		{name: "negative", input: "-1", expected: "negative numbers are not prime number"},
		{name: "typed", input: "three", expected: "Please enter a whole number!"},
		{name: "decimal", input: "1.1", expected: "Please enter a whole number!"},
		{name: "quit", input: "q", expected: ""},
		{name: "QUIT", input: "Q", expected: ""},
		{name: "greek", input: "αεφπ", expected: "Please enter a whole number!"},
	}
	for _, e := range tests {
		input := strings.NewReader(e.input) //convert to string so can be read
		reader := bufio.NewScanner(input)   //input from command
		res, _ := checkNumbers(reader)

		if !strings.EqualFold(res, e.expected) { //equalfold= compared s1 and s2 must all equal but no sensitive with upper/lowercase
			t.Errorf("%s: expected %s, but got %s", e.name, e.expected, res)
		}
	}
}
func Test_readUserInput(t *testing.T) {
	//make channel and need an instance of an io.Reader
	doneChan := make(chan bool)
	var stdin bytes.Buffer //create a reference to bytes.Buffer //bytes.Buffer buffer of read and write funcion
	stdin.Write([]byte("1\nq\n"))
	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}
