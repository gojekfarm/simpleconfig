package simpleconfig

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestExtractFieldsGetsAllConfigs(t *testing.T) {
	type singleLevel struct {
		F string
		G int
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
	testConfig := TestConfig{}
	listOfKeys := extractFields(testConfig)
	expectedList := []string{
		"A", "B", "C_F", "C_G", "D_E_F", "D_E_G", "D_H", "D_J",
	}

	assert.ElementsMatch(t, expectedList, listOfKeys, "Invalid keys found")
}

func TestPopulateValueFailsForNonPointers(t *testing.T) {
	type TestConfig struct {
		A string
		B int
	}

	testConfig := TestConfig{}
	err := populateValue("A", testConfig, "hello")
	assert.NotNil(t, err, "No error seen")

}

func TestPopulateValueSucceedsForPointers(t *testing.T) {
	type TestConfig struct {
		A string
		B int
	}

	testConfig := TestConfig{}
	err := populateValue("A", &testConfig, "hello")
	assert.Nil(t, err, "Error Seen")
	assert.Equal(t, "hello", testConfig.A, "Invalid A value")
	err = populateValue("B", &testConfig, "10")
	assert.Nil(t, err, "Error Seen")
	assert.Equal(t, 10, testConfig.B, "Invalid A value")

}

func TestPopulateValueSucceedsFor2Levels(t *testing.T) {
	type singleLevel struct {
		F string
		G int
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
	testConfig := TestConfig{}
	err := populateValue("C_F", &testConfig, "hello")
	assert.Nil(t, err, "Error Seen")
	assert.Equal(t, "hello", testConfig.C.F, "Invalid A value")
	err = populateValue("C_G", &testConfig, "10")
	assert.Nil(t, err, "Error Seen")
	assert.Equal(t, 10, testConfig.C.G, "Invalid A value")

}

func TestPopulateValueSucceedsFor3Levels(t *testing.T) {
	type singleLevel struct {
		F string
		G int
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
	testConfig := TestConfig{}
	err := populateValue("D_E_F", &testConfig, "hello")
	assert.Nil(t, err, "Error Seen")
	assert.Equal(t, "hello", testConfig.D.E.F, "Invalid A value")
	err = populateValue("D_E_G", &testConfig, "10")
	assert.Nil(t, err, "Error Seen")
	assert.Equal(t, 10, testConfig.D.E.G, "Invalid A value")

}

func TestTemp(t *testing.T) {
	type TestConfig struct {
		A string
		B int
	}
	testConfig := &TestConfig{}
	assert.True(t, reflect.ValueOf(testConfig).Elem().FieldByName("A").CanSet(), "Cant set")
	reflect.ValueOf(testConfig).Elem().FieldByName("A").SetString("Param")
	assert.Equal(t, "Param", testConfig.A, "Invalid A value")
}

func TestLoadConfig(t *testing.T) {
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

	testConfig := TestConfig{}

	os.Setenv("C_F", "hello")
	os.Setenv("C_G", "10")
	os.Setenv("C_K", "false")
	os.Setenv("D_E_F", "hello")
	os.Setenv("D_E_G", "15")
	os.Setenv("D_E_K", "true")

	LoadConfig(&testConfig)

	assert.Equal(t, "hello", testConfig.C.F, "Invalid A value")
	assert.Equal(t, 10, testConfig.C.G, "Invalid A value")
	assert.False(t, testConfig.C.K, "Invalid A value")
	assert.Equal(t, "hello", testConfig.D.E.F, "Invalid A value")
	assert.Equal(t, 15, testConfig.D.E.G, "Invalid A value")
	assert.True(t, testConfig.D.E.K, "Invalid A value")

	os.Unsetenv("C_F")
	os.Unsetenv("C_G")
	os.Unsetenv("C_K")
	os.Unsetenv("D_E_F")
	os.Unsetenv("D_E_G")
	os.Unsetenv("D_E_K")

}

func TestLoadRealisticConfig(t *testing.T) {
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

	LoadConfig(&testConfig)

	assert.Equal(t, "localhost", testConfig.Database.Host, "Incorrect Database Host")
	assert.Equal(t, 5432, testConfig.Database.Port, "Incorrect Database Host")
	assert.Equal(t, 35, testConfig.Database.MaxConns, "Incorrect Database Host")

	os.Unsetenv("DATABASE_HOST")
	os.Unsetenv("DATABASE_PORT")
	os.Unsetenv("DATABASE_MAXCONNS")

}

func TestLoadConfigTagDefaults(t *testing.T) {
	type DBConfig struct {
		Host     string `d:"localhost"`
		Port     int    `d:"5432"`
		MaxConns int    `d:"35"`
	}

	type TestConfig struct {
		Database DBConfig
	}

	testConfig := TestConfig{}

	LoadConfig(&testConfig)

	assert.Equal(t, "localhost", testConfig.Database.Host, "Incorrect Database Host")
	assert.Equal(t, 5432, testConfig.Database.Port, "Incorrect Database Host")
	assert.Equal(t, 35, testConfig.Database.MaxConns, "Incorrect Database Host")

}

// Test that Environment variables over-ride defaults even when empty.
func TestLoadConfigTagSetEnvironmentVsDefaults(t *testing.T) {
	type DBConfig struct {
		Host     string `d:"localhost"`
		Port     int    `d:"5432"`
		MaxConns int    `d:"35"`
	}

	type TestConfig struct {
		Database DBConfig
	}

	testConfig := TestConfig{}

	os.Setenv("DATABASE_HOST", "")
	os.Setenv("DATABASE_PORT", "9090")

	LoadConfig(&testConfig)

	assert.Equal(t, "", testConfig.Database.Host, "Incorrect Database Host")
	assert.Equal(t, 9090, testConfig.Database.Port, "Incorrect Database Host")
	assert.Equal(t, 35, testConfig.Database.MaxConns, "Incorrect Database Host")

	os.Unsetenv("DATABASE_HOST")
	os.Unsetenv("DATABASE_PORT")

}
