package main

import (
	"context"
	"fmt"

	service "github.com/ryhnfhrza/YoutubeSummerize/services/YtSummerizeService"
)

func main() {

	ytSummerizeService := service.NewYtSummerizaService()

	ctx := context.Background()
	result,err := ytSummerizeService.Summerize(ctx, "https://www.youtube.com/watch?v=Nfg-l-syLsA")
	if err != nil {
		fmt.Printf("Failed to summarize video: %v\n", err)
	}

	fmt.Println(result)
}
