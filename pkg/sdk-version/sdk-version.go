package sdkversion

import (
	"encoding/json"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/Masterminds/semver/v3"
	itbasisCoreLog "github.com/itbasis/go-tools-core/log"
	"github.com/pkg/errors"
)

const (
	_msgFailSemVer = "fail init semver for version [%s]"
)

var _reSemVer = regexp.MustCompile(`^(0|[1-9]\d*)(\.(0|[1-9]\d*))?(\.(0|[1-9]\d*))?((rc|beta)\d+)?$`)

type VersionType string

const (
	TypeStable   VersionType = "stable"
	TypeUnstable VersionType = "unstable"
	TypeArchived VersionType = "archived"
)

type SDKVersion struct {
	id          string
	versionType VersionType

	installed bool
	semVer    *semver.Version
}

type _jsonSdkVersion struct {
	ID   string      `json:"id"`
	Type VersionType `json:"type"`
}

var EmptySdkVersion = SDKVersion{}

func NewSDKVersion(id string, versionType VersionType, installed bool) SDKVersion {
	var sdkVersion = SDKVersion{id: id, versionType: versionType, installed: installed}

	if sv, err := makeSemVer(sdkVersion.id); err != nil {
		slog.Error("fail parsing semver", itbasisCoreLog.SlogAttrError(err))
	} else {
		sdkVersion.semVer = sv
	}

	return sdkVersion
}

func (r *SDKVersion) GetId() string             { return r.id }
func (r *SDKVersion) GetType() VersionType      { return r.versionType }
func (r *SDKVersion) HasInstalled() bool        { return r.installed }
func (r *SDKVersion) SetInstalled(flag bool)    { r.installed = flag }
func (r SDKVersion) GetSemVer() *semver.Version { return r.semVer }

func (r *SDKVersion) UnmarshalJSON(data []byte) error {
	var m = &_jsonSdkVersion{}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}

	r.id = m.ID
	r.versionType = m.Type
	if sv, err := makeSemVer(r.id); err != nil {
		return err
	} else {
		r.semVer = sv
	}

	return nil
}
func (r SDKVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(_jsonSdkVersion{ID: r.id, Type: r.versionType})
}

func makeSemVer(id string) (*semver.Version, error) {
	if id == "" {
		return nil, semver.ErrEmptyString
	}

	m := _reSemVer.FindStringSubmatch(id)

	if len(m) < 1 {
		return nil, semver.ErrInvalidSemVer
	}

	var (
		err                       error
		svMajor, svMinor, svPatch uint64
		svPre                     string
	)

	if svMajor, err = strconv.ParseUint(m[1], 10, 64); err != nil {
		return nil, errors.WithMessagef(err, _msgFailSemVer, id)
	}

	if s := m[3]; len(s) > 0 {
		if svMinor, err = strconv.ParseUint(m[3], 10, 64); err != nil {
			return nil, errors.WithMessagef(err, _msgFailSemVer, id)
		}
	}

	if s := m[5]; len(s) > 0 {
		if svPatch, err = strconv.ParseUint(s, 10, 64); err != nil {
			return nil, errors.WithMessagef(err, _msgFailSemVer, id)
		}
	}

	svPre = m[6]

	return semver.New(svMajor, svMinor, svPatch, svPre, ""), nil
}
