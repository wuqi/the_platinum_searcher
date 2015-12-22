package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func newEncodeReader(r io.Reader, encoding int) io.Reader {
	switch encoding {
	case EUCJP:
		return transform.NewReader(r, japanese.EUCJP.NewDecoder())
	case SHIFTJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	case GBK:
		return transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	case EUCKR:
		return transform.NewReader(r, korean.EUCKR.NewDecoder())
	case ISO2022JP:
		return transform.NewReader(r, japanese.ISO2022JP.NewDecoder())
	case BIG5:
		return transform.NewReader(r, traditionalchinese.Big5.NewDecoder())
	case UTF16BE:
		return transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewDecoder())
	case UTF16LE:
		return transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder())
	}
	return nil
}
