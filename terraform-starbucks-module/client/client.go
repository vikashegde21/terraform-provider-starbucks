package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type StarbucksClient struct {
    APIKey     string
    Endpoint   string
    Region     string
    HTTPClient *http.Client
}

func NewStarbucksClient(apiKey, endpoint, region string, timeout int64) *StarbucksClient {
    return &StarbucksClient{
        APIKey:   apiKey,
        Endpoint: endpoint,
        Region:   region,
        HTTPClient: &http.Client{
            Timeout: time.Duration(timeout) * time.Second,
        },
    }
}

func (c *StarbucksClient) DoRequest(method, path string, body interface{}) ([]byte, error) {
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("error marshaling request: %w", err)
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }

    req, err := http.NewRequest(method, c.Endpoint+path, reqBody)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+c.APIKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Region", c.Region)

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
    }

    return respBody, nil
}
