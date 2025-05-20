package main

import "html/template"

func parse() map[string]*template.Template {
    templates := map[string]*template.Template{}

    templates["about.html"] = template.Must(template.ParseFiles(
        "templates/layout/base.html",
        "templates/about.html",
        "templates/partials/footer.html",
        "templates/partials/header.html"))

    templates["posts.html"] = template.Must(template.ParseFiles(
        "templates/layout/base.html",
        "templates/posts.html",
        "templates/partials/footer.html",
        "templates/partials/header.html",
        "templates/partials/post.html"))

    templates["contact.html"] = template.Must(template.ParseFiles(
        "templates/layout/base.html",
        "templates/contact.html",
        "templates/partials/footer.html",
        "templates/partials/header.html"))

    templates["submit.html"] = template.Must(template.ParseFiles(
        "templates/partials/submit.html"))

    return templates
}
