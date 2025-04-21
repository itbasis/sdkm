package plugin

import "log/slog"

func SlogAttrPluginId(id ID) slog.Attr {
	return slog.String("pluginId", string(id))
}
