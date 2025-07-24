package posts

type Repository interface {
	GetPost(id int) (*Post, error)
	GetPosts() ([]Post, error)
	GetPostsPaginated(page int) ([]Post, PaginationInfo, error)
	GetTotalPostsCount() (int, error)
	DeletePost(id int) error
	UpdatePost(id int, title, description, body string) error
	CreatePost(title, description, body, author string) error
}
