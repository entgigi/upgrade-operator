package utils

import (
	"strings"

	ref "github.com/containers/image/v5/docker/reference"
)

type ImageInfo struct {
	ref.Reference
}

func NewImageInfo(imageUrl string) (ImageInfo, error) {
	ref, err := ref.ParseNormalizedNamed(imageUrl)
	if err != nil {
		return ImageInfo{}, err
	}
	im := ImageInfo{ref}

	return im, nil

}

func (i ImageInfo) Name() string {
	named, _ := i.Reference.(ref.Named)
	_, name := ref.SplitHostname(named)
	lastInd := strings.LastIndex(name, "/")
	return name[lastInd+1:]
}

func (i ImageInfo) Org() string {
	named, _ := i.Reference.(ref.Named)
	_, name := ref.SplitHostname(named)
	lastInd := strings.LastIndex(name, "/")
	return name[:lastInd]
}

func (i ImageInfo) Hostname() string {
	named, _ := i.Reference.(ref.Named)
	hostname, _ := ref.SplitHostname(named)
	return hostname
}

func (i ImageInfo) IsTag() bool {
	_, isTag := i.Reference.(ref.Tagged)
	return isTag
}

func (i ImageInfo) Tag() string {
	tag, _ := i.Reference.(ref.Tagged)
	return tag.String()
}

func (i ImageInfo) IsDigest() bool {
	_, isDigest := i.Reference.(ref.Digested)
	return isDigest
}

func (i ImageInfo) Digest() string {
	digest, _ := i.Reference.(ref.Digested)
	return digest.Digest().String()
}
