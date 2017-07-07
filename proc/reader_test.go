package proc_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/badapps/proc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	var (
		procPath   string
		procReader proc.Reader
	)

	BeforeEach(func() {
		var err error
		procPath, err = ioutil.TempDir("", "proc")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.Mkdir(filepath.Join(procPath, "100"), 0755)).To(Succeed())
		Expect(os.Mkdir(filepath.Join(procPath, "135"), 0755)).To(Succeed())
		Expect(os.Mkdir(filepath.Join(procPath, "400"), 0755)).To(Succeed())
		Expect(os.Mkdir(filepath.Join(procPath, "info"), 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(procPath, "410"), []byte{}, 0755)).To(Succeed())

		procReader = proc.NewReader(procPath)
	})

	Describe("Pids", func() {
		It("returns the list of pids in the proc fs", func() {
			pids, err := procReader.Pids()
			Expect(err).NotTo(HaveOccurred())

			Expect(pids).To(ConsistOf("100", "135", "400"))
		})
	})

	Describe("Env", func() {
		BeforeEach(func() {
			Expect(ioutil.WriteFile(filepath.Join(procPath, "100", "environ"), []byte("banana=monkey\x00carrot=rabbit"), 0755)).To(Succeed())
		})

		It("returns the environment of the specified pid", func() {
			env, err := procReader.Env("100")
			Expect(err).NotTo(HaveOccurred())

			Expect(env["banana"]).To(Equal("monkey"))
			Expect(env["carrot"]).To(Equal("rabbit"))
		})

		Context("when env contains a malformed var", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(filepath.Join(procPath, "100", "environ"), []byte("banana=monkey\x00wolf"), 0755)).To(Succeed())
			})
			It("ignores that var", func() {
				env, err := procReader.Env("100")
				Expect(err).NotTo(HaveOccurred())

				Expect(env).To(HaveKey("banana"))
				Expect(env).NotTo(HaveKey("wolf"))
			})
		})

		Context("when trying to read the env for a process which no longer exists", func() {
			It("does not fail", func() {
				_, err := procReader.Env("666")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when permission to read the proc dir is denied", func() {
			It("does not fail", func() {
				Expect(os.Chmod(filepath.Join(procPath, "100"), 0666))
				_, err := procReader.Env("100")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
