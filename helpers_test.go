package main_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func appEnvData(appName, appID string) string {
	return fmt.Sprintf(`{"application_id":"%s","application_name":"%s","application_uris":["%s.cf.example.com"],"application_version":"a08f6419-7b53-466f-ab89-235da7e9c6fe","cf_api":"https://api.cf.example.com","host":"0.0.0.0","instance_id":"f0c11c7c-dcaa-483a-5c9c-e202","instance_index":8,"limits":{"disk":64,"fds":16384,"mem":32},"name":"%s","port":8080,"space_id":"fb712285-0bfc-4da4-8d19-e91a10843527","space_name":"diagnostics","uris":["%s.cf.example.com"],"version":"a08f6419-7b53-466f-ab89-235da7e9c6fe"}`, appID, appName, appName, appName, appName)
}

func startApp(appName, appId string) *gexec.Session {
	cmd := exec.Command("sleep", "1000")
	cmd.Env = []string{fmt.Sprintf("VCAP_APPLICATION=%s", appEnvData(appName, appId))}
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
