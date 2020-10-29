package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	fileName = "problems.csv"
	shuffle  = flag.Bool("r", false, "for shuffling the quiz")
	timer    = flag.Int64("t", 30, "stoppage time for quiz")
)

type result struct {
	quizzes     []quiz
	userAnswer string
	score      int
	resultChan chan string
}

type quiz struct {
	question string
	answer   string
}

func reader(f io.Reader) []quiz {
	r := csv.NewReader(f)

	var quizzes []quiz

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		quizzes = append(quizzes, quiz{
			question: record[0],
			answer:   record[1],
		})
	}
	return quizzes
}

func (r *result) ask() {
	if *shuffle == true {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(r.quizzes), func(i, j int) {
			r.quizzes[i], r.quizzes[j] = r.quizzes[j], r.quizzes[i]
		})
	}

	for _, quiz := range r.quizzes {
		fmt.Printf("%s = ", quiz.question)
		fmt.Scan(&r.userAnswer)

		if r.userAnswer == quiz.answer {
			r.score++
		}

		r.resultChan <- r.userAnswer
	}
	close(r.resultChan)
}

func main() {
	fmt.Print(`Hit "Enter" to start!`)
	fmt.Scanln()

	flag.Parse()
	timeout := time.After(time.Second * time.Duration(*timer))


	var res result
	var wg sync.WaitGroup
	res.resultChan = make(chan string)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("error opening file")
	}
	defer f.Close()

	quizzes := reader(f)

	res.quizzes = quizzes

	wg.Add(1)
	go func() {
		defer wg.Done()

	Loop:
		for {
			select {
			case _, ok := <-res.resultChan:
				if !ok {
					break Loop
				}

			case <-timeout:
				fmt.Println("\nTimeout!")
				break Loop
			}
		}
	}()

	go func() {
		res.ask()
	}()

	wg.Wait()
	fmt.Printf("Total score = %d/%v\n", res.score, len(res.quizzes))
}
