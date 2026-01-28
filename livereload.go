//go:build !prod

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

var serverStartTime = fmt.Sprintf("%d", time.Now().UnixNano())

func setupLiveReload(e *echo.Echo) {
	e.GET("/dev/reload", liveReloadHandler)
}

func liveReloadHandler(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	w := c.Response()
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "data: %s\n\n", serverStartTime)
	flusher.Flush()

	<-c.Request().Context().Done()
	return nil
}

func liveReloadScript() string {
	return `<script>
(function() {
	let lastServerId = null;
	let retryCount = 0;
	
	function connect() {
		const es = new EventSource('/dev/reload');
		
		es.onmessage = function(e) {
			const serverId = e.data;
			if (lastServerId !== null && lastServerId !== serverId) {
				console.log('[livereload] server restarted, reloading...');
				location.reload();
			}
			lastServerId = serverId;
			retryCount = 0;
		};
		
		es.onerror = function() {
			es.close();
			retryCount++;
			const delay = Math.min(1000 * retryCount, 5000);
			console.log('[livereload] connection lost, retrying in ' + delay + 'ms...');
			setTimeout(connect, delay);
		};
	}
	
	connect();
})();
</script>`
}
