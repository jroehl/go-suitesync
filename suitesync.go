//+build !test

package main

import (
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jroehl/go-suitesync/sdf"
	"github.com/jroehl/go-suitesync/suitetalk"
	"github.com/kardianos/osext"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/urfave/cli"
)

type Exec struct {
}

var (
	bash        = new(Exec)
	errRequired = errors.New("required args missing")
)

type ExecCmd interface {
	Run() error
}

func (bash *Exec) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func before(c *cli.Context) error {
	lib.IsVerbose = c.GlobalBool("verbose") || c.GlobalBool("debug")
	lib.IsDebug = c.GlobalBool("debug")
	if lib.IsVerbose {
		lib.PrWarnf("\nRUNNING VERBOSE MODE\n")
	}
	if lib.IsDebug {
		lib.PrWarnf("RUNNING DEBUG MODE\n")
	}
	if os.Getenv("GO_ENV") == "local" {
		// get correct folder when running go run ...
		_, callerFile, _, _ := runtime.Caller(0)
		lib.CurrentDir = filepath.Dir(callerFile)
	} else {
		folderPath, _ := osext.ExecutableFolder()
		lib.CurrentDir = folderPath
	}
	return lib.InitEnv(c.Command.FullName() == "issuetoken")
}

func after(c *cli.Context) error {
	if lib.IsVerbose {
		lib.PrResultf("\nExecution successful\n\n")
	}

	if !lib.IsDebug {
		// cleanup of temporary files
		dir := lib.MkTempDir()
		res := lib.FindDir(filepath.Dir(dir), "suitesync_*")
		for _, r := range res {
			lib.Remove(r)
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	app.Name = "suitesync"
	app.Usage = "a netsuite filehandling cli"
	app.Version = "0.0.3"
	app.After = after
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose output",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debugging mode",
		},
	}

	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Johann Röhl",
			Email: "mail@johannroehl.de",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "init",
			Aliases:   []string{"i"},
			Before:    before,
			Usage:     "Initialize suitesync",
			ArgsUsage: "[src] (dest)",
			Action: func(c *cli.Context) error {
				lib.PrNoticef("\ninit successful\n")
				return nil
			},
		},
		{
			Name:      "issuetoken",
			Aliases:   []string{"it"},
			Before:    before,
			Usage:     "Issue sdf cli token",
			ArgsUsage: "[password]",
			Action: func(c *cli.Context) error {
				lib.PrNoticef("\ninit successful\n")
				args := c.Args()
				if len(args) == 0 {
					args[0] = lib.Credentials[lib.Password]
				}
				_, err := sdf.GenerateToken(bash, lib.CheckRequired(c.Args(), 0, "password"))
				return err
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Before:  before,
			Usage:   "Sync two directories (creates hash file)",
			// Flags: []cli.Flag{
			// 	cli.BoolFlag{
			// 		Name:  "bidirectional, b",
			// 		Usage: "Sync both ways (local-remote / remote-local filesystem)",
			// 	},
			// },
			ArgsUsage: "[src] [dest]",
			Action: func(c *cli.Context) error {
				_, _, _, _, err := sdf.Sync(
					bash,
					http.DefaultClient,
					lib.CheckRequired(c.Args(), 0, "src"),
					lib.CheckRequired(c.Args(), 1, "dest"),
					false, // c.Bool("bidirectional"),
					nil,
				)
				return err
			},
		},
		{
			Name:      "upload",
			Aliases:   []string{"u"},
			Before:    before,
			Usage:     "Upload files and/or directories to filecabinet directory",
			ArgsUsage: "[src...] [dest]",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if !args.Present() || len(args) < 2 {
					return errRequired
				}
				dest, srcs := args[len(args)-1:][0], args[:len(args)-1]
				_, err := sdf.Upload(bash, srcs, dest)
				return err
			},
		},
		{
			Name:      "download",
			Aliases:   []string{"d"},
			Before:    before,
			Usage:     "Download files and/or directories to local filesystem",
			ArgsUsage: "[src...] [dest]",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if !args.Present() || len(args) < 2 {
					return errRequired
				}
				dest, srcs := args[len(args)-1:][0], args[:len(args)-1]
				_, err := sdf.Download(bash, http.DefaultClient, srcs, dest)
				return err
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Before:    before,
			Usage:     "List files and directories from filecabinet",
			ArgsUsage: "[fcDir]",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if !args.Present() || len(args) < 1 {
					return errRequired
				}
				fcDir := args[len(args)-1:][0]
				err := sdf.List(bash, http.DefaultClient, fcDir)
				return err
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"del"},
			Before:    before,
			Usage:     "Delete files and/or directories in filecabinet",
			ArgsUsage: "[files...|dirs...]",
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					return errRequired
				}
				res, _ := suitetalk.DeleteRequest(http.DefaultClient, c.Args())
				lib.PrintResponse("Delete results", res)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		lib.PrFatalf("\n%s\n", err.Error())
	}
}
