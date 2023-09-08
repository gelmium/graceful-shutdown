# graceful-shutdown
Graceful shutdown helper for long running executable go program

[![go.yml<main>](https://github.com/gelmium/graceful-shutdown/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/gelmium/graceful-shutdown/actions/workflows/go.yml)

# Usage
``` go
triggerCtx, triggerFn := context.WithCancel(context.Background())
// calling triggerFn() will trigger graceful shutdown manually instead of
// waiting for syscall signal: SIGINT, SIGTERM, SIGHUP or SIGKILL
// if you don't need this feature, just pass context.Background() to GracefulShutdown

// Setup your program resources here, e.g. http server, database connection, etc.

// Graceful shutdown will exit with code 1 after 30s is passed
shutdownTimeout := 30 * time.Second
// Other while <- wait will give code 0
wait := gfshutdown.GracefulShutdown(triggerCtx, shutdownTimeout, map[string]gfshutdown.operation{
		"fiber-http-server": func(ctx context.Context) error {
			return app.ShutdownWithTimeout(shutdownTimeout - 1*time.Second)
		},
		// Add other cleanup operations here, db connection, etc.
	})

// Run your main program here

os.Exit(<-wait)
```
