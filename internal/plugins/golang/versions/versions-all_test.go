package versions_test

import (
	"context"
	"path"

	pluginGoVersions "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/versions"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe(
	"All Versions", func() {
		defer ginkgo.GinkgoRecover()

		ginkgo.DescribeTable(
			"different versions of answers", func(testResponseFilePath string, wantLen int, wantSDKVersions []sdkmSDKVersion.SDKVersion) {
				var server = initFakeServer(path.Join("testdata", "all-versions", testResponseFilePath))
				defer server.Close()

				sdkVersionList, err := pluginGoVersions.NewVersions(server.URL()).AllVersions(context.Background(), false)

				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sdkVersionList.AsList()).
					To(
						gomega.SatisfyAll(
							gomega.HaveLen(wantLen),

							gomega.ContainElements(
								sdkmSDKVersion.NewSDKVersion("1.23rc1", sdkmSDKVersion.TypeArchived, false),
								sdkmSDKVersion.NewSDKVersion("1.22.0", sdkmSDKVersion.TypeArchived, false),
								sdkmSDKVersion.NewSDKVersion("1.18", sdkmSDKVersion.TypeArchived, false),
								sdkmSDKVersion.NewSDKVersion("1.18.10", sdkmSDKVersion.TypeArchived, false),
								sdkmSDKVersion.NewSDKVersion("1.4beta1", sdkmSDKVersion.TypeArchived, false),
								sdkmSDKVersion.NewSDKVersion("1.3rc1", sdkmSDKVersion.TypeArchived, false),
							),

							gomega.ContainElements(wantSDKVersions),
						),
					)
			},
			ginkgo.Entry(
				nil, "001.html", 292, []sdkmSDKVersion.SDKVersion{
					sdkmSDKVersion.NewSDKVersion("1.23rc2", sdkmSDKVersion.TypeUnstable, false),
					sdkmSDKVersion.NewSDKVersion("1.22.5", sdkmSDKVersion.TypeStable, false),
				},
			),
			ginkgo.Entry(
				nil, "002.html", 295, []sdkmSDKVersion.SDKVersion{
					sdkmSDKVersion.NewSDKVersion("1.23.0", sdkmSDKVersion.TypeStable, false),
					sdkmSDKVersion.NewSDKVersion("1.23rc2", sdkmSDKVersion.TypeArchived, false),
					sdkmSDKVersion.NewSDKVersion("1.22.6", sdkmSDKVersion.TypeStable, false),
					sdkmSDKVersion.NewSDKVersion("1.22.5", sdkmSDKVersion.TypeArchived, false),
				},
			),
			ginkgo.Entry(
				nil, "003.html", 304, []sdkmSDKVersion.SDKVersion{
					sdkmSDKVersion.NewSDKVersion("1.24rc1", sdkmSDKVersion.TypeUnstable, false),
					sdkmSDKVersion.NewSDKVersion("1.23.4", sdkmSDKVersion.TypeStable, false),
					sdkmSDKVersion.NewSDKVersion("1.23.0", sdkmSDKVersion.TypeArchived, false),
					sdkmSDKVersion.NewSDKVersion("1.23rc2", sdkmSDKVersion.TypeArchived, false),
					sdkmSDKVersion.NewSDKVersion("1.22.10", sdkmSDKVersion.TypeStable, false),
					sdkmSDKVersion.NewSDKVersion("1.22.6", sdkmSDKVersion.TypeArchived, false),
					sdkmSDKVersion.NewSDKVersion("1.22.5", sdkmSDKVersion.TypeArchived, false),
				},
			),
		)
	},
)
