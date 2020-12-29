
package login

import (
	"fmt"
	"io/ioutil"
	"github.com/pkg/browser"
	"github.com/sorend/gokmp/pkg/flickr"
	"github.com/spf13/viper"
)


func Run(ApiKey string, ApiSecret string) error {
	client, err := flickr.NewFlickrLogin(ApiKey, ApiSecret)
	if err != nil {
		return err
	}
	loginUrl, err := client.GetLoginURL()
	if err != nil {
		return err
	}
	fmt.Printf("\nOpening browser for URL:\n   %s\n   (please do it manually if opening browser fails)\n", loginUrl)
	// avoid output from browser plugin
	browser.Stderr = ioutil.Discard
	browser.Stdout = ioutil.Discard
	if err = browser.OpenURL(loginUrl); err != nil {
		return err
	}
	fmt.Printf("\nEnter verifier PIN here: ")
	var verifier string
	_, err = fmt.Scanf("%s", &verifier)

	accessToken, err := client.GetAccessToken(verifier)
	if err != nil {
		return err
	}

	fmt.Printf("Access Token : %s\n", accessToken.Token)
	fmt.Printf("Access Secret: %s\n", accessToken.TokenSecret)

	viper.Set("accessToken", accessToken.Token)
	viper.Set("accessSecret", accessToken.TokenSecret)
	viper.WriteConfig()
	
	return nil
}
