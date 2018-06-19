package lib

import (
	"fmt"
	"log"
	"strings"
)

func print(str string, color int) {
	fmt.Println(str)
}

func printF(format string, color int, a ...interface{}) {
	fmt.Printf(format, a...)
}

func PrHeaderF(format string, a ...interface{}) {
	printF(format, 0, a...)
}

func PrNoticeF(format string, a ...interface{}) {
	printF(format, 0, a...)
}

func PrWarnF(format string, a ...interface{}) {
	printF(format, 0, a...)
}

func PrFatalf(format string, a ...interface{}) {
	log.Fatalf(format, a...)
	// os.Exit(1)
}

// PrettyList output list prettified for terminal
func PrettyList(str string, list []string) {
	if len(list) > 0 {
		fmt.Println(strings.Join([]string{str, "\n"}, ""))
		for _, s := range list {
			fmt.Printf("  - %s\n", s)
		}
	}
}

// PrettyHash output list prettified for terminal
func PrettyHash(str string, list []Hash) {
	PrNoticeF("%s\n", str)
	fmt.Println("  #   NAME\tHASH\tPATH")
	for i, s := range list {
		fmt.Printf("  %d   %s\t%s\t%s\n", i+1, s.Name, s.Hash, s.Path)
	}
}

// PrintResponse print restlet response prettified
func PrintResponse(s string, a []Response) {
	if len(a) > 0 {
		PrNoticeF("%s\n", s)
		fmt.Println("  #   ID\tTYPE\tCODE\tSTATUS\tPATH\tMESSAGE")
		for i, it := range a {
			fmt.Printf("  %d   %s\t%s\t%d\t%s\t%s\t%s\n", i+1, it.ID, it.Type, it.Code, it.Status, it.Path, it.Message)
		}
	}
}
