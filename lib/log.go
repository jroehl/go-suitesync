package lib

import (
	"fmt"
	"os"
	"strings"

	tm "github.com/buger/goterm"
)

func print(str string, color int) {
	tm.Println(tm.Color(str, color))
	tm.Flush()
}

func printF(format string, color int, a ...interface{}) {
	print(fmt.Sprintf(format, a...), color)
}

func PrHeaderF(format string, a ...interface{}) {
	printF(format, tm.MAGENTA, a...)
}

func PrNoticeF(format string, a ...interface{}) {
	printF(format, tm.GREEN, a...)
}

func PrWarnF(format string, a ...interface{}) {
	printF(format, tm.MAGENTA, a...)
}

func PrFatalf(format string, a ...interface{}) {
	print(tm.Bold(fmt.Sprintf(format, a)), tm.RED)
	os.Exit(1)
}

// PrettyList output list prettified for terminal
func PrettyList(str string, list []string) {
	if len(list) > 0 {
		tbl := tm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintln(tbl, tm.Color(strings.Join([]string{str, "\n"}, ""), tm.GREEN))
		for _, s := range list {
			fmt.Fprintf(tbl, "  - %s\n", s)
		}
		tm.Println(tbl)
		tm.Flush()
	}
}

// PrettyHash output list prettified for terminal
func PrettyHash(str string, list []Hash) {
	PrNoticeF(str)
	tbl := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintln(tbl, "  #   NAME\tHASH\tPATH")
	for i, s := range list {
		fmt.Fprintf(tbl, "  %d   %s\t%s\t%s\n", i+1, s.Name, s.Hash, s.Path)
	}
	tm.Println(tbl)
	tm.Flush()
}

// PrintResponse print restlet response prettified
func PrintResponse(s string, a []Response) {
	if len(a) > 0 {
		PrNoticeF(s)
		tbl := tm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintln(tbl, "  #   ID\tTYPE\tCODE\tSTATUS\tPATH\tMESSAGE")
		for i, it := range a {
			fmt.Fprintf(tbl, "  %d   %s\t%s\t%d\t%s\t%s\t%s\n", i+1, it.Id, it.Type, it.Code, it.Status, it.Path, it.Message)
		}
		tm.Println(tbl)
		tm.Flush()
	}
}
