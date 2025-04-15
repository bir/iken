package currency_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/currency"
)

func TestFixedFormatter_Format(t *testing.T) {
	tests := []struct {
		name        string
		symbol      string
		precision   int
		factor      int
		thousandSep string
		decimalSep  string
		format      string
		amount      int64
		want        string
	}{
		{"default format", "USD", 2, 100, ",", ".", "", 123456, "USD1,234.56"},
		{"basic dollar", "$", 2, 100, ",", ".", "%s%v", 123460, "$1,234.60"},
		{"basic dollar", "$", 2, 100, ",", ".", "%s%v", 123406, "$1,234.06"},
		{"no thousands dollar", "$", 2, 100, "", ".", "%s%v", 123406, "$1234.06"},
		{"basic euro", "€", 2, 100, ".", ",", "%v%s", 123406, "1.234,06€"},
		{"long euro", "€", 2, 100, ".", ",", "%v%s", 123456789012345, "1.234.567.890.123,45€"},
		{"no precision", "", 0, 1, ",", "", "%v%s", 123456789012345, "123,456,789,012,345"},
		{"high precision", "$", 3, 1000, ",", ".", "%s%v", 1234060, "$1,234.060"},
		{"over precision", "$", 3, 100, ",", ".", "%s%v", 123406, "$1,234.060"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := currency.NewFixedFormatter(tt.symbol, tt.precision, tt.factor, tt.thousandSep, tt.decimalSep, tt.format)
			if got := cf.Format(tt.amount); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPower10(t *testing.T) {
	assert.Equal(t, currency.Power10(0), 1)
	assert.Equal(t, currency.Power10(1), 10)
}
