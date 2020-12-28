
package backup

import (
	"fmt"
	"os"
	"github.com/sorend/gokmp/pkg/storage"
	"github.com/sorend/gokmp/pkg/flickr"
)


func Run(nsid string, destination string) error {


	client := flickr.NewFlickr(
		os.Getenv("FLICKR_API_KEY"),
		os.Getenv("FLICKR_API_SECRET"),
		os.Getenv("FLICKR_OAUTH_TOKEN"),
		os.Getenv("FLICKR_OAUTH_SECRET"),
	)
	fmt.Printf("Client %s\n", client)

	if err := doBackup(client, nsid, destination); err != nil {
		panic(err)
	}
	fmt.Println("All done :-)")
	return nil
}


func doBackup(client *flickr.Flickr, nsid string, destination string) error {
	fmt.Println("Backing up ...")
	existing, err := storage.New(destination)
	if err != nil {
		return err
	}

	fmt.Printf("%d photos found in %s\n", len(existing.Existing), destination)

	photosets, err := client.PhotosetsGetList(nsid)
	if err != nil {
		return err
	}

	fmt.Printf("Looking for photos to backup in %d photosets...\n", len(photosets))

	needBackup := []*QueueItem{}
	for _, photoset := range photosets {
		photos, err := client.PhotosetsGetPhotosWalk(photoset.Id)
		if err != nil {
			return err
		}
		for _, photo := range photos {
			if photo.Media == "photo" && !existing.Has(photoset.Id, photo.Id) {
				needBackup = append(needBackup, &QueueItem{
					photo: photo,
					photoset: photoset,
				})
			}
		}
		fmt.Printf(" .. %s .. queue %d\n", photoset.Title.Text, len(needBackup))
	}
	for _, item := range needBackup {
		sizes, err := client.PhotosGetSizes(item.photo.Id)
		if err != nil {
			return err
		}
		best := takeBestWidth(sizes)
		fmt.Printf("> Backing up %s ..\n", best.Source)
		destinationFile := existing.Filename(item.photoset.Id, item.photoset.Title.Text, item.photo.Id, item.photo.Title)
		if err = client.Download(best.Source, destinationFile); err != nil {
			return err
		}
	}
	return err
}


func takeBestWidth(sizes []*flickr.PhotosGetSizesSize) *flickr.PhotosGetSizesSize {
	best := sizes[0]
	for _, s := range sizes {
		if s.Width > best.Width {
			best = s
		}
	}
	return best
}


type QueueItem struct {
	photo *flickr.PhotosetsGetPhotosPhoto
	photoset *flickr.GetListPhotoset
}