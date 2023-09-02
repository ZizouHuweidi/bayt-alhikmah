package middleware

import (
	"os"
	"testing"

	"github.com/zizouhuweidi/bayt-alhikmah/config"
	"github.com/zizouhuweidi/bayt-alhikmah/ent"
	"github.com/zizouhuweidi/bayt-alhikmah/pkg/services"
	"github.com/zizouhuweidi/bayt-alhikmah/pkg/tests"
)

var (
	c   *services.Container
	usr *ent.User
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvTest)

	// Create a new container
	c = services.NewContainer()

	// Create a user
	var err error
	if usr, err = tests.CreateUser(c.ORM); err != nil {
		panic(err)
	}

	// Run tests
	exitVal := m.Run()

	// Shutdown the container
	if err = c.Shutdown(); err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}
