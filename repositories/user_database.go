package repositories

type UserDatabase interface {
	CreateUsers(users Users) error
	GetUserByUsername(username string) (*User, error)
	GetUsersByUsernames(usernames []string) (Users, error)
}

type userDatabase struct {
	client DbClient
}

func NewUserDatabase(client DbClient) UserDatabase {
	return &userDatabase{client: client}
}

func (p *userDatabase) CreateUsers(users Users) error {
	if err := p.client.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

func (p *userDatabase) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := p.client.Take(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *userDatabase) GetUsersByUsernames(usernames []string) (Users, error) {
	var users Users
	if err := p.client.Find(&users, "username IN ?", usernames).Error; err != nil {
		return nil, err
	}
	return users, nil
}
