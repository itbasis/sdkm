package sdkversion_test

import (
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe(
	"Print", func() {
		ginkgo.DescribeTable(
			"Strange models, but need to check", func(model sdkmSDKVersion.SDKVersion, expected string) {
				gomega.Expect(model.Print()).To(gomega.Equal(expected))
			},
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), ""),
		)

		ginkgo.DescribeTable(
			"correct models", func(model sdkmSDKVersion.SDKVersion, expected string) {
				gomega.Expect(model.Print()).To(gomega.Equal(expected))
			},
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), "1.2 [not installed]"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, true), "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), "1.2 (unstable) [not installed]"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, true), "1.2 (unstable)"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeArchived, false), "1.2 (archived) [not installed]"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeArchived, true), "1.2 (archived)"),
		)
	},
)

var _ = ginkgo.Describe(
	"Print with options", func() {
		ginkgo.DescribeTable(
			"Strange models, but need to check",
			func(model sdkmSDKVersion.SDKVersion, outType, outInstalled, outNotInstalled bool, expected string) {
				gomega.Expect(model.PrintWithOptions(outType, outInstalled, outNotInstalled)).To(gomega.Equal(expected))
			},

			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), false, false, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), false, true, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), false, true, true, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), true, false, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), true, true, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeStable, false), true, true, true, ""),

			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), false, false, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), false, true, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), false, true, true, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), true, false, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), true, true, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeUnstable, false), true, true, true, ""),

			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), false, false, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), false, true, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), false, true, true, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), true, false, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), true, true, false, ""),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("", sdkmSDKVersion.TypeArchived, false), true, true, true, ""),
		)

		ginkgo.DescribeTable(
			"correct",
			func(model sdkmSDKVersion.SDKVersion, outType, outInstalled, outNotInstalled bool, expected string) {
				gomega.Expect(model.PrintWithOptions(outType, outInstalled, outNotInstalled)).To(gomega.Equal(expected))
			},

			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), false, false, false, "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), false, true, false, "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), false, true, true, "1.2 [not installed]"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), true, false, false, "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), true, true, false, "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeStable, false), true, true, true, "1.2 [not installed]"),

			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), false, false, false, "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), false, true, false, "1.2"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), false, true, true, "1.2 [not installed]"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), true, false, false, "1.2 (unstable)"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), true, true, false, "1.2 (unstable)"),
			ginkgo.Entry(nil, sdkmSDKVersion.NewSDKVersion("1.2", sdkmSDKVersion.TypeUnstable, false), true, true, true, "1.2 (unstable) [not installed]"),
		)
	},
)
