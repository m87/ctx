package localstorage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"github.com/spf13/viper"
)

type LocalStorageMigratorV2_1_0_V_2_2_0 struct {
	statePath   string
	archivePath string
}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) Id() string {
	return "V2_1_0_V_2_2_0"
}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) Migrate() error {
	log.Println(`Migration v2.1.0 -> v2.2.0
	Migration plan:
	- Convert context comments to objects with ids
	`)

	log.Println("Migrating contexts in state...")
	var state map[string]any

	log.Println("Loading state from", migrator.statePath)
	f, err := os.Open(migrator.statePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	decoder.UseNumber()
	if err := decoder.Decode(&state); err != nil {
		return err
	}

	log.Println("Transforming context comments...")
	contexts, _ := state["contexts"].(map[string]any)
	for ctxId, ctxData := range contexts {
		ctxMap, ok := ctxData.(map[string]any)
		if !ok {
			continue
		}

		comments, _ := ctxMap["comments"].([]string)
		commentObjects := map[string]core.Comment{}
		for _, commentContent := range comments {
			commentId := uuid.New().String()
			commentObjects[commentId] = core.Comment{
				Id:      commentId,
				Content: commentContent,
			}
		}
		ctxMap["comments"] = commentObjects
		contexts[ctxId] = ctxMap
	}
	state["contexts"] = contexts

	log.Println("Saving migrated state to", migrator.statePath)
	data, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(migrator.statePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Context state migration completed.")
	return nil
}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) MigrateArchive() error {
	log.Println(`Archive migration v2.1.0 -> v2.2.0
	Migration plan:
	- Convert context comments to objects with ids
	`)

	log.Println("Migrating contexts in archive...")
	entries, err := os.ReadDir(migrator.archivePath)
	if err != nil {
		panic(err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".ctx") {
			files = append(files, filepath.Join(migrator.archivePath, entry.Name()))
		}
	}

	for _, filePath := range files {
		log.Println("Migrating archive file:", filePath)
		var context map[string]any
		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		decoder.UseNumber()
		if err := decoder.Decode(&context); err != nil {
			return err
		}

		log.Println("Transforming context comments...")
		comments, _ := context["comments"].([]string)
		commentObjects := map[string]core.Comment{}
		for _, commentContent := range comments {
			commentId := uuid.New().String()
			commentObjects[commentId] = core.Comment{
				Id:      commentId,
				Content: commentContent,
			}
		}
		context["comments"] = commentObjects

		log.Println("Saving migrated context to", filePath)
		data, err := json.Marshal(context)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Context state migration completed.")
	return nil
}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) MigrateConfig() error {
	log.Println(`Config migration v2.1.0 -> v2.2.0
	Migration plan:
	- Change version to 2.2.0
	`)

	viper.Set("version", "2.2.0")

	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			_ = viper.SafeWriteConfig()
		} else {
			log.Panicf("Updated config cannot be saved")
		}
	}

	log.Println("Config migration completed.")
	return nil
}
