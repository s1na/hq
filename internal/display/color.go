package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	Green  = color.New(color.FgGreen)
	Red    = color.New(color.FgRed)
	Yellow = color.New(color.FgYellow)
	Bold   = color.New(color.Bold)
	Cyan   = color.New(color.FgCyan)
)

// ColorizeDiff takes a raw diff/log string and prints it with color.
// Lines starting with ++ are green, lines starting with -- are red.
func ColorizeDiff(text string, noColor bool) {
	if noColor {
		fmt.Print(text)
		return
	}
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "++"):
			Green.Println(line)
		case strings.HasPrefix(line, "--"):
			Red.Println(line)
		case strings.HasPrefix(line, "@@"):
			Cyan.Println(line)
		default:
			fmt.Println(line)
		}
	}
}

// PassFail returns a colored pass/fail string.
func PassFail(pass bool) string {
	if pass {
		return Green.Sprint("PASS")
	}
	return Red.Sprint("FAIL")
}

// PassFailCount returns a colored string like "10/12".
func PassFailCount(passes, total int) string {
	if passes == total {
		return Green.Sprintf("%d/%d", passes, total)
	}
	return Yellow.Sprintf("%d/%d", passes, total)
}
