package tasks

import "testing"

func TestValidateTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{name: "valid", title: "Buy milk"},
		{name: "empty", title: "", wantErr: true},
		{name: "blank", title: "   ", wantErr: true},
		{name: "too long", title: string(make([]byte, maxTitleLength+1)), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateTitle(tt.title)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestNormalizeLimitOffset(t *testing.T) {
	t.Parallel()

	limit, offset := normalizeLimitOffset(0, -10)
	if limit != defaultLimit {
		t.Fatalf("expected default limit %d, got %d", defaultLimit, limit)
	}
	if offset != 0 {
		t.Fatalf("expected offset 0, got %d", offset)
	}

	limit, offset = normalizeLimitOffset(500, 20)
	if limit != maxLimit {
		t.Fatalf("expected max limit %d, got %d", maxLimit, limit)
	}
	if offset != 20 {
		t.Fatalf("expected offset 20, got %d", offset)
	}
}
