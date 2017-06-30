package cfapp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfapp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cfapp Suite")
}
