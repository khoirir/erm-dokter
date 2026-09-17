package master

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type icdSnapshot struct {
	icd10List []ICD10
	icd9List  []ICD9
	icd10Map  map[string]bool
	icd9Map   map[string]bool
	syncedAt  time.Time
}

type ICDCache struct {
	snapshot atomic.Pointer[icdSnapshot]
	syncMu   sync.Mutex
}

func NewICDCache() *ICDCache {
	return &ICDCache{}
}

func (c *ICDCache) IsLoaded() bool {
	return c.snapshot.Load() != nil
}

func (c *ICDCache) LastSync() time.Time {
	snap := c.snapshot.Load()
	if snap == nil {
		return time.Time{}
	}
	return snap.syncedAt
}

func (c *ICDCache) NeedsRefresh(ttl time.Duration) bool {
	snap := c.snapshot.Load()
	if snap == nil {
		return true
	}
	return time.Since(snap.syncedAt) > ttl
}

func (c *ICDCache) Load(icd10 []ICD10, icd9 []ICD9) {
	c.syncMu.Lock()
	defer c.syncMu.Unlock()

	map10 := make(map[string]bool, len(icd10))
	for _, item := range icd10 {
		map10[item.Kode] = true
	}

	map9 := make(map[string]bool, len(icd9))
	for _, item := range icd9 {
		map9[item.Kode] = true
	}

	newSnap := &icdSnapshot{
		icd10List: icd10,
		icd9List:  icd9,
		icd10Map:  map10,
		icd9Map:   map9,
		syncedAt:  time.Now(),
	}

	c.snapshot.Store(newSnap)
}

func (c *ICDCache) SearchICD10(filter FilterMasterICD) ([]ICD10, int) {
	snap := c.snapshot.Load()
	if snap == nil {
		return nil, 0
	}

	if filter.Keyword == "" {
		total := len(snap.icd10List)
		return paginateSlice(snap.icd10List, filter.Offset(), filter.Limit), total
	}

	kwLower := strings.ToLower(filter.Keyword)
	matched := make([]ICD10, 0, 64)
	for _, item := range snap.icd10List {
		if strings.Contains(strings.ToLower(item.Kode), kwLower) || strings.Contains(strings.ToLower(item.Nama), kwLower) {
			matched = append(matched, item)
		}
	}

	total := len(matched)
	return paginateSlice(matched, filter.Offset(), filter.Limit), total
}

func (c *ICDCache) SearchICD9(filter FilterMasterICD) ([]ICD9, int) {
	snap := c.snapshot.Load()
	if snap == nil {
		return nil, 0
	}

	if filter.Keyword == "" {
		total := len(snap.icd9List)
		return paginateSlice(snap.icd9List, filter.Offset(), filter.Limit), total
	}

	kwLower := strings.ToLower(filter.Keyword)
	matched := make([]ICD9, 0, 64)
	for _, item := range snap.icd9List {
		if strings.Contains(strings.ToLower(item.Kode), kwLower) || strings.Contains(strings.ToLower(item.Nama), kwLower) {
			matched = append(matched, item)
		}
	}

	total := len(matched)
	return paginateSlice(matched, filter.Offset(), filter.Limit), total
}

func (c *ICDCache) CheckICD10(listKode []string) map[string]bool {
	snap := c.snapshot.Load()
	if snap == nil {
		return nil
	}

	result := make(map[string]bool, len(listKode))
	for _, kode := range listKode {
		result[kode] = snap.icd10Map[kode]
	}
	return result
}

func (c *ICDCache) CheckICD9(listKode []string) map[string]bool {
	snap := c.snapshot.Load()
	if snap == nil {
		return nil
	}

	result := make(map[string]bool, len(listKode))
	for _, kode := range listKode {
		result[kode] = snap.icd9Map[kode]
	}
	return result
}

func (c *ICDCache) Counts() (int, int) {
	snap := c.snapshot.Load()
	if snap == nil {
		return 0, 0
	}
	return len(snap.icd10List), len(snap.icd9List)
}

func paginateSlice[T any](items []T, offset, limit int) []T {
	total := len(items)
	if offset >= total {
		return []T{}
	}
	end := offset + limit
	if end > total {
		end = total
	}
	result := make([]T, end-offset)
	copy(result, items[offset:end])
	return result
}
