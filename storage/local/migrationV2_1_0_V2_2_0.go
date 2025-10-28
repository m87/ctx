package localstorage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
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

	return nil
}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) MigrateConfig() error {
	log.Println(`Config migration v2.1.0 -> v2.2.0
	Migration plan:
	- Change version to 2.2.0
	`)

	return nil
}
