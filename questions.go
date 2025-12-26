package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Question struct {
	Id       int    `json:"id"`
	Topic    string `json:"topic"`
	Question string `json:"question"`
}

type State struct {
	Used []int `json:"used"`
}

func loadQuestions(pathToFile string) ([]Question, error) {
	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	var questions []Question
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func findQuestionById(questions []Question, id int) (*Question, error) {
	for i := range questions {
		if questions[i].Id == id {
			return &questions[i], nil
		}
	}
	return nil, fmt.Errorf("id %d not found in questions", id)
}

func loadState(pathToState string) (State, error) {
	data, err := os.ReadFile(pathToState)
	if err != nil {
		return State{}, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, err
	}
	return state, nil
}

func saveState(pathToState string, state State) error {
	if err := renameFile(pathToState, pathToState+".tmp"); err != nil {
		return err
	}
	defer os.Remove(pathToState + ".tmp")

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pathToState, data, 0644); err != nil {
		return err
	}
	return nil
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomUnused(state State, questions int) (int, error) {
	idsLeft := removeValuesFromRange(1, questions, state.Used)
	if len(idsLeft) == 0 {
		return -1, fmt.Errorf("no ids left, reset state")
	}

	return idsLeft[rng.Intn(len(idsLeft))], nil
}

func removeValuesFromRange(start, end int, remove []int) []int {
	toRemove := make(map[int]struct{}, len(remove))
	result := make([]int, 0, end-start+1)

	if len(remove) >= end {
		return result
	}

	for _, v := range remove {
		toRemove[v] = struct{}{}
	}

	for i := start; i <= end; i++ {
		if _, found := toRemove[i]; !found {
			result = append(result, i)
		}
	}

	return result
}

func randomInt(min, max int) int {
	return rng.Intn(max-min+1) + min
}

func renameFile(pathToFile, tmpName string) error {
	return os.Rename(pathToFile, tmpName)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false // other errors (permission, etc.)
}

func createEmptyJsonFile(path string) error {
	err := os.WriteFile(path, []byte("{}"), 0644)
	if err != nil {
		return err
	}
	return nil
}
