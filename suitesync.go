package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/jroehl/go-suitesync/restlet"
	"github.com/jroehl/go-suitesync/sdf"

	tm "github.com/buger/goterm"
	"github.com/jroehl/go-suitesync/lib"
	"github.com/urfave/cli"
)

func init() {
	lib.InitEnv()
}

// https://github.com/sanathkr/go-npm

func before(c *cli.Context) error {
	lib.IsVerbose = c.GlobalBool("verbose")
	if c.GlobalBool("verbose") {
		fmt.Println()
		tm.Println(tm.Color(tm.Bold("RUNNING VERBOSE MODE"), tm.RED))
		tm.Flush()
	}
	return nil
}

func after(c *cli.Context) error {
	if lib.IsVerbose {
		lib.PrNoticeF("Execution successful")
	}
	return nil
}

func checkRequired(s []string, i int, arg string) string {
	if i >= len(s) {
		lib.PrFatalf("Required arg \"%s\" is missing\n", arg)
	}
	return s[i]
}

func optArgDest(c *cli.Context, i int, def string) string {
	if i >= len(c.Args()) {
		return optDest(def, "")
	}
	return c.Args()[i]
}

func optDest(d string, def string) (s string) {
	s = d
	if s == "" {
		s = def
	}
	if s == "" {
		s, _ = os.Getwd()
	}
	return
}

var (
	promptRE = regexp.MustCompile("Type YES to continue.")
	passRE   = regexp.MustCompile("password:")
	doneRE   = regexp.MustCompile("Finished at:")
)

const (
	timeout = 10 * time.Minute
)

func main() {
	app := cli.NewApp()

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	app.Name = "suitesync"
	app.Usage = "a netsuite filehandling cli"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose output",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "sync two directories",
			Before:  before,
			After:   after,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "bidirectional, b",
					Usage: "Should it sync both ways (not only from local to remote filesystem set this flag)",
				},
			},
			ArgsUsage: "[src] (dest)",
			Action: func(c *cli.Context) error {
				sdf.Sync(
					checkRequired(c.Args(), 0, "src"),
					optArgDest(c, 1, ""),
					c.Bool("bidirectional"),
				)
				return nil
			},
		},
		{
			Name:    "upload",
			Aliases: []string{"ul"},
			Usage:   "upload files/directories or restlet to filecabinet",
			Before:  before,
			After:   after,
			Subcommands: []cli.Command{
				{
					Category: "upload",
					Name:     "restlet",
					Aliases:  []string{"r"},
					Usage:    "upload restlet to filecabinet",
					Action: func(c *cli.Context) error {
						sdf.UploadRestlet()
						return nil
					},
				},
				{
					Category:  "upload",
					Name:      "dir",
					Usage:     "upload directory to filecabinet",
					Aliases:   []string{"d"},
					ArgsUsage: "[src] [dest]",
					Action: func(c *cli.Context) error {
						sdf.UploadDir(
							checkRequired(c.Args(), 0, "src"),
							optArgDest(c, 1, "/SuiteScripts"),
						)
						return nil
					},
				},
				{
					Category: "upload",
					Name:     "files",
					Aliases:  []string{"f"},
					Usage:    "upload files to filecabinet",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "root, r",
							Usage: "root directory that is trimmed from local filepath",
						},
						cli.StringFlag{
							Name:  "dest, d",
							Usage: "destination directory in filecabinet",
						},
					},
					ArgsUsage: "[files]",
					Action: func(c *cli.Context) error {
						if c.String("root") == "" {
							lib.PrFatalf("Required flag \"root\" is missing\n")
						}
						if !c.Args().Present() {
							lib.PrFatalf("Required args (files) missing\n")
						}
						sdf.UploadFiles(
							c.String("root"),
							c.Args(),
							optDest(c.String("dest"), "/SuiteScripts"),
						)
						return nil
					},
				},
			},
		},
		{
			Name:    "download",
			Aliases: []string{"dl"},
			Usage:   "options for task templates",
			Before:  before,
			After:   after,
			Subcommands: []cli.Command{
				{
					Category:  "download",
					Name:      "dir",
					Usage:     "Download directory from filecabinet to local filesystem",
					ArgsUsage: "[src] [dest]",
					Aliases:   []string{"d"},
					Action: func(c *cli.Context) error {
						sdf.DownloadDir(
							checkRequired(c.Args(), 0, "src"),
							optArgDest(c, 1, ""),
						)
						return nil
					},
				},
				{
					Category:  "download",
					Name:      "files",
					Aliases:   []string{"f"},
					Usage:     "Download files from filecabinet to local filesystem",
					ArgsUsage: "[files]",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dest, d",
							Usage: "Destination directory in filecabinet",
						},
					},
					Action: func(c *cli.Context) error {
						if !c.Args().Present() {
							lib.PrFatalf("Required args (files) missing\n")
						}
						sdf.DownloadFiles(c.Args(), optDest(c.String("dest"), ""))
						return nil
					},
				},
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "delete files and/or directories in filecabinet",
			ArgsUsage: "[files|dirs]",
			Before:    before,
			After:     after,
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					lib.PrFatalf("Required args (files/dirs) missing\n")
				}
				restlet.Delete(c.Args())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
