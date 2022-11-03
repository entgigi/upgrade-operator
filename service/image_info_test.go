package service

import (
	"testing"
)

func TestParseRepositoryInfo(t *testing.T) {
	type tcase struct {
		TestString, Domain, Organization, Repository, Digest string
	}

	tcases := []tcase{
		{
			TestString:   "bar@sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
			Domain:       "docker.io",
			Organization: "library",
			Repository:   "bar",
			Digest:       "sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
		},
		{
			TestString:   "fooo/bar@sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
			Domain:       "docker.io",
			Organization: "fooo",
			Repository:   "bar",
			Digest:       "sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
		},
		{
			TestString:   "quay.io/fooo/bar@sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
			Domain:       "quay.io",
			Organization: "fooo",
			Repository:   "bar",
			Digest:       "sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c",
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
		if expected, actual := item.Digest, image.Digest(); !image.IsDigest() || expected != actual {
			t.Fatalf("Invalid Digest for %q. Expected %q, got %q", item.TestString, expected, actual)
		}
	}
}
