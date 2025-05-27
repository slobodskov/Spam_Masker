package presenter

import "os"

type Presenter struct {
	outputFileName string
}

type IPresenter interface {
	Present([]string) error
}

func NewPresenter(outputFileName string) *Presenter {
	return &Presenter{outputFileName: outputFileName}
}

func (pres Presenter) Present(data []string) error {
	name := pres.outputFileName
	if _, err := os.Stat(name); os.IsNotExist(err) {
		//файла с таким названием нет
	} else {
		os.Remove(name)
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(data); i++ {
		_, err := file.WriteString(data[i] + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
