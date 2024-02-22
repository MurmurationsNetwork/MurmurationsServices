package config

var Values = config{}

type config struct {
	Library libraryConf
	Mongo   mongoConf
	Redis   redisConf
	Github  githubConf
	IsLocal bool `env:"IS_LOCAL,required"`
}

type libraryConf struct {
	URL string `env:"LIBRARY_URL,required"`
}

type mongoConf struct {
	USERNAME string `env:"MONGO_USERNAME,required"`
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	HOST     string `env:"MONGO_HOST,required"`
	DBName   string `env:"MONGO_DB_NAME,required"`
}

type redisConf struct {
	URL string `env:"REDIS_URL,required"`
}

type githubConf struct {
	TOKEN     string `env:"GITHUB_TOKEN,required"`
	BranchURL string `env:"GITHUB_BRANCH_URL,required"`
	TreeURL   string `env:"GITHUB_TREE_URL,required"`
}
