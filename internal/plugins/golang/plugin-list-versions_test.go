package golang_test

import (
	"context"

	sdkmPluginGo "github.com/itbasis/go-tools-sdkm/internal/plugins/golang"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = ginkgo.Describe(
	"ListAllVersions", func() {
		defer ginkgo.GinkgoRecover()

	},
)

var _ = ginkgo.Describe(
	"ListAllVersionsByPrefix", func() {
		defer ginkgo.GinkgoRecover()

		var pluginGo sdkmPlugin.SDKMPlugin

		ginkgo.BeforeEach(
			func() {
				mockController := gomock.NewController(ginkgo.GinkgoT())

				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().
					AllVersions(gomock.Any(), false).
					Return(
						[]sdkmSDKVersion.SDKVersion{
							{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable},
							{ID: "1.21.12", Type: sdkmSDKVersion.TypeStable},
							{ID: "1.23rc2", Type: sdkmSDKVersion.TypeUnstable},
							{ID: "1.23rc1", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.22.4", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.22.3", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.22.0", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.21.11", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.21.10", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.21.0", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.21rc3", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.20.14", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.19.13", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.19.12", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.18", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.18.10", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.4beta1", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.3rc1", Type: sdkmSDKVersion.TypeArchived},
						},
						nil,
					).MaxTimes(1)

				mockBasePlugin := sdkmPlugin.NewMockBasePlugin(mockController)
				mockBasePlugin.EXPECT().GetSDKDir().Return("").AnyTimes()
				mockBasePlugin.EXPECT().HasInstalled(pluginGoConsts.PluginID, gomock.Any()).Return(false).AnyTimes()

				plugin, err := sdkmPluginGo.GetPlugin(mockBasePlugin)
				gomega.Expect(err).To(gomega.Succeed())
				pluginGo = plugin.WithVersions(mockSDKVersions)
			},
		)

		ginkgo.It(
			"", func() {
				gomega.Expect(pluginGo.ListAllVersions(context.Background(), false)).
					To(
						gomega.SatisfyAll(
							gomega.HaveLen(18),

							gomega.ContainElements(
								sdkmSDKVersion.SDKVersion{ID: "1.23rc2", Type: sdkmSDKVersion.TypeUnstable},
								sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable},
								sdkmSDKVersion.SDKVersion{ID: "1.22.0", Type: sdkmSDKVersion.TypeArchived},
								sdkmSDKVersion.SDKVersion{ID: "1.18", Type: sdkmSDKVersion.TypeArchived},
								sdkmSDKVersion.SDKVersion{ID: "1.18.10", Type: sdkmSDKVersion.TypeArchived},
								sdkmSDKVersion.SDKVersion{ID: "1.4beta1", Type: sdkmSDKVersion.TypeArchived},
								sdkmSDKVersion.SDKVersion{ID: "1.3rc1", Type: sdkmSDKVersion.TypeArchived},
							),
						),
					)
			},
		)

		ginkgo.When(
			"By Prefix", func() {
				ginkgo.It(
					"empty prefix", func() {
						gomega.Expect(pluginGo.ListAllVersionsByPrefix(context.Background(), false, "")).
							To(
								gomega.SatisfyAll(
									gomega.HaveLen(18),

									gomega.ContainElements(
										sdkmSDKVersion.SDKVersion{ID: "1.23rc2", Type: sdkmSDKVersion.TypeUnstable},
										sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable},
										sdkmSDKVersion.SDKVersion{ID: "1.22.0", Type: sdkmSDKVersion.TypeArchived},
										sdkmSDKVersion.SDKVersion{ID: "1.18", Type: sdkmSDKVersion.TypeArchived},
										sdkmSDKVersion.SDKVersion{ID: "1.18.10", Type: sdkmSDKVersion.TypeArchived},
										sdkmSDKVersion.SDKVersion{ID: "1.4beta1", Type: sdkmSDKVersion.TypeArchived},
										sdkmSDKVersion.SDKVersion{ID: "1.3rc1", Type: sdkmSDKVersion.TypeArchived},
									),
								),
							)
					},
				)

				ginkgo.DescribeTable(
					"success", func(prefix string, wantCount int, wantSDKVersions []sdkmSDKVersion.SDKVersion) {
						gomega.Expect(pluginGo.ListAllVersionsByPrefix(context.Background(), false, prefix)).
							To(
								gomega.SatisfyAll(
									gomega.HaveLen(wantCount),
									gomega.ContainElements(wantSDKVersions),
								),
							)
					},
					ginkgo.Entry(
						nil, "1.23", 2, []sdkmSDKVersion.SDKVersion{
							{ID: "1.23rc2", Type: sdkmSDKVersion.TypeUnstable},
							{ID: "1.23rc1", Type: sdkmSDKVersion.TypeArchived},
						},
					),
					ginkgo.Entry(
						nil, "1.21", 5, []sdkmSDKVersion.SDKVersion{
							{ID: "1.21.12", Type: sdkmSDKVersion.TypeStable},
							{ID: "1.21.11", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.21.0", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.21rc3", Type: sdkmSDKVersion.TypeArchived},
						},
					),
				)
			},
		)
	},
)
