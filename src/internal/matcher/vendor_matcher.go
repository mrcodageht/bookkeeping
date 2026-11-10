package matcher

import (
	"bytes"
	"strings"

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

				if !(bytes.EqualFold(begin, []byte{' '}) && isFinishMatch) {
					continue
				}

				t.SetVendor(v.Name)
				t.SetCategory(v.Category)
				t.SetStatus(model.Matched)
				isMatch = true
				break
			}
		}
		if !isMatch {
			t.SetVendor(unknown.Name)
			t.SetCategory(unknown.Category)
			t.SetStatus(model.Unknown)
		}
	}
}
