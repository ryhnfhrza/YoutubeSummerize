package helper

import (
	"fmt"
	"net/url"
	"os/exec"
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
    cmd := exec.Command("bin\\youtube-dl.exe", "--write-auto-sub", "--sub-lang", "en", "--skip-download", "https://www.youtube.com/watch?v="+videoID)
    err := cmd.Run()
    if err != nil {
        return "", err
    }

    subtitleFile := videoID + ".en.vtt"
    return subtitleFile, nil
}
