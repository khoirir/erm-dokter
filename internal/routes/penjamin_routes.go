package routes

func (c *RouteConfig) setupPenjaminRoutes() {
	auth := c.AuthMiddleware
	timeout := c.TimeoutMiddleware

	c.Mux.HandleFunc("GET /api/v1/penjamin", auth(timeout(c.PenjaminHandler.DaftarPenjamin)))
}
