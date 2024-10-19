package helper

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ryhnfhrza/YoutubeSummerize/exception"
)

func ExtractVideoID(ytLink string) (string, error) {
	parsedURL, err := url.Parse(ytLink)
	if err != nil {
		return "", exception.NewBadRequestError("invalid url format")
	}

	if strings.Contains(parsedURL.Host, "youtube.com") || strings.Contains(parsedURL.Host, "youtu.be") {
		queryParams := parsedURL.Query()
		videoID := queryParams.Get("v")
		if videoID == "" {
			pathParts := strings.Split(parsedURL.Path, "/")
			return pathParts[len(pathParts)-1], exception.NewBadRequestError("unable to find video id from URL")
		}
		return videoID, nil
	}
	return "", exception.NewBadRequestError("the url is not from youtube")
}

func GetSubtitle(videoID string) (string, error) {
	availableSubtitles, err := ListAvailableSubtitles(videoID)
	if err != nil {
		return "", exception.NewInternalServerError(fmt.Sprintf("Failed to get the subtitle list for the ID video: %s", videoID))
	}

	if len(availableSubtitles) == 0 {
		return "", exception.NewInternalServerError(fmt.Sprintf("no subtitles available for video ID: %s", videoID))
	}

	languagePriority := []string{"en", "id"} 
	var subtitleFiles []string  
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make(chan error, len(languagePriority)*2)

	
	for _, lang := range languagePriority {
		wg.Add(1)
		go func(lang string) {
			defer wg.Done()
			cmd := exec.Command("bin\\yt-dlp.exe", "--write-sub", "--sub-lang", lang, "--skip-download", "https://www.youtube.com/watch?v="+videoID)
			_, err := cmd.CombinedOutput()
			if err == nil {
				subtitlePattern := fmt.Sprintf("*%s*.%s.vtt", videoID, lang)
				matches, _ := filepath.Glob(subtitlePattern)
				if len(matches) > 0 {
					mu.Lock()
					subtitleFiles = append(subtitleFiles, matches...)
					mu.Unlock()
				}
			} else {
				errors <- exception.NewInternalServerError(fmt.Sprintf("Failed to write subtitles for languages: %s, error:%v", lang, err))
			}
		}(lang)
	}

	wg.Wait()
	if len(subtitleFiles) == 0 {
		for _, lang := range languagePriority {
			wg.Add(1)
			go func(lang string) {
				defer wg.Done()
				cmd := exec.Command("bin\\yt-dlp.exe", "--write-auto-sub", "--sub-lang", lang, "--skip-download", "https://www.youtube.com/watch?v="+videoID)
				_, err := cmd.CombinedOutput()
				if err == nil {
					subtitlePattern := fmt.Sprintf("*%s*.%s.vtt", videoID, lang)
					matches, _ := filepath.Glob(subtitlePattern)
					if len(matches) > 0 {
						mu.Lock()
						subtitleFiles = append(subtitleFiles, matches...)
						mu.Unlock()
					}
				} else {
					errors <- exception.NewInternalServerError(fmt.Sprintf("Failed to write auto-generated subtitles for languages: %s, error:%v", lang, err))
				}
			}(lang)
		}
	}

	wg.Wait()
	close(errors)

	// delete file after use
	var cleanedSubtitle string
	if len(subtitleFiles) > 0 {
		for _, subtitleFile := range subtitleFiles {
			cleaned, err := CleanSubtitle(subtitleFile)
			if err != nil {
				return "", exception.NewInternalServerError(fmt.Sprintf("Failed to clean up subtitles for ID video: %s, error: %v", videoID, err))
			}
			cleanedSubtitle += cleaned

			
			err = os.Remove(subtitleFile)
			if err != nil {
				return "", exception.NewInternalServerError(fmt.Sprintf("Failed to delete subtitle file: %s for ID video: %s, error: %v", subtitleFile, videoID, err))
			}
		}
		return cleanedSubtitle, nil
	}

	return "", exception.NewNotFoundError(fmt.Sprintf("Subtitle file not found for video ID: %s in the desired language", videoID))
}




