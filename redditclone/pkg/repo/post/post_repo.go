package post

import (
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/api/v1/services"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/comment"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/vote"
	"strconv"
	"sync"
)

func CopyDataFromPost(post Post) Post {
	post.Comments = CopyCommentsFields(post)
	post.Author = CopyUserFields(post)
	post.Votes = CopyVoteFields(post)
	return post
}

func CopyVoteFields(post Post) []*vote.Vote {
	copiedVotes := make([]*vote.Vote, 0, len(post.Votes))
	for _, curVote := range post.Votes {
		copiedVotes = append(copiedVotes, &vote.Vote{
			Value:  curVote.Value,
			UserID: curVote.UserID,
		})
	}
	return copiedVotes
}

func CopyCommentsFields(post Post) []*comment.Comment {
	copiedComments := make([]*comment.Comment, 0, len(post.Comments))
	for _, curComment := range post.Comments {
		copiedComments = append(copiedComments, &comment.Comment{
			TimeCreated: curComment.TimeCreated,
			Author:      curComment.Author,
			Body:        curComment.Body,
			ID:          curComment.ID,
		})
	}
	return copiedComments
}
func CopyUserFields(post Post) *user.User {
	copiedAuthor := &user.User{
		ID:    post.Author.ID,
		Login: post.Author.Login,
	}
	return copiedAuthor
}

func PercentageCount(post *Post) int {
	var positiveScore int
	for _, vote := range post.Votes {
		if vote.Value == 1 {
			positiveScore++
		}
	}
	return (100 * positiveScore) / len(post.Votes)
}

type PostMemoryRepository struct {
	data             []*Post
	mu               *sync.RWMutex
	previouslyPostID int
}

func NewMemoryRepo() *PostMemoryRepository {
	return &PostMemoryRepository{
		data: make([]*Post, 0),
		mu:   &sync.RWMutex{},
	}
}

func (p *PostMemoryRepository) GetAllPosts() ([]Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allPosts := make([]Post, 0)
	for _, post := range p.data {
		allPosts = append(allPosts, CopyDataFromPost(*post))
	}
	return allPosts, nil
}

func (p *PostMemoryRepository) CategoryPosts(category string) ([]Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allPostsByRating := make([]Post, 0)
	for _, post := range p.data {
		if post.Category == category {
			allPostsByRating = append(allPostsByRating, CopyDataFromPost(*post))
		}
	}
	return allPostsByRating, nil
}

func (p *PostMemoryRepository) UserPosts(userName string) ([]Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	userPosts := make([]Post, 0)
	for _, post := range p.data {
		if post.Author.Login == userName {
			userPosts = append(userPosts, CopyDataFromPost(*post))
		}
	}
	return userPosts, nil
}

func (p *PostMemoryRepository) PostByID(postID string) (Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, post := range p.data {
		if post.ID == postID {
			post.Views++
			postToReturn := CopyDataFromPost(*post)
			return postToReturn, nil
		}
	}
	return Post{}, ErrNoPost
}

func (p *PostMemoryRepository) AddPost(post *Post, user *user.User) (Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.previouslyPostID++
	post.ID = strconv.Itoa(p.previouslyPostID)
	post.Author = user
	post.Votes = make([]*vote.Vote, 0, 1)
	post.Votes = append(post.Votes, &vote.Vote{
		Value:  1,
		UserID: user.ID,
	})
	if post.Type == "text" {
		post.URL = ""
	} else {
		post.Text = ""
	}
	post.Views = 0
	post.Comments = make([]*comment.Comment, 0)
	post.Created = services.GetCreationTime()
	post.UpvotePercentage = 100
	post.Score = 1
	p.data = append(p.data, post)
	postToReturn := CopyDataFromPost(*post)
	return postToReturn, nil
}

