package autogen

import "testing"

func ptr(v uint64) *uint64 { return &v }

func TestPager(t *testing.T) {
	tests := []struct {
		name                                     string
		page, limit, count                       *uint64
		wantDBPage, wantPage, wantLimit, wantMax uint64
	}{
		{"no limit uses default", ptr(1), nil, ptr(25), 0, 1, 10, 3},
		{"zero limit uses default", ptr(1), ptr(0), ptr(25), 0, 1, 10, 3},
		{"limit capped at 100", ptr(1), ptr(500), ptr(250), 0, 1, 100, 3},
		{"explicit limit", ptr(2), ptr(5), ptr(12), 1, 2, 5, 3},
		{"no page is first page", nil, ptr(5), ptr(12), 0, 1, 5, 3},
		{"no count", ptr(1), ptr(5), nil, 0, 1, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbPage, page, limit, maxPage := Pager(tt.page, tt.limit, tt.count)
			if dbPage != tt.wantDBPage || page != tt.wantPage || limit != tt.wantLimit || maxPage != tt.wantMax {
				t.Errorf("Pager() = (%d, %d, %d, %d), want (%d, %d, %d, %d)",
					dbPage, page, limit, maxPage, tt.wantDBPage, tt.wantPage, tt.wantLimit, tt.wantMax)
			}
		})
	}
}
