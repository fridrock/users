package friend

import (
	"github.com/fridrock/users/api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FriendStorage interface {
	AddFriend(uuid.UUID, uuid.UUID) error
	DeleteFriend(uuid.UUID, uuid.UUID) error
	GetFriends(uuid.UUID) ([]api.GetUserDto, error)
}

type FriendStorageImpl struct {
	db *sqlx.DB
}

func (fs *FriendStorageImpl) AddFriend(userId uuid.UUID, friendId uuid.UUID) error {
	q := `INSERT INTO friends(fr1id, fr2id) VALUES ($1, $2),($2, $1)`
	_, err := fs.db.Exec(q, userId, friendId)
	return err
}
func (fs *FriendStorageImpl) DeleteFriend(userId uuid.UUID, friendId uuid.UUID) error {
	q := `DELETE FROM friends WHERE (fr1id = $1 OR fr2id = $2) OR (fr1id = $2 AND fr2id = $1)`
	_, err := fs.db.Exec(q, userId, friendId)
	return err
}
func (fs *FriendStorageImpl) GetFriends(userId uuid.UUID) ([]api.GetUserDto, error) {
	var friendsList []api.GetUserDto
	q := `SELECT
			users.id,
			users.email,
			users.username,
			users.name,
			users.surname
			FROM friends
			LEFT JOIN users ON friends.fr2id = users.id where fr1id = $1;`
	err := fs.db.Select(&friendsList, q, userId)
	return friendsList, err

}

func NewFriendStorage(db *sqlx.DB) FriendStorage {
	return &FriendStorageImpl{
		db: db,
	}
}
