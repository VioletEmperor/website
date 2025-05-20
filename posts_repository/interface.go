package posts_repository

type PostsRepository interface {
    GetPost(id int) (*Post, error)
    GetPosts() ([]Post, error)
}
