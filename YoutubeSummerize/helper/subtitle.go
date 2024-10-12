package helper

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)


func CleanSubtitle(subtitleFile string) (string, error) {
	file, err := os.Open(subtitleFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)

	timePattern := regexp.MustCompile(`\d{2}:\d{2}:\d{2}\.\d{3} --> \d{2}:\d{2}:\d{2}\.\d{3}`)

	previousLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		if timePattern.MatchString(line) || strings.TrimSpace(line) == "" {
			continue
		}

		cleanLine := removeVTTTags(line)

		if cleanLine == previousLine {
			continue
		}

		sb.WriteString(cleanLine + "\n")
		previousLine = cleanLine
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func removeVTTTags(text string) string {
	tagPattern := regexp.MustCompile(`<.*?>`)
	return tagPattern.ReplaceAllString(text, "")
}
