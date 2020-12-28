
package backup

import (
	"fmt"
	//"io/ioutil"
	//"flag"
	//"os"
	//"path"
	//"strings"
	//"strconv"
	"github.com/stvp/slug"
	//"github.com/sorend/gokmp/pkg/flickr"
)


func Run(nsid string, destination string) error {

	// config slug package
	slug.Replacement = '-'

	/*
	client := &flickr.Client{
		Key: os.Getenv("FLICKR_API_KEY"),
		Token: os.Getenv("FLICKR_OAUTH_TOKEN"),
	}
	err := doBackup(client, *backupNsid, *backupDestination)
	if err != nil {
		panic(err)
	}
	*/
	fmt.Println("All done :-)")
	return nil
}
