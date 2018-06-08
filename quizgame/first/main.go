package main

import (
	"flag"
	"os"
	"log"
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a CSV file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz seconds")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatalf("failed to open the CSV file: %s", *csvFilename)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("failed to parse the provided CSV file")
	}

	problems := parseCSVContent(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

loop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break loop
		case answer := <-answerChannel:
			if answer == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))

}
func parseCSVContent(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{question: line[0], answer: strings.TrimSpace(line[1])}
	}
	return problems

}
