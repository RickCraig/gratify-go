package helpers

import "os"

// Gets an environment variable but sets a
// default in case it hasn't been set
//
//	common.EnvVariable("MONGO_URI", "mongodb+srv://localhost:27017")
//
// Returns the string value of the environment variable
func EnvVariable(key string, d string) string {
	os.Setenv(key, d)
	return os.Getenv(key)
}
