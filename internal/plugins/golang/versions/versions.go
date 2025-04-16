package versions

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"

	itbasisCoreLog "github.com/itbasis/go-tools-core/log"
	sdkmCache "github.com/itbasis/go-tools-sdkm/internal/cache"
	sdkmHttp "github.com/itbasis/go-tools-sdkm/internal/http"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/pkg/errors"
)

type versions struct {
	urlReleases string

	httpClient *http.Client

	muParsing sync.Mutex

	contentReleases string

	reStableGroupVersions   *regexp.Regexp
	reUnstableGroupVersions *regexp.Regexp
	reArchivedGroupVersions *regexp.Regexp
	reGoVersion             *regexp.Regexp

	cache sdkmSDKVersion.Cache
}

func NewVersions(urlReleases string) sdkmSDKVersion.SDKVersions {
	return &versions{
		urlReleases: urlReleases,
		cache:       sdkmCache.NewCache(),

		httpClient: sdkmHttp.NewHTTPClient(sdkmHttp.DefaultTimeout),

		reStableGroupVersions:   regexp.MustCompile(`<h2 id="stable">.*?<h2`),
		reUnstableGroupVersions: regexp.MustCompile(`<h2 id="unstable">.*?<div.*?id="archive"`),
		reArchivedGroupVersions: regexp.MustCompile(`id="archive">.+?</article`),
		reGoVersion:             regexp.MustCompile(`id="go(.+?)"`),
	}
}

func (receiver *versions) WithCache(cache sdkmSDKVersion.Cache) sdkmSDKVersion.SDKVersions {
	slog.Debug("setting cache: " + cache.GoString())

	receiver.cache = cache

	return receiver
}

func (receiver *versions) getContent(url string) (string, error) {
	slog.Debug(fmt.Sprintf("getting content for url: %s", url))

	resp, errGet := receiver.httpClient.Get(url) //nolint:noctx // TODO
	if errGet != nil {
		return "", errGet //nolint:wrapcheck // TODO
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("AttrError closing body after receiving content", itbasisCoreLog.SlogAttrError(err))
		}
	}()

	body, errReadAll := io.ReadAll(resp.Body)

	content := strings.ReplaceAll(string(body), "\n", "")

	slog.Debug(fmt.Sprintf("received content size: %d", len(content)))

	return content, errors.Wrap(errReadAll, ErrMsgFailedGetVersion)
}

func (receiver *versions) parseVersions(
	ctx context.Context,
	versionType sdkmSDKVersion.VersionType,
	reGroupVersions *regexp.Regexp,
	cleanContent bool,
) {
	receiver.muParsing.Lock()
	defer receiver.muParsing.Unlock()

	if receiver.contentReleases == "" {
		receiver.contentReleases, _ = receiver.getContent(receiver.urlReleases)

		if cleanContent {
			defer receiver.cleanContent()
		}
	}

	var content = reGroupVersions.FindString(receiver.contentReleases)
	if content == "" {
		slog.Debug(fmt.Sprintf("content is empty for version: %s", versionType))

		return
	}

	slog.Debug(fmt.Sprintf("found groups for version type: %s", versionType))

	var (
		m           = receiver.reGoVersion.FindAllStringSubmatch(content, -1)
		sdkVersions = make([]sdkmSDKVersion.SDKVersion, len(m))
	)

	for i, row := range m {
		if row[1] != "" {
			sdkVersion := sdkmSDKVersion.NewSDKVersion(row[1], versionType, false)

			sdkVersions[i] = sdkVersion
		}
	}

	slog.Debug(fmt.Sprintf("found %d SDK versions for version type: %s", len(sdkVersions), versionType))

	receiver.cache.Store(ctx, versionType, sdkVersions)
}

func (receiver *versions) cleanContent() {
	receiver.contentReleases = ""
}

//nolint:unparam // the argument "stable" must be left
func (receiver *versions) updateCache(ctx context.Context, stable, unstable, archived bool) {
	if stable {
		receiver.parseVersions(ctx, sdkmSDKVersion.TypeStable, receiver.reStableGroupVersions, false)
	}

	if unstable {
		receiver.parseVersions(ctx, sdkmSDKVersion.TypeUnstable, receiver.reUnstableGroupVersions, false)
	}

	if archived {
		receiver.parseVersions(ctx, sdkmSDKVersion.TypeArchived, receiver.reArchivedGroupVersions, false)
	}

	receiver.cleanContent()
}

func (receiver *versions) GoString() string {
	return "versions{" + receiver.urlReleases + "}"
}
