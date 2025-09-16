package parser

import (
	"go-doc-parser/internal/entity"
	"strings"
	"unicode"
)

type NameParser struct {
	In       []rune
	Position int
}

func (p *NameParser) MatchChar(options ...string) string {
	// not needed if no recursion
	// because this is already covered by
	// the for loop condition
	// ???
	if len(p.In) <= p.Position {
		return ""
	}

	suitable := map[string]bool{}

	for _, o := range options {
		suitable[o] = true
	}

	longest := []rune{}

	for i := 0; i < len(p.In[p.Position:]); i++ {
		for key := range suitable {
			// skip already marked
			if !suitable[key] {
				continue
			}

			// we have gone past the option length
			// skip, do not mark
			option := []rune(key)
			if i >= len(option) {
				continue
			}

			// has wrong characters (case ignored), mark and skip
			if unicode.ToUpper(option[i]) != unicode.ToUpper(p.In[p.Position+i]) {
				suitable[key] = false
				continue
			}

			// continue if this is not the option's last character
			if i < len(option)-1 {
				continue
			}

			// use the option
			// the mark remains true from now on

			if len(option) > len(longest) {
				longest = option
			}
		}
	}

	// fmt.Printf("suitable: %#v\n", suitable)
	// fmt.Printf("longest: %s\n", string(longest))

	p.Position += len(longest)

	return string(longest)
}

func (p *NameParser) SkipSpace() {
	if len(p.In) <= p.Position {
		return
	}

	if unicode.IsSpace(rune(p.In[p.Position])) {
		p.Position++
		p.SkipSpace()
	}
}

func (p *NameParser) trimQuotes() (out []rune) {
	if len(p.In) <= p.Position {
		return
	}

	r := p.In[p.Position]

	// fmt.Printf("%c %t\n", r, isQuotation)

	if unicode.Is(unicode.Quotation_Mark, r) {
		p.Position++
	}

	out = p.In[p.Position:len(p.In)]

	r = out[len(out)-1]

	// assume the whole input ends in quote
	if unicode.Is(unicode.Quotation_Mark, r) {
		out = out[:len(out)-1]
	}

	p.Position = len(p.In)

	return out
}

func (p *NameParser) ParseName() entity.ShortID {
	// t = p.MatchChar("віпс", "впс")

	var t = []string{}

	for {
		matched := p.MatchChar("віпс", "впс", "ГОРВ", "ПОРВ", "ВОПР та ПБПС", "ВАЗ", "ВАК", "УОРД ПдРУ", "ВБТЗ", "2", "ПРИКЗ")

		if len(matched) > 0 {
			t = append(t, matched)
		}

		// fmt.Printf("%s %d %c\n", t, p.Position, p.In[p.Position])

		// fmt.Println(len(t), t)

		p.SkipSpace()

		if len(matched) < 1 {
			break
		}
	}

	// if len(t) < 1 {
	// 	return entity.ShortID{"", string(p.In[p.Position:])}
	// }

	// fmt.Printf("%d %c\n", p.Position, p.In[p.Position])

	// n := p.filterSymbols([]rune{})

	name := p.trimQuotes()

	// fmt.Printf("%s %d %c\n", string(n), p.Position, p.In[p.Position-1])

	return entity.ShortID{strings.Join(t, " "), string(name)}
}
