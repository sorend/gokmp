package main


import (
	"github.com/sorend/gokmp/cmd"
)


func main() {
	cmd.Execute()
}

/*
func doBackup(client *flickr.Client, nsid string, destination string) error {
	fmt.Println("Backing up ...")
	existing, err := findExistingPhotoIds(destination)
	if err != nil {
		return err
	}

	fmt.Printf("%d photos found in %s\n", len(existing), destination)

	photosets, err := photosetsGetList(client, nsid)
	if err != nil {
		return err
	}

	for _, photoset := range photosets {
		// backupDirname := path.Join(destination, slug.Clean(photoset.Title.Text) + "-" + photoset.Id)
		photos, err := photosetsGetPhotosWalk(client, photoset)
		if err != nil {
			return err
		}
		fmt.Printf("photos %s\n", photos)
		// fmt.Printf("photoset id:%s title:'%s' dirname:'%s' updated:%d\n", photoset.Id, photoset.Title, backupDirname, photoset.DateUpdate)
	}
	return err
}

func photosetsGetList(client *flickr.Client, nsid string) ([]*GetListPhotoset, error) {
	response, err := client.Request("photosets.getList", flickr.Params{"user_id": nsid})
	if err != nil {
		return nil, err
	}
	photosets := &GetListPhotosetsRaw{}
	err = flickr.Parse(response, &photosets)
	return photosets.Photosets.Photosets, err
}


func photosetsGetPhotosWalk(client *flickr.Client, photoset *GetListPhotoset) ([]RemotePhotoDetails, error) {
	response, err := photosetsGetPhotos(client, photoset.Id, 1)
	if err != nil {
		return nil, err
	}

	res := []RemotePhotoDetails{}
	for _, photo := range response.Photoset.Photos {
		farmNumber, _ := strconv.Atoi(photo.Farm)
		res = append(res, RemotePhotoDetails{
			Id: photo.Id,
			SetId: photoset.Id,
			Title: photo.Title,
			Url: flickr.GenerateURL("o", photo.Id, farmNumber, photo.Secret, photo.Server, "jpg"),
		})
	}

	pages := response.Photoset.Pages
	fmt.Printf("pages %d\n", pages)

	// res = append(res, response.Photos)

	for page := 2; page <= pages; page++ {
		response, err = photosetsGetPhotos(client, photoset.Id, page)
		if err != nil {
			return res, err // partial with error
		}
		for _, photo := range response.Photoset.Photos {
			farmNumber, _ := strconv.Atoi(photo.Farm)
			res = append(res, RemotePhotoDetails{
				Id: photo.Id,
				SetId: photoset.Id,
				Title: photo.Title,
				Url: flickr.GenerateURL("o", photo.Id, farmNumber, photo.Secret, photo.Server, "jpg"),
			})
		}
	}
	return res, nil
}

func photosetsGetPhotos(client *flickr.Client, photosetId string, page int) (*GetPhotosPhotoset, error) {
	res := &GetPhotosPhotoset{}
	response, err := client.Request("photosets.getPhotos", flickr.Params{"photoset_id": photosetId, "page": fmt.Sprintf("%d", page)})
	if err != nil {
		return res, err
	}
	err = flickr.Parse(response, &res)
	if  err != nil {
		return res, err
	}
	return res, nil
}


func photoBestSize(client *flickr.Client, photoId string) (string, error) {
	_, err := client.Request("photos.getInfo", flickr.Params{"photo_id": photoId})
	if err != nil {
		return "", err
	}
	return "o", err
}


func hasExisting(ids []LocalPhotoDetails, photoSetId string, photoId string) bool {
	for _, d := range ids {
		if d.photoSetId == photoSetId && d.photoId == photoId {
			return true
		}
	}
	return false
}

func findExistingPhotoIds(destination string) (ids []LocalPhotoDetails, err error) {

	photoset_folders, err := ioutil.ReadDir(destination)
	if err != nil {
		return nil, err
	}

	res := []LocalPhotoDetails{}
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
				res = append(res, LocalPhotoDetails { photoSetId: photoset_id, photoId: photo_id })
			}
		}
	}

	fmt.Printf("photos %q\n", res)

	return res, nil
}


type GetPhotosPhotoset struct {
	Photoset struct {
		Id string `json:"id"`
		Photos []*GetPhotosPhoto `json:"photos"`
		Page string `json:"page"`
		Pages int `json:"pages"`
		Title string `json:"title"`
		Total int `json:"total"`
	} `json:"photoset"`
}

type GetPhotosPhoto struct {
	Id string `json:"id"`
	Secret string `json:"secret"`
	Server string `json:"server"`
	Farm string `json:"farm"`
	Title string `json:"title"`
	Media string `json:"media"`
}

type GetListPhotosetsRaw struct {
	Photosets struct {
		CanCreate int `json:"cancreate"`
		Page int `json:"page"`
		Pages int `json:"pages"`
		Total int `json:"total"`
		Photosets []*GetListPhotoset `json:"photoset"`
	} `json:"photosets"`
}

type GetListPhotoset struct {
	Id string `json:"id"`
	Owner string `json:"owner"`
	Secret string `json:"secret"`
	Server string `json:"server"`
	Farm int `json:"farm"`
	CountPhotos int `json:"count_photos"`
	CountVideos int `json:"count_videos"`
	Title Content `json:"title"`
	Description Content `json:"description"`
	CanComment int `json:"can_comment"`
	DateCreate string `json:"date_create"`
	DateUpdate string `json:"date_update"`
}

type Content struct {
	Text string `json:"_content"`
}

type RemotePhotoDetails struct {
	Id string
	SetId string
	Title string
	Url string
}


type LocalPhotoDetails struct {
	photoSetId string
	photoId string
}
*/
