package main

import (
	"encoding/json"
	"os"

	"github.com/denisbrodbeck/machineid"
)

// @Incomplete: Add game? saves here

type GameSettings struct {
	// @Incomplete: Add all game settings here
	Fullscreen bool

	// @Incompltete: Go to mixer and search UserVolume
	MusicVolume float32

	Debug bool
}

type User struct {
	Name string
	GameSettings
	// @Incomplete: Add game? saves here
}

// @Incomplete: Handle multiple user saves, or just one
// Maybe it will be cool if each user save will be bound to a steam profile id

// @Incoplete: Manage it more nicely
// @Cleanup: Refactor this

type UserManager struct {
	CurrentUser *User
}

// @Cleanup: We even don't need this. Just use CurrentUser string and json files as saves
func NewUserManager() *UserManager {
	return &UserManager{}
}

// @Cleanup: Maybe move it to NewUserManager()
func (u *UserManager) AddUser() {
	// If already exists, load
	if u.CurrentUser != nil {
		u.LoadUser()
		return
	}

	name, err := machineid.ID()
	if err != nil {
		name = "unknown"
	}

	if _, err = os.Stat(SavesFolder + name + ".json"); err == nil {
		// @Cleanup: Refactor this horrible code
		u.CurrentUser = &User{
			Name: name,
		}

		u.LoadUser()
		return
	}

	// Create new
	u.CurrentUser = &User{
		Name: name,
	}

	// Load default settings
	u.UpdateSettings()
	u.ApplySettings()
	u.SaveUser()
}

func (u *UserManager) SaveUser() {
	if u.CurrentUser == nil {
		u.AddUser()
	}

	json, err := json.Marshal(u.CurrentUser)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(SavesFolder + u.CurrentUser.Name + ".json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write(json)
	if err != nil {
		panic(err)
	}
}

func (u *UserManager) LoadUser() {
	// if u.CurrentUser == nil {
	// 	u.AddUser()
	// }

	file, err := os.ReadFile(SavesFolder + u.CurrentUser.Name + ".json")
	if err != nil {
		panic(err)
	}

	user := User{}
	err = json.Unmarshal(file, &user)
	if err != nil {
		panic(err)
	}

	u.CurrentUser = &user

	// Load user settings
	u.ApplySettings()
}

func (u *UserManager) UpdateSettings() {
	if u.CurrentUser == nil {
		u.AddUser()
	}

	user := u.CurrentUser

	user.MusicVolume = UserVolume.Music
	user.Debug = Debug.Visible
	user.Fullscreen = GameDisplay.Fullscreen
	// @Incomplete: Add more stuff to saves

	u.CurrentUser = user
}

func (u *UserManager) ApplySettings() {
	if u.CurrentUser == nil {
		u.AddUser()
	}

	UserVolume.SetVolume(Music, u.CurrentUser.MusicVolume)
	Debug.Visible = u.CurrentUser.Debug
	GameDisplay.Fullscreen = u.CurrentUser.Fullscreen
	// @Incomplete: Add more stuff to saves
}
