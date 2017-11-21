package routers

import (
	"DockerPlatform/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login.html", &controllers.LoginController{}, "post:Login")
	beego.Router("/logout.html", &controllers.LoginController{}, "get:Logout")
	beego.Router("/dashboard.html", &controllers.DashboardController{})
	beego.Router("/namespaces.html", &controllers.NamespacesController{})
	beego.Router("/workload.html", &controllers.WorkloadController{})
	beego.Router("/deployment.html", &controllers.DeploymentController{})
	beego.Router("/service.html", &controllers.ServiceController{})
	beego.Router("/pods.html", &controllers.PodsController{})
	beego.Router("/ingress.html", &controllers.IngressController{})

	//  harbor
	beego.Router("/images.html", &controllers.RepoController{}, "get:GetAll")
	beego.Router("/repos.html", &controllers.RepoController{})
	beego.Router("/tags.html", &controllers.TagsController{})

	beego.Router("/bitbucket.html", &controllers.BitbucketController{})
	beego.Router("/bitbucketrepos.html", &controllers.BitbucketController{}, "get:BitbucketRepos")

	// jenkins controller
	// beego.Router("/jenkins.html", &controllers.JenkinsController{}, "get:GetAllJobs")
	// beego.Router("/jkbuildids.html", &controllers.JenkinsController{}, "get:GetAllBuildIds")
	// beego.Router("/jkbuildresp.html", &controllers.JenkinsController{}, "get:GetConsoleOutput")
	// beego.Router("/api/build", &controllers.JenkinsController{}, "post:Build")

	//buildforadmin.html
	beego.Router("/buildforadmin.html", &controllers.BuildController{}, "get:GetForAdmin")

	beego.Router("/buildlist.html", &controllers.BuildController{}, "get:List")
	beego.Router("/build.html", &controllers.BuildController{})
	beego.Router("/buildconfig.html", &controllers.BuildController{}, "get:BuildConfigForGet")
	beego.Router("/buildconfig.html", &controllers.BuildController{}, "post:BuildConfigForPost")
	beego.Router("/api/build/add", &controllers.BuildController{})

	// api/repo/delete
	beego.Router("/api/repo/delete", &controllers.BuildController{}, "post:RepoDelete")

	beego.Router("/api/deploy", &controllers.DeployController{})

	beego.Router("/api/pipelinetoclone", &controllers.BuildController{}, "post:PipelineToClone")
	beego.Router("/api/pipelinetobuild", &controllers.BuildController{}, "post:PipelineToBuild")
	beego.Router("/api/pipelinetopush", &controllers.BuildController{}, "post:PipelineToPush")

	beego.Router("/api/buildmsg/get", &controllers.DashboardController{}, "get:GetBuildMsg")
	beego.Router("/api/podsmsg/get", &controllers.DashboardController{}, "get:GetPodsMsg")
	beego.Router("/api/kubernetes/init", &controllers.KubernetesController{}, "get:INIT")
	beego.Router("/api/kubenetesmsg/get", &controllers.KubernetesController{}, "get:GetKubenetesMsg")
}
