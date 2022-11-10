package backup

import (
	"fmt"
	"strings"

	"github.com/sorend/gokmp/pkg/flickr"
	"github.com/sorend/gokmp/pkg/storage"
	"go.uber.org/zap"
)

var logger = zap.NewExample().Sugar()

func Run(ApiKey string, ApiSecret string, accessToken string, accessSecret string, nsid string, destination string) error {
	defer logger.Sync()
	client := flickr.NewFlickr(
		ApiKey,
		ApiSecret,
		accessToken,
		accessSecret,
	)
	if err := doBackup(client, nsid, destination); err != nil {
		panic(err)
	}
	logger.Infow("All done :-)")
	fmt.Println("All done :-)")
	return nil
}

func doBackup(client *flickr.Flickr, nsid string, destination string) error {
	fmt.Println("| Backing up ...")
	existing, err := storage.New(destination)
	if err != nil {
		return err
	}
	fmt.Printf("| %d photos found in %s\n", len(existing.Existing), destination)
	photosets, err := client.PhotosetsGetList(nsid)
	if err != nil {
		return err
	}
	fmt.Printf("| Looking for photos to backup without a set...\n")
	photos, err := client.PhotosGetNotInSetWalk()
	if err != nil {
		return err
	}
	notInSet := &flickr.PhotosetsGetListPhotoset{
		Id: storage.NotInSetId,
		Title: flickr.Content{
			Text: storage.NotInSet,
		},
	}
	if err = doBackupPhotos(client, existing, notInSet, photos); err != nil {
		return err
	}
	fmt.Printf("| Looking for photos to backup in %d photosets...\n", len(photosets.Photosets.Photosets))
	for _, photoset := range photosets.Photosets.Photosets {
		photos, err := client.PhotosetsGetPhotosWalk(photoset.Id)
		if err != nil {
			return err
		}
		if err = doBackupPhotos(client, existing, photoset, photos); err != nil {
			return err
		}
	}
	return err
}

func doBackupPhotos(client *flickr.Flickr, existing *storage.Storage, photoset *flickr.PhotosetsGetListPhotoset, photos []*flickr.PhotosetsGetPhotosPhoto) error {
	needBackup := []*QueueItem{}
	for _, photo := range photos {
		if photo.Media == "photo" && !existing.Has(photoset.Id, photo.Id) {
			needBackup = append(needBackup, &QueueItem{
				photo:    photo,
				photoset: photoset,
			})
		}
	}
	fmt.Printf("`-> %s .. to-backup %d\n", photoset.Title.Text, len(needBackup))
	for _, item := range needBackup {
		if err := doBackupPhoto(client, existing, item); err != nil {
			return err
		}
	}
	return nil
}

func doBackupPhoto(client *flickr.Flickr, existing *storage.Storage, item *QueueItem) error {
	info, err := client.PhotosGetInfo(item.photo.Id)
	if err != nil {
		return err
	}
	sizes, err := client.PhotosGetSizes(item.photo.Id)
	if err != nil {
		return err
	}
	best := takeBest(sizes.Sizes.Size)
	takenDate := takenDate(info)
	destinationFile := existing.Filename(item.photoset.Id, item.photoset.Title.Text, item.photo.Id, item.photo.Title, takenDate)
	fmt.Printf("  `-> backing up %s ..\n", best.Source)
	fmt.Printf("    >      .. to %s ..\n", destinationFile)
	if err = client.Download(best.Source, destinationFile); err != nil {
		return err
	}
	return nil
}

func takenDate(info *flickr.PhotosGetInfo) string {
	taken := info.Photo.Dates.Taken
	return strings.Split(taken, " ")[0]
}

func takeBest(sizes []*flickr.PhotosGetSizesSize) *flickr.PhotosGetSizesSize {
	best := sizes[0]
	for _, s := range sizes {
		if s.Label == "Original" {
			return s
		}
		if s.Width > best.Width && s.Height > best.Height {
			best = s
		}
	}
	return best
}

type QueueItem struct {
	photo    *flickr.PhotosetsGetPhotosPhoto
	photoset *flickr.PhotosetsGetListPhotoset
}
