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

```

To set TestConfig.A you set the environment variable `A`
To set TestConfig.C.F you set environment variable `C_F`
To set TestConfig.D.E.K you set environment variable `D_E_K`

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
