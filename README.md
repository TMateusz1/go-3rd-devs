# Go-3rd-Devs

A Go application that provides a chat interface with AI, maintaining conversation context through summaries. <br>
Source in TS: https://github.com/i-am-alice/3rd-devs

## Prerequisites

- Go 1.24 or higher
- OpenAI API key

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/TMateusz1/go-3rd-devs.git
   cd go-3rd-devs
   ```

2. Create a `.env` file in the root directory with your OpenAI API key:

   You can use OpenRouter as well because it's api is same as OpenAI
   ```
   OPENAI_API_KEY=your_api_key_here (required)
   OPENAI_MODEL=openai/gpt-4o-mini (not required OpenAI gpt-4o-mini default)
   OPENAI_BASE_URL=https://openrouter.ai/api/v1 (not required OpenAI default)
   FIRECRAWL_API_KEY=Firecrawl_api_key (required for websearch)
   ```

3. Install dependencies:
   ```
   go mod download
   ```

## Running the Application

Run the examples using the Makefile ```make {example}```

Current examples:

```
make thread
make websearch
```

This will start a server on port 8080 with endpoints /api/{example}: /api/thread