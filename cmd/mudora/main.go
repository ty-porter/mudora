package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ty-porter/mudora/internal/sni"
	"github.com/ty-porter/mudora/internal/tracker"
	"github.com/ty-porter/mudora/internal/ui"
)

func main() {
	var (
		sniAddr      = flag.String("sni", "localhost:8191", "address of the SNI gRPC server")
		pollInterval = flag.Duration("poll", 500*time.Millisecond, "how often to poll SNES memory (ms)")
		// settingsPath = flag.String("settings", "mudora.json", "path to persisted settings")
	)

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	client, err := sni.Dial(ctx, *sniAddr)
	if err != nil {
		log.Fatalf("connecting to SNI at %s: %v", *sniAddr, err)
	}
	defer stop()

	t := tracker.New(client, *pollInterval)
	w := ui.New()

	// Bridge the tracker's poll goroutine to the UI. These callbacks only
	// stash values; the window applies them on the Tk event loop.
	t.OnStatus(w.SetStatus)
	t.OnEvent(func(e tracker.Event) { w.ApplyState(e.State) })

	go func() {
		if err := t.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("tracker stopped: %v", err)
		}
	}()

	w.Run()
}
