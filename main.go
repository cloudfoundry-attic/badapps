package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code.cloudfoundry.org/badapps/cfapp"
	"code.cloudfoundry.org/badapps/proc"

	"github.com/pkg/errors"
)

const procPath = "/proc"

func main() {
	procReader := proc.NewReader(procPath)
	pids, err := procReader.Pids()
	if err != nil {
		exitWithError(errors.Wrapf(err, "Failed to get pids %s"))
	}

	output := make(map[string]cfapp.Info)

	for _, pid := range pids {
		env, err := procReader.Env(pid)
		if err != nil {
			if os.IsPermission(errors.Cause(err)) {
				fmt.Fprintf(os.Stderr, "WARNING: Permission denied to read %s/%s/environ. Ignoring\n", procPath, pid)
				continue
			}

			exitWithError(errors.Wrapf(err, "Failed to read env for pid %s", pid))
		}

		if vcapApplication, ok := env["VCAP_APPLICATION"]; ok {
			info, err := cfapp.Parse(vcapApplication)
			if err != nil {
				exitWithError(errors.Wrap(err, "Failed to parse vcap_application"))
			}

			output[pid] = info
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		exitWithError(errors.Wrap(err, "Failed to encode output"))
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
