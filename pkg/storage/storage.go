package storage

import (
	"io/ioutil"
	"strings"
	"path"
	"fmt"
	"github.com/gosimple/slug"
)

const NotInSet = "_not_in_set"
const NotInSetId = ""

type Storage struct {
	Existing map[string]*LocalPhoto
	Destination string
}

type LocalPhoto struct {
	PhotosetId string
	PhotoId string
	PhotoFilename string
	PhotosetDirname string
}

func New(destination string) (*Storage, error) {
	existing, err := findLocalPhotos(destination)
	if err != nil {
		return nil, err
	}
	return &Storage{
		Existing: existing,
		Destination: destination,
	}, nil
	
}

func (s *Storage) Has(photosetId string, photoId string) bool {
	key := photosetId + "/" + photoId
	_, ok := s.Existing[key]
	return ok
}

func (s *Storage) Filename(photosetId string, photosetTitle string, photoId string, photoTitle string, photoTakenDate string) string {
	var dirname = NotInSet
	if photosetId != NotInSetId {
		dirname = fmt.Sprintf("%s-%s", slug.Make(photosetTitle), photosetId)
	}
	filename := fmt.Sprintf("%s-%s-%s.jpg", photoTakenDate, slug.Make(photoTitle), photoId)
	return path.Join(s.Destination, dirname, filename)
}

func findLocalPhotos(destination string) (map[string]*LocalPhoto, error) {
	photosetDirs, err := ioutil.ReadDir(destination)
	if err != nil {
		return nil, err
	}

	res := map[string]*LocalPhoto{}
	for _, photosetDir := range photosetDirs {
		if !photosetDir.IsDir() {
			continue
		}
		splits := strings.Split(photosetDir.Name(), "-")
		var photosetId = "unknown"
		if photosetDir.Name() == NotInSet {
			photosetId = NotInSetId
		} else if len(splits) > 1 {
			photosetId = splits[len(splits) - 1]
		}
		
		photoFiles, err := ioutil.ReadDir(path.Join(destination, photosetDir.Name()))
		if err != nil {
			return res, err // partial
		}
		for _, photoFile := range photoFiles {
			if !strings.HasSuffix(photoFile.Name(), ".jpg") {
				continue
			}
			splitsPhoto := strings.Split(photoFile.Name(), "-")
			photoId := strings.Split(splitsPhoto[len(splitsPhoto) - 1], ".")[0]
			res[photosetId + "/" + photoId] = &LocalPhoto{
				PhotosetId: photosetId,
				PhotoId: photoId,
				PhotoFilename: photoFile.Name(),
				PhotosetDirname: path.Join(destination, photosetDir.Name()),
			}
		}
	}
	return res, nil
}


