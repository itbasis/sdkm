package sdkversion

import (
	"context"
	"errors"
	"fmt"
)

type MapSdkVersionGroupType = map[VersionType][]SDKVersion

var (
	ErrCacheNotFoundFile = errors.New("cache file not found")
	ErrCacheInvalidFile  = errors.New("fail cache file")
	ErrCacheExpired      = errors.New("cache expired")
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=cache-storage.mock.go
type CacheStorage interface {
	fmt.GoStringer

	Validate(ctx context.Context) bool
	Load(ctx context.Context) (MapSdkVersionGroupType, error)
	Store(ctx context.Context, versions MapSdkVersionGroupType)
}
