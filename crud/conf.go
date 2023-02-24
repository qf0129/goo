package crud

var defaultConf = &Config{
	PrimaryKey:       "id",
	DefaultPageIndex: 1,
	defaultPageSize:  10,
}

type Config struct {
	PrimaryKey       string
	DefaultPageIndex int
	defaultPageSize  int
}

func (conf *Config) Update(new *Config) {
	if new.PrimaryKey != "" {
		conf.PrimaryKey = new.PrimaryKey
	}
	if new.DefaultPageIndex == 0 {
		conf.DefaultPageIndex = new.DefaultPageIndex
	}
	if new.defaultPageSize == 0 {
		conf.defaultPageSize = new.defaultPageSize
	}
}
