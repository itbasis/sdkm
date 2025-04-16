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
	"ListAllVersionsByPrefix", func() {
		defer ginkgo.GinkgoRecover()

		var pluginGo sdkmPlugin.SDKMPlugin

		ginkgo.BeforeEach(
			func() {
				mockController := gomock.NewController(ginkgo.GinkgoT())

				mockSDKVersions := sdkmSDKVersion.NewMockSDKVersions(mockController)
				mockSDKVersions.EXPECT().AllVersions(gomock.Any(), false).Return(testSdkList, nil).MaxTimes(1)

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
				var actualSdkList, err = pluginGo.ListAllVersions(context.Background(), false)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(actualSdkList.AsList()).
					To(
						gomega.SatisfyAll(
							gomega.HaveLen(testSdkList.Len()),
							gomega.ContainElements(go1_23_rc2, go1_22_5, go1_22_0, go1_18, go1_18_10, go1_4_beta1, go1_3_rc1),
						),
					)
			},
		)

		// 		ginkgo.When(
		// 			"By Prefix", func() {
		// 				ginkgo.It(
		// 					"empty prefix", func() {
		// 						gomega.Expect(pluginGo.ListAllVersionsByPrefix(context.Background(), false, "")).
		// 							To(
		// 								gomega.SatisfyAll(
		// 									gomega.HaveLen(len(sdkList)),
		// 									gomega.ContainElements(wantGo1_23_rc2, wantGo1_22_5, wantGo1_22_0, wantGo1_18, wantGo1_18_10, wantGo1_4_beta1, wantGo1_3_rc1),
		// 								),
		// 							)
		// 					},
		// 				)

		// 				ginkgo.DescribeTable(
		// 					"success", func(prefix string, wantCount int, wantSDKVersions []sdkmSDKVersion.SDKVersion) {
		// 						gomega.Expect(pluginGo.ListAllVersionsByPrefix(context.Background(), false, prefix)).
		// 							To(
		// 								gomega.SatisfyAll(
		// 									gomega.HaveLen(wantCount),
		// 									gomega.ContainElements(wantSDKVersions),
		// 								),
		// 							)
		// 					},
		// 					ginkgo.Entry(nil, "1.23", 2, []sdkmSDKVersion.SDKVersion{wantGo1_23_rc2, wantGo1_23_rc1}),
		// 					ginkgo.Entry(nil, "1.21", 5, []sdkmSDKVersion.SDKVersion{wantGo1_21_12, wantGo1_21_11, wantGo1_21_0, wantGo1_21_rc3}),
		// 				)
		// 			},
		// 		)
	},
)
