# Go Redis Ollama

## tl;dr

```powershell
git clone "https://github.com/davidhintelmann/red-ollama.git"
go build -o "llm.exe"
.\llm.exe
```
use p flag to enter a prompt

```powershell
.\llm.exe -p "tell me a joke"
```

use m flag to enter which model you want to use

```powershell
.\llm.exe -p "tell me a joke" -m "phi3"
```

## Summary

This repo is an example of how one can use the go programming language to send prompts to [Ollama](https://ollama.com/) server hosted locally.
Using Ollama one can request prompts from LLM or SLM hosted locally.

For example you can download and serve:
- Microsoft's Phi3 SLM
- Meta's Llama3.1 LLM

Additionally, using Redis to cache prompts along with their responses

## Before Using this Repo

Prerequisites:
1. [Download go](https://go.dev/dl/)
2. [Install Ollama](https://ollama.com/)
3. [Install Redis on Windows](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-on-windows/)

## Cache Responses

A simple approach to using Redis is to cache prompts along with their response, and then if a user enters the same prompt twice then the cached result will be returned instead.

This was developed on Windows 11 and one can use WSL 2 to [install Redis on Windows](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-on-windows/).

### Caution 
This example only uses [Redis Strings](https://redis.io/docs/latest/develop/data-types/#strings) to cache data.

From [Redis docs](https://redis.io/docs/latest/develop/get-started/data-store/)
> Similar to byte arrays, Redis strings store sequences of bytes, including text, serialized objects, counter values, and binary arrays.

There are other types as well, for example:
- Hash
- List
- Geospatial

If you [install Redis Stack](https://redis.io/docs/latest/operate/oss_and_stack/install/install-stack/) you can also store data as JSON, [more info here](https://redis.io/docs/latest/develop/data-types/json/).

LLMs often output their responses in JSON and caching the data in the same format would be the ideal approach to take.