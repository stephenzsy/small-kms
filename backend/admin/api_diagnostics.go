package admin

import (
	"runtime"

	"github.com/gin-gonic/gin"
)

func (s *adminServer) GetDiagnosticsV1(c *gin.Context) {
	r := RequestDiagnostics{
		ServiceRuntime: RequestDiagnostics_ServiceRuntime{
			GoVersion: runtime.Version(),
		},
	}
	for k, v := range c.Request.Header {
		r.RequestHeaders = append(r.RequestHeaders, RequestHeaderEntry{
			Key:   k,
			Value: v,
		})
	}
	c.JSON(200, r)
}
