package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type action int
type match int

const (
	keep         action = 0
	command      action = 1
	hotkey       action = 2
	replaceOne   action = 3
	replaceTwo   action = 4
	replaceThree action = 5
	overwrite    action = 6

	matchFalse   match = 0
	matchTrue    match = 1
	matchHotkey  match = 2
	matchCommand match = 3
)

// Expression encapusulates a regular expression and an associated action
type Expression struct {
	action action
	regex  *regexp.Regexp
}

// matches returns MatchTrue if the line matches the regex or a more specific match (MatchCommand/MatchHotkey)
func (e Expression) matches(line string) match {

	if e.regex.MatchString(line) {
		switch e.action {

		case command:
			return matchCommand

		case hotkey:
			return matchHotkey

		default:
			return matchTrue
		}
	}

	return matchFalse
}

func (e Expression) extract(line string) string {
	return e.regex.ReplaceAllString(line, "$1")
}

// replace returns a string modified according to the regex and the action
func (e Expression) replace(line string, key string) string {

	switch e.action {

	case keep:
		return line

	case command:
		return line

	case hotkey:
		return line

	case overwrite:
		return fmt.Sprintf("%s%s", e.regex.ReplaceAllString(line, "$1"), key)

	case replaceOne:
		return fmt.Sprintf(e.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)"), key) + e.regex.ReplaceAllString(line, "$5")

	case replaceTwo:
		return fmt.Sprintf(e.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)$5$6$7$8 (|cffffcc00%s|r)"), key, key)

	case replaceThree:
		return fmt.Sprintf(e.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)$5$6$7$8 (|cffffcc00%s|r)$9$10$11$12 (|cffffcc00%s|r)$13$14"), key, key, key)
	}

	return "<< ERROR >>"
}

// NewExpressions returns a set of regex expressions with correspodning actions
func NewExpressions() []Expression {

	expressions := []Expression{}

	expressions = append(expressions, Expression{ // command
		action: command,
		regex:  regexp.MustCompile(`^\[[\w]*\][ \t]*$`),
	})
	expressions = append(expressions, Expression{ // hotkey
		action: hotkey,
		regex:  regexp.MustCompile(`^Hotkey=(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	expressions = append(expressions, Expression{ // researchhotkey
		action: hotkey,
		regex:  regexp.MustCompile(`^Researchhotkey=(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	expressions = append(expressions, Expression{ // comment
		action: keep,
		regex:  regexp.MustCompile(`^\/\/.*$`),
	})
	expressions = append(expressions, Expression{ // empty
		action: keep,
		regex:  regexp.MustCompile(`^[ \t]*$`),
	})

	expressions = append(expressions, Expression{ // unhotkey
		action: overwrite,
		regex:  regexp.MustCompile(`^(?P<name>Unhotkey=)(?P<hotkey>[\w \!\.]*)$`),
	})

	expressions = append(expressions, Expression{ // Awakentip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Awakentip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	expressions = append(expressions, Expression{ // Awakentip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Awakentip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	expressions = append(expressions, Expression{ // Awakentip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Awakentip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})

	expressions = append(expressions, Expression{ // Researchtip=t(i)p [Level %d]
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1> - \[\|cffffcc00Level %d\|r\])"?[ \t]*$`),
	})
	expressions = append(expressions, Expression{ // Researchtip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Researchtip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	expressions = append(expressions, Expression{ // Researchtip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	expressions = append(expressions, Expression{ // Researchtip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})

	expressions = append(expressions, Expression{ // Revivetip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Revivetip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	expressions = append(expressions, Expression{ // Revivetip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Revivetip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	expressions = append(expressions, Expression{ // Revivetip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Revivetip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})

	expressions = append(expressions, Expression{ // Untip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Untip="?\|cffc3dbff[\w \-\!\.]{2,}\|r"?$`),
	})
	expressions = append(expressions, Expression{ // Untip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Untip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	expressions = append(expressions, Expression{ // Untip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Untip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	expressions = append(expressions, Expression{ // Untip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Untip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})

	expressions = append(expressions, Expression{ // Tip=t(i)p1,t(i)p2,t(i)p3
		action: replaceThree,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1> - \[\|cffffcc00Level 1\|r\],)` +
			`(?P<p3>[\w \-\!\.]*)\|cffffcc00(?P<key2>\w)\|r(?P<p4>[\w \-\!\.]*)(?P<l2> - \[\|cffffcc00Level 2\|r\],)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3> - \[\|cffffcc00Level 3\|r\])"?[ \t]*$`),
	})
	expressions = append(expressions, Expression{ // Tip=t(i)p1,t(i)p2,t(i)p3
		action: replaceThree,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1>,)` +
			`(?P<p3>[\w \-\!\.]*)\|cffffcc00(?P<key2>\w)\|r(?P<p4>[\w \-\!\.]*)(?P<l2>,)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3>)"?$`),
	})

	expressions = append(expressions, Expression{ // Tip=t(i)p1,t(i)p2
		action: replaceTwo,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1> - \[\|cffffcc00Level 1\|r\],)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3> - \[\|cffffcc00Level 2\|r\])"?[ \t]*$`),
	})
	expressions = append(expressions, Expression{ // Tip=t(i)p1,t(i)p2
		action: replaceTwo,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1>,)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3>)"?$`),
	})

	expressions = append(expressions, Expression{ // Tip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Tip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	expressions = append(expressions, Expression{ // Tip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Tip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	expressions = append(expressions, Expression{ // Tip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Tip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})

	return expressions
}

// Group denotes a single hotkey group within the CustomKeys.txt file
type Group struct {
	Hotkey string
	Lines  []string
}

// Adjust will replace all the lines with their regex equivalents
func (g Group) Adjust(expressions []Expression) {
	for i, line := range g.Lines {
		for _, e := range expressions {
			m := e.matches(line)
			if m == matchCommand || m == matchHotkey || m == matchTrue {
				g.Lines[i] = e.replace(line, g.Hotkey)
				break
			}
		}
	}
}

// Print outputs the group to stdout
func (g Group) Print() {
	for _, line := range g.Lines {
		fmt.Printf("%s\n", line)
	}
}

func main() {

	expressions := NewExpressions()

	current := Group{}

	f, err := os.Open("CustomKeys.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

	innerloop:
		for _, e := range expressions {
			switch e.matches(line) {

			case matchCommand:
				current.Adjust(expressions)
				current.Print()
				current = Group{Lines: []string{line}}
				break innerloop

			case matchHotkey:
				current.Lines = append(current.Lines, line)
				current.Hotkey = e.extract(line)
				break innerloop

			case matchTrue:
				current.Lines = append(current.Lines, line)
				break innerloop
			}
		}
	}

	current.Adjust(expressions)
	current.Print()

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
