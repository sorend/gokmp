package storage

import (
	"io/ioutil"
	"strings"
	"path"
	"fmt"
	"github.com/gosimple/slug"
)

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

func (s *Storage) Filename(photosetId string, photosetTitle string, photoId string, photoTitle string) string {
	dirname := fmt.Sprintf("%s-%s", slug.Make(photosetTitle), photosetId)
	filename := fmt.Sprintf("%s-%s.jpg", slug.Make(photoTitle), photoId)
	return path.Join(s.Destination, dirname, filename)
}

func findLocalPhotos(destination string) (map[string]*LocalPhoto, error) {
	photoset_folders, err := ioutil.ReadDir(destination)
	if err != nil {
		return nil, err
	}

	res := map[string]*LocalPhoto{}
	for _, photoset_folder := range photoset_folders {
		splits := strings.Split(photoset_folder.Name(), "-")
		if photoset_folder.IsDir() && len(splits) > 1 {
			photoset_id := splits[len(splits) - 1]
			photo_files, err := ioutil.ReadDir(path.Join(destination, photoset_folder.Name()))
			if err != nil {
				return res, err // partial
			}
			for _, photo_file := range photo_files {
				photo_splits := strings.Split(photo_file.Name(), "-")
				photo_id := strings.Split(photo_splits[len(photo_splits) - 1], ".")[0]
				res[photoset_id + "/" + photo_id] = &LocalPhoto{
					PhotosetId: photoset_id,
					PhotoId: photo_id,
					PhotoFilename: photo_file.Name(),
					PhotosetDirname: path.Join(destination, photoset_folder.Name()),
				}
			}
		}
	}
	return res, nil
}


