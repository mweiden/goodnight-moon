package text_metrics

import (
	"regexp"
	"strings"
	"unicode"
)

/*
 * Packages here are copied for security reasons
 */

//https://github.com/BluntSporks/readability/blob/master/fk.go
// CntWords counts the number of words in a text by counting the spaces.
func CntWords(text string) int {
	cnt := 0
	wasSpace := false
	text = strings.TrimSpace(text) + " "
	for _, char := range text {
		if unicode.IsSpace(char) {
			if !wasSpace {
				cnt++
				wasSpace = true
			}
		} else {
			wasSpace = false
		}
	}
	return cnt
}

// CntSents counts the number of sentences in a text by counting ending marks.
func CntSents(text string) int {
	cnt := 0
	wasEnd := false
	text = strings.TrimSpace(text) + " "
	for _, char := range text {
		if char == '.' || char == '?' || char == '!' {
			wasEnd = true
		} else if unicode.IsLetter(char) {
			wasEnd = false
		} else if wasEnd && unicode.IsSpace(char) {
			cnt++
			wasEnd = false
		}
	}
	return cnt
}

func FleschKincaid(text string) (float64, float64) {
	syllableCount := float64(SyllablesIn(text))
	wordCnt := float64(CntWords(text))
	sentCnt := float64(CntSents(text))
	gradeLevel := 0.39*(wordCnt/sentCnt) + 11.8*(syllableCount/wordCnt) - 15.59
	score := 206.835 - 1.015 * (wordCnt / sentCnt) - 84.6 * (syllableCount / wordCnt)
	return score, gradeLevel
}

// https://github.com/mtso/syllables/blob/master/count.go

// Returns the integer count of syllables in the input string.
type counter struct {
	count, index, length int
	singular             string
	parts                []string
}

// Returns the integer count of syllables in the input byte array.
func InBytes(b []byte) int {
	s := string(b[:len(b)])
	return SyllablesIn(s)
}

func SyllablesIn(text string) int {

	// Prepare input text by converting to lowercase
	// and removing all non-alphabetic runes
	text = strings.ToLower(text)
	text = expressionNonalphabetic.ReplaceAllString(text, "")

	// Return early when possible
	if len(text) < 1 {
		return 0
	}
	if len(text) < 3 {
		return 1
	}

	// If value is part of cornercases,
	// return hardcoded value
	if syllables, ok := cornercases[text]; ok {
		return syllables
	}

	// Initialize counter
	c := counter{}

	// Count and remove matched prefixes and suffixes
	text = expressionTriple.ReplaceAllStringFunc(text, c.countAndRemove(3))
	text = expressionDouble.ReplaceAllStringFunc(text, c.countAndRemove(2))
	text = expressionSingle.ReplaceAllStringFunc(text, c.countAndRemove(1))

	// Count multiple consanants
	c.parts = consanants.Split(text, -1)
	c.index = 0
	c.length = len(c.parts)

	for ; c.index < c.length; c.index++ {
		if c.parts[c.index] != "" {
			c.count++
		}
	}

	// Subtract one for maches which should be
	// counted as one but are counted as two
	subtractOne := c.countInPlace(-1)
	expressionMonosyllabicOne.ReplaceAllStringFunc(text, subtractOne)
	expressionMonosyllabicTwo.ReplaceAllStringFunc(text, subtractOne)

	// Add one for maches which should be
	// counted as two but are counted as one
	addOne := c.countInPlace(1)
	expressionDoubleSyllabicOne.ReplaceAllStringFunc(text, addOne)
	expressionDoubleSyllabicTwo.ReplaceAllStringFunc(text, addOne)
	expressionDoubleSyllabicThree.ReplaceAllStringFunc(text, addOne)
	expressionDoubleSyllabicFour.ReplaceAllStringFunc(text, addOne)

	if c.count < 1 {
		return 1
	}
	return c.count
}

func (c *counter) countAndRemove(increment int) func(string) string {
	return func(in string) string {
		c.count += increment
		return ""
	}
}

func (c *counter) countInPlace(increment int) func(string) string {
	return func(in string) string {
		c.count += increment
		return in
	}
}

var cornercases = map[string]int{
	"abalone":     4,
	"abare":       3,
	"abed":        2,
	"abruzzese":   4,
	"abbruzzese":  4,
	"aborigine":   5,
	"acreage":     3,
	"adame":       3,
	"adieu":       2,
	"adobe":       3,
	"anemone":     4,
	"apache":      3,
	"aphrodite":   4,
	"apostrophe":  4,
	"ariadne":     4,
	"cafe":        2,
	"calliope":    4,
	"catastrophe": 4,
	"chile":       2,
	"chloe":       2,
	"circe":       2,
	"coyote":      3,
	"epitome":     4,
	"facsimile":   4,
	"forever":     3,
	"gethsemane":  4,
	"guacamole":   4,
	"hyperbole":   4,
	"jesse":       2,
	"jukebox":     2,
	"karate":      3,
	"machete":     3,
	"maybe":       2,
	"people":      2,
	"recipe":      3,
	"sesame":      3,
	"shoreline":   2,
	"simile":      3,
	"syncope":     3,
	"tamale":      3,
	"yosemite":    4,
	"daphne":      2,
	"eurydice":    4,
	"euterpe":     3,
	"hermione":    4,
	"penelope":    4,
	"persephone":  4,
	"phoebe":      2,
	"zoe":         2,
}

