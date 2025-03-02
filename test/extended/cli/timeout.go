package cli

import (
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	exutil "github.com/openshift/origin/test/extended/util"
	k8simage "k8s.io/kubernetes/test/utils/image"
)

var _ = g.Describe("[sig-cli] oc --request-timeout", func() {
	defer g.GinkgoRecover()

	oc := exutil.NewCLI("oc-request-timeout")

	g.It("works as expected", func() {
		busyBoxImage := k8simage.GetE2EImage(k8simage.BusyBox)
		err := oc.Run("create").Args("deploymentconfig", "testdc", "--image="+busyBoxImage).Execute()
		o.Expect(err).NotTo(o.HaveOccurred())

		out, err := oc.Run("get", "dc/testdc").Args("-w", "-v=5", "--request-timeout=1s").Output()
		o.Expect(err).NotTo(o.HaveOccurred())
		// timeout is set for both the request and on context in request
		// seek8s.io/client-go/rest/request.go#request so if we get timeout
		// from server or from context it's ok
		o.Expect(out).Should(o.SatisfyAny(o.ContainSubstring("request canceled"), o.ContainSubstring("context deadline exceeded")))

		out, err = oc.Run("get", "dc/testdc").Args("--request-timeout=1s").Output()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(out).To(o.ContainSubstring("testdc"))

		out, err = oc.Run("get", "dc/testdc").Args("--request-timeout=1").Output()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(out).To(o.ContainSubstring("testdc"))

		out, err = oc.Run("get", "pods").Args("--watch", "-v=5", "--request-timeout=1s").Output()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(out).Should(o.SatisfyAny(
			o.ContainSubstring("request canceled"),
			o.ContainSubstring("context deadline exceeded"),
			o.ContainSubstring("Client.Timeout exceeded while awaiting headers")))
	})
})
