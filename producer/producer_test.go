import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Produce(t *testing.T) {
	prod := &Producer{
		inputFileName: "test.txt",
	}

	file, err := os.Create(prod.inputFileName)
	assert.NoError(t, err)
	defer os.Remove(prod.inputFileName)

	_, err = file.WriteString("str 1\nstr 2\n")
	assert.NoError(t, err)
	file.Close()

	result, err := prod.Produce()
	assert.NoError(t, err)
	assert.Equal(t, []string{"str 1", "str 2"}, result)
}
