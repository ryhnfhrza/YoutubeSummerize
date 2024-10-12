package main

import (
	"context"
	"fmt"

	service "github.com/ryhnfhrza/YoutubeSummerize/services/YtSummerizeService"
)

func main() {

	ytSummerizeService := service.NewYtSummerizaService()

	ctx := context.Background()
	result := ytSummerizeService.Summerize(ctx,"https://www.youtube.com/watch?v=OA5rrLNvNGg")

	fmt.Println(result)
}