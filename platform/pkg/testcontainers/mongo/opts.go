package mongo

type Option func(*Config)

func WithNetworkName(name string) Option {
	return func(c *Config) {
		c.NetworkName = name
	}
}

func WithContainerName(name string) Option {
	return func(c *Config) {
		c.ContainerName = name
	}
}

func WithImageName(name string) Option {
	return func(c *Config) {
		c.ImageName = name
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

func WithUsername(username string) Option {
	return func(c *Config) {
		c.Username = username
	}
}

func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

func WithAuthDB(authDB string) Option {
	return func(c *Config) {
		c.AuthDB = authDB
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}