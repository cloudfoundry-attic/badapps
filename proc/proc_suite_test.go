package proc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestProc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proc Suite")
}
