package golang_test

import (
	"context"
	"sort"

	sdkmPluginGo "github.com/itbasis/go-tools-sdkm/internal/plugins/golang"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = ginkgo.Describe(
	"LatestVersion", func() {
		defer ginkgo.GinkgoRecover()

		var pluginGo sdkmPlugin.SDKMPlugin

		sort.Sort(testSdkList)

		ginkgo.BeforeEach(
			func() {
				mockController := gomock.NewController(ginkgo.GinkgoT())

				sort.Sort(testSdkList)
				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().AllVersions(gomock.Any(), false).Return(testSdkList, nil)

				mockBasePlugin := sdkmPlugin.NewMockBasePlugin(mockController)
				mockBasePlugin.EXPECT().GetSDKDir().Return("").AnyTimes()

				plugin, err := sdkmPluginGo.GetPlugin(mockBasePlugin)
				gomega.Expect(err).To(gomega.Succeed())
				pluginGo = plugin.WithVersions(mockSDKVersions)

			},
		)

		ginkgo.DescribeTable(
			"LatestVersion", func(onlyInstalled bool, wantSDKVersion sdkmSDKVersion.SDKVersion) {
				gomega.Expect(pluginGo.LatestVersion(context.Background(), false, onlyInstalled)).
					To(gomega.Equal(wantSDKVersion))
			},
			ginkgo.Entry(nil, true, go1_21_12),
			ginkgo.Entry(nil, false, go1_22_8),
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

				sort.Sort(testSdkList)
				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().AllVersions(gomock.Any(), false).Return(testSdkList, nil).MaxTimes(1)

				mockBasePlugin := sdkmPlugin.NewMockBasePlugin(mockController)
				mockBasePlugin.EXPECT().GetSDKDir().Return("").AnyTimes()
				mockBasePlugin.EXPECT().HasInstalled(pluginGoConsts.PluginID, gomock.Any()).DoAndReturn(func(_ sdkmPlugin.ID, version string) bool {
					v, ok := sdkMap[version]
					gomega.Expect(ok).To(gomega.BeTrue())

					return v.HasInstalled()
				}).AnyTimes()

				plugin, err := sdkmPluginGo.GetPlugin(mockBasePlugin)
				gomega.Expect(err).To(gomega.Succeed())
				pluginGo = plugin.WithVersions(mockSDKVersions)
			},
		)

		ginkgo.DescribeTable(
			"success", func(prefix string, onlyInstalled bool, wantSDKVersion sdkmSDKVersion.SDKVersion) {
				gomega.Expect(pluginGo.LatestVersionByPrefix(context.Background(), false, onlyInstalled, prefix)).
					To(gomega.Equal(wantSDKVersion))
			},
			ginkgo.Entry("empty prefix", "", true, go1_21_12),
			ginkgo.Entry("empty prefix", "", false, go1_22_8),
			ginkgo.Entry(nil, "1.23", true, go1_23_rc1),
			ginkgo.Entry(nil, "1.23", false, go1_23_rc2),
			ginkgo.Entry(nil, "1.22", true, go1_22_5),
			ginkgo.Entry(nil, "1.22", false, go1_22_8),
			ginkgo.Entry(nil, "1.21", true, go1_21_12),
			ginkgo.Entry(nil, "1.21", false, go1_21_12),
			ginkgo.Entry(nil, "1.20", false, go1_20_14),
			ginkgo.Entry(nil, "1.19", true, go1_19_13),
			ginkgo.Entry(nil, "1.19", false, go1_19_13),
		)

		ginkgo.DescribeTable(

			"fail", func(prefix string, onlyInstalled bool, wantErr string) {
				gomega.Expect(pluginGo.LatestVersionByPrefix(context.Background(), false, onlyInstalled, prefix)).
					Error().To(
					gomega.MatchError(
						sdkmSDKVersion.ErrSDKVersionNotFound,
						gomega.ContainSubstring(wantErr),
					),
				)
			},
			ginkgo.Entry("", "1.24", true, "version by prefix 1.24"),
			ginkgo.Entry("", "1.24", false, "version by prefix 1.24"),
			ginkgo.Entry("", "1.20", true, "version by prefix 1.20"),
		)
	},
)
