package main

import (
	"bufio"
	"fmt"
	"os"
)

func stringInArray(str string, array []string) bool {
	for _, b := range array {
		if b == str {
			return true
		}
	}
	return false
}

func stringInSet(str string, d Set) bool {
	if _, ok := d[str]; ok {
		return true
	} else {
		return false
	}
}

func saveToFile(filePath string, lines []string) error {
	f, err := os.Create(filePath)
	if os.IsExist(err) {
		f, err = os.Open(filePath)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range lines {
		fmt.Fprintln(f, value)
	}
	return nil
}

func readFromFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if os.IsExist(err) {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
