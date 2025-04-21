package filestorage

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	itbasisCoreLog "github.com/itbasis/go-tools-core/log"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	sdkmPlugin "github.com/itbasis/go-tools-sdkm/pkg/plugin"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/jonboulle/clockwork"
)

const (
	cacheExpirationDuration = 24 * time.Hour
)

type _fileStorage struct {
	lock sync.Mutex

	filePath string
}

func NewFileCacheStorage(rootDir string, pluginId sdkmPlugin.ID) sdkmSDKVersion.CacheStorage {
	slog.Debug("init file cache storage",
		slog.String("rootDir", rootDir),
		sdkmPlugin.SlogAttrPluginId(pluginId),
	)

	var filePath = path.Join(rootDir, string(pluginId)+".json")
	slog.Debug("using cache with file path: " + filePath)

	return &_fileStorage{filePath: filePath}
}

func (r *_fileStorage) GoString() string {
	return "FileCacheStorage[" + r.filePath + "]"
}

func (r *_fileStorage) Validate(ctx context.Context) bool {
	var _, err = r.load(ctx)

	return err == nil
}

func (r *_fileStorage) Load(ctx context.Context) (sdkmSDKVersion.MapSdkVersionGroupType, error) {
	return r.load(ctx)
}

func (r *_fileStorage) Store(ctx context.Context, versions sdkmSDKVersion.MapSdkVersionGroupType) {
	r.lock.Lock()
	defer r.lock.Unlock()

	filePath := r.filePath

	slog.Debug("storing cache to file: " + filePath)

	var bytes, errMarshal = json.Marshal(
		model{
			Updated:  updated(clockwork.FromContext(ctx).Now()),
			Versions: versions,
		},
	)

	if errMarshal != nil {
		slog.Error(
			"error marshalling cache file",
			itbasisCoreLog.SlogAttrError(errMarshal),
			itbasisCoreLog.SlogAttrFilePath(filePath),
		)

		return
	}

	dir := filepath.Dir(filePath)
	if errMkdir := os.MkdirAll(dir, itbasisCoreOs.DefaultDirMode); errMkdir != nil {
		slog.Error("error creating cache dir: "+dir, itbasisCoreLog.SlogAttrError(errMkdir))

		return
	}

	if errWriteFile := os.WriteFile(filePath, bytes, itbasisCoreOs.DefaultFileMode); errWriteFile != nil {
		slog.Error("error writing cache file: "+filePath, itbasisCoreLog.SlogAttrError(errWriteFile))
	}
}

func (r *_fileStorage) load(ctx context.Context) (sdkmSDKVersion.MapSdkVersionGroupType, error) {
	var bytes, errRead = r.readFile()
	if errRead != nil {
		return nil, errRead
	}

	var m model
	if err := json.Unmarshal(bytes, &m); err != nil {
		slog.Error("fail unmarshalling cache file",
			itbasisCoreLog.SlogAttrFilePath(r.filePath),
			itbasisCoreLog.SlogAttrError(err),
		)

		return nil, sdkmSDKVersion.ErrCacheInvalidFile
	}

	if !r.validate(ctx, m) {
		return nil, sdkmSDKVersion.ErrCacheExpired
	}

	return m.Versions, nil
}

func (r *_fileStorage) readFile() ([]byte, error) {
	var filePath = r.filePath

	slog.Debug("reading cache file", itbasisCoreLog.SlogAttrFilePath(filePath))

	if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
		slog.Debug("cache file not found", itbasisCoreLog.SlogAttrFilePath(filePath))

		return nil, sdkmSDKVersion.ErrCacheNotFoundFile
	} else if err != nil {
		slog.Error("fail accessing cache file", itbasisCoreLog.SlogAttrError(err))

		return nil, sdkmSDKVersion.ErrCacheInvalidFile
	}

	var bytes, errReadFile = os.ReadFile(filePath)
	if errReadFile != nil {
		slog.Error("fail reading cache file", itbasisCoreLog.SlogAttrFilePath(filePath), itbasisCoreLog.SlogAttrError(errReadFile))

		return nil, sdkmSDKVersion.ErrCacheInvalidFile
	}

	return bytes, nil
}

func (r *_fileStorage) validate(ctx context.Context, m model) bool {
	var (
		now    = clockwork.FromContext(ctx).Now()
		sub    = now.Sub(time.Time(m.Updated))
		result = sub <= cacheExpirationDuration
	)

	slog.Debug("validate",
		slog.Time("cache date", time.Time(m.Updated)),
		slog.Duration("expiration", cacheExpirationDuration),
		slog.Time("current date", now),
		slog.Duration("sub", sub),
		slog.Bool("valid", result),
	)

	return result
}
