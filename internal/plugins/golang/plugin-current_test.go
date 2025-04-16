package golang_test

import (
	"context"
	"log/slog"
	"path"

	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	sdkmPluginGo "github.com/itbasis/go-tools-sdkm/internal/plugins/golang"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = ginkgo.Describe(
	"Current", func() {
		var (
			mockBasePlugin *sdkmPlugin.MockBasePlugin
			pluginGo       sdkmPlugin.SDKMPlugin
		)

		ginkgo.BeforeEach(
			func() {
				mockController := gomock.NewController(ginkgo.GinkgoT())

				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().AllVersions(gomock.Any(), false).Return(testSdkList, nil)

				mockBasePlugin = sdkmPlugin.NewMockBasePlugin(mockController)
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

		var getBaseDir = func(dir string) string {
			baseDir := path.Join(itbasisCoreOs.Pwd(), "testdata/current", dir)
			slog.Debug("baseDir: " + baseDir)
			gomega.Expect(baseDir).To(gomega.BeADirectory())

			return baseDir
		}

		ginkgo.DescribeTable(
			"success", func(testDir string, onlyInstalled bool, wantSDK sdkmSDKVersion.SDKVersion) {
				baseDir := getBaseDir(testDir)

				gomega.Expect(pluginGo.Current(context.Background(), false, onlyInstalled, baseDir)).To(gomega.Equal(wantSDK))
			},
			ginkgo.Entry(nil, "001", true, go1_21_12),
			ginkgo.Entry(nil, "002", true, go1_22_5),
			ginkgo.Entry(nil, "003", true, go1_22_5),
			ginkgo.Entry(nil, "003", false, go1_22_5),
			ginkgo.Entry(nil, "004", false, go1_22_3),
			ginkgo.Entry(nil, "005", true, go1_23_rc1),
			ginkgo.Entry(nil, "006", true, go1_23_rc1),
			ginkgo.Entry(nil, "007", true, go1_23_rc1),
			ginkgo.Entry(nil, "008", true, go1_22_5),
			ginkgo.Entry(nil, "008", false, go1_22_8),
		)

		ginkgo.DescribeTable("fail", func(testDir string, onlyInstalled bool) {
			baseDir := getBaseDir(testDir)

			gomega.Expect(pluginGo.Current(context.Background(), false, onlyInstalled, baseDir)).
				Error().To(gomega.MatchError(sdkmSDKVersion.ErrSDKVersionNotFound))
		},
			ginkgo.Entry(nil, "004", true),
		)
	},
)
