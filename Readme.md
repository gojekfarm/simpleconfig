# simple config

A library to remove the messy way we handle configs with viper cobra or other vile creatures

# Config Structure

It is assumed that your app has a config structure like this

```
type Config struct {
  Config1  ConfigA
  Config2  ConfigB
  Config3  ConfigC
}
```

Where ConfigA is a structure holding the information for Config1.

# Setting Config

For example

```
	type DBConfig struct {
		Host     string
		Port     int
		MaxConns int
	}

	type TestConfig struct {
		Database DBConfig
	}

	testConfig := TestConfig{}

	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("DATABASE_MAXCONNS", "35")

```

Multilevel example

```
	type singleLevel struct {
		F string
		G int
		K bool
	}

	type doubleLevel struct {
		E singleLevel
		H int
		J string
	}

	type TestConfig struct {
		A string
		B int
		C singleLevel
		D doubleLevel
	}

To set TestConfig.A you set the environment variable `A`
To set TestConfig.C.F you set environment variable `C_F`
To set TestConfig.D.E.K you set environment variable `D_E_K`
```
# Defaults

To add a default just add it as an annotation
```
	type DBConfig struct {
		Host     string `d:"localhost"`
		Port     int    `d:"5432"`
		MaxConns int    `d:"35"`
	}
```

Please Note the syntax here `d:"default"`. The `""` are important in the value as that is the syntax supported by reflect tag library.

# Supported data types

Presently
1. string
2. int
3. bool

More types can be added by modifying the functions `populateDefaultValue` and `populateValue` Look for 
`/* Add a new type here! */`

# Code

```
import "simpleconfig"

// Define your structure
	type singleLevel struct {
		F string
		G int
		K bool
	}

	type doubleLevel struct {
		E singleLevel
		H int
		J string
	}

	type TestConfig struct {
		A string
		B int
		C singleLevel
		D doubleLevel
	}

func main() {
  config := TestConfig{}
  simpleconfig.LoadConfig(&config)
}

```
_voila_

Use a library like godotenv or config tool to set the environment variables on docker

# Status

All tests pass
