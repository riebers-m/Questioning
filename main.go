package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func setupLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, // Days,
		Compress:   true,
	})
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

const stateFile = "state.json"
const questionsFile = "questions.json"

func main() {
	setupLogger()
	log.Println("running Questioning...")

	if !IsWeekDay(time.Now()) {
		log.Println("no weekday have a nice weekend")
		return
	}

	if !fileExists(stateFile) {
		log.Println("state file does not exists, creating new one")
		if err := createEmptyJsonFile(stateFile); err != nil {
			log.Fatalf("creating state file failed: %v", err)
		}
	}
	state, err := loadState(stateFile)
	if err != nil {
		log.Fatalf("loading state file failed: %v", err)
	}

	questions, err := loadQuestions(questionsFile)
	if err != nil {
		log.Fatalf("reading questions.json failed: %v", err)
	}
	id, err := getRandomUnused(state, len(questions))

	if err != nil {
		log.Println("no remaining id, resetting state...")
		if err := os.Remove(stateFile); err != nil {
			log.Fatalf("could not reset state file: %v", err)
		}
		if err := createEmptyJsonFile(stateFile); err != nil {
			log.Fatalf("could not create empty state file: %v", err)
		}
		state = State{}
		id, err = getRandomUnused(state, 5)
	}

	log.Printf("next question choosed with id: %d", id)

	question, err := findQuestionById(questions, id)
	if err != nil {
		log.Fatal(err)
	}
	// notifer send question
	botMessage := fmt.Sprintf("Today question comes from category %s.\n%d: %s", question.Topic, question.Id, question.Question)
	if err := sendMessage(botMessage); err != nil {
		log.Fatalf("sending question failed: %v", err)
	}

	state.Used = append(state.Used, id)

	if err := saveState(stateFile, state); err != nil {
		log.Fatalf("could not save state: %v", err)
	}
}
