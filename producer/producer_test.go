package producer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Produce(t *testing.T) {
	serv := &Producer{
		inputFileName: "test.txt",
	}

	file, err := os.Create(serv.inputFileName)
	assert.NoError(t, err)
	defer os.Remove(serv.inputFileName)

	_, err = file.WriteString("str 1\nstr 2\n")
	assert.NoError(t, err)
	file.Close()

	result, err := serv.Produce()
	assert.NoError(t, err)
	assert.Equal(t, []string{"str 1", "str 2"}, result)
}
