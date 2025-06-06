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

	s.fileText = make([]string, len(data))
	for i, line := range data {
		s.fileText[i] = s.SpamMasker(line)
	}

	s.pres.Present(s.fileText)
	return nil
}

func (s Service) SpamMasker(input string) string { //функция переделана в метод сервиса

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

	return string(link)
}
