package sdkversion

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

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

type SDKVersion interface {
	GetId() string
	GetType() VersionType
	HasInstalled() bool
	SetInstalled(flag bool)

	Print() string
	PrintWithOptions(outType, outInstalled, outNotInstalled bool) string

	GetSemVer() *semver.Version
}

type _sdkVersion struct {
	ID   string
	Type VersionType

	installed bool            `json:"-"`
	semVer    *semver.Version `json:"-"`
}

func NewSDKVersion(id string, versionType VersionType, installed bool) SDKVersion {
	var sdkVersion = &_sdkVersion{ID: id, Type: versionType, installed: installed}

	if sv, err := makeSemVer(sdkVersion.ID); err != nil {
		slog.Error("fail parsing semver", itbasisCoreLog.SlogAttrError(err))
	} else {
		sdkVersion.semVer = sv
	}

	return sdkVersion
}

func (r *_sdkVersion) GetId() string             { return r.ID }
func (r *_sdkVersion) GetType() VersionType      { return r.Type }
func (r *_sdkVersion) HasInstalled() bool        { return r.installed }
func (r *_sdkVersion) SetInstalled(flag bool)    { r.installed = flag }
func (r _sdkVersion) GetSemVer() *semver.Version { return r.semVer }

func makeSemVer(id string) (*semver.Version, error) {
	if id == "" {
		return nil, semver.ErrEmptyString
	}

	slog.Debug("sdk version id=" + id)

	m := _reSemVer.FindStringSubmatch(id)

	slog.Debug(fmt.Sprintf("matches (len=%d): %s", len(m), strings.Join(m, ",")))

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
