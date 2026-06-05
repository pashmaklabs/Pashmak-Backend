package services_auth

import (
	"os"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"pashmak.com/pashmak/bootstrap"
)

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

	appConfig := bootstrap.LoadEnvVars()
	redisClient := bootstrap.SetupRedis()
	defer redisClient.Close()

	as := NewAuthService(nil, redisClient, appConfig)

	err = as.StoreOTPAndSendEmail("test@example.com")
	if err != nil {
		as.CaptureAuthError(err, "send_email", "test@example.com", map[string]interface{}{
			"test_run": true,
		})
		t.Logf("error captured and sent to Sentry: %v", err)
	} else {
		t.Log("email sent successfully — check inbox")
	}

	sentry.Flush(5 * time.Second)
	t.Log("Event sent — check your Sentry dashboard")
}
