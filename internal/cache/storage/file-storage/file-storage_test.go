package filestorage_test

import (
	"context"
	"path"
	"time"

	filestorage "github.com/itbasis/go-tools-sdkm/internal/cache/storage/file-storage"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/jonboulle/clockwork"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.DescribeTable("GoString", func(dir, id, want string) {
	gomega.Expect(filestorage.NewFileCacheStorage(dir, sdkmPlugin.ID(id)).GoString()).To(gomega.Equal(want))
},
	ginkgo.Entry(nil, "rootDir", "go", "FileCacheStorage[rootDir/go.json]"),
	ginkgo.Entry(nil, "rootDir", "jvm", "FileCacheStorage[rootDir/jvm.json]"),
	ginkgo.Entry(nil, "rootDir", "blah", "FileCacheStorage[rootDir/blah.json]"),
)

var _ = ginkgo.DescribeTable(
	"Fail", func(testDir string, testDate, pluginId string, wantErr error) {
		var (
			ctx              = initContext(testDate)
			rootDir          = path.Join("testdata", testDir)
			fileCacheStorage = filestorage.NewFileCacheStorage(rootDir, sdkmPlugin.ID(pluginId))
		)

		gomega.Expect(fileCacheStorage.Load(ctx)).Error().To(gomega.MatchError(wantErr))
	},
	ginkgo.Entry(nil, "000", "2025-04-18 00:00:00", "go", sdkmSDKVersion.ErrCacheInvalidFile),
	ginkgo.Entry(nil, "001", "2025-04-18 00:00:02", "go", sdkmSDKVersion.ErrCacheExpired),
)

var _ = ginkgo.DescribeTable("Load", func(testDir, testDate, pluginId string, want sdkmSDKVersion.MapSdkVersionGroupType) {
	var (
		ctx              = initContext(testDate)
		rootDir          = path.Join("testdata", testDir)
		fileCacheStorage = filestorage.NewFileCacheStorage(rootDir, sdkmPlugin.ID(pluginId))
	)
	gomega.Expect(fileCacheStorage.Load(ctx)).To(gomega.Equal(want))
},
	ginkgo.Entry(nil, "001", "2025-04-17 00:00:00", "go", sdkmSDKVersion.MapSdkVersionGroupType{
		sdkmSDKVersion.TypeArchived: []sdkmSDKVersion.SDKVersion{
			sdkmSDKVersion.NewSDKVersion("1.22.5", sdkmSDKVersion.TypeArchived, false),
		},
		sdkmSDKVersion.TypeStable: []sdkmSDKVersion.SDKVersion{
			sdkmSDKVersion.NewSDKVersion("1.22.8", sdkmSDKVersion.TypeStable, false),
		},
	}),
)

func initContext(testDate string) context.Context {
	var t, err = time.Parse(filestorage.UpdatedFormat, testDate)
	gomega.Expect(err).To(gomega.Succeed())

	return clockwork.AddToContext(context.Background(), clockwork.NewFakeClockAt(t))
}
