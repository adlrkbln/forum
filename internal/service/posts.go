package service

import (
	"forum/internal/models"
)

func (s *service) InsertPost(form models.PostCreateForm, data *models.TemplateData) (int, error) {
	id, err := s.repo.InsertPost(data.User.Id, form.Title, form.Content)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *service) GetPost(id int) (*models.Post, error) {
	post, err := s.repo.GetPost(id)
	if err != nil {
		return nil, err
	}

	comments, err := s.repo.GetCommentsForPost(post.Id)
	if err != nil {
		return nil, err
	}

	post.Comments = comments
	return post, nil
}

func (s *service) GetAllPosts() ([]*models.Post, error) {
	posts, err := s.repo.GetAllPosts()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

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

func (s *service) GetAllComments() ([]*models.Comment, error) {
	comments, err := s.repo.GetAllComments()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *service) InsertComment(post_id int, user_id int, content string) error {
	err := s.repo.InsertComment(post_id, user_id, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetCreatedPosts(user_id int) ([]*models.Post, error) {
	posts, err := s.repo.GetCreatedPosts(user_id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) GetLikedPosts(user_id int) ([]*models.Post, error) {
	posts, err := s.repo.GetLikedPosts(user_id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) AddLikePost(post_id int, user_id int) error {
	reaction, err := s.repo.CheckUserReactionsPost(post_id, user_id)
	if err != nil {
		return err
	}
	if reaction != 0 {
		err = s.repo.RemoveUserReactionPost(post_id, user_id, reaction)
		if err != nil {
			return err
		}
	}

	if reaction == 1 {
		return nil
	}

	err = s.repo.AddLikePost(post_id, user_id)
	if err != nil {
		return err
	}

	err = s.repo.InsertUserReactionPost(post_id, user_id, 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddDislikePost(post_id int, user_id int) error {
	reaction, err := s.repo.CheckUserReactionsPost(post_id, user_id)
	if err != nil {
		return err
	}
	if reaction != 0 {
		err = s.repo.RemoveUserReactionPost(post_id, user_id, reaction)
		if err != nil {
			return err
		}
	}
	if reaction == -1 {
		return nil
	}
	err = s.repo.AddDislikePost(post_id, user_id)
	if err != nil {
		return err
	}
	err = s.repo.InsertUserReactionPost(post_id, user_id, -1)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddLikeComment(comment_id, user_id int) error {
	reaction, err := s.repo.CheckUserReactionComment(comment_id, user_id)
	if err != nil {
		return err
	}
	if reaction != 0 {
		err = s.repo.RemoveUserReactionComment(comment_id, user_id, reaction)
		if err != nil {
			return err
		}
	}

	if reaction == 1 {
		return nil
	}

	err = s.repo.AddLikeComment(comment_id, user_id)
	if err != nil {
		return err
	}

	err = s.repo.InsertUserReactionComment(comment_id, user_id, 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddDislikeComment(comment_id, user_id int) error {
	reaction, err := s.repo.CheckUserReactionComment(comment_id, user_id)
	if err != nil {
		return err
	}
	if reaction != 0 {
		err = s.repo.RemoveUserReactionComment(comment_id, user_id, reaction)
		if err != nil {
			return err
		}
	}
	if reaction == -1 {
		return nil
	}
	err = s.repo.AddDislikeComment(comment_id, user_id)
	if err != nil {
		return err
	}
	err = s.repo.InsertUserReactionComment(comment_id, user_id, -1)
	if err != nil {
		return err
	}

	return nil
}
