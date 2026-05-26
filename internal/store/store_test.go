package store_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/deploylog/deploylog/internal/event"
	"github.com/deploylog/deploylog/internal/store"
)

func makeEvent(id, source string) event.Event {
	return event.Event{
		ID:        id,
		Source:    source,
		Service:   "api",
		Status:    event.StatusSuccess,
		Timestamp: time.Now().UTC().Truncate(time.Second),
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.json")
	s := store.New(path)

	events := []event.Event{makeEvent("1", "github"), makeEvent("2", "gitlab")}
	if err := s.Save(events); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 events, got %d", len(loaded))
	}
	if loaded[0].ID != "1" || loaded[1].ID != "2" {
		t.Errorf("unexpected event IDs: %v", loaded)
	}
}

func TestStore_Load_FileNotExist(t *testing.T) {
	dir := t.TempDir()
	s := store.New(filepath.Join(dir, "missing.json"))

	events, err := s.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty slice, got %d events", len(events))
	}
}

func TestStore_Append(t *testing.T) {
	dir := t.TempDir()
	s := store.New(filepath.Join(dir, "events.json"))

	if err := s.Save([]event.Event{makeEvent("1", "github")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Append([]event.Event{makeEvent("2", "jenkins")}); err != nil {
		t.Fatal(err)
	}

	loaded, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 events after append, got %d", len(loaded))
	}
}

func TestStore_Save_InvalidPath(t *testing.T) {
	s := store.New("/nonexistent/dir/events.json")
	err := s.Save([]event.Event{makeEvent("1", "github")})
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestStore_Load_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	if err := os.WriteFile(path, []byte("not-valid-json{"), 0o644); err != nil {
		t.Fatal(err)
	}

	s := store.New(path)
	_, err := s.Load()
	if err == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
}
