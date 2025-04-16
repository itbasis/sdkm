package sdkversion_test

import (
	"testing"

	"github.com/itbasis/go-test-utils/v5/ginkgo"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
)

func TestSDKVersionSuite(t *testing.T) {
	ginkgo.InitGinkgoSuite(t, "SDK Version Suite")
}

var (
	go1_23_rc2  = sdkmSDKVersion.NewSDKVersion("1.23rc2", sdkmSDKVersion.TypeArchived, false)
	go1_23_rc1  = sdkmSDKVersion.NewSDKVersion("1.23rc1", sdkmSDKVersion.TypeArchived, true)
	go1_22_8    = sdkmSDKVersion.NewSDKVersion("1.22.8", sdkmSDKVersion.TypeStable, false)
	go1_22_6    = sdkmSDKVersion.NewSDKVersion("1.22.6", sdkmSDKVersion.TypeArchived, false)
	go1_22_5    = sdkmSDKVersion.NewSDKVersion("1.22.5", sdkmSDKVersion.TypeArchived, true)
	go1_18_10   = sdkmSDKVersion.NewSDKVersion("1.18.10", sdkmSDKVersion.TypeArchived, false)
	go1_4_beta1 = sdkmSDKVersion.NewSDKVersion("1.4beta1", sdkmSDKVersion.TypeArchived, false)
	go1_3_rc1   = sdkmSDKVersion.NewSDKVersion("1.3rc1", sdkmSDKVersion.TypeArchived, false)
)
