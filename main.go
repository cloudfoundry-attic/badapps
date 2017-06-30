package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code.cloudfoundry.org/badapps/cfapp"
	"code.cloudfoundry.org/badapps/proc"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "badapps"
	app.Usage = "list running apps from the current diego-cell"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "proc",
			Value: "/proc",
			Usage: "procfs path",
		},
	}

	app.Action = func(c *cli.Context) error {
		procPath := c.String("proc")
		procReader := proc.NewReader(procPath)
		pids, err := procReader.Pids()
		if err != nil {
			return errors.Wrapf(err, "Failed to get pids %s")
		}

		output := make(map[string]cfapp.Info)

		for _, pid := range pids {
			env, err := procReader.Env(pid)
			if err != nil {
				if os.IsPermission(errors.Cause(err)) {
					fmt.Fprintf(os.Stderr, "WARNING: Permission denied to read %s/%s/environ. Ignoring\n", procPath, pid)
					continue
				}

				return errors.Wrapf(err, "Failed to read env for pid %s", pid)
			}

			if vcapApplication, ok := env["VCAP_APPLICATION"]; ok {
				info, err := cfapp.Parse(vcapApplication)
				if err != nil {
					return errors.Wrap(err, "Failed to parse vcap_application")
				}

				output[pid] = info
			}
		}

		if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
			return errors.Wrap(err, "Failed to encode output")
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
