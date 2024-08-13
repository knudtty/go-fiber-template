package state

import "my_project/platform/database"

type AppState struct {
	DB *database.Stores
}
