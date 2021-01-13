
- always discard body e.g. io.Copy(ioutil.Discard, resp.Body) if you don't use it
  - HTTP client's Transport will not reuse connections unless the body is read to completion and closed

```go
res, _ := client.Do(req)
io.Copy(ioutil.Discard, res.Body)
defer res.Body.Close()
```

```go
// https://github.com/hashicorp/go-retryablehttp/blob/master/client.go
// Try to read the response body so we can reuse this connection.
func (c *Client) drainBody(body io.ReadCloser) {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
	if err != nil {
		if c.logger() != nil {
			switch v := c.logger().(type) {
			case LeveledLogger:
				v.Error("error reading response body", "error", err)
			case Logger:
				v.Printf("[ERR] error reading response body: %v", err)
			}
		}
	}
}
```
