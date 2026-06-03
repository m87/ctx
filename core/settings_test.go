package core

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/m87/nod"
)

type mockSettingsRepository struct {
	settings *Settings
	loadErr  error
	saveErr  error

	loadCalls int
	saveCalls int
	saved     *Settings
}

func (r *mockSettingsRepository) Load() (*Settings, error) {
	r.loadCalls++
	if r.loadErr != nil {
		return nil, r.loadErr
	}
	return r.settings, nil
}

func (r *mockSettingsRepository) Save(settings *Settings) error {
	r.saveCalls++
	r.saved = settings
	return r.saveErr
}

func TestSettingsManagerInitSettingsIfNotExistsCreatesDefaults(t *testing.T) {
	repo := &mockSettingsRepository{loadErr: gorm.ErrRecordNotFound}
	manager := NewSettingsManager(repo)

	err := manager.InitSettingsIfNotExists()

	require.NoError(t, err)
	require.Equal(t, 1, repo.loadCalls)
	require.Equal(t, 1, repo.saveCalls)
	want := map[string]string{
		"client.general.theme":    "light",
		"client.general.firstDay": "Monday",
	}
	require.Equal(t, want, repo.saved.raw)
	require.Same(t, repo.saved, manager.cache)
}

func TestSettingsManagerInitSettingsIfNotExistsDoesNotOverrideExisting(t *testing.T) {
	existing := &Settings{raw: map[string]string{"client.general.theme": "dark"}}
	repo := &mockSettingsRepository{
		settings: existing,
	}
	manager := NewSettingsManager(repo)

	err := manager.InitSettingsIfNotExists()

	require.NoError(t, err)
	require.Equal(t, 0, repo.saveCalls)
	require.Same(t, existing, manager.cache)
}

func TestSettingsManagerInitSettingsIfNotExistsReturnsLoadError(t *testing.T) {
	wantErr := errors.New("load failed")
	repo := &mockSettingsRepository{loadErr: wantErr}
	manager := NewSettingsManager(repo)

	err := manager.InitSettingsIfNotExists()

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 1, repo.loadCalls)
	require.Equal(t, 0, repo.saveCalls)
	require.Nil(t, manager.cache)
}

func TestSettingsManagerGetClientLoadsFiltersAndCaches(t *testing.T) {
	repo := &mockSettingsRepository{
		settings: &Settings{raw: map[string]string{
			"client.general.theme":    "dark",
			"client.general.firstDay": "Sunday",
			"database.path":           "/tmp/ctx.db",
		}},
	}
	manager := NewSettingsManager(repo)

	got, err := manager.GetClient()
	require.NoError(t, err)
	_, err = manager.GetClient()
	require.NoError(t, err)

	want := map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
	}
	require.Equal(t, want, got)
	require.Equal(t, 1, repo.loadCalls)
}

func TestSettingsManagerGetClientReturnsLoadError(t *testing.T) {
	wantErr := errors.New("load failed")
	manager := NewSettingsManager(&mockSettingsRepository{loadErr: wantErr})

	got, err := manager.GetClient()

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}

func TestSettingsManagerSaveClientSavesAndUpdatesCache(t *testing.T) {
	repo := &mockSettingsRepository{loadErr: gorm.ErrRecordNotFound}
	manager := NewSettingsManager(repo)
	settings := map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
	}

	err := manager.SaveClient(settings)

	require.NoError(t, err)
	require.Equal(t, 1, repo.saveCalls)
	require.Equal(t, settings, repo.saved.raw)
	require.Same(t, repo.saved, manager.cache)
	require.Equal(t, "dark", repo.saved.general.theme)
	require.Equal(t, "Sunday", repo.saved.general.firstDay)
}

func TestSettingsManagerSaveClientMergesWithExistingSettings(t *testing.T) {
	repo := &mockSettingsRepository{
		settings: &Settings{raw: map[string]string{
			"client.general.theme":    "light",
			"client.general.firstDay": "Sunday",
		}},
	}
	manager := NewSettingsManager(repo)

	err := manager.SaveClient(map[string]string{"client.general.theme": "dark"})

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
	}, repo.saved.raw)
	require.Equal(t, "dark", repo.saved.general.theme)
	require.Equal(t, "Sunday", repo.saved.general.firstDay)
}

