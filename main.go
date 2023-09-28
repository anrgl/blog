package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Post struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type PostsStorage struct {
	posts map[int64]Post
}

type PostsResponse struct {
	Posts []Post `json:"posts"`
}

type PostResponse struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func main() {
	app := fiber.New()

	storage := &PostsStorage{
		posts: map[int64]Post{
			1: {
				Id:        1,
				Title:     "Post #1",
				Content:   "Posts content #1",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			2: {
				Id:        2,
				Title:     "Post #2",
				Content:   "Posts content #2",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
		},
	}

	// Show newest posts
	app.Get("/posts", func(ctx *fiber.Ctx) error {
		var resp PostsResponse
		posts := storage.GetAllPosts()
		for _, post := range posts {
			resp.Posts = append(resp.Posts, post)
		}
		return ctx.JSON(resp.Posts)
	})

	app.Get("/posts/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		post, err := storage.GetPostById(int64(id))
		if err != nil {
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		return ctx.JSON(PostResponse{
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	})

	log.Fatal(app.Listen(":9988"))
}

func (s *PostsStorage) GetAllPosts() map[int64]Post {
	return s.posts
}

func (s *PostsStorage) GetPostById(id int64) (Post, error) {
	post, exists := s.posts[id]
	if !exists {
		return Post{}, errors.New("Post not found")
	}

	return post, nil
}
