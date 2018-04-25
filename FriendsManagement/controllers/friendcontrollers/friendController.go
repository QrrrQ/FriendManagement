package friendcontrollers

import (
	"FriendsManagement/models/friendmodels"
	"encoding/json"
	"fmt"
)

//0 none  1 friends 2 applying（user1 proposer，user2 proposed）3 block（user1 block user2） 4 delete
const (
	none = iota
	friends
	applying
	block
	delete
)

type CretateUserController struct {
	BaseController
}

type requestCretateUserController struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CretateUserController) Post() {
	fmt.Println("CretateUserController POST start")
	var model requestCretateUserController
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	if !c.IsEmail(model.Email) {
		c.echoError(-1, "email address format not available")
		return
	}
	if model.Password == "" {
		c.echoError(-1, "password can not be empty")
		return
	}
	var user friendmodels.User
	user.Email = model.Email
	user.Password = model.Password
	if user.CheckUserEmailExist() {
		c.echoError(-1, "email address already exist")
		return
	}
	err := user.CreateUser()
	if err != nil {
		c.echoError(-1, "can not create user")
		fmt.Println(err)
		return
	}
	c.echoSuccess(nil)
}

// Question 1  make friends with 2 email addresses
type CretateConnectionController struct {
	BaseController
}

type requestCretateConnectionController struct {
	Friends []string `json:"friends"`
}

func (c *CretateConnectionController) Post() {
	fmt.Println("CretateConnectionController POST start")
	var model requestCretateConnectionController
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	//c.ParseForm(&model)
	if len(model.Friends) != 2 {
		c.echoError(-1, "need 2 emails to make friends, more or less is not allowed")
		fmt.Printf("received %d emails", len(model.Friends))
		return
	}
	for _, s := range model.Friends {
		if !c.IsEmail(s) {
			c.echoError(-1, "email address format not available")
			return
		}
	}
	var user1 friendmodels.User
	var user2 friendmodels.User
	user1.Email = model.Friends[0]
	user2.Email = model.Friends[1]
	err := user1.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "Can not get user info by first email address")
		return
	}
	err = user2.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "Can not get user info by second email address")
		return
	}
	var connection friendmodels.Connection
	connection.User1 = user1.UserId
	connection.User2 = user2.UserId
	connection.Relation = friends
	connection.Deleted = false
	err = connection.CreateOrUpdateConnection()
	if err != nil {
		c.echoError(-1, "make friends failed")
		return
	}
	c.echoSuccess(nil)
}

//Q2  get friend list by email
type GetFriendList struct {
	BaseController
}

type requestGetFriendList struct {
	Email string `json:"email"`
}

type responseGetFriendList struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func (c *GetFriendList) Post() {
	fmt.Println("GetFriendList POST start")
	var model requestGetFriendList
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	if !c.IsEmail(model.Email) {
		c.echoError(-1, "email address format not available")
		return
	}
	var user friendmodels.User
	user.Email = model.Email
	err := user.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "Can not get user info by email address")
		return
	}
	friends := user.GetFriendEmailList()
	var res responseGetFriendList
	res.Success = true
	res.Friends = friends
	res.Count = len(friends)
	c.echoData(res)
}

//Q3  get common friend list by emails
type GetCommonFriendList struct {
	BaseController
}

type requestCommonGetFriendList struct {
	Friends []string `json:"friends"`
}

type responseGetCommonFriendList struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func (c *GetCommonFriendList) Post() {
	fmt.Println("GetCommonFriendList POST start")
	var model requestCommonGetFriendList
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	//c.ParseForm(&model)
	if len(model.Friends) != 2 {
		c.echoError(-1, "need 2 emails addresses, more or less is not allowed")
		fmt.Printf("received %d emails", len(model.Friends))
		return
	}
	for _, s := range model.Friends {
		if !c.IsEmail(s) {
			c.echoError(-1, "email address format not available")
			return
		}
	}
	var user1 friendmodels.User
	var user2 friendmodels.User
	user1.Email = model.Friends[0]
	user2.Email = model.Friends[1]
	err := user1.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "Can not get user info by first email address")
		return
	}
	err = user2.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "Can not get user info by second email address")
		return
	}
	friends := user1.GetCommonList(user2)
	var res responseGetCommonFriendList
	res.Success = true
	res.Friends = friends
	res.Count = len(friends)
	c.echoData(res)
}

