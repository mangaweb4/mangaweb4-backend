package browse

import "github.com/julienschmidt/httprouter"

func Register(router *httprouter.Router) {
	router.GET(PathThumbnail, GetThumbnailHandler)
	router.POST(PathBrowse, Handler)
	router.POST(PathHistory, historyHandler)
}
