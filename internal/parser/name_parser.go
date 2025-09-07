package parser

import (
	"go-doc-parser/internal/entity"
	"unicode"
)

type NameParser struct {
	In       []rune
	Position int
}

// narrow down options character by character
func (p *NameParser) matchChar(options [][]rune, i int) [][]rune {
	if len(options) == 1 && len(options[0]) == i { // reached the end of the word?
		p.Position += i
		return options
	}

	out := [][]rune{}

	if p.Position+i >= len(p.In) {
		return nil
	}

	for _, option := range options {
		if i >= len(option) {
			continue
		}

		if p.In[p.Position+i] == option[i] {
			out = append(out, option)
		}
	}

	if len(out) < 1 {
		return nil
	}

	i++

	return p.matchChar(out, i)
}

func (p *NameParser) MatchChar(options ...string) string {
	runeOptions := [][]rune{}

	for _, option := range options {
		runeOptions = append(runeOptions, []rune(option))
	}

	out := p.matchChar(runeOptions, 0)

	if len(out) > 0 {
		return string(out[0])
	}

	// longest := ""

	// for _, option := range out {
	// 	if len(option) > len(longest) {
	// 		longest = string(option)
	// 	}
	// }

	// p.Position += len(longest)

	// return longest

	return ""
}

func (p *NameParser) SkipSpace() {
	if len(p.In) == p.Position {
		return
	}

	if unicode.IsSpace(rune(p.In[p.Position])) {
		p.Position++
		p.SkipSpace()
	}
}

func (p *NameParser) trimQuotes() (out []rune) {
	if len(p.In) == p.Position {
		return
	}

	r := p.In[p.Position]

	// fmt.Printf("%c %t\n", r, isQuotation)

	if unicode.Is(unicode.Quotation_Mark, r) {
		p.Position++
	}

	out = p.In[p.Position:len(p.In)]

	r = out[len(out)-1]

	if unicode.Is(unicode.Quotation_Mark, r) {
		out = out[:len(out)-1]
	}

	p.Position = len(p.In)

	return out
}

func (p *NameParser) ParseName() entity.ShortID {
	// t := p.MatchChar("віпс", "впс", "ГОРВ", "ПОРВ", "ВОПР та ПБПС", "ГОРВ ВАЗ") //
	t := p.MatchChar("віпс", "впс")

	// fmt.Printf("%s %d %c\n", t, p.Position, p.In[p.Position])

	p.SkipSpace()
	// fmt.Printf("%d %c\n", p.Position, p.In[p.Position])

	// n := p.filterSymbols([]rune{})

	name := p.trimQuotes()

	// fmt.Printf("%s %d %c\n", string(n), p.Position, p.In[p.Position-1])

	return entity.ShortID{t, string(name)}
}
