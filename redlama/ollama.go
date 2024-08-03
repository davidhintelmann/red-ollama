/*
ollama package is for interacting with local Ollama server.

Local instance of Redis is running on WSL 2 to return cached responses.
*/
package redlama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

/*
struct for decoding json response returned
from Ollama server.
*/
type OllamaOutput struct {
	Model             string `json:"model"`
	CreatedAt         string `json:"created_at"`
	Response          string `json:"response"`
	Done              bool   `json:"done"`
	DoneReason        string `json:"done_reason"`
	Context           []int  `json:"context"`
	TotalDuration     int    `json:"total_duration"`
	LoadDuration      int    `json:"load_duration"`
	PromptValCount    int    `json:"prompt_val_count"`
	PromptValDuration int    `json:"prompt_eval_duration"`
	EvalCount         int    `json:"eval_count"`
	EvalDuration      int    `json:"eval_duration"`
}

/*
CheckLocalConnetion function will return a string
for Ollama server status.

Must be running Ollama (hosted locally).

Otherwise this function returns an error.

Returns:
  - string: response for Ollama server status
  - int: http status code
  - error: error message if there is one
*/
func CheckLocalConnetion() (string, int, error) {
	resp, err := http.Get("http://localhost:11434")
	if err != nil {
		return "", -1, errors.New("ollama is not running")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", resp.StatusCode, err
		}
		return string(bodyBytes), resp.StatusCode, nil
	} else {
		return "", resp.StatusCode, nil
	}
}

/*
PromptOllama function will return a response,
encoded as json, and a status code.

Must be running Ollama (hosted locally).

Otherwise this function returns an error.

Parameters:
  - prompt: string, message or question you wish to ask
  - model: string, model that ollama will use for prompt
  - cache: bool, set to true to used cached response from redis
    or set to false to reset cached response
  - redisClient: *redis.Client, use output of RedisClient function for this functions parameter

Returns:
  - *OllamaOutput: response returned as json struct
  - int: http status code
  - error: error message if there is one
*/
func PromptOllama(ctx context.Context, prompt string, model string, cache bool, redisClient *redis.Client) (*OllamaOutput, int, error) {
	// check if prompt is in redis cache
	// then return prompt if it exists and reset
	// parameter is set to false
	// if reset is set to true by pass the cache
	checkRedis := fmt.Sprintf("%s:prompt:%s", strings.ToLower(model), strings.ToLower(prompt))
	val, err := redisClient.Get(ctx, checkRedis).Result()
	if err == redis.Nil || !cache {
		return postOllama(ctx, prompt, model, checkRedis, redisClient)
	} else { // return prompt output if it is cached
		return &OllamaOutput{Response: val}, -2, err
	}
}

// internal function for PromptOllama function
func postOllama(ctx context.Context, prompt string, model string, checkRedis string, redisClient *redis.Client) (*OllamaOutput, int, error) {
	// HTTP endpoint
	postURL := "http://localhost:11434/api/generate"

	// JSON body
	input := fmt.Sprintf(`{"model": "%s", "prompt": "%s", "stream": %t}`, model, prompt, false)
	body := []byte(input)

	// Create a HTTP post request
	request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, -1, errors.New("can not create a new POST request")
	}

	// Add headers
	request.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, -1, errors.New("can not add headers to POST request")
	}
	defer resp.Body.Close()

	// check if response is okay
	// then return response as json
	if resp.StatusCode == http.StatusOK {
		output := &OllamaOutput{}
		outputCheck := json.NewDecoder(resp.Body).Decode(output)
		if outputCheck != nil {
			return nil, resp.StatusCode, errors.New("can not decode json from struct")
		}
		// set key value pair into redis
		// then return output
		setErr := redisClient.Set(ctx, checkRedis, output.Response, 0).Err()
		if setErr != nil {
			return nil, -1, errors.New("can not set value in redis")
		}
		return output, resp.StatusCode, nil
	} else {
		return nil, resp.StatusCode, errors.New("error with PromptOllama function")
	}
}

/*
RedisClient function will connect to local Redis server,
hosted on WSL2, and return a Redis client.

The client returned is used by passing it to functions requiring
Redis to cache data.

Parameters:
  - db_num: int, redis db number [0-15]

Returns:
  - *redis.Client
  - error: error message if there is one
*/
func RedisClient(db_num int) (*redis.Client, error) {
	url, port := "localhost", "6379"
	addr := fmt.Sprintf("%s:%s", url, port)
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
		// Username: "",     // no username set
		Password: "",     // no password set
		DB:       db_num, // use default DB
	})

	addrCheck := strings.Split(redisClient.Options().Addr, ":")
	if len(addrCheck) > 2 {
		errorText := fmt.Sprintf("address to redis local host is incorrect\nexpected %s\ngot %s\n", addr, addrCheck)
		return nil, errors.New(errorText)
	} else if addrCheck[0] != url || addrCheck[1] != port {
		errorText := fmt.Sprintf("address to redis local host is incorrect\nexpected url: %s\ngot url: %s\nexpected port: %s\ngot port: %s\n", url, addrCheck[0], port, addrCheck[1])
		return nil, errors.New(errorText)
	}

	return redisClient, nil
}
