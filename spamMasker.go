package main

import (
	"fmt"
)

func SpamMasker(input string) string {

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

func Test(text string) {
	output := SpamMasker(text)
	fmt.Println(output)
}

func main() {
	Test("Here's my spammy page: http://hehefouls.netHAHAHA see you.")
	Test("http://hehefouls.netHAHAHA")
	Test("one more link hTTp://cherry")
	Test("http/new.laptop")
	Test("1 http://the_last-one")
}
