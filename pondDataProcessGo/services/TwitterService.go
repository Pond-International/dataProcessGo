package services

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"pondDataProcessGo/models"
	"pondDataProcessGo/utils"
)

const (
	apikey    = "UdlPCOv9sDgdzTKoI3ruF7D7Ma80NbndlhTHXNpxZ42OO%7C1709993771088134144-63jz0h0BYZqoHuinV9KedljFruaiya"
	bearToken = "AAAAAAAAAAAAAAAAAAAAADK5qwEAAAAAz3pfOmYNJxBheDDRl1XFbQ8DCHY%3DyAzPYE5vLFXL4afdO1i7hLCnbv6PtD8XNzfzlDwCoFFpMZuXgr"
)

type TwitterService struct {
}

func NewTwitterService() *TwitterService {
	return &TwitterService{}
}

func (s *TwitterService) GetUserInfoByID(userIDSlice []int64) []models.User {
	userIDs := utils.Int64SliceToString(userIDSlice)
	//使用twitter官方api
	url := fmt.Sprintf("https://api.twitter.com/2/users?ids=%s&user.fields=id,name,username,created_at,description,profile_image_url,public_metrics", userIDs)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+bearToken)
	req.Header.Set("User-Agent", "v2FullArchiveSearchPython")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	bodyStr := string(body)
	var response models.Response
	err = json.Unmarshal([]byte(bodyStr), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	return response.Data
}

func (s *TwitterService) GetFollowIdsByUserId(userId string, fType bool) []int64 {
	//fType = true 表示请求follower， 反之表示following
	//请求twittertools 获取到followerids
	url := fmt.Sprintf("https://twitter.utools.me/api/base/apitools/followersIds?apiKey=%s&userId=%s", apikey, userId)
	if fType == false {
		url = fmt.Sprintf("https://twitter.utools.me/api/base/apitools/followingsIds?apiKey=%s&userId=%s", apikey, userId)
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "*/*")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	bodystr := string(body)
	var idsRes models.TwitterToolResponse
	err := json.Unmarshal([]byte(bodystr), &idsRes)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	var idsData models.TwitterIds
	err = json.Unmarshal([]byte(idsRes.Data), &idsData)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return idsData.Ids
}
