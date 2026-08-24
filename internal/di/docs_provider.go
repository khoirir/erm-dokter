package di

import "erm-dokter/internal/docs"

func provideDocs() *docs.Handler {
	return docs.NewHandler()
}
