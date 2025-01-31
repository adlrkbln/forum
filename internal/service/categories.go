package service

import (
	"fmt"
	"forum/internal/models"
)

func (s *service) GetCategories() ([]*models.Category, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) GetPostByCategory(id int) ([]*models.Post, error) {
	posts, err := s.repo.GetPostByCategory(id)
	if err != nil {
		return nil, err
	}
	for i, post := range posts {
		categories, err := s.repo.GetCategoriesForPost(post.Id)
		if err != nil {
			return nil, err
		}
		posts[i].Categories = categories
	}
	return posts, nil
}

func (s *service) PostCategoryPost(post_id int, categoryIds []int) error {
	for _, categoryId := range categoryIds {
		err := s.repo.PostCategoryPost(post_id, categoryId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CreateCategory(form models.CategoryCreateForm) error {
	if form.Name == "" {
		return fmt.Errorf("Blank category name")
	}
	return s.repo.CreateCategory(form.Name)
}

func (s *service) DeleteCategory(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid category ID")
	}
	posts, err := s.repo.GetPostByCategory(id)
	if err != nil {
		return err
	}
	for _, post := range posts {
		err = s.repo.DeletePost(post.Id)
		if err != nil {
			return err
		}
	}
	err = s.repo.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
