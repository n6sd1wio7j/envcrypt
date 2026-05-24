package env

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// DefaultTagsPath is the default file path for storing environment tags.
const DefaultTagsPath = ".envcrypt/tags.json"

// Tag represents a named marker for a specific snapshot or version of the env file.
type Tag struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SnapshotRef string    `json:"snapshot_ref"` // path or identifier of the snapshot
	CreatedAt   time.Time `json:"created_at"`
}

// TagFile holds all tags.
type TagFile struct {
	Tags []Tag `json:"tags"`
}

// LoadTagFile loads the tag file from the given path.
// Returns an empty TagFile if the file does not exist.
func LoadTagFile(path string) (TagFile, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return TagFile{}, nil
	}
	if err != nil {
		return TagFile{}, err
	}
	var tf TagFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return TagFile{}, err
	}
	return tf, nil
}

// SaveTagFile writes the TagFile to the given path, creating directories as needed.
func SaveTagFile(path string, tf TagFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// AddTag appends a new tag to the TagFile. Returns an error if a tag with the same name exists.
func AddTag(tf *TagFile, tag Tag) error {
	for _, t := range tf.Tags {
		if t.Name == tag.Name {
			return errors.New("tag already exists: " + tag.Name)
		}
	}
	if tag.CreatedAt.IsZero() {
		tag.CreatedAt = time.Now().UTC()
	}
	tf.Tags = append(tf.Tags, tag)
	return nil
}

// RemoveTag removes a tag by name. Returns an error if not found.
func RemoveTag(tf *TagFile, name string) error {
	for i, t := range tf.Tags {
		if t.Name == name {
			tf.Tags = append(tf.Tags[:i], tf.Tags[i+1:]...)
			return nil
		}
	}
	return errors.New("tag not found: " + name)
}

// FindTag returns a tag by name, or an error if not found.
func FindTag(tf TagFile, name string) (Tag, error) {
	for _, t := range tf.Tags {
		if t.Name == name {
			return t, nil
		}
	}
	return Tag{}, errors.New("tag not found: " + name)
}
