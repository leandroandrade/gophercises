package main

import (
	"flag"
	"os"
	"strings"
	"fmt"
	"log"
	"time"
	"math/rand"
	"encoding/csv"
	"io"
)

type problem struct {
	question string
	answer   string
}

func readFile(csvFilename *string) ([]problem, error) {
	file, err := os.Open(*csvFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open the CSV file: %s", *csvFilename)
	}

	problems := make([]problem, 0)

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		processData(record, &problems)
	}
	return problems, nil
}

func processData(line []string, problems *[]problem) {
	if len(line) != 0 {
		*problems = append(*problems,
			problem{question: line[0], answer: strings.TrimSpace(line[1])})
	}
}

func answer(answerChannel chan string) {
	var answer string
	fmt.Scanf("%v\n", &answer)

	answerChannel <- answer
}

func countableAnswer(timer *time.Timer, answerChannel chan string, p problem, correct *int) bool {
	select {
	case <-timer.C:
		fmt.Println()

		return true

	case answer := <-answerChannel:
		if answer == p.answer {
			*correct++
		}
		return false
	}
}

func shuffleQuestions(problems []problem) {
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a CSV file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz seconds")
	random := flag.Bool("random", false, "shuffle the quiz order each time it is run")

	flag.Parse()

	problems, err := readFile(csvFilename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if *random {
		shuffleQuestions(problems)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	answerChannel := make(chan string)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		go answer(answerChannel)

		if timeout := countableAnswer(timer, answerChannel, p, &correct); timeout {
			break
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}
