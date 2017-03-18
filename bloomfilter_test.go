package bloomfilter

import "testing"

func TestInsert(t *testing.T) {
	vals := []string{
		"9192901",
		"oakdoiaskdok",
		"asdkmaskdmas",
		"9jiasdiao",
		"ads000000d",
		"adiaid109390120 0109 902",
		"oijvijvioijˆˆøåß∂¬…0a9sdkao",
		"¨ßˆ∂ªº0qkomamc;m",
		"xxxx",
		"",
	}

	f := New(1000)

	for i := range vals {
		f.Insert([]byte(vals[i]))

		if !f.Has([]byte(vals[i])) {
			t.Error("Bloom filter is not working!")
		}
	}
}
