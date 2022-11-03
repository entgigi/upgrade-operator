package service

import (
	"testing"
)

func TestParseRepositoryInfo(t *testing.T) {
	type tcase struct {
		TestString, Domain, Organization, Repository, Version string
	}

	tcases := []tcase{
		{
			TestString:   "bar@sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
			Domain:       "docker.io",
			Organization: "library",
			Repository:   "bar",
			Version:      "sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
		},
		{
			TestString:   "fooo/bar@sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
			Domain:       "docker.io",
			Organization: "fooo",
			Repository:   "bar",
			Version:      "sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
		},
		{
			TestString:   "quay.io/fooo/bar@sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
			Domain:       "quay.io",
			Organization: "fooo",
			Repository:   "bar",
			Version:      "sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
		},
	}

	for _, item := range tcases {
		image, err := NewImageInfo(item.TestString)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if expected, actual := item.Domain, image.Hostname(); expected != actual {
			t.Fatalf("Invalid Domain for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
		if expected, actual := item.Organization, image.Org(); expected != actual {
			t.Fatalf("Invalid Organization for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
		if expected, actual := item.Repository, image.Name(); expected != actual {
			t.Fatalf("Invalid Repository for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
		if expected, actual := item.Version, image.Digest(); !image.IsDigest() || expected != actual {
			t.Fatalf("Invalid Digest for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
	}

	tcases2 := []tcase{
		{
			TestString:   "bar:7.0.1",
			Domain:       "docker.io",
			Organization: "library",
			Repository:   "bar",
			Version:      "7.0.1",
		},
	}

	for _, item := range tcases2 {
		image, err := NewImageInfo(item.TestString)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if expected, actual := item.Domain, image.Hostname(); expected != actual {
			t.Fatalf("Invalid Domain for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
		if expected, actual := item.Organization, image.Org(); expected != actual {
			t.Fatalf("Invalid Organization for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
		if expected, actual := item.Repository, image.Name(); expected != actual {
			t.Fatalf("Invalid Repository for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
		if expected, actual := item.Version, image.Tag(); !image.IsTag() || expected != actual {
			t.Fatalf("Invalid Digest for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
	}
}
