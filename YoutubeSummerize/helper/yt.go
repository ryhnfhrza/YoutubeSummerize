package helper

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExtractVideoID(ytLink string) (string, error) {
    parsedURL, err := url.Parse(ytLink)
    if err != nil {
        return "", err
    }

    if strings.Contains(parsedURL.Host, "youtube.com") || strings.Contains(parsedURL.Host, "youtu.be") {
        queryParams := parsedURL.Query()
        videoID := queryParams.Get("v")
        if videoID == "" {
            pathParts := strings.Split(parsedURL.Path, "/")
            return pathParts[len(pathParts)-1], nil 
        }
        return videoID, nil
    }
    return "", fmt.Errorf("the url is not from youtube")
}

func GetSubtitle(videoID string) (string, error) {
    cmd := exec.Command("bin\\yt-dlp.exe", "--write-auto-sub", "--sub-lang", "en", "--skip-download", "https://www.youtube.com/watch?v="+videoID)

    fmt.Println("Running command:", cmd.String())

    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error running command:", err)
        fmt.Println("Command output:", string(output)) 
        return "", err
    }

    subtitlePattern := fmt.Sprintf("*%s*.en.vtt", videoID)
    matches, err := filepath.Glob(subtitlePattern)
    if err != nil || len(matches) == 0 {
        return "", fmt.Errorf("subtitle file not found for video ID: %s", videoID)
    }

    subtitleFile := matches[0]
    fmt.Println("Subtitle file found:", subtitleFile)

    
    cleanedSubtitle, err := CleanSubtitle(subtitleFile)
	if err != nil {
		return "", err
	}

    //delete file subtitle
    err = os.Remove(subtitleFile)
	if err != nil {
		fmt.Println("Error deleting subtitle file:", err)
		return "", err
	}
	fmt.Println("Subtitle file deleted:", subtitleFile)

    return cleanedSubtitle, nil
}
