package presenter

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Present(t *testing.T) {
	serv := &Presenter{
		outputFileName: "output.txt",
	}

	data := []string{"str 1", "str 2"}

	err := serv.Present(data)
	assert.NoError(t, err)

	file, err := os.Open(serv.outputFileName)
	assert.NoError(t, err)
	defer os.Remove(serv.outputFileName)

	scanner := bufio.NewScanner(file)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	assert.NoError(t, scanner.Err())
	assert.Equal(t, data, text)
}
