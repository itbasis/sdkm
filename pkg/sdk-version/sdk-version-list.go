package sdkversion

import (
	"encoding/json"
	"iter"
	"sort"

	"github.com/pkg/errors"
)

type SdkVersionList interface {
	sort.Interface

	IsEmpty() bool
	Add(sdkVersions ...SDKVersion)

	First() (SDKVersion, error)
	Seq() iter.Seq[SDKVersion]

	AsList() []SDKVersion
	AsMap() map[string]SDKVersion

	json.Unmarshaler
	json.Marshaler
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

func (r *_sdkVersionList) IsEmpty() bool { return len(r.list) == 0 }

func (r *_sdkVersionList) First() (SDKVersion, error) {
	if len(r.list) == 0 {
		return EmptySdkVersion, errors.WithMessage(ErrSDKVersionNotFound, "list is empty")
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
	if len(r.list) == 0 {
		return false
	}

	var (
		sdkI = r.list[i]
		sdkJ = r.list[j]

		typeI = sdkI.GetType()
		typeJ = sdkJ.GetType()
	)

	if sdkI.GetType() == sdkJ.GetType() {
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

func (r *_sdkVersionList) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.list)
}

func (r *_sdkVersionList) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.list)
}