func TestSettingsManagerSaveClientIgnoresNonClientSettings(t *testing.T) {
	repo := &mockSettingsRepository{loadErr: gorm.ErrRecordNotFound}
	manager := NewSettingsManager(repo)

	err := manager.SaveClient(map[string]string{
		"client.general.theme": "dark",
		"database.path":        "/tmp/ctx.db",
	})

	require.NoError(t, err)
	require.Equal(t, map[string]string{"client.general.theme": "dark"}, repo.saved.raw)
}

func TestSettingsManagerSaveClientReturnsSaveError(t *testing.T) {
	wantErr := errors.New("save failed")
	repo := &mockSettingsRepository{loadErr: gorm.ErrRecordNotFound, saveErr: wantErr}
	manager := NewSettingsManager(repo)

	err := manager.SaveClient(map[string]string{"client.general.theme": "dark"})

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, manager.cache)
}

func TestSettingsManagerGetClientKeyOnlyAllowsClientKeys(t *testing.T) {
	repo := &mockSettingsRepository{
		settings: &Settings{raw: map[string]string{"client.general.theme": "dark"}},
	}
	manager := NewSettingsManager(repo)

	got, err := manager.GetClientKey("database.path")
	require.NoError(t, err)
	require.Empty(t, got)
	require.Equal(t, 0, repo.loadCalls)

	got, err = manager.GetClientKey("client.general.theme")
	require.NoError(t, err)
	require.Equal(t, "dark", got)
}

func TestSettingsManagerGetKeyUsesCacheThenViperFallback(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	err := os.WriteFile(configPath, []byte("database:\n  path: /tmp/ctx.db\n"), 0600)
	require.NoError(t, err)
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	repo := &mockSettingsRepository{
		settings: &Settings{raw: map[string]string{"client.general.theme": "dark"}},
	}
	manager := NewSettingsManager(repo)

	got, err := manager.GetKey("client.general.theme")
	require.NoError(t, err)
	require.Equal(t, "dark", got)

	got, err = manager.GetKey("database.path")
	require.NoError(t, err)
	require.Equal(t, "/tmp/ctx.db", got)
	require.Equal(t, 1, repo.loadCalls)
}

func TestSettingsManagerSaveSavesAndUpdatesCache(t *testing.T) {
	repo := &mockSettingsRepository{}
	manager := NewSettingsManager(repo)
	settings := map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
		"database.path":           "/tmp/ctx.db",
	}

	err := manager.Save(settings)

	require.NoError(t, err)
	require.Equal(t, 1, repo.saveCalls)
	require.Equal(t, settings, repo.saved.raw)
	require.Same(t, repo.saved, manager.cache)
}

func TestSettingsMapperFromNodeRestoresRawSettings(t *testing.T) {
	theme := "dark"
	firstDay := "Sunday"
	mapper := NewSettingsMapper()
	node, err := mapper.ToNode(NewSettings(map[string]string{
		"client.general.theme":    theme,
		"client.general.firstDay": firstDay,
	}))
	require.NoError(t, err)

	settings, err := mapper.FromNode(node)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"client.general.theme":    theme,
		"client.general.firstDay": firstDay,
	}, settings.raw)
	require.Equal(t, theme, settings.general.theme)
	require.Equal(t, firstDay, settings.general.firstDay)
}

func TestSettingsMapperFromNodeHandlesMissingAndNilKV(t *testing.T) {
	theme := "dark"
	mapper := NewSettingsMapper()
	node := &nod.Node{
		KV: map[string]*nod.KV{
			"client.general.theme":    {Key: "client.general.theme", ValueText: &theme},
			"client.general.firstDay": {Key: "client.general.firstDay"},
		},
	}

	settings, err := mapper.FromNode(node)

	require.NoError(t, err)
	require.Equal(t, map[string]string{"client.general.theme": theme}, settings.raw)
	require.Equal(t, theme, settings.general.theme)
	require.Empty(t, settings.general.firstDay)
}
