package service

import (
	"context"
	"regexp"

	"github.com/ryhnfhrza/YoutubeSummerize/exception"
	"github.com/ryhnfhrza/YoutubeSummerize/helper"
	"github.com/ryhnfhrza/YoutubeSummerize/model/web"
)

type YtSummerizaServiceImpl struct {
	
}

func NewYtSummerizaService() YtSummerizeService {
	return &YtSummerizaServiceImpl{
	}
}



func(service *YtSummerizaServiceImpl)Summerize(ctx context.Context,ytLink string)web.SummerizeResponses{
	var youtubeRegex = regexp.MustCompile(`^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$`)
	
	if !youtubeRegex.MatchString(ytLink) {
		panic(exception.NewBadRequestError("Invalid YouTube link format"))
	}


	
	 videoID, err := helper.ExtractVideoID(ytLink)
	 if err != nil {
			 helper.PanicIfError(err)
	 }

	 
	 subtitleFile, err := helper.GetSubtitle(videoID)
	 if err != nil {
			 helper.PanicIfError(err)
	 }

	 

	 return response
}

