package minify

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/tdewolff/test"
)

func TestMediatype(t *testing.T) {
	mediatypeTests := []struct {
		mediatype string
		expected  string
	}{
		{"text/html", "text/html"},
		{"text/html; charset=UTF-8", "text/html;charset=utf-8"},
		{"text/html; charset=UTF-8 ; param = \" ; \"", "text/html;charset=utf-8;param=\" ; \""},
		{"text/html, text/css", "text/html,text/css"},
	}
	for _, tt := range mediatypeTests {
		t.Run(tt.mediatype, func(t *testing.T) {
			mediatype := Mediatype([]byte(tt.mediatype))
			test.Minify(t, tt.mediatype, nil, string(mediatype), tt.expected)
		})
	}
}

func TestDataURI(t *testing.T) {
	dataURITests := []struct {
		dataURI  string
		expected string
	}{
		{"datx:x", "datx:x"},
		{"data:,text", "data:,text"},
		{"data:text/plain;charset=us-ascii,text", "data:,text"},
		{"data:TEXT/PLAIN;CHARSET=US-ASCII,text", "data:,text"},
		{"data:text/plain;charset=us-asciiz,text", "data:;charset=us-asciiz,text"},
		{"data:;base64,dGV4dA==", "data:,text"},
		{"data:text/svg+xml;base64,IyMjIyMj", "data:text/svg+xml;base64,IyMjIyMj"},
		{"data:text/xml;version=2.0,content", "data:text/xml;version=2.0,content"},
		{"data:text/xml; version = 2.0,content", "data:text/xml;version=2.0,content"},
		{"data:,%23%23%23%23%23", "data:,%23%23%23%23%23"},
		{"data:,%23%23%23%23%23%23", "data:;base64,IyMjIyMj"},
		{"data:text/x,<?xx?>", "data:text/x,%3C?xx?%3E"},
		{"data:text/other,\"<\u2318", "data:text/other,%22%3C%E2%8C%98"},
		{"data:text/other,\"<\u2318>", "data:text/other;base64,IjzijJg+"},
		{`data:text/svg+xml,%3Csvg height="100" width="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" /></svg>`, `data:text/svg+xml,%3Csvg height="100" width="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" /></svg>`},
	}
	m := New()
	m.AddFunc("text/x", func(_ *M, w io.Writer, r io.Reader, _ map[string]string) error {
		b, _ := ioutil.ReadAll(r)
		test.String(t, string(b), "<?xx?>")
		w.Write(b)
		return nil
	})
	for _, tt := range dataURITests {
		t.Run(tt.dataURI, func(t *testing.T) {
			dataURI := DataURI(m, []byte(tt.dataURI))
			test.Minify(t, tt.dataURI, nil, string(dataURI), tt.expected)
		})
	}
}

func TestDecimal(t *testing.T) {
	numberTests := []struct {
		number   string
		expected string
	}{
		{"", ""},
		{"0", "0"},
		{".0", "0"},
		{"1.0", "1"},
		{"0.1", ".1"},
		{"+1", "1"},
		{"-1", "-1"},
		{"-0.1", "-.1"},
		{"10", "10"},
		{"100", "100"},
		{"1000", "1000"},
		{"0.001", ".001"},
		{"0.0001", ".0001"},
		{"0.252", ".252"},
		{"1.252", "1.252"},
		{"-1.252", "-1.252"},
		{"0.075", ".075"},
		{"789012345678901234567890123456789e9234567890123456789", "789012345678901234567890123456789e9234567890123456789"},
		{".000100009", ".000100009"},
		{".0001000009", ".0001000009"},
		{".0001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009", ".0001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009"},
		{"E\x1f", "E\x1f"}, // fuzz
	}
	for _, tt := range numberTests {
		t.Run(tt.number, func(t *testing.T) {
			number := Decimal([]byte(tt.number), -1)
			test.Minify(t, tt.number, nil, string(number), tt.expected)
		})
	}
}

func TestDecimalTruncate(t *testing.T) {
	numberTests := []struct {
		number   string
		truncate int
		expected string
	}{
		{"0.1", 1, ".1"},
		{"0.0001", 1, ".0001"},
		{"0.111", 1, ".1"},
		{"0.111", 0, ".111"},
		{"1.111", 1, "1"},
		{"0.075", 1, ".08"},
		{"0.025", 1, ".03"},
		{"0.105", 2, ".11"},
		{"0.104", 2, ".1"},
		{"9.99", 2, "10"},
		{"9.99", 1, "10"},
		{"8.88", 2, "8.9"},
		{"8.88", 1, "9"},
		{"8.00", 1, "8"},
		{".88", 1, ".9"},
		{"1.234", 2, "1.2"},
		{"33.33", 2, "33"},
		{"29.666", 2, "30"},
		{"1.51", 2, "1.5"},
		{"1.01", 2, "1"},
		{".99", 1, "1"},
		{"-16.400000000000006", 3, "-16.4"}, // #233
		{"1.00000000000001", 15, "1.00000000000001"},
		{"1.000000000000001", 15, "1"},
		{"1.000000000000009", 15, "1.00000000000001"},
		{"100000000000009", 15, "100000000000009"},
		{"1000000000000009", 15, "1000000000000009"},
		{"10000000000000009", 15, "10000000000000009"},
		{"139.99999999", 8, "140"},
	}
	for _, tt := range numberTests {
		t.Run(tt.number, func(t *testing.T) {
			number := Decimal([]byte(tt.number), tt.truncate)
			test.Minify(t, tt.number, nil, string(number), tt.expected, "truncate to", tt.truncate)
		})
	}
}