//Q4 As	a user,I need an API to subscribe to updates from an email address.
type SubscribController struct {
	BaseController
}

//"requestor": "lisa@example.com",   "target": "john@example.com"
type requestSubscribController struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

func (c *SubscribController) Post() {
	fmt.Println("requestSubscribController POST start")
	var model requestSubscribController
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	if !c.IsEmail(model.Requestor) {
		c.echoError(-1, "Requestor email address format not available")
		return
	}
	if !c.IsEmail(model.Target) {
		c.echoError(-1, "Target email address format not available")
		return
	}
	var requestor friendmodels.User
	requestor.Email = model.Requestor
	var target friendmodels.User
	target.Email = model.Target
	err := requestor.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "requestor not exist")
		return
	}
	err = target.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "target not exist")
		return
	}
	var con friendmodels.Connection
	con.User1 = target.UserId
	con.User2 = requestor.UserId
	con.CheckRalation()
	if con.Relation == block || con.Relation == delete {
		c.echoError(-1, "can not send request , you are deleted or blocked")
	}
	var connection friendmodels.Connection
	connection.User1 = requestor.UserId
	connection.User2 = target.UserId
	connection.Relation = applying
	connection.Deleted = false
	err = connection.CreateOrUpdateConnection()
	if err != nil {
		c.echoError(-1, "apply failed")
		return
	}
	c.echoSuccess(nil)
}

//Q5 As	a user,I need an API to block to updates from an email address.
type BlockController struct {
	BaseController
}

//"requestor": "lisa@example.com",   "target": "john@example.com"
type requestBlockController struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

func (c *BlockController) Post() {
	fmt.Println("BlockController POST start")
	var model requestBlockController
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	if !c.IsEmail(model.Requestor) {
		c.echoError(-1, "Requestor email address format not available")
		return
	}
	if !c.IsEmail(model.Target) {
		c.echoError(-1, "Target email address format not available")
		return
	}
	var requestor friendmodels.User
	requestor.Email = model.Requestor
	var target friendmodels.User
	target.Email = model.Target
	err := requestor.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "requestor not exist")
		return
	}
	err = target.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "target not exist")
		return
	}
	var connection friendmodels.Connection
	connection.User1 = requestor.UserId
	connection.User2 = target.UserId
	connection.Relation = block
	connection.Deleted = false
	err = connection.CreateOrUpdateConnection()
	if err != nil {
		c.echoError(-1, "block failed")
		return
	}
	c.echoSuccess(nil)
}

//Q6 	retrieveall	email addresses	that can receive updates
type CanRetriveUpdatesListController struct {
	BaseController
}

//"requestor": "lisa@example.com",   "target": "john@example.com"
type requestCanRetriveUpdatesListController struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

type responseAvailabelEmails struct {
	Success    bool     `json:"success"`
	Recipients []string `json:"recipients"`
}

func (c *CanRetriveUpdatesListController) Post() {
	fmt.Println("CanRetriveUpdatesListController POST start")
	var model requestCanRetriveUpdatesListController
	json.Unmarshal(c.Ctx.Input.RequestBody, &model)
	fmt.Println(model)
	if !c.IsEmail(model.Sender) {
		c.echoError(-1, "sender email address format not available")
		return
	}

	var sender friendmodels.User
	sender.Email = model.Sender
	err := sender.GetUserByEmail()
	if err != nil {
		c.echoError(-1, "sender not exist")
		return
	}
	emails := sender.GetAvailabelEmails()
	var res responseAvailabelEmails
	res.Success = true
	res.Recipients = emails
	c.echoData(res)
}
