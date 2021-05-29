
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

- golang的net如何对于epoll进行封装，使用上看似通过过程呢？

- 所有的网络操作以网络描述符netFD为中心实现，netFD与底层PollDesc结构绑定，当在一个netFD上读写遇到EAGAIN的错误的时候，把当前goroutine存储到这个netFD对应的PollDesc中，同时把goroutine给park
- 直到这个netFD发生读写事件，才将该goroutine激活重新开始。在底层通知goroutine再次发生读写事件的方式就是epoll事件驱动机制。
