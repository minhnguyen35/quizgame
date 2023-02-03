package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

type Game struct {
	listQuiz      []Quiz
	correctAnswer int
}

type Quiz struct {
	problem string
	answer  string
}

type SafeIndex struct {
	mutex sync.Mutex
	index int
}

func (idx *SafeIndex) increaseIndex() {
	idx.mutex.Lock()
	idx.index++
	idx.mutex.Unlock()
}

func (idx *SafeIndex) getIndex() int {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()
	return idx.index
}

func convertCsvToQuize(data [][]string) []Quiz {
	results := make([]Quiz, len(data))
	for i, row := range data {
		results[i] = Quiz{row[0], row[1]}
	}
	return results
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "Please provide csv file with this format (question,answer)")
	flag.Parse()
	f, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file %s\n", *csvFileName))
	}

	defer f.Close()

	inputFromFile := csv.NewReader(f)
	data, err := inputFromFile.ReadAll()
	if err != nil {
		exit("Failed to parse csv file.")
	}
	quizList := convertCsvToQuize(data)

	quizGame := Game{quizList, 0}
	startGame(quizGame)

}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func startTimeout(i int, c chan int) {
	time.Sleep(30 * time.Second)
	i++
	c <- i
}
func startGame(game Game) {
	totalQuiz := len(game.listQuiz)
	safeIndex := SafeIndex{index: 0}
	i := 0
	for i < totalQuiz {
		fmt.Printf("Please answer %s ?\n", game.listQuiz[safeIndex.getIndex()].problem)
		var answer string

		fmt.Scanln(&answer)
		if answer == game.listQuiz[safeIndex.getIndex()].answer {
			game.correctAnswer++
			fmt.Println("Correct answer ", game.correctAnswer)
		}
		i++
	}
	fmt.Println("Correct answer ", game.correctAnswer)
	fmt.Println("Total Quiz ", totalQuiz)
}
