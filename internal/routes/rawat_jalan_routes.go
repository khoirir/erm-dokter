package routes

func (c *RouteConfig) setupRawatJalanRoutes() {
	auth := c.AuthMiddleware
	timeout := c.TimeoutMiddleware

	c.Mux.HandleFunc("GET /api/v1/rawat-jalan/referensi-filter", auth(timeout(c.RawatJalanHandler.GetReferensiFilter)))
	c.Mux.HandleFunc("POST /api/v1/rawat-jalan/antrean", auth(timeout(c.RawatJalanHandler.DaftarAntreanDokter)))
	c.Mux.HandleFunc("GET /api/v1/rawat-jalan/detail/{no_rawat}", auth(timeout(c.RawatJalanHandler.DetailKunjungan)))
}
