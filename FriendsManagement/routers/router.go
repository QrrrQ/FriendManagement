package routers

import (
	"FriendsManagement/controllers"
	"FriendsManagement/controllers/friendcontrollers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/1.0/createuser", &friendcontrollers.CretateUserController{})
	beego.Router("/api/1.0/makefriends", &friendcontrollers.CretateConnectionController{})
	beego.Router("/api/1.0/getfriendlist", &friendcontrollers.GetFriendList{})
	beego.Router("/api/1.0/getcommonfriends", &friendcontrollers.GetCommonFriendList{})
	beego.Router("/api/1.0/subscribeupdates", &friendcontrollers.SubscribController{})
	beego.Router("/api/1.0/blocksomeone", &friendcontrollers.BlockController{})
	beego.Router("/api/1.0/getcansendmaillist", &friendcontrollers.CanRetriveUpdatesListController{})
}
