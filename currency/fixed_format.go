package currency

import (
	"strconv"
	"strings"
)

const (
	defaultFormatPattern = "%s%v"
)

// FixedFormatter holds currency formatting configuration.
type FixedFormatter struct {
	symbol      string
	precision   int
	factor      int
	thousandSep string
	decimalSep  string
	format      string
}

// NewFixedFormatter creates a new currency formatter with the specified configuration.
func NewFixedFormatter(symbol string, precision, factor int, thousandSep, decimalSep, format string) FixedFormatter {
	return FixedFormatter{
		symbol:      symbol,
		precision:   precision,
		factor:      factor,
		thousandSep: thousandSep,
		decimalSep:  decimalSep,
		format:      format,
	}
}

// Format formats a monetary amount based on fixed point integer of the currency.
func (cf FixedFormatter) Format(amount int64) string {
	return cf.formatAmount(amount/int64(cf.factor), amount%int64(cf.factor))
}

// formatAmount handles the common formatting logic for both Format and FormatFixed.
func (cf FixedFormatter) formatAmount(whole, fractional int64) string {
	wholeStr := cf.formatWholePartWithSeparators(whole)
	fractionalStr := cf.formatFractionalPart(fractional)

	return cf.combinePartsAndFormat(wholeStr, fractionalStr)
}

// formatWholePartWithSeparators formats the whole number part with thousand separators.
func (cf FixedFormatter) formatWholePartWithSeparators(wholePart int64) string {
	wholeStr := strconv.FormatInt(wholePart, 10)
	if cf.thousandSep == "" {
		return wholeStr
	}

	var result strings.Builder

	length := len(wholeStr)
	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(cf.thousandSep)
		}

		result.WriteByte(wholeStr[i])
	}

	return result.String()
}

const base = 10

// formatFractionalPart formats the fractional part with proper precision.
func (cf FixedFormatter) formatFractionalPart(fractionalPart int64) string {
	if cf.precision <= 0 {
		return ""
	}

	fractionalStr := strconv.FormatInt(fractionalPart, base)

	offset := Power10(cf.precision) / cf.factor

	for offset > 1 {
		fractionalStr += "0"
		offset /= base
	}

	if len(fractionalStr) < cf.precision {
		return strings.Repeat("0", cf.precision-len(fractionalStr)) + fractionalStr
	}

	return fractionalStr
}

// combinePartsAndFormat combines whole and fractional parts and applies the final formatting.
func (cf FixedFormatter) combinePartsAndFormat(wholeStr, fractionalStr string) string {
	valueStr := wholeStr
	if cf.precision > 0 {
		valueStr += cf.decimalSep + fractionalStr
	}

	format := cf.format
	if format == "" {
		format = defaultFormatPattern
	}

	result := strings.ReplaceAll(format, "%s", cf.symbol)

	return strings.ReplaceAll(result, "%v", valueStr)
}

// Power10 calculates 10 to the mth power.
func Power10(m int) int {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return base
	}

	result := 10
	for i := 2; i <= m; i++ {
		result *= base
	}

	return result
}
