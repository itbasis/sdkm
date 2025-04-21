package cache

import (
	"context"
	"sync"

	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"golang.org/x/exp/maps"
)

type cache struct {
	storeLock sync.Mutex

	cacheStorage sdkmSDKVersion.CacheStorage
	cache        sdkmSDKVersion.MapSdkVersionGroupType
}

func NewCache() sdkmSDKVersion.Cache {
	return &cache{
		cache: make(sdkmSDKVersion.MapSdkVersionGroupType),
	}
}

func (receiver *cache) WithExternalStore(cacheStorage sdkmSDKVersion.CacheStorage) sdkmSDKVersion.Cache {
	receiver.cacheStorage = cacheStorage
	maps.Clear(receiver.cache)

	return receiver
}

func (receiver *cache) Valid(ctx context.Context) bool {
	if receiver.cacheStorage != nil {
		m, _ := receiver.cacheStorage.Load(ctx)
		receiver.cache = m
	}

	return len(receiver.cache) > 0
}

func (receiver *cache) Load(ctx context.Context, versionType sdkmSDKVersion.VersionType) sdkmSDKVersion.SdkVersionList {
	receiver.storeLock.Lock()
	defer receiver.storeLock.Unlock()

	if cacheStorage := receiver.cacheStorage; cacheStorage != nil {
		if len(receiver.cache) == 0 {
			if m, err := cacheStorage.Load(ctx); err != nil {
				receiver.cache = m
			}
		}
	}

	var list, ok = receiver.cache[versionType]
	if !ok {
		return sdkmSDKVersion.NewSdkVersionList()
	}

	return sdkmSDKVersion.NewSdkVersionList(list...)
}

func (receiver *cache) Store(ctx context.Context, versionType sdkmSDKVersion.VersionType, sdkVersionList sdkmSDKVersion.SdkVersionList) {
	receiver.storeLock.Lock()
	defer receiver.storeLock.Unlock()

	receiver.cache[versionType] = sdkVersionList.AsList()

	if receiver.cacheStorage != nil {
		receiver.cacheStorage.Store(ctx, receiver.cache)
	}
}

func (receiver *cache) GoString() string {
	var result = "SDKVersionCache"

	if receiver.cacheStorage != nil {
		result = result + " (" + receiver.cacheStorage.GoString() + ")"
	}

	return result
}
