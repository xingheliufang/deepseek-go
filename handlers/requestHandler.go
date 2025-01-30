package handlers

import (
	"fmt"
	"net/http"
	"time"
  "os"
)
func HandleTimeout() (time.Duration, error){
	if timeout := os.Getenv("DEEPSEEK_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			timeout = d
      return timeout, nil
		} 
	}
    return time.ParseDuration("5m"), err //Default timeout behavior if no timeout is set
}

func HandleSendChatCompletionRequest(req *http.Request) (*http.Response, error) {
  timeout, err = HandleTimeout()
  if err != nil{
    return nil, fmt.Errorf("Error sending request %w", err)
  }
	client := &http.Client{
    Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil{
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}

func HandelNormalRequest(req *http.Request) (*http.Response, error) {
  timeout, err = HandleTimeout()
  if err != nil{
    return nil, fmt.Errorf("Error sending request %w", err)
  }
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err !=nil{
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}


}
