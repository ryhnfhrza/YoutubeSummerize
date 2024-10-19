package helper

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ryhnfhrza/YoutubeSummerize/exception"
)

func ListAvailableSubtitles(videoID string) (map[string]string, error) {
	cmd := exec.Command("bin\\yt-dlp.exe", "--list-subs", "https://www.youtube.com/watch?v="+videoID)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, exception.NewInternalServerError(fmt.Sprintf("Failed to list subtitles for video ID %s: %v", videoID, err))
	}

	subtitles := make(map[string]string)
	lines := strings.Split(string(output), "\n")

	langPattern := regexp.MustCompile(`^(\w+)\s+[\w\s]+vtt`)

	for _, line := range lines {
		matches := langPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			langCode := matches[1] 
			subtitles[langCode] = "[auto-generated]"
		}
	}

	if len(subtitles) == 0 {
		return nil, exception.NewNotFoundError(fmt.Sprintf("No subtitles available for video ID: %s", videoID))
	}

	return subtitles, nil
}

func CleanSubtitle(subtitleFile string) (string, error) {
	file, err := os.Open(subtitleFile)
	if err != nil {
		return "", exception.NewInternalServerError(fmt.Sprintf("Failed to open subtitle file %s: %v", subtitleFile, err))
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
		return "", exception.NewInternalServerError(fmt.Sprintf("Error reading subtitle file %s: %v", subtitleFile, err))
	}

	return sb.String(), nil
}

func removeVTTTags(text string) string {
	tagPattern := regexp.MustCompile(`<.*?>`)
	return tagPattern.ReplaceAllString(text, "")
}
