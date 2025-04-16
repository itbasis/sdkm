package golang_test

import (
	"testing"

	itbasisTestUtils "github.com/itbasis/go-test-utils/v5/ginkgo"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
)

func TestVersionSuite(t *testing.T) {
	itbasisTestUtils.InitGinkgoSuite(t, "Golang Plugin Suite")
}

var (
	// stable
	go1_22_8  = sdkmSDKVersion.NewSDKVersion("1.22.8", sdkmSDKVersion.TypeStable, false)
	go1_21_12 = sdkmSDKVersion.NewSDKVersion("1.21.12", sdkmSDKVersion.TypeStable, true)
	// unstable
	go1_23_rc2 = sdkmSDKVersion.NewSDKVersion("1.23rc2", sdkmSDKVersion.TypeUnstable, false)
	// archived
	go1_23_rc1  = sdkmSDKVersion.NewSDKVersion("1.23rc1", sdkmSDKVersion.TypeArchived, true)
	go1_22_5    = sdkmSDKVersion.NewSDKVersion("1.22.5", sdkmSDKVersion.TypeArchived, true)
	go1_22_4    = sdkmSDKVersion.NewSDKVersion("1.22.4", sdkmSDKVersion.TypeArchived, true)
	go1_22_3    = sdkmSDKVersion.NewSDKVersion("1.22.3", sdkmSDKVersion.TypeArchived, false)
	go1_22_0    = sdkmSDKVersion.NewSDKVersion("1.22.0", sdkmSDKVersion.TypeArchived, false)
	go1_21_11   = sdkmSDKVersion.NewSDKVersion("1.21.11", sdkmSDKVersion.TypeArchived, true)
	go1_21_10   = sdkmSDKVersion.NewSDKVersion("1.21.10", sdkmSDKVersion.TypeArchived, true)
	go1_21_0    = sdkmSDKVersion.NewSDKVersion("1.21.0", sdkmSDKVersion.TypeArchived, true)
	go1_21_rc3  = sdkmSDKVersion.NewSDKVersion("1.21rc3", sdkmSDKVersion.TypeArchived, false)
	go1_20_14   = sdkmSDKVersion.NewSDKVersion("1.20.14", sdkmSDKVersion.TypeArchived, false)
	go1_19_13   = sdkmSDKVersion.NewSDKVersion("1.19.13", sdkmSDKVersion.TypeArchived, true)
	go1_19_12   = sdkmSDKVersion.NewSDKVersion("1.19.12", sdkmSDKVersion.TypeArchived, true)
	go1_18      = sdkmSDKVersion.NewSDKVersion("1.18", sdkmSDKVersion.TypeArchived, false)
	go1_18_10   = sdkmSDKVersion.NewSDKVersion("1.18.10", sdkmSDKVersion.TypeArchived, false)
	go1_4_beta1 = sdkmSDKVersion.NewSDKVersion("1.4beta1", sdkmSDKVersion.TypeArchived, false)
	go1_3_rc1   = sdkmSDKVersion.NewSDKVersion("1.3rc1", sdkmSDKVersion.TypeArchived, false)

	testSdkList = sdkmSDKVersion.NewSdkVersionList(
		go1_22_8,
		go1_22_5,
		go1_22_4,
		go1_22_3,
		go1_22_0,
		go1_21_12,
		go1_21_11,
		go1_21_10,
		go1_21_0,
		go1_21_rc3,
		go1_23_rc2,
		go1_23_rc1,
		go1_20_14,
		go1_19_13,
		go1_19_12,
		go1_18,
		go1_18_10,
		go1_4_beta1,
		go1_3_rc1,
	)
	sdkMap = testSdkList.AsMap()
)
