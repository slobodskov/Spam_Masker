package service

import (
	"fmt"
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

func TestService_Run(t *testing.T) {
	mockProd := new(MockServ)
	mockPres := new(MockServ)
	serv := &Service{
		prod: mockProd,
		pres: mockPres,
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
