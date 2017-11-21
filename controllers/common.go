package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	var FilterUser = func(ctx *context.Context) {
		username := ctx.Input.Session("username")
		if username == nil {
			ctx.Redirect(302, "/")
			return
		}
	}
	// beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/dashboard.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/workload.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/deployment.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/namespaces.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/service.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/pods.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/repos.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/tags.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/bitbucket.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/bitbucketrepos.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/build.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/buildforadmin.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/buildlist.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/ingress.html", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/buildconfig.html", beego.BeforeRouter, FilterUser)
	// beego.InsertFilter("/jenkins.html", beego.BeforeRouter, FilterUser)

}
