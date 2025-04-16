package sdkversion

import (
	"github.com/Masterminds/semver/v3"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.DescribeTable("makeSemVer", func(version string, wantSemVer *semver.Version) {
	defer ginkgo.GinkgoRecover()

	gomega.Expect(makeSemVer(version)).To(gomega.Equal(wantSemVer))
},
	ginkgo.Entry(nil, "1", semver.New(1, 0, 0, "", "")),
	ginkgo.Entry(nil, "1.23.3", semver.New(1, 23, 3, "", "")),
	ginkgo.Entry(nil, "1.23", semver.New(1, 23, 0, "", "")),
	ginkgo.Entry(nil, "1.23rc2", semver.New(1, 23, 0, "rc2", "")),
	ginkgo.Entry(nil, "1.3beta1", semver.New(1, 3, 0, "beta1", "")),
)
