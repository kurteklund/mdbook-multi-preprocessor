package main

import (
	"regexp"
)

const suitSpade = `<span class="suit-s">♠</span>`
const suitHeart = `<span class="suit-h">♥</span>`
const suitDiamond = `<span class="suit-d">♦</span>`
const suitClub = `<span class="suit-c">♣</span>`

var spadeRegExp = regexp.MustCompile(`!s`)
var heartRegExp = regexp.MustCompile(`!h`)
var diamondRegExp = regexp.MustCompile(`(!d)|(!r)`) // Swedish, r => ruter
var clubRegExp = regexp.MustCompile(`(!c)|(!k)`)    // Swedish, k => klöver

func processSections(bookSections *MdBookTopItem) {
	for i := range bookSections.Sections {
		var section = &bookSections.Sections[i]
		processSection(section)
	}
}

func processSection(section *MdBookSection) {
	processChapter(&section.Chapter)
	for i := range section.Chapter.SubItems {
		subSection := &section.Chapter.SubItems[i]
		processSection(subSection)
	}
}

func processChapter(chapter *MdBookChapter) {
	chapter.Content = replaceCardSuitStrings(chapter.Content)
}

func replaceCardSuitStrings(text string) string {
	text = spadeRegExp.ReplaceAllString(text, suitSpade)
	text = heartRegExp.ReplaceAllString(text, suitHeart)
	text = diamondRegExp.ReplaceAllString(text, suitDiamond)
	text = clubRegExp.ReplaceAllString(text, suitClub)

	return text
}
