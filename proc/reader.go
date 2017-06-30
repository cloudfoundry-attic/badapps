package proc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Reader struct {
	procPath string
}

func NewReader(path string) Reader {
	return Reader{path}
}

func (p Reader) Pids() ([]string, error) {
	entries, err := ioutil.ReadDir(p.procPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read proc dir %s")
	}

	pids := []string{}
	for _, info := range entries {

		if info.IsDir() && isNumeric(info) {
			pids = append(pids, info.Name())
		}
	}
	return pids, nil
}

func (p Reader) Env(pid string) (map[string]string, error) {
	environ, err := ioutil.ReadFile(filepath.Join(p.procPath, pid, "environ"))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read environ for pid %s", pid)
	}

	envMap := make(map[string]string)
	for _, envVar := range strings.Split(string(environ), "\x00") {
		fields := strings.Split(envVar, "=")
		if len(fields) != 2 {
			continue
		}
		envMap[fields[0]] = fields[1]
	}

	return envMap, nil
}

func isNumeric(info os.FileInfo) bool {
	_, err := strconv.Atoi(info.Name())
	if err != nil {
		return false
	}
	return true
}
