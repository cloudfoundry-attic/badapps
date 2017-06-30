package cfapp_test

import (
	"code.cloudfoundry.org/badapps/cfapp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
	It("parses a the VCAP_APPLICATION string into a Info struct", func() {
		vcapApplication := `{"application_id":"123","application_name":"test-app","application_uris":["test-app.cf.example.com"],"application_version":"a08f6419-7b53-466f-ab89-235da7e9c6fe","cf_api":"https://api.cf.example.com","host":"0.0.0.0","instance_id":"f0c11c7c-dcaa-483a-5c9c-e202","instance_index":8,"limits":{"disk":64,"fds":16384,"mem":32},"name":"test-app","port":8080,"space_id":"fb712285-0bfc-4da4-8d19-e91a10843527","space_name":"diagnostics","uris":["test-app.cf.example.com"],"version":"a08f6419-7b53-466f-ab89-235da7e9c6fe"}`

		info, err := cfapp.Parse(vcapApplication)
		Expect(err).NotTo(HaveOccurred())

		Expect(info).To(Equal(cfapp.Info{
			ApplicationID:   "123",
			ApplicationName: "test-app",
			ApplicationUris: []string{"test-app.cf.example.com"},
			InstanceID:      "f0c11c7c-dcaa-483a-5c9c-e202",
			Limits:          map[string]int{"disk": 64, "fds": 16384, "mem": 32},
			SpaceId:         "fb712285-0bfc-4da4-8d19-e91a10843527",
			SpaceName:       "diagnostics",
		}))
	})
})
