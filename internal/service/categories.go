package service

import "forum/internal/models"

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
	return s.repo.CreateCategory(form.Name)
}

func (s *service) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}
