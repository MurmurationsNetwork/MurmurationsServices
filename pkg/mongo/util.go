package mongo

func GetURI(username string, password string, host string) string {
	return "mongodb://" + username + ":" + password + "@" + host
}
