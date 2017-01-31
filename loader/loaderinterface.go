package loader

type Loader interface {
	Load(config interface{}, file string) error
	Dump(config interface{}, file string) error
}
