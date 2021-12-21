package schema

import (
	"fmt"
	"regexp"
	"strings"
)

type LinterRule struct {
	when *regexp.Regexp
	then *regexp.Regexp
	message string
}

var activeRules = []LinterRule{{
	when: regexp.MustCompile(`(?is)CREATE[\s]+(FUNCTION|VIEW)[\s]+.*AS`),
	then: regexp.MustCompile(`(?is)CREATE[\s]+(FUNCTION|VIEW)[\s]+.*WITH[\s]+SCHEMABINDING.*AS`),
	message: "Functions and views are only deterministic when defined with WITH SCHEMABINDING"}}
/*
func withoutComments(batch []byte) []byte {
	return regexp.MustCompile(`--.*\n`).ReplaceAll(batch, []byte{})
}*/

func Lint(batch []byte, reportLinesWithOffset int) {
	for _, rule := range activeRules {
		loc := rule.when.FindIndex(batch)
		if loc != nil && !rule.then.Match(batch[loc[0]:loc[1]]) {
			atLine := strings.Count(string(batch[0:loc[0]]), "\n") + reportLinesWithOffset +1
			fmt.Println("  line", atLine, ":", rule.message)
		}
	}
}
