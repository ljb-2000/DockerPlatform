package routers

import (
	"DockerPlatform/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/dashboard.html", &controllers.DashboardController{})
	beego.Router("/namespaces.html", &controllers.NamespacesController{})
	beego.Router("/deployment.html", &controllers.DeploymentController{})
	beego.Router("/service.html", &controllers.ServiceController{})
	beego.Router("/pods.html", &controllers.PodsController{})

	beego.Router("/image.html", &controllers.ImageController{})
	beego.Router("/repos.html", &controllers.RepoController{})
}
