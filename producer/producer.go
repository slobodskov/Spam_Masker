package producer

import (
	"bufio"
	"os"
)

type Producer struct {
	inputFileName string
	fileText      []string
}

type iProducer interface {
	Produce() ([]string, error)
}

func (prod Producer) Produce() ([]string, error) {
	file, err := os.Open(prod.inputFileName)
	if err != nil {
		return nil, err
	}

	var text []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	error := scanner.Err()
	if error != nil {
		return nil, error
	}
	prod.fileText = text
	return text, nil
}
