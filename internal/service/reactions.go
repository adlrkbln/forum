package service

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
	
	err = s.NotifyUser(author.Id, post_id, "like", user.Name + " liked your post.")
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
	
	err = s.NotifyUser(author.Id, post_id, "dislike", user.Name + " disliked your post.")
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
