package fans

import (
	"fmt"

	"github.com/fansforyou/fan-gopher/model"
)

type PostNotFoundErr struct {
	PostID      int64
	CreatorName string
}

func (e *PostNotFoundErr) Error() string {
	return fmt.Sprintf("unable to find post ID %d for creator '%s'", e.PostID, e.CreatorName)
}

type Gopher interface {
	// VerifyExists verifies if a post for the given ID and creator name exists
	VerifyExists(postID int64, creatorName string) (bool, error)

	// GetPostDetails gets the details of the given post by the given creator
	// Returns PostNotFoundErr if the post cannot be found
	GetPostDetails(postID int64, creatorName string) (*model.Post, error)
}
