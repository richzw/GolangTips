
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

- Q: golang的net如何对于epoll进行封装，使用上看似通过过程呢？

- 所有的网络操作以网络描述符netFD为中心实现，netFD与底层PollDesc结构绑定，当在一个netFD上读写遇到EAGAIN的错误的时候，把当前goroutine存储到这个netFD对应的PollDesc中，同时把goroutine给park 直到这个netFD发生读写事件，才将该goroutine激活重新开始。在底层通知goroutine再次发生读写事件的方式就是epoll事件驱动机制。

- Q: 为什么 gnet 会比 Go 原生的 net 包更快？
- Multi-Reactors 模型相较于 Go 原生模型在以下场景具有性能优势： 
  - 1. 高频创建新连接：我们从源码里可以知道 Go 模式下所有事件都是在一个 epoll 实例来管理的，接收新连接和 IO 读写；而在 Reactors 模式下，accept 新连接和 IO 读写分离，它们在各自独立的 goroutines 里用自己的 epoll 实例来管理网络事件。 
  - 2. 海量网络连接：Go net 处理网络请求的模式是 goroutine per connection，甚至是 multiple goroutines per connection，而 gnet 一般使用与机器 CPU 核心数相同的 goroutines 来处理网络请求，所以在海量网络连接的场景下 gnet 更节省系统资源，进而提高性能。 
  - 3. 时间窗口内连接总数大而活跃连接数少：这种场景下，Go 原生网络模型因为 goroutine per connection 模式，依然需要维持大量的 goroutines 去等待 IO 事件(保持 1:1 的关系)，Go scheduler 对大量 idle goroutines 的调度势必会损耗系统整体性能；而 gnet 模式下需要维护的仅仅是与 CPU 核心数相同的 goroutines，而且得益于 Reactors 模型和 epoll/kqueue，可以确保每个 goroutine 在大多数时间里都是在处理活跃连接。 
  - 4. 短连接场景：gnet 内部维护了一个内存池，在短连接这种场景下，可以大量复用内存，进一步节省资源和提高性能。

- Q: Go 的网络模型有『惊群效应』吗？
- 没有。 我们看下源码里是怎么初始化 listener 的 epoll 示例的：

```go
var serverInit sync.Once

func (pd *pollDesc) init(fd *FD) error {
	serverInit.Do(runtime_pollServerInit)
	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
```

   这里用了 sync.Once 来确保初始化一次 epoll 实例，这就表示一个 listener 只持有一个 epoll 实例来管理网络连接，既然只有一个 epoll 实例，当然就不存在『惊群效应』了

- gnet - Servers can utilize the SO_REUSEPORT 3 option which allows multiple sockets on the same host to bind to the same port and the OS kernel takes care of the load balancing for you, it wakes one socket per accpet event coming to resolved the thundering herd


