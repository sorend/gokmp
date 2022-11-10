package flickr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/dghubble/oauth1"
)

func NewFlickr(ApiKey string, ApiSecret string, OAuthToken string, OAuthSecret string) *Flickr {
	config := oauth1.NewConfig(ApiKey, ApiSecret)
	token := oauth1.NewToken(OAuthToken, OAuthSecret)
	oauthClient := config.Client(oauth1.NoContext, token)

	return &Flickr{
		ApiKey: ApiKey,
		client: oauthClient,
	}
}

type Flickr struct {
	ApiKey string
	client *http.Client
}

type Params map[string]string // Params datatype

func (f *Flickr) PhotosetsGetList(nsid string) (*PhotosetsGetList, error) {
	res := &PhotosetsGetList{}
	if err := f.Request("photosets.getList", Params{"user_id": nsid}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Wraps PhotosetsGetPhotosPage and returns only the photos from each page
func (f *Flickr) PhotosetsGetPhotosWalk(photosetId string) ([]*PhotosetsGetPhotosPhoto, error) {
	thePage, err := f.PhotosetsGetPhotosPage(photosetId, 1)
	if err != nil {
		return nil, err
	}

	res := []*PhotosetsGetPhotosPhoto{}
	for _, photo := range thePage.Photoset.Photos {
		res = append(res, photo)
	}

	pages := thePage.Photoset.Pages

	for page := 2; page <= pages; page++ {
		thePage, err = f.PhotosetsGetPhotosPage(photosetId, page)
		if err != nil {
			return res, err // partial with error
		}
		for _, photo := range thePage.Photoset.Photos {
			res = append(res, photo)
		}
	}
	return res, nil
}

func (f *Flickr) PhotosetsGetPhotosPage(photosetId string, page int) (*PhotosetsGetPhotos, error) {
	if page < 1 {
		page = 1
	}
	res := &PhotosetsGetPhotos{}
	if err := f.Request("photosets.getPhotos", Params{"photoset_id": photosetId, "page": fmt.Sprintf("%d", page), "extras": "media"}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Wraps PhotosetsGetPhotosPage and returns only the photos from each page
func (f *Flickr) PhotosGetNotInSetWalk() ([]*PhotosetsGetPhotosPhoto, error) {
	thePage, err := f.PhotosGetNotInSetPage(1)
	if err != nil {
		return nil, err
	}
	res := []*PhotosetsGetPhotosPhoto{}
	for _, photo := range thePage.Photos.Photo {
		res = append(res, photo)
	}
	pages := thePage.Photos.Pages
	for page := 2; page <= pages; page++ {
		thePage, err = f.PhotosGetNotInSetPage(page)
		if err != nil {
			return res, err // partial with error
		}
		for _, photo := range thePage.Photos.Photo {
			res = append(res, photo)
		}
	}
	return res, nil
}

func (f *Flickr) PhotosGetNotInSetPage(page int) (*PhotosGetNotInSet, error) {
	if page < 1 {
		page = 1
	}
	res := &PhotosGetNotInSet{}
	if err := f.Request("photos.getNotInSet", Params{"page": fmt.Sprintf("%d", page), "extras": "media", "per_page": "500"}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (f *Flickr) PhotosGetSizes(photoId string) (*PhotosGetSizes, error) {
	res := &PhotosGetSizes{}
	if err := f.Request("photos.getSizes", Params{"photo_id": photoId}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (f *Flickr) PhotosGetInfo(photoId string) (*PhotosGetInfo, error) {
	res := &PhotosGetInfo{}
	if err := f.Request("photos.getInfo", Params{"photo_id": photoId}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (f *Flickr) Request(method string, params Params, v interface{}) error {
	url := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.%s&api_key=%s&format=json&nojsoncallback=1", method, f.ApiKey)
	for key, value := range params {
		url = fmt.Sprintf("%s&%s=%s", url, key, value)
	}
	response, err := f.client.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Error, got response %d %s", response.StatusCode, response.Status))
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = Parse(bytes, &v); err != nil {
		return err
	}
	return nil
}

func (f *Flickr) Download(url string, destinationFile string) error {
	tempFile := destinationFile + ".tmp"
	destinationDir := path.Dir(destinationFile)
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		if err := os.Mkdir(destinationDir, 0755); err != nil {
			return err
		}
	}
	resp, err := f.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	out.Close()
	return os.Rename(tempFile, destinationFile)
}

type FailResponse struct {
	Stat    string
	Code    int
	Message string
}

func Parse(data []byte, v interface{}) error {
	fail := Fail(data)
	if fail != nil {
		return fail
	}
	err := json.Unmarshal(data, v)
	if err != nil {
		fail := Fail(data)
		if fail != nil {
			return fail
		}
		return err
	}
	return nil
}

func Fail(data []byte) error {
	fail := &FailResponse{}
	err := json.Unmarshal(data, fail)

	if err == nil && fail.Stat == "fail" {
		return errors.New(fail.Message)
	}

	return nil
}

type PhotosetsGetPhotos struct {
	Photoset struct {
		Id     string                     `json:"id"`
		Photos []*PhotosetsGetPhotosPhoto `json:"photo"`
		Page   string                     `json:"page"`
		Pages  int                        `json:"pages"`
		Title  string                     `json:"title"`
		Total  int                        `json:"total"`
	} `json:"photoset"`
	Stat string `json:"stat"`
}

type PhotosGetNotInSet struct {
	Photos struct {
		Page    int                        `json:"page"`
		Pages   int                        `json:"pages"`
		Perpage int                        `json:"perpage"`
		Total   int                        `json:"total"`
		Photo   []*PhotosetsGetPhotosPhoto `json:"photo"`
	} `json:"photos"`
	Stat string `json:"stat"`
}

type PhotosetsGetPhotosPhoto struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
	Server string `json:"server"`
	Farm   int    `json:"farm"`
	Title  string `json:"title"`
	Media  string `json:"media"`
}

type PhotosetsGetList struct {
	Photosets struct {
		CanCreate int                         `json:"cancreate"`
		Page      int                         `json:"page"`
		Pages     int                         `json:"pages"`
		Total     int                         `json:"total"`
		Photosets []*PhotosetsGetListPhotoset `json:"photoset"`
	} `json:"photosets"`
}

type PhotosetsGetListPhotoset struct {
	Id          string  `json:"id"`
	Owner       string  `json:"owner"`
	Secret      string  `json:"secret"`
	Server      string  `json:"server"`
	Farm        int     `json:"farm"`
	CountPhotos int     `json:"count_photos"`
	CountVideos int     `json:"count_videos"`
	Title       Content `json:"title"`
	Description Content `json:"description"`
	CanComment  int     `json:"can_comment"`
	DateCreate  string  `json:"date_create"`
	DateUpdate  string  `json:"date_update"`
}

type Content struct {
	Text string `json:"_content"`
}

type PhotosGetSizes struct {
	Sizes struct {
		CanBlog     int                   `json:"canblog"`
		CanPrint    int                   `json:"canprint"`
		CanDownload int                   `json:"candownload"`
		Size        []*PhotosGetSizesSize `json:"size"`
	} `json:"sizes"`
	Stat string `json:"stat"`
}

type PhotosGetSizesSize struct {
	Label  string `json:"label"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Source string `json:"source"`
	Url    string `json:"url"`
	Media  string `json:"media"`
}

type PhotosGetInfo struct {
	Photo *PhotosGetInfoPhoto `json:"photo"`
	Stat  string              `json:"stat"`
}

type PhotosGetInfoPhoto struct {
	Id           string `json:"id"`
	Secret       string `json:"secret"`
	Server       string `json:"server"`
	Farm         int    `json:"farm"`
	DateUploaded string `json:"farm"`
	IsFavorite   int    `json:"isfavorite"`
	// License int `json:"license"`
	// SafetyLevel int `json:"safety_level"`
	Rotation   int     `json:"rotation"`
	Title      Content `json:"title"`
	Visibility struct {
		IsPublic int `json:"ispublic"`
		IsFriend int `json:"isfriend"`
		IsFamily int `json:"isfamily"`
	} `json:"visibility"`
	Dates struct {
		Posted     string `json:"posted"`
		Taken      string `json:"taken"`
		Lastupdate string `json:"lastupdate"`
	} `json:"dates"`
	Media string `json:"media"`
}

/*
{
  "photo": {
    "id": "50686688091",
    "secret": "cd8310012b",
    "server": "65535",
    "farm": 66,
    "dateuploaded": "1607263874",
    "isfavorite": 0,
    "license": 0,
    "safety_level": 0,
    "rotation": 0,
    "originalsecret": "71f7de4a92",
    "originalformat": "jpg",
    "owner": {
      "nsid": "59362368@N00",
      "username": "AtmakuriDavidsen",
      "realname": "Soren Atmakuri Davidsen",
      "location": "",
      "iconserver": "8273",
      "iconfarm": 9,
      "path_alias": "sorend"
    },
    "title": {
      "_content": "2020-11-01 11.16.46"
    },
    "description": {
      "_content": ""
    },
    "visibility": {
      "ispublic": 0,
      "isfriend": 0,
      "isfamily": 0
    },
    "dates": {
      "posted": "1607263874",
      "taken": "2020-11-01 11:16:46",
      "takengranularity": 0,
      "takenunknown": 0,
      "lastupdate": "1607263875"
    },
    "permissions": {
      "permcomment": 3,
      "permaddmeta": 2
    },
    "views": 0,
    "editability": {
      "cancomment": 1,
      "canaddmeta": 1
    },
    "publiceditability": {
      "cancomment": 1,
      "canaddmeta": 0
    },
    "usage": {
      "candownload": 1,
      "canblog": 1,
      "canprint": 1,
      "canshare": 1
    },
    "comments": {
      "_content": 0
    },
    "notes": {
      "note": []
    },
    "people": {
      "haspeople": 0
    },
    "tags": {
      "tag": []
    },
    "urls": {
      "url": [
        {
          "type": "photopage",
          "_content": "https://www.flickr.com/photos/sorend/50686688091/"
        }
      ]
    },
    "media": "photo"
  },
  "stat": "ok"
  }
*/
