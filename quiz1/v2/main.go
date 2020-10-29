package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	csvFile = flag.String("f", "problems.csv", "file name that contains the quiz")
	shuffle = flag.Bool("r", false, "for shuffling the quiz")
	timer   = flag.Int64("t", 30, "stoppage time for the quiz in seconds")
)

type result struct {
	quizzes    []problem
	score      int
	timeout    <-chan time.Time
	answerChan chan string
}

type problem struct {
	question string
	answer   string
}

func reader(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))

	for i, record := range records {
		problems[i] = problem{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}
	return problems
}

func (r *result) ask() {
	if *shuffle == true {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(r.quizzes), func(i, j int) {
			r.quizzes[i], r.quizzes[j] = r.quizzes[j], r.quizzes[i]
		})
	}

	for i, problem := range r.quizzes {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			r.answerChan <- answer
		}()

		select {
		case answer := <-r.answerChan:
			if answer == problem.answer {
				r.score++
			}

		case <-r.timeout:
			fmt.Println("\nTimeout!")
			return
		}
	}
	close(r.answerChan)
}

func main() {
	var res result
	res.answerChan = make(chan string)

	flag.Parse()

	records, err := reader(*csvFile)
	if err != nil {
		log.Fatalf("failed to open csv file %s\n", *csvFile)
	}

	problems := parseRecords(records)

	res.quizzes = problems

	fmt.Print(`Hit "Enter" to start!`)
	fmt.Scanln()

	res.timeout = time.After(time.Duration(*timer) * time.Second)

	res.ask()

	fmt.Printf("You scored %d out of %v.\n", res.score, len(res.quizzes))
}
