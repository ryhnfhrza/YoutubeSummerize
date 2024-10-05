package service

import (
	"context"

	"github.com/ryhnfhrza/YoutubeSummerize/model/web"
)

type YtSummerizeService interface {
	Summerize(ctx context.Context, ytLink string) web.SummerizeResponses
}