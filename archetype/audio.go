package archetype

import (
	"amaru/assets"
	"amaru/component"
)

func PlayAudioMenu() {
	if !assets.MenuAdioPlayer.IsPlaying() {
		assets.MenuAdioPlayer.Rewind()
		assets.MenuAdioPlayer.Play()
	}
}

func StopAudioMenu() {
	if assets.MenuAdioPlayer.IsPlaying() {
		assets.MenuAdioPlayer.Pause()
		assets.MenuAdioPlayer.Rewind()
	}
}

func PlayAudioGame() {
	if !assets.GameAudioPlayer.IsPlaying() {
		assets.GameAudioPlayer.Rewind()
		assets.GameAudioPlayer.Play()
	}
}

func StopAudioGame() {
	if assets.GameAudioPlayer.IsPlaying() {
		assets.GameAudioPlayer.Pause()
		assets.GameAudioPlayer.Rewind()
	}
}

func PlayWinnerAudio(gameData *component.GameData) {
	PlayWavesAudio()
	if gameData.Session.RemoteClient.GameData.Counter >= 5 {
		PlayHeronAudio()
	}
	if gameData.Session.RemoteClient.GameData.Counter >= 10 {
		PlayShipAudio()
	}
}

func StopWinnerAudio() {
	if assets.WavesAudioPlayer.IsPlaying() {
		assets.WavesAudioPlayer.Pause()
		assets.WavesAudioPlayer.Rewind()
	}
	if assets.HeronAudioPlayer.IsPlaying() {
		assets.HeronAudioPlayer.Pause()
		assets.HeronAudioPlayer.Rewind()
	}
	if assets.ShipAudioPlayer.IsPlaying() {
		assets.ShipAudioPlayer.Pause()
		assets.ShipAudioPlayer.Rewind()
	}
}

func PlayWavesAudio() {
	if !assets.WavesAudioPlayer.IsPlaying() {
		assets.WavesAudioPlayer.Rewind()
		assets.WavesAudioPlayer.Play()
	}
}

func PlayHeronAudio() {
	if !assets.HeronAudioPlayer.IsPlaying() {
		assets.HeronAudioPlayer.Rewind()
		assets.HeronAudioPlayer.Play()
	}
}

func PlayShipAudio() {
	if !assets.ShipAudioPlayer.IsPlaying() {
		assets.ShipAudioPlayer.Rewind()
		assets.ShipAudioPlayer.Play()
	}
}

func PlayCollectedAudio() {
	if !assets.CollectedAudioPlayer.IsPlaying() {
		assets.CollectedAudioPlayer.Rewind()
		assets.CollectedAudioPlayer.Play()
	}
}

func PlayButtonClickAudio() {
	assets.ButtonClickPlayer.Rewind()
	assets.ButtonClickPlayer.Play()
}
