package sdkversion

import (
	"fmt"
	"iter"
	"log/slog"
	"sort"

	"github.com/pkg/errors"
)

type SdkVersionList interface {
	sort.Interface

	Add(sdkVersions ...SDKVersion)
	First() (SDKVersion, error)
	Seq() iter.Seq[SDKVersion]

	AsList() []SDKVersion
	AsMap() map[string]SDKVersion
}

type _sdkVersionList struct {
	list []SDKVersion
}

func NewSdkVersionList(list ...SDKVersion) SdkVersionList {
	return &_sdkVersionList{list: list}
}

func (r *_sdkVersionList) Add(sdkVersion ...SDKVersion) {
	r.list = append(r.list, sdkVersion...)
}

func (r *_sdkVersionList) First() (SDKVersion, error) {
	if len(r.list) == 0 {
		return nil, errors.WithMessage(ErrSDKVersionNotFound, "list is empty")
	}

	return r.list[0], nil
}

func (r *_sdkVersionList) AsList() []SDKVersion {
	return r.list
}

func (r *_sdkVersionList) AsMap() map[string]SDKVersion {
	var m = make(map[string]SDKVersion, r.Len())

	for _, el := range r.list {
		m[el.GetId()] = el
	}

	return m
}

func (r *_sdkVersionList) Seq() iter.Seq[SDKVersion] {
	return func(yield func(SDKVersion) bool) {
		for i := range r.list {
			if !yield(r.list[i]) {
				return
			}
		}
	}
}

func (r _sdkVersionList) Len() int {
	return len(r.list)
}

func (r *_sdkVersionList) Less(i, j int) (result bool) {
	var (
		sdkI = r.list[i]
		sdkJ = r.list[j]

		idI   = sdkI.GetId()
		idJ   = sdkJ.GetId()
		typeI = sdkI.GetType()
		typeJ = sdkJ.GetType()
	)

	slog.Debug("comparison to less", slog.String("sdk[i]", idI), slog.String("sdk[j]", idJ))

	defer func() {
		slog.Debug(fmt.Sprintf("%s less %s: %t", idI, idJ, result))
	}()

	if sdkI.GetType() == sdkJ.GetType() {
		slog.Debug(fmt.Sprintf("sort by semver: i=%s, j=%s", sdkI.GetId(), sdkJ.GetId()))

		return sdkI.GetSemVer().GreaterThan(sdkJ.GetSemVer())
	}

	if typeI == TypeStable {
		return true
	}

	if typeI == TypeUnstable {
		return typeJ == TypeArchived
	}

	return false
}

func (r *_sdkVersionList) Swap(i, j int) {
	r.list[i], r.list[j] = r.list[j], r.list[i]
}
