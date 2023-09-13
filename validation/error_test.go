package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name     string
		Message  string
		Source   error
		want     string
		wantUser string
	}{
		{"Message Only", "public", nil, "public", "public"},
		{"Err Only", "", errors.New("private"), "private", "private"},
		{"Both", "public", errors.New("private"), "public: private", "public"},
		{"Neither", "", nil, "", ""},
		{"Nest", "public", Error{"public2", errors.New("PRIVATE")}, "public: public2: PRIVATE", "public"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Message: tt.Message,
				Source:  tt.Source,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
			assert.Equalf(t, tt.wantUser, e.UserError(), "Error()")

			if tt.Source != nil {
				assert.ErrorIsf(t, e, tt.Source, "Error()")
			}
		})
	}
}
