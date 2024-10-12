package service

import (
	"context"
)

type YtSummerizeService interface {
	Summerize(ctx context.Context, ytLink string) string
}