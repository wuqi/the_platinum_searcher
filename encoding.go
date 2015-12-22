package the_platinum_searcher

import (
	"github.com/gogits/chardet"
)

const (
	UNKNOWN = iota
	ERROR
	BINARY
	ASCII
	UTF8
	EUCJP
	SHIFTJIS
	GBK
	EUCKR
	BIG5
	ISO2022JP
	UTF16BE
	UTF16LE
)

func detectEncoding(bs []byte) int {

	var (
		suspiciousBytes = 0
		likelyUtf8      = 0
		likelyEucjp     = 0
		likelyShiftjis  = 0
	)
	var detector = chardet.NewTextDetector()
	length := len(bs)

	if length == 0 {
		return ASCII
	}

	if length >= 3 && bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		// UTF-8 BOM. This isn't binary.
		return UTF8
	}

	if length >= 5 && bs[0] == 0x25 && bs[1] == 0x50 && bs[2] == 0x44 && bs[3] == 0x46 && bs[4] == 0x2D {
		/*  %PDF-. This is binary. */
		return BINARY
	}

	result, err := detector.DetectBest(bs)
	if err != nil {
		return ERROR
	}

	if result.Charset == "UTF-8" {
		return UTF8
	} else if result.Charset == "GB18030" {
		return GBK
	} else if result.Charset == "EUC-JP" {
		return EUCJP
	} else if result.Charset == "Shift_JIS" {
		return SHIFTJIS
	} else if result.Charset == "EUC-KR" {
		return EUCKR
	} else if result.Charset == "Big5" {
		return BIG5
	} else if result.Charset == "ISO-2022-JP" {
		return ISO2022JP
	} else if result.Charset == "UTF-16BE" {
		return UTF16BE
	} else if result.Charset == "UTF-16LE" {
		return UTF16LE
	} else {
		for i := 0; i < length; i++ {
			if bs[i] == 0x00 {
				/* NULL char. It's binary */
				return BINARY
			} else if (bs[i] < 7 || bs[i] > 14) && (bs[i] < 32 || bs[i] > 127) {
				/* UTF-8 detection */
				if bs[i] > 193 && bs[i] < 224 && i+1 < length {
					i++
					if bs[i] > 127 && bs[i] < 192 {
						likelyUtf8++
						continue
					}

				} else if bs[i] > 223 && bs[i] < 240 && i+2 < length {
					i++
					if bs[i] > 127 && bs[i] < 192 && bs[i+1] > 127 && bs[i+1] < 192 {
						i++
						likelyUtf8++
						continue
					}
				}

				/* EUC-JP detection */
				if bs[i] == 142 && i+1 < length {
					i++
					if bs[i] > 160 && bs[i] < 224 {
						likelyEucjp++
						continue
					}
				} else if bs[i] > 160 && bs[i] < 255 && i+1 < length {
					i++
					if bs[i] > 160 && bs[i] < 255 {
						likelyEucjp++
						continue
					}
				}

				/* Shift-JIS detection */
				if bs[i] > 160 && bs[i] < 224 {
					likelyShiftjis++
					continue
				} else if ((bs[i] > 128 && bs[i] < 160) || (bs[i] > 223 && bs[i] < 240)) && i+1 < length {
					i++
					if (bs[i] > 63 && bs[i] < 127) || (bs[i] > 127 && bs[i] < 253) {
						likelyShiftjis++
						continue
					}
				}

				suspiciousBytes++
				if i >= 32 && (suspiciousBytes*100)/length > 10 {
					return BINARY
				}

			}
		}
	}
	if (suspiciousBytes*100)/length > 10 {
		return BINARY
	}

	// fmt.Printf("Detected points[utf8/eucjp/shiftjis] is %d/%d/%d.\n", likelyUtf8, likelyEucjp, likelyShiftjis)

	if likelyUtf8 == 0 && likelyEucjp == 0 && likelyShiftjis == 0 {
		return ASCII
	} else if likelyUtf8 >= likelyEucjp && likelyUtf8 >= likelyShiftjis {
		return UTF8
	} else if likelyEucjp >= likelyUtf8 && likelyEucjp >= likelyShiftjis {
		return EUCJP
	} else if likelyShiftjis >= likelyUtf8 && likelyShiftjis >= likelyEucjp {
		return SHIFTJIS
	}

	return ASCII

}
