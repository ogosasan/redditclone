package main

import (
	"github.com/gorilla/mux"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/api/v1/handlers"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/middleware"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/post"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/session"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"

	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
)

func main() {
	templates := template.Must(template.ParseGlob("./05_web_app/99_hw/redditclone/static/html/*"))

	sessionManager := session.NewSessionManager()
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Error in create logger %s", err)
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Printf("all entries in logger buffer have been deleted")
		}
	}(zapLogger)
	logger := zapLogger.Sugar()

	userRepository := user.NewMemoryRepo()
	postRepository := post.NewMemoryRepo()

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepository,
		Logger:   logger,
		Sessions: sessionManager,
	}

	postHandler := &handlers.PostHandler{
		Tmpl:     templates,
		PostRepo: postRepository,
		Logger:   logger,
	}

	r := mux.NewRouter()
	staticRouter := r.PathPrefix("/static/").Subrouter()
	static := "./05_web_app/99_hw/redditclone/static"
	staticRouter.PathPrefix("/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/posts/", postHandler.GetAllPosts).Methods("GET")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postHandler.CategoryPosts).Methods("GET")
	r.HandleFunc("/api/user/{USER_LOGIN}", postHandler.UserPosts).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postHandler.PostByID).Methods("GET")

	routerAuth := mux.NewRouter()

	r.Handle("/api/posts", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("POST")
	r.Handle("/api/post/{POST_ID}", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("POST")
	r.Handle("/api/post/{POST_ID}/{COMMENT_ID}", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("DELETE")
	r.Handle("/api/post/{POST_ID}/upvote", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("GET")
	r.Handle("/api/post/{POST_ID}/downvote", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("GET")
	r.Handle("/api/post/{POST_ID}/unvote", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("GET")
	r.Handle("/api/post/{POST_ID}", middleware.AuthMiddleware(logger, sessionManager, routerAuth)).Methods("DELETE")

	routerAuth.HandleFunc("/api/posts", postHandler.AddPost).Methods("POST")
	routerAuth.HandleFunc("/api/post/{POST_ID}", postHandler.AddComment).Methods("POST")
	routerAuth.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postHandler.DeleteComment).Methods("DELETE")
	routerAuth.HandleFunc("/api/post/{POST_ID}/upvote", postHandler.Vote).Methods("GET")
	routerAuth.HandleFunc("/api/post/{POST_ID}/downvote", postHandler.Vote).Methods("GET")
	routerAuth.HandleFunc("/api/post/{POST_ID}/unvote", postHandler.Vote).Methods("GET")
	routerAuth.HandleFunc("/api/post/{POST_ID}", postHandler.DeletePost).Methods("DELETE")

	accessLogMiddleware := middleware.AccessLogMiddleware(logger, r)
	mux := middleware.PanicMiddleware(logger, accessLogMiddleware)

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = templates.ExecuteTemplate(w, "index.html", struct{}{})
		if err != nil {
			log.Printf("Error in template transmission %s", err)
		}
	})

	addr := ":8081"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Printf("Error in start service %s", err)
	}
}
