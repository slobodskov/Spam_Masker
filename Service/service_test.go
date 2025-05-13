package service

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServ struct {
	mock.Mock
}

func (m *MockServ) Produce() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockServ) Present(data []string) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestService_Produce(t *testing.T) {
	serv := &Service{
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

func TestService_Present(t *testing.T) {
	serv := &Service{
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

func TestService_Run(t *testing.T) {
	mockProd := new(MockServ)
	mockPres := new(MockServ)
	serv := &Service{
		inputFileName:  "File1.txt",
		outputFileName: "File2.txt",
		prod:           mockProd,
		pres:           mockPres,
	}

	inputData := []string{
		"Line 1",
		"line 2",
		"Here's my spammy page: http://hehefouls.netHAHAHA see you.",
		"http/new.laptop",
	}

	expectedOutput := []string{
		"Line 1",
		"line 2",
		"Here's my spammy page: http://******************* see you.",
		"http/new.laptop",
	}

	mockProd.On("Produce").Return(inputData, nil)
	mockPres.On("Present", expectedOutput).Return(nil)

	er := serv.Run()
	assert.NoError(t, er)

	mockProd.AssertExpectations(t)
	mockPres.AssertExpectations(t)

	fmt.Println("FileText:", serv.fileText)

	assert.Equal(t, []string{"Line 1", "line 2", "Here's my spammy page: http://******************* see you.", "http/new.laptop"}, serv.fileText)

}

func TestService_SpamMasker(t *testing.T) {
	serv := &Service{}

	tests := []struct {
		input    string
		expected string
	}{
		{input: "Here's my spammy page: http://hehefouls.netHAHAHA see you.", expected: "Here's my spammy page: http://******************* see you."},
		{input: "http://hehefouls.netHAHAHA", expected: "http://*******************"},
		{input: "one more link hTTp://cherry", expected: "one more link hTTp://cherry"},
		{input: "http/new.laptop", expected: "http/new.laptop"},
		{input: "1 http://the_last-one", expected: "1 http://************"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("SpamMasker(%s)", test.input), func(t *testing.T) {
			result := serv.SpamMasker(test.input)
			assert.Equal(t, test.expected, result, "Expected %v for input %s, but got %v", test.expected, test.input, result)
		})
	}
}

func TestNewService(t *testing.T) {

	tests := []struct {
		input1, input2  string
		expectedService *Service
	}{
		{input1: "File", input2: "File1", expectedService: &Service{inputFileName: "File", outputFileName: "File1"}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("NewService(%s, %s)", test.input1, test.input2), func(t *testing.T) {
			result := NewService(test.input1, test.input2)
			assert.Equal(t, test.expectedService, result, "Expected %v for input %s and %s, but got %v", test.expectedService, test.input1, test.input2, result)
		})
	}

}

func TestGetProdAndPres(t *testing.T) {
	mockProd := new(MockServ)
	mockPres := new(MockServ)

	service := getProdAndPres(mockProd, mockPres)

	assert.NotNil(t, service)
	assert.Equal(t, mockProd, service.prod)
	assert.Equal(t, mockPres, service.pres)
}