var (
	expressionMonosyllabicOne = regexp.MustCompile(
		"cia(l|$)|" +
			"tia|" +
			"cius|" +
			"cious|" +
			"[^aeiou]giu|" +
			"[aeiouy][^aeiouy]ion|" +
			"iou|" +
			"sia$|" +
			"eous$|" +
			"[oa]gue$|" +
			".[^aeiuoycgltdb]{2,}ed$|" +
			".ely$|" +
			"^jua|" +
			"uai|" +
			"eau|" +
			"^busi$|" +
			"(" +
			"[aeiouy]" +
			"(" +
			"b|" +
			"c|" +
			"ch|" +
			"dg|" +
			"f|" +
			"g|" +
			"gh|" +
			"gn|" +
			"k|" +
			"l|" +
			"lch|" +
			"ll|" +
			"lv|" +
			"m|" +
			"mm|" +
			"n|" +
			"nc|" +
			"ng|" +
			"nch|" +
			"nn|" +
			"p|" +
			"r|" +
			"rc|" +
			"rn|" +
			"rs|" +
			"rv|" +
			"s|" +
			"sc|" +
			"sk|" +
			"sl|" +
			"squ|" +
			"ss|" +
			"th|" +
			"v|" +
			"y|" +
			"z" +
			")" +
			"ed$" +
			")|" +
			"(" +
			"[aeiouy]" +
			"(" +
			"b|" +
			"ch|" +
			"d|" +
			"f|" +
			"gh|" +
			"gn|" +
			"k|" +
			"l|" +
			"lch|" +
			"ll|" +
			"lv|" +
			"m|" +
			"mm|" +
			"n|" +
			"nch|" +
			"nn|" +
			"p|" +
			"r|" +
			"rn|" +
			"rs|" +
			"rv|" +
			"s|" +
			"sc|" +
			"sk|" +
			"sl|" +
			"squ|" +
			"ss|" +
			"st|" +
			"t|" +
			"th|" +
			"v|" +
			"y" +
			")" +
			"es$" +
			")",
	)

	expressionMonosyllabicTwo = regexp.MustCompile(
		"[aeiouy]" +
			"(" +
			"b|" +
			"c|" +
			"ch|" +
			"d|" +
			"dg|" +
			"f|" +
			"g|" +
			"gh|" +
			"gn|" +
			"k|" +
			"l|" +
			"ll|" +
			"lv|" +
			"m|" +
			"mm|" +
			"n|" +
			"nc|" +
			"ng|" +
			"nn|" +
			"p|" +
			"r|" +
			"rc|" +
			"rn|" +
			"rs|" +
			"rv|" +
			"s|" +
			"sc|" +
			"sk|" +
			"sl|" +
			"squ|" +
			"ss|" +
			"st|" +
			"t|" +
			"th|" +
			"v|" +
			"y|" +
			"z" +
			")" +
			"e$",
	)

	expressionDoubleSyllabicOne = regexp.MustCompile(
		"(" +
			"(" +
			"[^aeiouy]" +
			// Remove unsupported backreference
			// Replace with {1,2} instead (one or two repeated consanants)
			// Will probably need to figure out a better way to do this
			// Original with backreference: `")\\2l|" +`
			"){1,2}l|" +
			"[^aeiouy]ie" +
			"(" +
			"r|" +
			"st|" +
			"t" +
			")|" +
			"[aeiouym]bl|" +
			"eo|" +
			"ism|" +
			"asm|" +
			"thm|" +
			"dnt|" +
			"uity|" +
			"dea|" +
			"gean|" +
			"oa|" +
			"ua|" +
			"eings?|" +
			"[aeiouy]sh?e[rsd]" +
			")$",
	)

	expressionDoubleSyllabicTwo = regexp.MustCompile(
		"[^gq]ua[^auieo]|" +
			"[aeiou]{3}|" +
			"^(" +
			"ia|" +
			"mc|" +
			"coa[dglx]." +
			")",
	)

	expressionDoubleSyllabicThree = regexp.MustCompile(
		"[^aeiou]y[ae]|" +
			"[^l]lien|" +
			"riet|" +
			"dien|" +
			"iu|" +
			"io|" +
			"ii|" +
			"uen|" +
			"real|" +
			"iell|" +
			"eo[^aeiou]|" +
			"[aeiou]y[aeiou]",
	)

	expressionDoubleSyllabicFour = regexp.MustCompile(
		"[^s]ia",
	)

	expressionSingle = regexp.MustCompile(
		"^" +
			"(" +
			"un|" +
			"fore|" +
			"ware|" +
			"none?|" +
			"out|" +
			"post|" +
			"sub|" +
			"pre|" +
			"pro|" +
			"dis|" +
			"side" +
			")" +
			"|" +
			"(" +
			"ly|" +
			"less|" +
			"some|" +
			"ful|" +
			"ers?|" +
			"ness|" +
			"cians?|" +
			"ments?|" +
			"ettes?|" +
			"villes?|" +
			"ships?|" +
			"sides?|" +
			"ports?|" +
			"shires?|" +
			"tion(ed)?" +
			")",
	)

	expressionDouble = regexp.MustCompile(
		"^" +
			"(" +
			"above|" +
			"anti|" +
			"ante|" +
			"counter|" +
			"hyper|" +
			"afore|" +
			"agri|" +
			"infra|" +
			"intra|" +
			"inter|" +
			"over|" +
			"semi|" +
			"ultra|" +
			"under|" +
			"extra|" +
			"dia|" +
			"micro|" +
			"mega|" +
			"kilo|" +
			"pico|" +
			"nano|" +
			"macro" +
			")" +
			"|" +
			"(" +
			"fully|" +
			"berry|" +
			"woman|" +
			"women" +
			")",
	)

	expressionTriple = regexp.MustCompile(
		"(ology|ologist|onomy|onomist)",
	)

	expressionNonalphabetic = regexp.MustCompile(
		"[^a-z]",
	)

	consanants = regexp.MustCompile(
		"[^aeiouy]+",
	)
)
