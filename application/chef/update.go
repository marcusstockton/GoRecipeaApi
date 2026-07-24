package chef

func (s *ChefService) UpdateChef(req UpdateChefRequest) (*ChefResponse, error) {
	chefEntity, err := s.Repo.FindByID(req.ID)
	if err != nil {
		return nil, err
	}

	if err := chefEntity.UpdateProfile(req.FirstName, req.LastName, req.Email); err != nil {
		return nil, err
	}

	if req.Password != "" {
		if err := chefEntity.SetPassword(req.Password); err != nil {
			return nil, err
		}
	}

	if err := s.Repo.Save(chefEntity); err != nil {
		return nil, err
	}

	return &ChefResponse{ID: chefEntity.ID, FirstName: chefEntity.FirstName, LastName: chefEntity.LastName, Email: chefEntity.Email}, nil
}
