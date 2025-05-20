package posts_repository

import "time"

type Post struct {
    ID      int
    Title   string
    Author  string
    Created time.Time
    Edited  time.Time
    Body    string
}
