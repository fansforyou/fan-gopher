package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"

	fangopherrod "github.com/fansforyou/fan-gopher/fans/rod"
	"github.com/fansforyou/fan-gopher/model"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type outputResult struct {
	Error        bool            `json:"error"`
	ErrorCode    string          `json:"errorCode"`
	ErrorMessage string          `json:"errorMessage"`
	ActorDetails []*actorDetails `json:"actorDetails"`
	VideoDetails *videoDetails   `json:"videoDetails"`
}

type actorDetails struct {
	ActorName       string `json:"actorName"`
	ProfileImageURL string `json:"profileImageUrl"`
}

type videoDetails struct {
	VideoDescription string `json:"videoDescription"`
}

func main() {
	logger := buildLogger().Sugar()
	defer logger.Sync()

	var postID int64
	flag.Int64Var(&postID, "postID", 0, "the ID of the OnlyFans post whose data is to be retrieved")

	var creatorName string
	flag.StringVar(&creatorName, "creatorName", "", "the username of the creator whose post data is to be retrieved")

	var verify bool
	flag.BoolVar(&verify, "verify", false, "true if the existence of the post should be verified")

	var getDetails bool
	flag.BoolVar(&getDetails, "get-details", false, "true if the details of the given post should be retrieved")

	flag.Parse()

	// Disable Leakless usage because Windows flags it as a virus
	rodLauncherURL := launcher.New().Leakless(false).MustLaunch()
	// Quiet the logger so that the only output is the eventual JSON
	browser := rod.New().ControlURL(rodLauncherURL).Logger(utils.LoggerQuiet)
	if connectErr := browser.Connect(); connectErr != nil {
		logger.Errorf("Failed to start browser while resolving post %d for creator '%s': %v", postID, creatorName, connectErr)
		printError(connectErr)
		return
	}
	defer browser.MustClose()

	gopher := fangopherrod.NewGopher(browser)

	if verify {
		logger.Infof("Verifying existence of post %d for creator '%s'", postID, creatorName)
		postExists, err := gopher.VerifyExists(postID, creatorName)
		if err != nil {
			logger.Errorf("Failed to verify existence of post %d for creator '%s': %v", postID, creatorName, err)
			printError(err)
		}

		if !postExists {
			logger.Infof("Post %d for creator '%s' does not exist", postID, creatorName)
			fmt.Printf(`{ "error": true, "errorCode": "POST_NOT_FOUND", "errorMessage": "%s"\n }`, fmt.Sprintf("No post for ID %d found for creator '%s'", postID, creatorName))
			return
		} else {
			logger.Infof("Post %d creator '%s' does exist", postID, creatorName)
			fmt.Println(`{ "error": false, "errorCode": "", "errorMessage": "" }`)
			return
		}
	} else if getDetails {
		logger.Info("Retrieving details of post %d for creator '%s'", postID, creatorName)
		postDetails, err := gopher.GetPostDetails(postID, creatorName)
		if err != nil {
			logger.Errorf("Failed to retrieve details of post %d for creator '%s': %v", postID, creatorName, err)
			printError(err)
			return
		}

		outputResult := toOutputResult(postDetails)
		jsonBytes, err := json.Marshal(outputResult)
		if err != nil {
			logger.Errorf("Failed to marshal post details for post %d for creator '%s' to JSON: %v", postID, creatorName, err)
			printError(err)
			return
		}

		fmt.Println(string(jsonBytes))
	} else {
		logger.Warnf("Unknown action requested for post %d for creator '%s'; no actual work will be done", postID, creatorName)
		printError(errors.New("no valid action requested"))
	}
}

func buildLogger() *zap.Logger {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "fan-gopher.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	return zap.New(core)
}

func printError(err error) {
	fmt.Printf(`{ "error": true, "errorCode": "GENERAL", "errorMessage": "%s" }\n`, err.Error())
}

func toOutputResult(postDetails *model.Post) *outputResult {
	result := &outputResult{}
	result.ActorDetails = make([]*actorDetails, len(postDetails.Actors))
	for actorIndex, actor := range postDetails.Actors {
		result.ActorDetails[actorIndex] = &actorDetails{
			ActorName:       actor.ActorName,
			ProfileImageURL: actor.ProfileImageURL,
		}
	}
	result.VideoDetails = &videoDetails{
		VideoDescription: postDetails.VideoDetails.VideoDescription,
	}
	return result
}
