package util

import (
	"bufio"
	"fmt"
	"strings"
)

func HandleRequiredFlag(fieldValue, fieldName string, userInputReader *bufio.Reader) (string, error) {
	if fieldValue != "" {
		return fieldValue, nil
	}
	result, err := requestUserInput(userInputReader, fieldName)
	if err != nil {
		return "", err
	}
	if result != "" {
		return result, nil
	}
	return "", fmt.Errorf("%s must be provided", fieldName)
}

func requestUserInput(userInputReader *bufio.Reader, fieldName string) (response string, err error) {
	fmt.Printf("Please enter %s: ", fieldName)
	response, err = userInputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(response, "\n"), nil
}
