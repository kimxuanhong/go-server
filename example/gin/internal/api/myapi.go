package api

import (
	"github.com/kimxuanhong/go-server/core"
)

// MyApiHandler
// @BaseUrl /iloveu/
type MyApiHandler struct {
}

// SayHi API
// @Api GET /say-hi
func (h *MyApiHandler) SayHi(c core.Context) {
	c.JSON(200, map[string]string{
		"message": "SayHi Api",
	})
}

// Print555 API
// @Api GET /print/555
func (h *MyApiHandler) Print555(c core.Context) {
	c.JSON(200, map[string]string{
		"message": "Print555 Api",
	})
}

func (h *MyApiHandler) SayHello(c core.Context) {
	c.JSON(200, map[string]string{
		"message": "SayHello",
	})
}

func (h *MyApiHandler) Routes() []core.RouteConfig {
	return []core.RouteConfig{
		{
			Method:  "GET",
			Path:    "/say-hello",
			Handler: h.SayHello,
		},
	}
}
