package main

import (
	"io"
	"os"
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
	if string(out) != "-> " { //perform test //use string func cuz return as slice of byte
		t.Errorf("incorrect prompt: expected -> but got %s", string(out))
	}
}
