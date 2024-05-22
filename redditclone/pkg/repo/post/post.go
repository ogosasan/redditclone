package post

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/comment"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/vote"
)

var (
	ErrNoPost   = errors.New("the post was not found")
	ErrNoAccess = errors.New("the action cannot be performed, authorization is required")
)

type Post struct {
	Author           *user.User         `json:"author"`
	Category         string             `json:"category"`
	Comments         []*comment.Comment `json:"comments"`
	Created          string             `json:"created"`
	ID               string             `json:"id"`
	Score            int                `json:"score"`
	Text             string             `json:"text,omitempty"`
	Title            string             `json:"title"`
	Type             string             `json:"type"`
	UpvotePercentage int                `json:"upvotePercentage"`
	Views            int                `json:"views"`
	Votes            []*vote.Vote       `json:"votes"`
	URL              string             `json:"url,omitempty" valid:"url"`
}

type PostRepo interface {
	GetAllPosts() ([]Post, error)
	CategoryPosts(category string) ([]Post, error)
	UserPosts(userID string) ([]Post, error)
	PostByID(postID string) (Post, error)
	AddPost(post *Post, user *user.User) (Post, error)
	AddComment(user *user.User, commentBody string, postID string) (Post, error)
	DeleteComment(userID, postID, commentID string) (Post, error)
	UpVote(postID, userID string) (Post, error)
	DownVote(postID, userID string) (Post, error)
	UnVote(postID, userID string) (Post, error)
	DeletePost(userID, postID string) (bool, error)
}

func (p *Post) UnmarshalJSON(data []byte) error {
	type Alias Post
	post := &Alias{}
	if err := json.Unmarshal(data, post); err != nil {
		return err
	}
	*p = Post(*post)
	return nil
}

func (p *Post) Validate() []string {
	_, err := govalidator.ValidateStruct(p)
	if err == nil {
		return nil
	}
	errors := make([]string, 0)
	if p.Type == "url" && p.Text == "" {
		errors = append(errors, "No URL address")
	}
	if p.Type == "text" && p.URL == "" {
		errors = append(errors, "No text")
	}
	if errs, ok := err.(govalidator.Errors); ok {
		for _, e := range errs {
			errors = append(errors, e.Error())
		}
	}
	return errors
}
