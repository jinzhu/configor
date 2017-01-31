package configor

type loader interface {
	Load(config interface{}, file string) error
	Dump(config interface{}, file string) error
}
