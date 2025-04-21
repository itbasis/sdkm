package log

import "log/slog"

func SlogAttrRootDir(sdkRootDir string) slog.Attr {
	return slog.String("root_dir", sdkRootDir)
}

func SlogAttrFilePath(filePath string) slog.Attr {
	return slog.String("file_path", filePath)
}
