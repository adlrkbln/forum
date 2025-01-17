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

	user, err := s.repo.GetUserByID(user_id)
	if err != nil {
		return err
	}
	author, err := s.repo.GetPostAuthor(post_id)
	if err != nil {
		return err
	}
	if author.Id == user_id {
		return nil
	}

	err = s.NotifyUser(author.Id, post_id, "comment", user.Name + " commented on your post.")
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

func (s *service) GetDislikedPosts(user_id int) ([]*models.Post, error) {
	posts, err := s.repo.GetDislikedPosts(user_id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) DeletePost(post_id int) error {
	err := s.repo.DeletePost(post_id)
	if err != nil {
		return err
	}

	reports, err := s.repo.FindReportsForPost(post_id)
	for _, report := range reports {
		err = s.repo.ChangeReportStatus(report.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) IgnoreReport(report_id int) error {
	err := s.repo.ChangeReportStatus(report_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteComment(commentID int) error {
	return s.repo.DeleteComment(commentID)
}

func (s *service) GetCommentedPosts(userId int) ([]*models.CommentWithPost, error) {
    return s.repo.GetCommentedPostsByUser(userId)
}