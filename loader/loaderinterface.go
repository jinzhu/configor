package loader

type Loader interface {
	Load(config interface{}, file string) error      // fills the config with unmarshalled data from the file
	PlainLoad(config interface{}, file string) error // fills the config with unmarshalled data from the file, without extension checking
	Dump(config interface{}, file string) error      // serializes config to file
}
