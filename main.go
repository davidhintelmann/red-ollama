package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/davidhintelmann/red-ollama/redlama"
)

func main() {
	ctx := context.Background()

	defaultPrompt := "i forget how to prompt"
	defaultModel := "llama3.1"
	pFlag := flag.String("p", defaultPrompt, "enter your prompt for LLM")
	mFlag := flag.String("m", defaultModel, "enter your LLM model name served by Ollama")
	flag.Parse()
	if *pFlag == defaultPrompt {
		fmt.Println("Use -p flag followed by your text prompt, ie -p \"tell me a dirty joke\"")
		fmt.Println()
		fmt.Printf("-=RESPONSE=-\n\n")
	}

	_, _, err := redlama.CheckLocalConnetion()
	if err != nil {
		log.Fatalf("error with ollama connection: %v\n", err)
	}

	redisClient, err := redlama.RedisClient(ctx, 2)
	if err != nil {
		log.Fatalf("error with redis connection: %v\n", err)
	}

	jsonResp, _, err := redlama.PromptOllama(ctx, *pFlag, *mFlag, true, redisClient)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Println(jsonResp.Response)
}
