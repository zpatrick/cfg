package cfg

type FileFormat string

const (
	FormatINI  FileFormat = "ini"
	FormatTOML FileFormat = "toml"
	FormatYAML FileFormat = "yaml"
	FormatJSON FileFormat = "json"
)

func File(format FileFormat, path string) Provider {
	return nil
}
