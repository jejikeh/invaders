package main

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Volumes = map[SoundCategory]float32{
	Props:          0.5,
	Ambience:       0.5,
	PlayerMovement: 0.5,
	Ui:             0.5,
	Music:          0.5,
}

func SetVolume(category SoundCategory, volume float32) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}

	Volumes[category] = roundFloat(volume, 1)
	AudioManager.UpdateVolume()
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}

type SoundCategory int

const (
	Props SoundCategory = iota
	Ambience
	PlayerMovement
	Ui
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
	Volume float32
}

func NewMixer() *Mixer {
	rl.InitAudioDevice()
	m := &Mixer{
		Sounds: loadSoundFiles(),
		Musics: loadMusicFiles(),
		Volume: 1.0,
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

func (m *Mixer) UpdateVolume() {
	for _, sound := range m.Sounds {
		rl.SetSoundVolume(sound.Sound, Volumes[sound.Category]*m.Volume)
	}

	for _, music := range m.Musics {
		rl.SetMusicVolume(music.Music, Volumes[music.Category]*m.Volume)
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
	m.Volume = volume
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
