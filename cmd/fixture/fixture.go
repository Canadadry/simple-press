package fixture

import (
	"app/cmd/migration"
	"app/config"
	"app/fixtures"
	"app/pkg/clock"
	"app/pkg/http/httpcaller"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	Action = "fixture"
)

var FixedNow = "2025-05-02T17:51:53+02:00"

func Run(c config.Parameters) error {
	_ = os.Remove(c.DatabaseUrl)
	err := migration.Run(c)
	if err != nil {
		return fmt.Errorf("migration failed : %w", err)
	}
	users := []fixtures.Article{}
	now, err := time.Parse(config.SerializerDateTimeFormat, FixedNow)
	if err != nil {
		return fmt.Errorf("invalid fixed clock date '%s': %w", FixedNow, err)
	}
	env, err := fixtures.Run(
		httpcaller.New(fmt.Sprintf("http://localhost:%d", c.Port), http.DefaultClient),
		&clock.Fixed{At: now},
		users,
	)
	if err != nil {
		return fmt.Errorf("fixture failed : %w", err)
	}
	env.SaveAt("tests/fixture.json")
	return os.Rename(c.DatabaseUrl, "tests/"+filepath.Base(c.DatabaseUrl))
}