func TestNumber(t *testing.T) {
	numberTests := []struct {
		number   string
		expected string
	}{
		{"", ""},
		{"0", "0"},
		{".0", "0"},
		{"1.0", "1"},
		{"0.1", ".1"},
		{"+1", "1"},
		{"-1", "-1"},
		{"-0.1", "-.1"},
		{"10", "10"},
		{"100", "100"},
		{"1000", "1e3"},
		{"0.001", ".001"},
		{"0.0001", "1e-4"},
		{"100e1", "1e3"},
		{"1e10", "1e10"},
		{"1e-10", "1e-10"},
		{"1000e-7", "1e-4"},
		{"1000e-6", ".001"},
		{"1.1e+1", "11"},
		{"1.1e-1", ".11"},
		{"1.1e6", "11e5"},
		{"1.1e", "1.1e"},   // broken number, don't parse
		{"1.1e+", "1.1e+"}, // broken number, don't parse
		{"0.252", ".252"},
		{"1.252", "1.252"},
		{"-1.252", "-1.252"},
		{"0.075", ".075"},
		{"789012345678901234567890123456789e9234567890123456789", "789012345678901234567890123456789e9234567890123456789"},
		{".000100009", "100009e-9"},
		{".0001000009", ".0001000009"},
		{".0001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009", ".0001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009"},
		{".6000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003e-9", ".6000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003e-9"},
		{"E\x1f", "E\x1f"}, // fuzz
		{"1e9223372036854775807", "1e9223372036854775807"},
		{"11e9223372036854775807", "11e9223372036854775807"},
		{".01e-9223372036854775808", ".01e-9223372036854775808"},
		{".011e-9223372036854775808", ".011e-9223372036854775808"},

		{".12345e8", "12345e3"},
		{".12345e7", "1234500"},
		{".12345e6", "123450"},
		{".12345e5", "12345"},
		{".012345e6", "12345"},
		{".12345e4", "1234.5"},
		{"-.12345e4", "-1234.5"},
		{".12345e0", ".12345"},
		{".12345e-1", ".012345"},
		{".12345e-2", ".0012345"},
		{".12345e-3", "12345e-8"},
		{".12345e-4", "12345e-9"},
		{".12345e-5", ".12345e-5"},

		{".123456e-3", "123456e-9"},
		{".123456e-2", ".00123456"},
		{".1234567e-4", ".1234567e-4"},
		{".1234567e-3", ".0001234567"},

		{"12345678e-1", "1234567.8"},
		{"72.e-3", ".072"},
		{"7640e-2", "76.4"},
		{"10.e-3", ".01"},
		{".0319e3", "31.9"},
		{"39.7e-2", ".397"},
		{"39.7e-3", ".0397"},
		{".01e1", ".1"},
		{".001e1", ".01"},
		{"39.7e-5", "397e-6"},
	}
	for _, tt := range numberTests {
		t.Run(tt.number, func(t *testing.T) {
			number := Number([]byte(tt.number), -1)
			test.Minify(t, tt.number, nil, string(number), tt.expected)
		})
	}
}

