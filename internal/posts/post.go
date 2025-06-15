package posts

import "time"

type Post struct {
    ID      int       `db:"id"`
    Title   string    `db:"title"`
    Author  string    `db:"author"`
    Created time.Time `db:"created"`
    Edited  time.Time `db:"edited"`
    Body    string    `db:"body"`
}
