package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/middleware"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/comment"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/post"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/utils"
	"go.uber.org/zap"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type PostHandler struct {
	Tmpl     *template.Template
	PostRepo post.PostRepo
	Logger   *zap.SugaredLogger
}

type RequestBody interface {
	Validate() []string
	UnmarshalJSON([]byte) error
}

func RequestPostHandler(logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) (interface{}, *user.User, error) {
	curUser, ok := r.Context().Value(middleware.CurUserKey).(*user.User)
	if !ok {
		textErr := []byte(`{"message": "Context can not converted to user"}`)
		utils.MakeResponse(logger, w, textErr, 500)
		return nil, nil, errors.New("context can not converted to user")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Request can not be readed %s"}`, err))
		utils.MakeResponse(logger, w, textErr, 500)
		return nil, nil, err
	}
	var bodyUnmarshal RequestBody
	if mux.Vars(r)["POST_ID"] != "" {
		bodyUnmarshal = &comment.BodyComment{}
	} else {
		bodyUnmarshal = &post.Post{}
	}
	err = bodyUnmarshal.UnmarshalJSON(body)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "The body cannot be unmarshaled %s"}`, err))
		utils.MakeResponse(logger, w, textErr, 400)
		return nil, nil, err
	}
	logger.Infof("body %v", bodyUnmarshal)
	if validationErrs := bodyUnmarshal.Validate(); len(validationErrs) != 0 {
		var validationErrsJSON []byte
		validationErrsJSON, err = json.Marshal(validationErrs)
		if err != nil {
			textErr := []byte(fmt.Sprintf(`{"message": "The errors cannot be marshaled %s"}`, err))
			utils.MakeResponse(logger, w, textErr, 500)
			return nil, nil, err
		}
		utils.MakeResponse(logger, w, validationErrsJSON, 500)
		return nil, nil, err
	}
	return bodyUnmarshal, curUser, nil
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostRepo.GetAllPosts()
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Posts cannot be ruturned %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	postInJSON, err := json.Marshal(posts)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Cannot marshal to JSON %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	utils.MakeResponse(h.Logger, w, postInJSON, 200)
}

func (h *PostHandler) CategoryPosts(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["CATEGORY_NAME"]
	posts, err := h.PostRepo.CategoryPosts(category)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Posts cannot be ruturned %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	postInJSON, err := json.Marshal(posts)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Cannot marshal to JSON %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	utils.MakeResponse(h.Logger, w, postInJSON, 200)
}

func (h *PostHandler) UserPosts(w http.ResponseWriter, r *http.Request) {
	login := mux.Vars(r)["USER_LOGIN"]
	posts, err := h.PostRepo.UserPosts(login)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Posts cannot be ruturned %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	postInJSON, err := json.Marshal(posts)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Cannot marshal to JSON %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	utils.MakeResponse(h.Logger, w, postInJSON, 200)
}

func (h *PostHandler) PostByID(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	posts, err := h.PostRepo.PostByID(postID)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Posts cannot be ruturned %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	postInJSON, err := json.Marshal(posts)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Cannot marshal to JSON %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	utils.MakeResponse(h.Logger, w, postInJSON, 200)
}

func (h *PostHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	bodyUnmarshal, curUser, err := RequestPostHandler(h.Logger, w, r)
	unmarshalBody, ok := bodyUnmarshal.(*post.Post)
	if !ok {
		textErr := []byte(fmt.Sprintf(`{"message": "Invalid body %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPost, err := h.PostRepo.AddPost(unmarshalBody, curUser)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "The post can not be added %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPostJSON, err := json.Marshal(newPost)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "New post can not be marhaled %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	h.Logger.Infof("The new post added %v", newPostJSON)
	utils.MakeResponse(h.Logger, w, newPostJSON, 201)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	bodyUnmarshal, curUser, err := RequestPostHandler(h.Logger, w, r)
	unmarshalBody, ok := bodyUnmarshal.(*comment.BodyComment)
	if !ok {
		textErr := []byte(fmt.Sprintf(`{"message": "Invalid body %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPost, err := h.PostRepo.AddComment(curUser, unmarshalBody.Body, postID)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "The comment can not be added %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPostJSON, err := json.Marshal(newPost)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "New post can not be marhaled %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	h.Logger.Infof("The new comment added %v", newPostJSON)
	utils.MakeResponse(h.Logger, w, newPostJSON, 201)
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	postID := values["POST_ID"]
	commentID := values["COMMENT_ID"]
	curUser, ok := r.Context().Value(middleware.CurUserKey).(*user.User)
	if !ok {
		textErr := []byte(`{"message": "Context can not converted to user"}`)
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPost, err := h.PostRepo.DeleteComment(curUser.ID, postID, commentID)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Comment can not be deleted %v"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPostJSON, err := json.Marshal(newPost)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "New post can not be marhaled %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	h.Logger.Infof("new comment added %v", newPostJSON)
	utils.MakeResponse(h.Logger, w, newPostJSON, 200)
}

func (h *PostHandler) Vote(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	curUser, ok := r.Context().Value(middleware.CurUserKey).(*user.User)
	if !ok {
		textErr := []byte(`{"message": "Context can not converted to user"}`)
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	segmentsURL := strings.Split(r.URL.Path, "/")
	voteCommand := segmentsURL[len(segmentsURL)-1]
	var curPost post.Post
	var err error
	switch voteCommand {
	case "upvote":
		curPost, err = h.PostRepo.UpVote(postID, curUser.ID)
	case "downvote":
		curPost, err = h.PostRepo.DownVote(postID, curUser.ID)
	default:
		curPost, err = h.PostRepo.UnVote(postID, curUser.ID)
	}
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Vote can not be added %v"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	newPostJSON, err := json.Marshal(curPost)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "New post can not be marhaled %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	h.Logger.Infof("new vote added %v", newPostJSON)
	utils.MakeResponse(h.Logger, w, newPostJSON, 200)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	curUser, ok := r.Context().Value(middleware.CurUserKey).(*user.User)
	if !ok {
		textErr := []byte(`{"message": "Context can not converted to user"}`)
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	isDeleted, err := h.PostRepo.DeletePost(curUser.ID, postID)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Post can not be deleted %s"}`, err))
		utils.MakeResponse(h.Logger, w, textErr, 500)
		return
	}
	h.Logger.Infof("Post successfully deleted")
	if isDeleted {
		utils.MakeResponse(h.Logger, w, []byte(`{"message": "success"}`), 200)
		return
	}
}