func (p *PostMemoryRepository) AddComment(user *user.User, commentBody string, postID string) (Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, postInBase := range p.data {
		if postInBase.ID == postID {
			p.previouslyPostID++
			newComment := &comment.Comment{
				TimeCreated: services.GetCreationTime(),
				ID:          strconv.Itoa(p.previouslyPostID),
				Author:      user,
				Body:        commentBody,
			}
			postInBase.Comments = append(postInBase.Comments, newComment)
			return CopyDataFromPost(*postInBase), nil
		}
	}
	return Post{}, ErrNoAccess
}

func (p *PostMemoryRepository) DeleteComment(userID, postID, commentID string) (Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, postInBase := range p.data {
		if postInBase.ID == postID {
			for i, commentInBase := range postInBase.Comments {
				if commentInBase.ID == commentID {
					if userID != commentInBase.Author.ID {
						return *postInBase, ErrNoAccess
					}
					postInBase.Comments = append(postInBase.Comments[:i], postInBase.Comments[i+1:]...)
					postToReturn := CopyDataFromPost(*postInBase)
					return postToReturn, nil
				}
			}
		}
	}
	return Post{}, ErrNoPost
}

func (p *PostMemoryRepository) UpVote(postID, userID string) (Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, postInBase := range p.data {
		if postInBase.ID == postID {
			for _, voteInBase := range postInBase.Votes {
				if voteInBase.UserID == userID {
					if voteInBase.Value == -1 {
						postInBase.Score += 2
						voteInBase.Value = 1
						postInBase.UpvotePercentage = PercentageCount(postInBase)
					}
					postToReturn := CopyDataFromPost(*postInBase)
					return postToReturn, nil
				}
			}
			postInBase.Votes = append(postInBase.Votes, vote.CreateVote(userID, 1))
			postInBase.UpvotePercentage = PercentageCount(postInBase)
			postInBase.Score++
			postToReturn := CopyDataFromPost(*postInBase)
			return postToReturn, nil
		}
	}
	return Post{}, ErrNoPost
}

func (p *PostMemoryRepository) DownVote(postID, userID string) (Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, postInBase := range p.data {
		if postInBase.ID == postID {
			for _, voteInBase := range postInBase.Votes {
				if voteInBase.UserID == userID {
					if voteInBase.Value == 1 {
						postInBase.Score -= 2
						voteInBase.Value = -1
						postInBase.UpvotePercentage = PercentageCount(postInBase)
					}
					postToReturn := CopyDataFromPost(*postInBase)
					return postToReturn, nil
				}
			}
			postInBase.Votes = append(postInBase.Votes, vote.CreateVote(userID, -1))
			postInBase.UpvotePercentage = PercentageCount(postInBase)
			postInBase.Score--
			postToReturn := CopyDataFromPost(*postInBase)
			return postToReturn, nil
		}
	}
	return Post{}, ErrNoPost
}

func (p *PostMemoryRepository) UnVote(postID, userID string) (Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, postInBase := range p.data {
		if postInBase.ID == postID {
			for i, voteInBase := range postInBase.Votes {
				if voteInBase.UserID == userID {
					postInBase.Votes = append(postInBase.Votes[:i], postInBase.Votes[i+1:]...)
					if len(postInBase.Votes) == 0 {
						postInBase.Score = 0
						postInBase.UpvotePercentage = 0
					} else if voteInBase.Value == 1 {
						postInBase.Score--
						postInBase.UpvotePercentage = PercentageCount(postInBase)
					}
					postToReturn := CopyDataFromPost(*postInBase)
					return postToReturn, nil
				}
			}
			postToReturn := CopyDataFromPost(*postInBase)
			return postToReturn, nil
		}
	}
	return Post{}, ErrNoPost
}

func (p *PostMemoryRepository) DeletePost(userID, postID string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, postInBase := range p.data {
		if postInBase.ID == postID {
			if userID != postInBase.Author.ID {
				return false, ErrNoAccess
			}
			p.data = append(p.data[:i], p.data[i+1:]...)
			return true, nil
		}
	}
	return false, ErrNoPost
}
