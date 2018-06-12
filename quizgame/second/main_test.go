package main

import (
	"testing"
	"time"
)

func TestProcessData(t *testing.T) {
	var tests = []struct {
		want     int
		line     []string
		problems []problem
	}{
		{1, []string{"5+5", "10"}, make([]problem, 0)},
		{0, []string{}, make([]problem, 0)},
	}

	for _, test := range tests {
		processData(test.line, &test.problems)
		if test.want != len(test.problems) {
			t.Errorf("FAIL: processData = %d, want %d", len(test.problems), test.want)
		}
	}
}

func TestReadFile(t *testing.T) {
	var tests = []struct {
		filename string
		want     int
	}{
		{"testdata/problems.csv", 2},
		{"problem.csv", 0},
	}

	for _, test := range tests {
		problems, err := readFile(&test.filename)
		if err != nil && len(problems) != test.want {
			t.Errorf("FAIL: readFile = %d, want %d", len(problems), test.want)
		}
	}
}

func TestCountableAnswerCorrect(t *testing.T) {
	timer := time.NewTimer(time.Duration(5) * time.Second)

	answerChannel := make(chan string, 1)
	answerChannel <- "2"

	p := problem{question: "1+1", answer: "2"}
	correct := 0

	//c := make(chan time.Time, 1)
	//c <- time.Now()

	timeout := countableAnswer(timer, answerChannel, p, &correct)
	if timeout {
		t.Errorf("FAIL: countableAnswer = %t, want %t", timeout, false)
	}

	if correct != 1 {
		t.Errorf("FAIL: countableAnswer = %v, want %v", correct, 1)
	}
}
