package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/ty-porter/mudora/internal/game"
	"github.com/ty-porter/mudora/internal/sni"
)

type Event struct {
	Description string
	State       game.State
	At          time.Time
}

type Tracker struct {
	client   *sni.Client
	interval time.Duration

	eventCallbacks  []func(Event)
	statusCallbacks []func(string)

	deviceURI string
	tracking  bool

	prev    game.State
	hasPrev bool
}

func New(client *sni.Client, interval time.Duration) *Tracker {
	return &Tracker{
		client:   client,
		interval: interval,
	}
}

// OnEvent registers a callback invoked with each tracker Event (initial and
// changed game state). Register before Run; callbacks run on the tracker's
// poll goroutine, so handlers that touch the UI must marshal onto the UI
// thread themselves.
func (t *Tracker) OnEvent(fn func(Event)) {
	t.eventCallbacks = append(t.eventCallbacks, fn)
}

// OnStatus registers a callback invoked with human-readable status updates
// (device discovery, tracking on/off, connection loss). Same threading caveat
// as OnEvent.
func (t *Tracker) OnStatus(fn func(string)) {
	t.statusCallbacks = append(t.statusCallbacks, fn)
}

func (t *Tracker) Run(ctx context.Context) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	t.status("searching for SNES device")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			t.tick(ctx)
		}
	}
}

func (t *Tracker) tick(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if t.deviceURI == "" {
		uri, err := t.client.FirstDevice(ctx)
		if err != nil {
			return
		}

		t.deviceURI = uri
		t.tracking = false
		t.hasPrev = false
		t.status(fmt.Sprintf("device: %s - waiting for game", uri))
	}

	if err := t.poll(ctx); err != nil {
		if ctx.Err() != nil {
			return
		}

		t.status(fmt.Sprintf("device lost (%v) - searching", err))
		t.deviceURI = ""
	}
}

func (t *Tracker) poll(ctx context.Context) error {
	mode, err := t.client.ReadMemory(ctx, t.deviceURI, game.GameModeAddr, 1)
	if err != nil {
		return err
	}

	if !game.InGame(mode[0]) {
		if t.tracking {
			t.tracking = false
			t.status(fmt.Sprintf("device: %s — waiting for game…", t.deviceURI))
		}

		return nil
	}

	if !t.tracking {
		t.tracking = true
		t.status(fmt.Sprintf("device: %s — tracking", t.deviceURI))
	}

	data, err := t.client.ReadMemory(ctx, t.deviceURI, game.SaveDataAddr, game.SaveDataSize)
	if err != nil {
		return err
	}

	state, err := game.ParseSaveData(data)
	if err != nil {
		return err
	}

	if t.hasPrev {
		t.diff(t.prev, state)
	} else {
		t.emit(Event{Description: fmt.Sprintf("Initial state: %+v", state), State: state, At: time.Now()})
	}

	t.prev = state
	t.hasPrev = true

	return nil
}

func (t *Tracker) diff(old, new game.State) {
	if old == new {
		return
	}
	t.emit(Event{Description: fmt.Sprintf("state changed: %+v", new), State: new, At: time.Now()})
}

func (t *Tracker) emit(event Event) {
	for _, fn := range t.eventCallbacks {
		fn(event)
	}
}

func (t *Tracker) status(msg string) {
	for _, fn := range t.statusCallbacks {
		fn(msg)
	}
}
