package utility

import (
	"log"
	"os"
)

func ParseTextFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var lines []string
	for {
		var line string
		_, err := file.Read([]byte(line))
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	return lines
}