func TestNumberTruncate(t *testing.T) {
	numberTests := []struct {
		number   string
		truncate int
		expected string
	}{
		{"0.1", 1, ".1"},
		{"0.01", 1, ".01"},
		{"0.001", 1, ".001"},
		{"0.0001", 1, "1e-4"},
		{"1000", 0, "1e3"},
		{"1234", 0, "1234"},
		{"0.111", 1, ".1"},
		{"0.111", 0, ".111"},
		{"0.075", 1, ".08"},
		{"0.025", 1, ".03"},
		{"0.105", 2, ".11"},
		{"0.104", 2, ".1"},
		{"9.99", 2, "10"},
		{"9.99", 1, "10"},
		{"99", 1, "99"},
		{"999", 1, "1e3"},
		{"99e1", 1, "1e3"},
		{"99.9", 1, "100"},
		{"999.9", 1, "1e3"},
		{"999.99", 1, "1e3"},
		{"111.99", 4, "112"},
		{"8.88", 2, "8.9"},
		{"8.88", 1, "9"},
		{"8.00", 1, "8"},
		{".88", 1, ".9"},
		{"1.234", 2, "1.2"},
		{"33.33", 2, "33"},
		{"29.666", 2, "30"},
		{"1.51", 2, "1.5"},
		{"1.51", 1, "2"},
		{"1.01", 2, "1"},
		{"1.01", 3, "1.01"},
		{"1.01", 4, "1.01"},
		{".99", 1, "1"},
		{"-16.400000000000006", 3, "-16.4"}, // #233
		{"1.00000000000001", 15, "1.00000000000001"},
		{"1.000000000000001", 15, "1"},
		{"1.000000000000009", 15, "1.00000000000001"},
		{"100000000000009", 15, "100000000000009"},
		{"1000000000000009", 15, "1000000000000009"},
		{"10000000000000009", 15, "1e16"},
		{"0.0000100000000000009", 15, ".100000000000009e-4"},
		{".000333336", 0, "333336e-9"},
		{".0003333337", 0, ".0003333337"},
		{".000033335", 0, "33335e-9"},
		{".0000333336", 0, ".333336e-4"},
		{".00003333337", 0, ".3333337e-4"},
		{".0000000003333337", 0, ".3333337e-9"},
		{".00000000003333338", 0, "3333338e-17"},
		{".333336e-3", 0, "333336e-9"},
		{".3333337e-3", 0, ".0003333337"},
		{".33335e-4", 0, "33335e-9"},
		{".333336e-4", 0, ".333336e-4"},
		{".3333337e-4", 0, ".3333337e-4"},
		{".000033333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333395", 0, ".33333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333395e-4"},
		{".0000333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333396", 0, ".333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333396e-4"},
		{".00003333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397", 0, ".3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-4"},
		{".0000033333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333395", 0, ".33333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333395e-5"},
		{".33333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333395e-4", 0, ".33333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333395e-4"},
		{".333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333396e-4", 0, ".333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333396e-4"},
		{".3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-4", 0, ".3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-4"},
		{".3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-902", 0, "3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-999"},
		{".3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-903", 0, ".3333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333397e-903"},
		{"29.629775e-9", 0, ".29629775e-7"},
		{"e-9223372036854775808", 0, "e-9223372036854775808"},
		{"139.99999999", 8, "140"},
	}
	for _, tt := range numberTests {
		t.Run(tt.number, func(t *testing.T) {
			number := Number([]byte(tt.number), tt.truncate)
			test.Minify(t, tt.number, nil, string(number), tt.expected, "truncate to", tt.truncate)
		})
	}
}

func TestDecimalRandom(t *testing.T) {
	N := int(1e4)
	if testing.Short() {
		N = 0
	}
	for i := 0; i < N; i++ {
		b := RandNumBytes(false)
		f, _ := strconv.ParseFloat(string(b), 64)

		b2 := make([]byte, len(b))
		copy(b2, b)
		b2 = Decimal(b2, -1)
		f2, _ := strconv.ParseFloat(string(b2), 64)
		if math.Abs(f-f2) > 1e-6 {
			fmt.Println("Bad:", f, "!=", f2, "in", string(b), "to", string(b2))
		}
	}
}

func TestNumberRandom(t *testing.T) {
	N := int(1e4)
	if testing.Short() {
		N = 0
	}
	for i := 0; i < N; i++ {
		b := RandNumBytes(true)
		f, _ := strconv.ParseFloat(string(b), 64)

		b2 := make([]byte, len(b))
		copy(b2, b)
		b2 = Number(b2, -1)
		f2, _ := strconv.ParseFloat(string(b2), 64)
		if math.Abs(f-f2) > 1e-6 {
			fmt.Println("Bad:", f, "!=", f2, "in", string(b), "to", string(b2))
		}
	}
}

////////////////

var n = 100
var numbers [][]byte

func TestMain(t *testing.T) {
	numbers = make([][]byte, 0, n)
	for j := 0; j < n; j++ {
		numbers = append(numbers, RandNumBytes(true))
	}
}

func RandNumBytes(withExp bool) []byte {
	var b []byte
	n := rand.Int() % 10
	for i := 0; i < n; i++ {
		b = append(b, byte(rand.Int()%10)+'0')
	}
	if rand.Int()%2 == 0 {
		b = append(b, '.')
		n = rand.Int() % 10
		for i := 0; i < n; i++ {
			b = append(b, byte(rand.Int()%10)+'0')
		}
	}
	if withExp && rand.Int()%2 == 0 {
		b = append(b, 'e')
		if rand.Int()%2 == 0 {
			b = append(b, '-')
		}
		n = 1 + rand.Int()%4
		for i := 0; i < n; i++ {
			b = append(b, byte(rand.Int()%10)+'0')
		}
	}
	return b
}

func BenchmarkNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			Number(numbers[j], -1)
		}
	}
}

func BenchmarkNumber2(b *testing.B) {
	num := []byte("1.2345e-6")
	for i := 0; i < b.N; i++ {
		Number(num, -1)
	}
}
