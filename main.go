package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Welcome to our quiz game!")
	fmt.Println("Please enter your name:")
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')

	if err != nil {
		panic("HEEEELP I'M DEAD")
	}

	g.Name = name
	fmt.Printf("Let's begin our quiz %s", g.Name)
}

func (g *GameState) ProcessCSV() {
	file, err := os.Open("quizgo.csv")

	if err != nil {
		panic("Failed to open CSV file")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Failed to read CSV file")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := strconv.Atoi(strings.TrimSpace(record[5]))
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}
			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) AskQuestion() {
	for index, q := range g.Questions {
		fmt.Printf("\033[33m %d. %s \033[0m\n", index+1, q.Text)
		for index, option := range q.Options {
			fmt.Printf("[%d] - %s\n", index+1, option)
		}
		fmt.Println("Please enter the number of your answer:")
		for {
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')

			answer, err := validateAnswer(input)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if answer == q.Answer {
				fmt.Println("Correct!")
				g.Points += 10
			} else {
				fmt.Println("Wrong!")
			}
			break
		}
	}
}

func main() {
	game := &GameState{}
	go game.ProcessCSV()

	game.Init()
	game.AskQuestion()

	fmt.Printf("Game over! Your final score is: %d\n", game.Points)
}

func validateAnswer(i string) (int, error) {
	a, err := strconv.Atoi(strings.TrimSpace(i))
	if err != nil {
		return 0, errors.New("Invalid input. Please enter a number.")
	}
	return a, nil
}
