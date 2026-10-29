package tui

import (
	"testing"
	"time"

	"github.com/museslabs/kyma/internal/config"
	"github.com/museslabs/kyma/internal/tui/transitions"
)

func TestSyncServer_GoTo(t *testing.T) {
	server, err := listenSync("127.0.0.1:0")
	if err != nil {
		t.Fatalf("listenSync() failed: %v", err)
	}
	server.Start()
	t.Cleanup(server.Stop)

	address := server.listener.Addr().String()

	// Speaker notes stay connected, and keep hearing about slides.
	notes, err := dialSync(address)
	if err != nil {
		t.Fatalf("dialSync() failed: %v", err)
	}
	t.Cleanup(notes.Close)

	for _, want := range []int{5, 2} {
		remote, err := dialSync(address)
		if err != nil {
			t.Fatalf("dialSync() failed: %v", err)
		}
		if err := remote.GoTo(want); err != nil {
			t.Fatalf("GoTo(%d) failed: %v", want, err)
		}
		remote.Close()

		select {
		case got := <-server.GoTos():
			if got != want {
				t.Errorf("server got slide %d, want %d", got, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("server never got slide %d", want)
		}
	}

	slides := make(chan int, 2)
	go notes.ListenForSlideChanges(slides)
	server.BroadcastSlideChange(3)
	for _, want := range []int{0, 3} {
		select {
		case got := <-slides:
			if got != want {
				t.Errorf("speaker notes got slide %d, want %d", got, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("speaker notes never got slide %d", want)
		}
	}
}

func TestModel_GoToSlideMsg(t *testing.T) {
	props := config.Properties{Transition: transitions.Get("swipeLeft", transitions.Fps)}

	var slides []*Slide
	for _, body := range []string{"# One", "# Two", "# Three"} {
		slide, err := NewSlide(body, props)
		if err != nil {
			t.Fatalf("NewSlide() failed: %v", err)
		}
		if n := len(slides); n > 0 {
			slides[n-1].Next, slide.Prev = slide, slides[n-1]
		}
		slides = append(slides, slide)
	}

	m := model{
		slide:        slides[1],
		rootSlide:    slides[0],
		keys:         keys,
		width:        80,
		height:       24,
		timerDisplay: NewTimerDisplay(),
		globalTimer:  NewTimer(),
	}

	// Halfway through a transition, with the go-to overlay open.
	m.slide.ActiveTransition = props.Transition.Start(m.width, m.height, transitions.Forwards)
	goTo := NewGoTo(len(slides)).SetShowing(true)
	m.goTo = &goTo

	updated, _ := m.Update(GoToSlideMsg{Number: 3})
	m = updated.(model)
	if m.slide != slides[2] {
		t.Fatalf("on %q, want %q", m.slide.Data, slides[2].Data)
	}
	if m.goTo != nil {
		t.Error("the go-to overlay stayed open")
	}
	if slides[1].ActiveTransition != nil {
		t.Error("the transition left behind is still running")
	}

	// The frame the transition had asked for still arrives.
	updated, _ = m.Update(transitions.FrameMsg(time.Now()))
	m = updated.(model)

	for _, number := range []int{0, 4, -1} {
		updated, _ = m.Update(GoToSlideMsg{Number: number})
		m = updated.(model)
		if m.slide != slides[2] {
			t.Errorf("slide %d moved to %q", number, m.slide.Data)
		}
	}

	updated, _ = m.Update(GoToSlideMsg{Number: 1})
	m = updated.(model)
	if m.slide != slides[0] {
		t.Errorf("on %q, want %q", m.slide.Data, slides[0].Data)
	}
}
