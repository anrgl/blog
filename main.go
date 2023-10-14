package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var storage = &PostsStorage{
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

var maxId = int64(len(storage.posts))

func main() {
	app := fiber.New()

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
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusBadRequest,
				"Message": err.Error(),
			}).Error("Invalid :id params")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		post, err := storage.GetPostById(int64(id))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusNotFound,
				"Message": err.Error(),
			}).Error("Post not found")
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		return ctx.JSON(PostResponse{
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	})

	app.Post("/posts", func(ctx *fiber.Ctx) error {
		var req CreatePostRequest
		if err := ctx.BodyParser(&req); err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusBadRequest,
				"Message": err.Error(),
			}).Error("Create post request error")
		}

		id := getMaxId(&maxId)
		post := Post{
			Id:        id,
			Title:     req.Title,
			Content:   req.Content,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		storage.AddNewPost(post)
		return ctx.SendStatus(fiber.StatusCreated)
	})

	app.Put("/posts/:id", func(ctx *fiber.Ctx) error {
		var req UpdatePostRequest
		if err := ctx.BodyParser(&req); err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusBadRequest,
				"Message": err.Error(),
			}).Error("Update post request error")
		}

		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusBadRequest,
				"Message": err.Error(),
			}).Error("Invalid :id params")
		}

		post, err := storage.GetPostById(int64(id))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusNotFound,
				"Message": err.Error(),
			}).Error("Post not found")
		}

		post.Title = req.Title
		post.Content = req.Content
		post.UpdatedAt = time.Now().Unix()

		storage.UpdatePost(int64(id), post)

		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Delete("/posts/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusBadRequest,
				"Message": err.Error(),
			}).Error("Invalid :id params")
		}

		if _, err := storage.GetPostById(int64(id)); err != nil {
			logrus.WithFields(logrus.Fields{
				"Status":  fiber.StatusNotFound,
				"Message": err.Error(),
			}).Error("Post not found")
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		storage.DeletePost(int64(id))

		return ctx.SendStatus(fiber.StatusNoContent)
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

func (s *PostsStorage) AddNewPost(p Post) {
	s.posts[p.Id] = p
}

func (s *PostsStorage) UpdatePost(id int64, p Post) {
	s.posts[id] = p
}

func (s *PostsStorage) DeletePost(id int64) {
	delete(s.posts, id)
}

func getMaxId(id *int64) int64 {
	*id++
	return *id
}
