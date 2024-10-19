package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ryhnfhrza/YoutubeSummerize/exception"
	"github.com/ryhnfhrza/YoutubeSummerize/helper"
	//"github.com/ryhnfhrza/YoutubeSummerize/model/web"
)

type YtSummerizaServiceImpl struct {
}

func NewYtSummerizaService() YtSummerizeService {
	return &YtSummerizaServiceImpl{}
}

func (service *YtSummerizaServiceImpl) Summerize(ctx context.Context, ytLink string) (string,error) {
	var youtubeRegex = regexp.MustCompile(`^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$`)

	if !youtubeRegex.MatchString(ytLink) {
		return "", exception.NewBadRequestError("Invalid YouTube link format")
	}

	videoID, err := helper.ExtractVideoID(ytLink)
	if err != nil {
		fmt.Println("Attempting to fetch subtitles for video ID:", videoID)
		fmt.Println("Error extracting video ID:", err)
		return "", fmt.Errorf("error extracting video ID: %v", err)
	}

	subtitleFile, err := helper.GetSubtitle(videoID)
	if err != nil {
		fmt.Println("Error getting subtitle:", err)
		return "", fmt.Errorf("error getting subtitle: %v", err)
	}

	return subtitleFile,nil
}
