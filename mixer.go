package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Volume struct {
	All      float32
	Props    float32
	Ambience float32
	Movement float32
	UI       float32
	Music    float32
}

var GameVolume = Volume{
	All:      1.0,
	Props:    1.0,
	Ambience: 1.0,
	Movement: 1.0,
	UI:       1.0,
	Music:    1.0,
}

// UserVolume is a volume, setted by the user, but not saved.
// @Incomplete: Save it in saves. It is possible to store that in GameSettings
// And access in here, but them is need to set proper load order, so global game settings
// will be loaded first. For now it other way.
var UserVolume = Volume{
	All:      1.0,
	Props:    1.0,
	Ambience: 1.0,
	Movement: 1.0,
	UI:       1.0,
	Music:    1.0,
}

func (v *Volume) SetVolume(category SoundCategory, volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	// @Cleanup: Cant make Volume map, because vars.go does not support maps yet
	switch category {
	case All:
		v.All = volume
	case Props:
		v.Props = volume
	case Ambience:
		v.Ambience = volume
	case PlayerMovement:
		v.Movement = volume
	case UI:
		v.UI = volume
	case Music:
		v.Music = volume
	}

	AudioManager.UpdateVolume()
}

func (v *Volume) GetVolume(category SoundCategory) float32 {
	switch category {
	case All:
		return v.All
	case Props:
		return v.Props
	case Ambience:
		return v.Ambience
	case PlayerMovement:
		return v.Movement
	case UI:
		return v.UI
	case Music:
		return v.Music
	}

	return 0
}

type SoundCategory int

const (
	All SoundCategory = iota
	Props
	Ambience
	PlayerMovement
	UI
	Music
)

type AudioFile struct {
	Name     string
	Category SoundCategory
}

type SoundFile struct {
	AudioFile
	Sound rl.Sound
}

type MusicFile struct {
	AudioFile
	Music rl.Music
}

type Mixer struct {
	Sounds map[string]SoundFile
	Musics map[string]MusicFile
}

func NewMixer() *Mixer {
	rl.InitAudioDevice()
	m := &Mixer{
		Sounds: loadSoundFiles(),
		Musics: loadMusicFiles(),
	}

	m.UpdateVolume()

	return m
}

func loadSoundFiles() map[string]SoundFile {
	return map[string]SoundFile{}
}

func loadMusicFiles() map[string]MusicFile {
	return map[string]MusicFile{
		"music": {
			AudioFile: AudioFile{
				Name:     "music",
				Category: Music,
			},
			Music: rl.LoadMusicStream(MusicFolder + "music.mp3"),
		},
	}
}

func GetGameUserVolume(category SoundCategory) float32 {
	userCategoryVolume := UserVolume.GetVolume(category)
	gameCategoryVolume := GameVolume.GetVolume(category)

	return userCategoryVolume * gameCategoryVolume * UserVolume.All * GameVolume.All
}

func (m *Mixer) UpdateVolume() {
	for _, sound := range m.Sounds {
		rl.SetSoundVolume(sound.Sound, GetGameUserVolume(sound.Category))
	}

	for _, music := range m.Musics {
		rl.SetMusicVolume(music.Music, GetGameUserVolume(music.Category))
	}
}

func (m *Mixer) PlaySound(sound string) {
	if _, ok := m.Sounds[sound]; !ok {
		log.Printf("Sound %s not found", sound)
		return
	}

	rl.PlaySound(m.Sounds[sound].Sound)
}

func (m *Mixer) PlayMusic(music string) {
	if _, ok := m.Musics[music]; !ok {
		log.Printf("Music %s not found", music)
		return
	}

	rl.PlayMusicStream(m.Musics[music].Music)
}

func (m *Mixer) SetVolume(volume float32) {
	// Do we not clamp value here on purpose?
	UserVolume.All = volume
	m.UpdateVolume()
}

func (m *Mixer) SetMusicLoop(music string, loop bool) {
	if _, ok := m.Musics[music]; !ok {
		log.Printf("Music %s not found", music)
		return
	}

	file := m.Musics[music]
	file.Music.Looping = loop

	m.Musics[music] = file
}

func (m *Mixer) UpdateMusic() {
	rl.UpdateMusicStream(m.Musics["music"].Music)
}

func (m *Mixer) Unload() {
	for _, sound := range m.Sounds {
		rl.UnloadSound(sound.Sound)
	}

	for _, music := range m.Musics {
		rl.UnloadMusicStream(music.Music)
	}

	rl.CloseAudioDevice()
}
