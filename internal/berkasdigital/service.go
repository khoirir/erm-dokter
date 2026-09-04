package berkasdigital

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	GetMasterBerkas(ctx context.Context) ([]MasterBerkasDigital, error)
	GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]BerkasDigitalPerawatan, error)
	GetBerkasStream(ctx context.Context, targetURL string) (*BerkasStream, error)
	BuildFullURL(lokasiFile string) string
}

type service struct {
	repo             Repository
	urlBerkasDigital string
	log              *logger.Logger
	httpClient       *http.Client
}

func NewService(repo Repository, urlBerkasDigital string, log *logger.Logger) Service {
	baseURL := ""
	if urlBerkasDigital != "" {
		baseURL = strings.TrimRight(urlBerkasDigital, "/") + "/"
	}
	return &service{
		repo:             repo,
		urlBerkasDigital: baseURL,
		log:              log,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *service) GetMasterBerkas(ctx context.Context) ([]MasterBerkasDigital, error) {
	list, err := s.repo.FetchMasterBerkas(ctx)
	if err != nil {
		s.log.Error("gagal mengambil master berkas digital: %v", err)
		return nil, err
	}
	return list, nil
}

func (s *service) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]BerkasDigitalPerawatan, error) {
	list, err := s.repo.FetchBerkasByNoRawat(ctx, noRawat, kodeList)
	if err != nil {
		s.log.Error("gagal mengambil berkas digital no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return list, nil
}

func (s *service) GetBerkasStream(ctx context.Context, targetURL string) (*BerkasStream, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, apperror.NewBusinessError("URL berkas digital tidak valid")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		s.log.Error("gagal membuat request stream berkas %s: %v", targetURL, err)
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.log.Error("gagal fetch berkas dari storage %s: %v", targetURL, err)
		return nil, apperror.NewNotFoundError("Berkas digital tidak dapat diakses dari server penyimpanan")
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		s.log.Warn("server berkas mengembalikan status %d untuk URL %s", resp.StatusCode, targetURL)
		return nil, apperror.NewNotFoundError("Berkas digital tidak ditemukan di server penyimpanan")
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(parsedURL.Path))
		switch ext {
		case ".pdf":
			contentType = "application/pdf"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		default:
			contentType = "application/pdf"
		}
	}

	filename := filepath.Base(parsedURL.Path)
	if filename == "" || filename == "." {
		filename = "dokumen.pdf"
	}

	return &BerkasStream{
		Body:          resp.Body,
		ContentType:   contentType,
		ContentLength: resp.Header.Get("Content-Length"),
		Filename:      filename,
	}, nil
}

func (s *service) BuildFullURL(lokasiFile string) string {
	if strings.HasPrefix(lokasiFile, "http://") || strings.HasPrefix(lokasiFile, "https://") {
		return lokasiFile
	}
	cleanPath := strings.TrimPrefix(lokasiFile, "/")
	return s.urlBerkasDigital + cleanPath
}

