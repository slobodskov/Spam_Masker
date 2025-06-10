package service

import (
	"github.com/slobodskov/spamMasker/presenter"
	"github.com/slobodskov/spamMasker/producer"
)

type Service struct {
	prod     producer.IProducer
	pres     presenter.IPresenter
	fileText []string
}

func NewService(prod producer.IProducer, pres presenter.IPresenter) *Service {
	return &Service{prod: prod, pres: pres} //конструктор структуры
}

func (s *Service) Run() error { //метод Run с использованием SpamMasker
	data, err := s.prod.Produce()
	if err != nil {
		return err
	}
	inputChan := make(chan string, 10) //для избежания deadlock нужен буфер; также нужен для ограничения запущенных spamMasker
	outputChan := make(chan string)
	defer close(outputChan)

	s.fileText = make([]string, len(data))

	for i, line := range data {
		inputChan <- line
		s.SpamMasker(inputChan, outputChan)
		s.fileText[i] = <-outputChan
	}

	close(inputChan)

	s.pres.Present(s.fileText)
	return nil
}

func (s Service) SpamMasker(inputChan <-chan string, outputChan chan<- string) { //функция переделана в метод сервиса
	input := <-inputChan
	go func() {
		data := []byte(input)
		size := len(data)
		var link []byte

		httpPrefix := []byte("http://")
		httpPrefixLen := len(httpPrefix)

		i := 0
		for i < size {
			if i <= size-httpPrefixLen && string(data[i:i+httpPrefixLen]) == string(httpPrefix) {
				link = append(link, data[i:i+httpPrefixLen]...)
				i += httpPrefixLen
				for i < size && (data[i] != ' ' && data[i] != '\n' && data[i] != '\r') {
					data[i] = '*'
					link = append(link, data[i])
					i++
				}
			} else {
				link = append(link, data[i])
				i++
			}
		}
		outputChan <- string(link)
	}()
}
