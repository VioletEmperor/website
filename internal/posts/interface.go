package posts

type Repository interface {
	GetPost(id int) (*Post, error)
	GetPosts() ([]Post, error)
}
