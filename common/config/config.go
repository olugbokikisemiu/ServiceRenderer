package config

import (
	"github.com/lalamove/konfig"
	"github.com/lalamove/konfig/loader/klfile"
	"github.com/lalamove/konfig/parser/kpyaml"
)

var configFiles = []klfile.File{
	{
		Path:   "config.yaml",
		Parser: kpyaml.Parser,
	},
	// {
	// 	Path:   "../secrets.json",
	// 	Parser: kpjson.Parser,
	// },
}

func init() {
	konfig.Init(konfig.DefaultConfig())
}

func LoadAndWatch() error {
	// load from yaml file
	konfig.RegisterLoaderWatcher(
		klfile.New(&klfile.Config{
			Files: configFiles,
			Watch: true,
		}),
	)

	return konfig.LoadWatch()
}

// Load config relative to specified path mostly for test
func LoadFromPath(path string) error {
	// load from yaml file
	konfig.RegisterLoaderWatcher(
		klfile.New(&klfile.Config{
			Files: []klfile.File{
				{
					Path:   path,
					Parser: kpyaml.Parser,
				},
			},
			Watch: true,
		}),
	)
	return konfig.LoadWatch()
}

func Exists(k string) bool {
	return konfig.Exists(k)
}

func Get(k string) interface{} {
	return konfig.Get(k)
}

func String(k string) string {
	return konfig.String(k)
}

func MustString(k string) string {
	return konfig.MustString(k)
}

func Int(k string) int {
	return konfig.Int(k)
}

func Float(k string) float64 {
	return konfig.Float(k)
}

func Bool(k string) bool {
	return konfig.Bool(k)
}

func MustBool(k string) bool {
	return konfig.MustBool(k)
}
