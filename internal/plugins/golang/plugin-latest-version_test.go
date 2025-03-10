package golang_test

import (
	"context"
	"fmt"

	sdkmPluginGo "github.com/itbasis/go-tools-sdkm/internal/plugins/golang"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"go.uber.org/mock/gomock"
)

var _ = ginkgo.Describe(
	"Plugin Latest Version", func() {
		defer ginkgo.GinkgoRecover()

		var pluginGo sdkmPlugin.SDKMPlugin

		ginkgo.BeforeEach(
			func() {
				mockController := gomock.NewController(ginkgo.GinkgoT())

				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().LatestVersion(gomock.Any(), false).Return(
					sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable},
					nil,
				)

				mockBasePlugin := sdkmPlugin.NewMockBasePlugin(mockController)
				mockBasePlugin.EXPECT().GetSDKDir().Return("").AnyTimes()
				mockBasePlugin.EXPECT().HasInstalled(pluginGoConsts.PluginID, gomock.Any()).Return(false)

				plugin, err := sdkmPluginGo.GetPlugin(mockBasePlugin)
				gomega.Expect(err).To(gomega.Succeed())
				pluginGo = plugin.WithVersions(mockSDKVersions)

			},
		)

		ginkgo.DescribeTable(
			"LatestVersion", func(wantSDKVersion sdkmSDKVersion.SDKVersion) {
				gomega.Expect(pluginGo.LatestVersion(context.Background(), false)).
					To(
						gomega.HaveValue(
							gstruct.MatchFields(
								gstruct.IgnoreExtras, gstruct.Fields{
									"ID":   gomega.Equal(wantSDKVersion.ID),
									"Type": gomega.Equal(wantSDKVersion.Type),
								},
							),
						),
					)
			},
			ginkgo.Entry(nil, sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable}),
		)
	},
)

var _ = ginkgo.Describe(
	"LatestVersionByPrefix", func() {
		defer ginkgo.GinkgoRecover()

		var pluginGo sdkmPlugin.SDKMPlugin

		ginkgo.BeforeEach(
			func() {
				mockController := gomock.NewController(ginkgo.GinkgoT())

				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().
					LatestVersion(gomock.Any(), false).
					Return(
						sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable},
						nil,
					).
					MaxTimes(1)
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
							{ID: "1.21.11", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.20.14", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.19.13", Type: sdkmSDKVersion.TypeArchived},
							{ID: "1.19.12", Type: sdkmSDKVersion.TypeArchived},
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

		ginkgo.DescribeTable(
			"success", func(prefix string, wantSDKVersion sdkmSDKVersion.SDKVersion) {
				gomega.Expect(pluginGo.LatestVersionByPrefix(context.Background(), false, prefix)).
					To(
						gomega.HaveValue(
							gstruct.MatchFields(
								gstruct.IgnoreExtras, gstruct.Fields{
									"ID":   gomega.Equal(wantSDKVersion.ID),
									"Type": gomega.Equal(wantSDKVersion.Type),
								},
							),
						),
					)
			},
			ginkgo.Entry("empty prefix", "", sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable}),
			ginkgo.Entry(nil, "1.23", sdkmSDKVersion.SDKVersion{ID: "1.23rc2", Type: sdkmSDKVersion.TypeUnstable}),
			ginkgo.Entry(nil, "1.22", sdkmSDKVersion.SDKVersion{ID: "1.22.5", Type: sdkmSDKVersion.TypeStable}),
			ginkgo.Entry(nil, "1.21", sdkmSDKVersion.SDKVersion{ID: "1.21.12", Type: sdkmSDKVersion.TypeStable}),
			ginkgo.Entry(nil, "1.20", sdkmSDKVersion.SDKVersion{ID: "1.20.14", Type: sdkmSDKVersion.TypeArchived}),
			ginkgo.Entry(nil, "1.19", sdkmSDKVersion.SDKVersion{ID: "1.19.13", Type: sdkmSDKVersion.TypeArchived}),
		)

		ginkgo.DescribeTable(
			"fail", func(prefix string) {
				gomega.Expect(pluginGo.LatestVersionByPrefix(context.Background(), false, prefix)).Error().To(
					gomega.MatchError(fmt.Sprintf("version by prefix %s: SDK version not found", prefix)),
				)
			},
			ginkgo.Entry("", "1.24"),
		)
	},
)
