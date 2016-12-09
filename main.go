package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattermost/platform/model"
)

var mmClient *model.Client
var mmUser *model.User
var mmTeam *model.Team
var initialLoad *model.InitialLoad
var mmChannel *model.Channel

var url string
var email string
var password string
var teamName string
var channelName string
var message string

func main() {

	flag.StringVar(&url, "url", "http://mattermost.com", "url of mattermost")
	flag.StringVar(&email, "email", "login@mattermost.com", "email address to login")
	flag.StringVar(&password, "password", "P4Ssw0rd", "user password")
	flag.StringVar(&teamName, "team", "", "mattermost team name")
	flag.StringVar(&channelName, "channel", "", "mattermost channel name")
	flag.StringVar(&message, "message", "", "message")

	flag.Parse()

	fmt.Printf("Url: %v\n", url)
	fmt.Printf("Email: %v\n", email)
	fmt.Printf("Password: %v\n", password)
	fmt.Printf("Team: %v\n", teamName)
	fmt.Printf("Channel: %v\n", channelName)
	fmt.Printf("Message: %v\n", message)

	mmClient = model.NewClient(url)
	makeSureServerIsRunning()
	loginAsTheUser()
	initialLoadModel()
	sendMessageToChannel(channelName, message)
}

func makeSureServerIsRunning() {
	if props, err := mmClient.GetPing(); err != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["version"])
	}
}

func loginAsTheUser() {
	if loginResult, err := mmClient.Login(email, password); err != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		os.Exit(1)
	} else {
		mmUser = loginResult.Data.(*model.User)
	}
}

func initialLoadModel() {
	if initialLoadResults, err := mmClient.GetInitialLoad(); err != nil {
		println("We failed to get the initial load")
		os.Exit(1)
	} else {
		initialLoad = initialLoadResults.Data.(*model.InitialLoad)
		for _, team := range initialLoad.Teams {
			if team.Name == teamName {
				//Team found
				fmt.Printf("Team found Id: %v\n", team.Id)
				mmClient.SetTeamId(team.Id)
				return
			}
		}

		//Team not found
		println("We do not appear to be a member of the team '" + teamName + "'")
		os.Exit(1)
	}
}

func sendMessageToChannel(channelName, msg string) (string, error) {
	fmt.Printf("Message: %v\n", msg)
	channelsResult, err := mmClient.GetChannels(channelName)
	if err != nil {
		return "", err
	}

	channelList := channelsResult.Data.(*model.ChannelList)
	for _, channel := range *channelList {
		if channel.Name == channelName {
			mmChannel = channel
			break
		}
	}

	post := &model.Post{ChannelId: mmChannel.Id, Message: msg, UserId: mmUser.Id}

	createPostResult, createPostErr := mmClient.CreatePost(post)

	if createPostErr != nil {
		fmt.Printf("Error: %v\n", createPostErr.Error())
		return "", createPostErr
	}
	fmt.Printf("RequestId: %v\n", createPostResult.RequestId)

	return "", nil
}
