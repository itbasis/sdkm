package golang_test

import (
	"context"
	"os"
	"path"

	itbasisCoreEnv "github.com/itbasis/go-tools-core/env"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	pluginBase "github.com/itbasis/go-tools-sdkm/internal/plugins/base"
	sdkmPluginGo "github.com/itbasis/go-tools-sdkm/internal/plugins/golang"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"go.uber.org/mock/gomock"
)

var _ = ginkgo.Describe(
	"EnvByVersion", func() {
		defer ginkgo.GinkgoRecover()

		ginkgo.DescribeTable(
			"success", func(sdkVersion, wantSDKPath, wantGoCachePath string) {
				var (
					originPath      = os.Getenv(itbasisCoreEnv.KeyPath)
					originPaths     = itbasisCoreOs.SplitPathList(originPath)
					countOriginPath = len(originPaths)

					mockController = gomock.NewController(ginkgo.GinkgoT())
					mockBasePlugin = sdkmPlugin.NewMockBasePlugin(mockController)

					sdkVersionDir = path.Join("sdk", sdkVersion)
				)

				mockBasePlugin.EXPECT().GetSDKDir().Return(sdkVersionDir).AnyTimes()
				mockBasePlugin.EXPECT().GetSDKVersionDir(pluginGoConsts.PluginID, sdkVersion).Return(sdkVersionDir).AnyTimes()
				mockBasePlugin.EXPECT().HasInstalled(pluginGoConsts.PluginID, sdkVersion).Return(true).AnyTimes()
				mockBasePlugin.EXPECT().
					PrepareEnvironment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						itbasisCoreEnv.Map{
							itbasisCoreEnv.KeyPath:                                    originPath,
							itbasisCoreEnv.KeyGoRoot:                                  os.Getenv(itbasisCoreEnv.KeyGoRoot),
							itbasisCoreEnv.KeyGoPath:                                  os.Getenv(itbasisCoreEnv.KeyGoPath),
							itbasisCoreEnv.KeyGoBin:                                   os.Getenv(itbasisCoreEnv.KeyGoBin),
							pluginBase.EnvSdkmOriginPrefix + itbasisCoreEnv.KeyPath:   originPath,
							pluginBase.EnvSdkmOriginPrefix + itbasisCoreEnv.KeyGoRoot: os.Getenv(itbasisCoreEnv.KeyGoRoot),
							pluginBase.EnvSdkmOriginPrefix + itbasisCoreEnv.KeyGoPath: os.Getenv(itbasisCoreEnv.KeyGoPath),
							pluginBase.EnvSdkmOriginPrefix + itbasisCoreEnv.KeyGoBin:  os.Getenv(itbasisCoreEnv.KeyGoBin),
						},
					).
					AnyTimes()

				var pluginGo, errGetPlugin = sdkmPluginGo.GetPlugin(mockBasePlugin)
				gomega.Expect(errGetPlugin).To(gomega.Succeed())

				var envs, err = pluginGo.EnvByVersion(context.Background(), sdkVersion)

				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(envs).To(
					gomega.SatisfyAll(
						gomega.HaveLen(8),

						gomega.HaveKeyWithValue(pluginBase.EnvSdkmOriginPrefix+itbasisCoreEnv.KeyPath, originPath),
					),
				)

				gomega.Expect(originPaths).To(gomega.HaveLen(countOriginPath))
				gomega.Expect(itbasisCoreOs.SplitPathList(envs[itbasisCoreEnv.KeyPath])).To(
					gomega.SatisfyAll(
						gomega.HaveLen(countOriginPath+3),
						gstruct.MatchElementsWithIndex(
							gstruct.IndexIdentity, gstruct.IgnoreExtras, gstruct.Elements{
								"0": gomega.Equal(wantSDKPath),
								"1": gomega.Equal(wantGoCachePath),
								"2": gomega.Equal(itbasisCoreOs.FixPath(itbasisCoreOs.ExecutableDir())),
							},
						),
					),
				)
			},
			ginkgo.Entry(nil, "1.23.0", path.Join("sdk", "1.23.0", "bin"), path.Join("1.23.0", "bin")),
			ginkgo.Entry(nil, "1.20.1", path.Join("sdk", "1.20.1", "bin"), path.Join("1.20.1", "bin")),
		)
	},
)
