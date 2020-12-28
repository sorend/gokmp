package flickr

import (
	"fmt"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"os"
	"github.com/dghubble/oauth1"
)

func NewFlickr(ApiKey string, ApiSecret string, OAuthToken string, OAuthSecret string) *Flickr {
	config := oauth1.NewConfig(ApiKey, ApiSecret)
	token := oauth1.NewToken(OAuthToken, OAuthSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return &Flickr{
		ApiKey: ApiKey,
		client: httpClient,
	}
}

type Flickr struct {
	ApiKey string
	client *http.Client
}

type Params map[string]string  // Params datatype

func (f *Flickr) PhotosetsGetList(nsid string) ([]*GetListPhotoset, error) {
	response, err := f.Request("photosets.getList", Params{"user_id": nsid})
	if err != nil {
		return nil, err
	}
	photosets := &GetListPhotosetsRaw{}
	err = Parse(response, &photosets)
	return photosets.Photosets.Photosets, err
}

func (f *Flickr) PhotosetsGetPhotosWalk(photosetId string) ([]*PhotosetsGetPhotosPhoto, error) {

	response, err := f.photosetsGetPhotosPage(photosetId, 1)
	if err != nil {
		return nil, err
	}

	res := []*PhotosetsGetPhotosPhoto{}
	for _, photo := range response.Photoset.Photos {
		res = append(res, photo)
	}

	pages := response.Photoset.Pages

	for page := 2; page <= pages; page++ {
		response, err = f.photosetsGetPhotosPage(photosetId, page)
		if err != nil {
			return res, err // partial with error
		}
		for _, photo := range response.Photoset.Photos {
			res = append(res, photo)
		}
	}
	return res, nil
}

func (f *Flickr) photosetsGetPhotosPage(photosetId string, page int) (*PhotosetsGetPhotos, error) {
	res := &PhotosetsGetPhotos{}
	response, err := f.Request("photosets.getPhotos", Params{"photoset_id": photosetId, "page": fmt.Sprintf("%d", page), "extras":"media"})
	if err != nil {
		return res, err
	}
	if err = Parse(response, &res); err != nil {
		return res, err
	}
	return res, nil
}

func (f *Flickr) PhotosGetSizes(photoId string) ([]*PhotosGetSizesSize, error) {
	res := &PhotosGetSizes{}
	response, err := f.Request("photos.getSizes", Params{"photo_id": photoId})
	if err != nil {
		return nil, err
	}
	if err = Parse(response, &res); err != nil {
		return nil, err
	}
	return res.Sizes.Size, nil
}


func (f *Flickr) Request(method string, params Params) ([]byte, error) {
	url := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.%s&api_key=%s&format=json&nojsoncallback=1", method, f.ApiKey)
	for key, value := range params {
		url = fmt.Sprintf("%s&%s=%s", url, key, value)
	}
	response, err := f.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)	
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
	Stat string
	Code int
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

func Fail (data []byte) error {
	fail := &FailResponse{}
	err := json.Unmarshal(data, fail)

	if err == nil && fail.Stat == "fail" {
		return errors.New(fail.Message)
	}

	return nil
}

type PhotosetsGetPhotos struct {
	Photoset struct {
		Id string `json:"id"`
		Photos []*PhotosetsGetPhotosPhoto `json:"photo"`
		Page string `json:"page"`
		Pages int `json:"pages"`
		Title string `json:"title"`
		Total int `json:"total"`
	} `json:"photoset"`
}

type PhotosetsGetPhotosPhoto struct {
	Id string `json:"id"`
	Secret string `json:"secret"`
	Server string `json:"server"`
	Farm int `json:"farm"`
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

type PhotosGetSizes struct {
	Sizes struct {
		CanBlog int `json:"canblog"`
		CanPrint int `json:"canprint"`
		CanDownload int `json:"candownload"`
		Size []*PhotosGetSizesSize `json:"size"`
	} `json:"sizes"`
	Stat string `json:"stat"`
}


type PhotosGetSizesSize struct {
	Label string `json:"label"`
	Width int `json:"width"`
	Height int `json:"height"`
	Source string `json:"source"`
	Url string `json:"url"`
	Media string `json:"media"`
}

/*
   { "sizes": { "canblog": 1, "canprint": 1, "candownload": 1, 
    "size": [
      { "label": "Square", "width": 75, "height": 75, "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_s.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/sq\/", "media": "photo" },
      { "label": "Large Square", "width": "150", "height": "150", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_q.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/q\/", "media": "photo" },
      { "label": "Thumbnail", "width": 75, "height": 100, "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_t.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/t\/", "media": "photo" },
      { "label": "Small", "width": "180", "height": "240", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_m.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/s\/", "media": "photo" },
      { "label": "Small 320", "width": "240", "height": "320", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_n.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/n\/", "media": "photo" },
      { "label": "Small 400", "width": "300", "height": "400", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_w.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/w\/", "media": "photo" },
      { "label": "Medium", "width": "375", "height": "500", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/m\/", "media": "photo" },
      { "label": "Medium 640", "width": "480", "height": "640", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_z.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/z\/", "media": "photo" },
      { "label": "Medium 800", "width": "600", "height": "800", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_c.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/c\/", "media": "photo" },
      { "label": "Large", "width": "768", "height": "1024", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_cd8310012b_b.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/l\/", "media": "photo" },
      { "label": "Large 1600", "width": "1200", "height": "1600", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_b82a7389fb_h.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/h\/", "media": "photo" },
      { "label": "Large 2048", "width": "1536", "height": "2048", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_56c7e90d4b_k.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/k\/", "media": "photo" },
      { "label": "X-Large 3K", "width": "2304", "height": "3072", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_11bccc3722_3k.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/3k\/", "media": "photo" },
      { "label": "X-Large 4K", "width": "3072", "height": "4096", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_b68a015b96_4k.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/4k\/", "media": "photo" },
      { "label": "Original", "width": "3072", "height": "4096", "source": "https:\/\/live.staticflickr.com\/65535\/50686688091_71f7de4a92_o.jpg", "url": "https:\/\/www.flickr.com\/photos\/sorend\/50686688091\/sizes\/o\/", "media": "photo" }
    ] }, "stat": "ok" }
*/

