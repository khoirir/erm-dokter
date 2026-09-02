package berkasdigital

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	GetMasterBerkas(ctx context.Context) ([]MasterBerkasDigital, error)
	GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]BerkasDigitalDB, error)
	StreamBerkasDigital(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error
	BuildBerkasItem(kode string, namaBerkas string, lokasiFile string) (*BerkasDigitalItem, error)
	BuildFullURL(lokasiFile string) string
}

type service struct {
	repo             Repository
	encryptionKey    string
	urlBerkasDigital string
	log              *logger.Logger
	httpClient       *http.Client
}

func NewService(repo Repository, encryptionKey string, urlBerkasDigital string, log *logger.Logger) Service {
	baseURL := ""
	if urlBerkasDigital != "" {
		baseURL = strings.TrimRight(urlBerkasDigital, "/") + "/"
	}
	return &service{
		repo:             repo,
		encryptionKey:    encryptionKey,
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

func (s *service) GetBerkasByNoRawat(ctx context.Context, noRawat string, kodeList []string) ([]BerkasDigitalDB, error) {
	list, err := s.repo.FetchBerkasByNoRawat(ctx, noRawat, kodeList)
	if err != nil {
		s.log.Error("gagal mengambil berkas digital no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	return list, nil
}

func (s *service) StreamBerkasDigital(ctx context.Context, encryptedIdBerkas string, w http.ResponseWriter) error {
	rawTargetURL, err := crypto.Decrypt(encryptedIdBerkas, s.encryptionKey)
	if err != nil {
		return apperror.NewBusinessError("Token berkas digital tidak valid")
	}

	parsedURL, err := url.Parse(rawTargetURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return apperror.NewBusinessError("URL berkas digital tidak valid")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawTargetURL, nil)
	if err != nil {
		s.log.Error("gagal membuat request stream berkas %s: %v", rawTargetURL, err)
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.log.Error("gagal fetch berkas dari storage %s: %v", rawTargetURL, err)
		return apperror.NewNotFoundError("Berkas digital tidak dapat diakses dari server penyimpanan")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.log.Warn("server berkas mengembalikan status %d untuk URL %s", resp.StatusCode, rawTargetURL)
		return apperror.NewNotFoundError("Berkas digital tidak ditemukan di server penyimpanan")
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

	w.Header().Set("Content-Type", contentType)
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		w.Header().Set("Content-Length", cl)
	}

	filename := filepath.Base(parsedURL.Path)
	if filename == "" || filename == "." {
		filename = "dokumen.pdf"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))

	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, resp.Body)
	return err
}

func (s *service) BuildBerkasItem(kode string, namaBerkas string, lokasiFile string) (*BerkasDigitalItem, error) {
	if lokasiFile == "" {
		return nil, nil
	}
	fullURL := s.BuildFullURL(lokasiFile)
	encIdBerkas, err := crypto.Encrypt(fullURL, s.encryptionKey)
	if err != nil {
		s.log.Error("gagal mengenkripsi URL berkas: %v", err)
		return nil, err
	}

	return &BerkasDigitalItem{
		Kode:       kode,
		NamaBerkas: namaBerkas,
		IdBerkas:   encIdBerkas,
		UrlBerkas:  "/api/v1/berkas-digital/" + encIdBerkas,
	}, nil
}

func (s *service) BuildFullURL(lokasiFile string) string {
	if strings.HasPrefix(lokasiFile, "http://") || strings.HasPrefix(lokasiFile, "https://") {
		return lokasiFile
	}
	cleanPath := strings.TrimPrefix(lokasiFile, "/")
	return s.urlBerkasDigital + cleanPath
}
