package tests

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var globalServer *MockServer

func TestMain(m *testing.M) {
	var err error
	globalServer, err = NewMockServer()
	if err != nil {
		panic("failed to start mock HEOS server: " + err.Error())
	}
	defer globalServer.Close()

	opts := godog.Options{
		Format:    "pretty",
		Paths:     []string{"../features"},
		Output:    colors.Colored(os.Stdout),
		Randomize: 0,
	}

	status := godog.TestSuite{
		Name:                "heos-cli",
		ScenarioInitializer: initializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}
