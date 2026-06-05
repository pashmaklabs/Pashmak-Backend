package services_auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
)

func TestCaptureAuthError_SentryEventShape(t *testing.T) {
	var captured []*sentry.Event

	_ = sentry.Init(sentry.ClientOptions{
		Dsn: "https://examplePublicKey@o0.ingest.sentry.io/0", // fake DSN
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			captured = append(captured, event)
			return nil
		},
	})
	defer sentry.Flush(time.Second)

	as := &AuthService{}
	as.CaptureAuthError(fmt.Errorf("redis timeout"), "validate_otp", "u@example.com", nil)
	sentry.Flush(time.Second)

	if len(captured) != 1 {
		t.Fatalf("expected 1 event, got %d", len(captured))
	}
	ev := captured[0]
	if ev.Tags["operation"] != "validate_otp" {
		t.Errorf("unexpected operation tag: %s", ev.Tags["operation"])
	}
	if ev.Tags["service"] != "auth" {
		t.Errorf("expected service=auth tag")
	}
}

func TestCaptureAuthError_RealSentry(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real Sentry test in short mode")
	}

	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		t.Skip("SENTRY_DSN not set, skipping real Sentry test")
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      "test",
		Debug:            true,
	})
	if err != nil {
		t.Fatalf("sentry.Init failed: %v", err)
	}
	defer sentry.Flush(5 * time.Second)

	as := &AuthService{}
	as.CaptureAuthError(
		fmt.Errorf("real test: redis timeout"),
		"validate_otp",
		"test@example.com",
		map[string]interface{}{
			"test_run": true,
			"password": "should_not_appear",
		},
	)

	sentry.Flush(5 * time.Second)
	t.Log("Event sent — check your Sentry dashboard")
}
