package sdkversion

import (
	"log/slog"
)

type PrintFormatOptions struct {
	OutputType bool

	OutputInstalled    bool
	OutputNotInstalled bool
}

func (receiver *SDKVersion) Print() string {
	return receiver.PrintWithOptions(true, false, true)
}

func (receiver *SDKVersion) PrintWithOptions(outType, outInstalled, outNotInstalled bool) string {
	if receiver == nil || receiver.GetId() == "" {
		return ""
	}

	var out = receiver.GetId()

	if outType {
		out += _sprintType(receiver.GetType())
	}

	if outInstalled && receiver.installed {
		out += " [installed]"
	} else if outNotInstalled && !receiver.installed {
		out += " [not installed]"
	}

	return out
}

func _sprintType(versionType VersionType) string {
	switch versionType {
	case TypeStable:
		slog.Debug("skip SDK Version Type: " + string(versionType))

	case TypeUnstable, TypeArchived:
		return " (" + string(versionType) + ")"

	default:
		slog.Error("unknown SDK version type: " + string(versionType))
	}

	return ""
}
