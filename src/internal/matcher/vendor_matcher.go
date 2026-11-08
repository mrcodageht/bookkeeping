package matcher

import (
	"bytes"
	"fmt"
	"strings"

	"mrccommunity.com/bookkeping/src/internal"
	"mrccommunity.com/bookkeping/src/internal/model"
)

func Match(transactions []*model.Transaction, vendors *model.VendorsYAML) {
	unknown := model.UnknownVendor()
	for _, t := range transactions {
		isMatch := false
		for _, v := range vendors.Vendors {
			for _, pattern := range v.Patterns {
				pattern = strings.TrimSpace(pattern)

				patternWithoutAsterix := strings.TrimSuffix(pattern, "*")
				indexSubstr := strings.Index(t.GetDescription(), patternWithoutAsterix)
				if indexSubstr == -1 {
					continue
				}

				patternSize := len(patternWithoutAsterix)

				begin := []byte{t.GetDescription()[indexSubstr-1]}
				patternEnd := []byte{pattern[len(pattern)-1]}
				isOutOfRange := len(t.GetDescription()) <= (indexSubstr + patternSize)

				var finish []byte
				if isOutOfRange {
					finish = []byte{t.GetDescription()[(indexSubstr+patternSize)-1]}
				} else {
					finish = []byte{t.GetDescription()[(indexSubstr + patternSize)]}
				}

				if bytes.EqualFold(begin, []byte{'*'}) {
					continue
				}

				isFinishMatch := bytes.EqualFold(finish, []byte{' '}) || bytes.EqualFold(patternEnd, []byte{'*'}) || isOutOfRange
				if internal.MustVerbose {
					fmt.Println(":: DEBUG :: is finish match:", isFinishMatch)
				}
				if !(bytes.EqualFold(begin, []byte{' '}) && isFinishMatch) {
					continue
				}

				t.Vendor = v.Name
				t.Category = v.Category
				t.Status = model.Matched
				isMatch = true
				break
			}
		}
		if !isMatch {
			t.Vendor = unknown.Name
			t.Category = unknown.Category
			t.Status = model.Unknown
		}
	}
}
