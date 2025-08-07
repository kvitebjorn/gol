package gui

var playing bool
var paused bool
var playStopCh chan struct{}

func stopPlayback() {
	if playing {
		if playStopCh != nil {
			close(playStopCh)
			playStopCh = nil
		}
		playing = false
		paused = false
	}
}
