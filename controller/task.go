package controller

import (
	"context"
	// "encoding/json"
	// "errors"
	"math/rand"
	"fmt"
	"time"
	"net/http"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"twitter-analytics/util"
	"twitter-analytics/util/constant"
)

func favoriteTweets(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	newCtx, _ := context.WithTimeout(ctx, time.Duration(300*time.Second))
	fmt.Println(newCtx)

	// Authorize Twitter API
	config := oauth1.NewConfig(constant.TwitterApiKey, constant.TwitterApiSecretKey)
	twitterToken := oauth1.NewToken(constant.TwitterAccessToken, constant.TwitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, twitterToken)
	twitterClient := twitter.NewClient(httpClient)

	queries := []string{"筋トレ", "筋肉エンジニア", "ベンチプレス", "デッドリフト", "スクワット", "筋肉痛"}
	shuffle(queries)

	searchTweets, resp, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{
		Query: queries[0],
		Lang: "ja",
		ResultType: "latest",
		Count: 200,
	})

	if err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(util.ResponseJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
	}

	if resp.StatusCode != 200 {
		log.Errorf(ctx, fmt.Sprintf("twitter api (GET search/tweets) response error: %s=", resp.Status))
		return err
	}

	var tweetIDsFiltered []int64
	for _, tweet := range searchTweets.Statuses {
		if tweet.Favorited == true {
			fmt.Println("already favorited the tweet")
			continue
		}
		if tweet.InReplyToStatusID != 0 || tweet.InReplyToUserIDStr != "" {
			fmt.Println("reply tweet to someone or include url")
			continue
		}

		if tweet.FavoriteCount > 30 || tweet.User.FollowersCount > 500 {
			fmt.Println("over 30 favorited or over 500 followers")
			continue
		}
		fmt.Println(tweet.Text)
		tweetIDsFiltered = append(tweetIDsFiltered, tweet.ID)
	}

	for _, id := range tweetIDsFiltered {
		favoriteTweets, resp, err := twitterClient.Favorites.Create(&twitter.FavoriteCreateParams{
			ID: id,
		})

		if err != nil {
			log.Errorf(ctx, err.Error())
			continue
		}

		if resp.StatusCode != 200 {
			log.Errorf(ctx, fmt.Sprintf("twitter api (GET favorites/create) response error: %s=", resp.Status))
			return err
		}
		fmt.Println(favoriteTweets)
		time.Sleep(1 * time.Second)
	}
	return c.JSON(util.ResponseJSON(http.StatusOK, http.StatusText(http.StatusOK)))
}

func shuffle(a []string) {
    rand.Seed(time.Now().UnixNano())
    for i := range a {
        j := rand.Intn(i + 1)
        a[i], a[j] = a[j], a[i]
    }
}
