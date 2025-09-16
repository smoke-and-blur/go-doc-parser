package parser

import (
	"fmt"
	"testing"
)

func TestMatcher(t *testing.T) {
	data := []string{
		"впс",
		"впс Кодима",
		"віпс",
		"віпс Шершенці",
		"ГОРВ ВАЗ",
		"ГОРВ",
		"ГОРВ 123",
		"ГОРВ ВАЗ 123",
	}

	for _, item := range data {
		p := NameParser{
			In: []rune(item),
		}

		out := p.MatchChar("впс", "віпс", "ГОРВ", "ГОРВ ВАЗ")

		t.Log(out)
	}
}

func TestSkipper(t *testing.T) {
	p := NameParser{
		In: []rune("  \t\n\rHello"),
	}

	p.SkipSpace()

	t.Logf("%c", p.In[p.Position])
}

var parsingData = []string{
	"впс «Кодима»",
	"віпс «Загнітків»",
	"віпс «Шершенці»",
	"впс «Станіславка»",
	"віпс «Тимкове»",
	"віпс «Чорна»",
	"впс «Окни»",
	"віпс «Ткаченкове»",
	"віпс «Гулянка»",
	"віпс «Новосеменівка»",
	"впс «Великокомарівка»",
	"віпс «Павлівка»",
	"впс «Велика Михайлівка»",
	"віпс «Слов’яносербка»",
	"віпс «Гребеники»",
	"впс «Степанівка»",
	"ВПС «Степанівка»",
	"віпс «Лучинське»",
	"віпс «Кучурган»",
	"віпс «Лиманське»",
	"ВОПР та ПБПС 123",
	"ГОРВ ВАК",
	"ГОРВ ВАЗ 123",
	"ГОРВ",
	"впс",
	"ГОРВ ВАК 2ПРИКЗ 2 ПРИКЗ \"Тест\"",
}

func TestParseName(t *testing.T) {
	for _, name := range parsingData {
		p := NameParser{
			In: []rune(name),
		}

		q := p.ParseName()

		_ = q

		fmt.Printf("%36.36q %q\n", q.Type, q.Name)
	}

	p := NameParser{
		In: []rune("   івивф  \tш 's'f'f asdb josoj oo o "),
	}

	q := p.ParseName()

	fmt.Printf("%6.6q %q\n", q.Type, q.Name)
}
