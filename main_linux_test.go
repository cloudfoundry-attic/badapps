package main_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"code.cloudfoundry.org/badapps/cfapp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	var (
		app1, app2       *gexec.Session
		app1Pid, app2Pid string
	)

	BeforeEach(func() {
		app1 = startApp("app-1", "1")
		app1Pid = fmt.Sprintf("%d", app1.Command.Process.Pid)

		app2 = startApp("app-2", "2")
		app2Pid = fmt.Sprintf("%d", app2.Command.Process.Pid)
	})

	AfterEach(func() {
		app1.Kill()
		app2.Kill()
	})

	It("returns the app id of all applications running", func() {
		cmd := exec.Command(badappsBinPath)
		stdoutBuffer := gbytes.NewBuffer()
		session, err := gexec.Start(cmd, stdoutBuffer, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		var applications map[string]cfapp.Info
		Expect(json.NewDecoder(stdoutBuffer).Decode(&applications)).To(Succeed())

		Expect(applications).To(Equal(map[string]cfapp.Info{
			app1Pid: cfapp.Info{
				ApplicationID:   "1",
				ApplicationName: "app-1",
				ApplicationUris: []string{"app-1.cf.example.com"},
				InstanceID:      "f0c11c7c-dcaa-483a-5c9c-e202",
				Limits:          map[string]int{"disk": 64, "fds": 16384, "mem": 32},
				SpaceId:         "fb712285-0bfc-4da4-8d19-e91a10843527",
				SpaceName:       "diagnostics",
			},
			app2Pid: cfapp.Info{
				ApplicationID:   "2",
				ApplicationName: "app-2",
				ApplicationUris: []string{"app-2.cf.example.com"},
				InstanceID:      "f0c11c7c-dcaa-483a-5c9c-e202",
				Limits:          map[string]int{"disk": 64, "fds": 16384, "mem": 32},
				SpaceId:         "fb712285-0bfc-4da4-8d19-e91a10843527",
				SpaceName:       "diagnostics",
			},
		}))
	})
})

func startApp(appName, appId string) *gexec.Session {
	cmd := exec.Command("sleep", "1000")
	cmd.Env = []string{fmt.Sprintf("VCAP_APPLICATION=%s", appEnvData(appName, appId))}
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
