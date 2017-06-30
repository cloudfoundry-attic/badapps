package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	badappsBinPath string
)

func TestBadapps(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		binPath, err := gexec.Build("code.cloudfoundry.org/badapps")
		Expect(err).NotTo(HaveOccurred())

		return []byte(binPath)
	}, func(data []byte) {
		badappsBinPath = string(data)
	})

	RunSpecs(t, "Badapps Suite")
}
