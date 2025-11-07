package integrationtest

import (
	"context"
	"os"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	accessToken := os.Getenv("TEST_TOKEN")

	info, err := GetUserInfo(context.Background(), accessToken)
	if err != nil {
		t.Errorf("failed to get user info, error = %v", err)
	}

	t.Logf("ID: %s\nName: %s\nPictureURL: %s\n", info.ID, info.Name, info.PictureURL)
}
