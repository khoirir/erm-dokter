package routes

import (
	"database/sql"
	"net/http"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/middleware"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/usecase"
)

func RegisterRawatJalanRoutes(mux *http.ServeMux, db *sql.DB, cfg *config.Config) {
	repo := repository.NewRawatJalanRepository(db)
	uc := usecase.NewRawatJalanUsecase(repo, cfg.EncryptionKey)
	h := handler.NewRawatJalanHandler(uc)

	authMiddleware := middleware.JWTMiddleware(cfg.JWTSecret)

	mux.HandleFunc("GET /api/v1/rawat-jalan/referensi-filter", authMiddleware(h.GetReferensiFilter))
	mux.HandleFunc("POST /api/v1/rawat-jalan/antrean", authMiddleware(h.DaftarAntreanDokter))
	mux.HandleFunc("GET /api/v1/rawat-jalan/detail/{no_rawat}", authMiddleware(h.DetailKunjungan))
}
