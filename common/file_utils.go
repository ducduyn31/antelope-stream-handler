package common

import (
	"bufio"
	"fmt"
	"os"
)

func FileHaveValue(path string, line string) (bool, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return StrArrContains(line, lines), nil
}

func OpenFileAndReadLines(path string) ([]string, *os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0600)
	if err != nil {
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		return nil, nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, f, scanner.Err()
}

func AppendLineAndClose(file *os.File, line string)  {
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if _, err := file.WriteString(line + "\n"); err != nil {
		panic(err)
	}
}

func WriteOverAndClose(file *os.File, lines []string) error {
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	err := file.Truncate(0)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = fmt.Fprintf(writer, "%s\n", line)
	}
	return writer.Flush()
}