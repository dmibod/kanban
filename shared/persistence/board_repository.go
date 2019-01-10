package persistence

import (
	"context"

	"github.com/dmibod/kanban/shared/domain"
	"github.com/dmibod/kanban/shared/kernel"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// BoardEntity entity
type BoardEntity struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Layout   string        `bson:"layout"`
	Name     string        `bson:"name"`
	Children []string      `bson:"children"`
	Owner    string        `bson:"owner"`
	Shared   bool          `bson:"shared"`
}

// BoardRepository interface
type BoardRepository interface {
	db.RepositoryEntity
	db.Repository
	DomainRepository(context.Context) domain.Repository
}

type boardRepository struct {
	db.Repository
}

func (r *boardRepository) CreateInstance() interface{} {
	return &BoardEntity{}
}

func (r *boardRepository) GetID(entity interface{}) string {
	return entity.(*BoardEntity).ID.Hex()
}

func (r *boardRepository) DomainRepository(ctx context.Context) domain.Repository {
	return &boardDomainRepository{
		ctx:        ctx,
		Repository: r.Repository,
	}
}

// CreateBoardRepository creates new repository
func CreateBoardRepository(f db.RepositoryFactory) BoardRepository {
	r := &boardRepository{}
	r.Repository = f.CreateRepository("boards", r)
	return r
}

type boardDomainRepository struct {
	ctx context.Context
	db.Repository
}

func (r *boardDomainRepository) Fetch(id kernel.Id) (interface{}, error) {
	persistent, err := r.Repository.FindByID(r.ctx, string(id))
	if err != nil {
		return nil, err
	}
	entity, ok := persistent.(*BoardEntity)
	if !ok {
		return nil, domain.ErrInvalidType
	}
	children := []kernel.Id{}
	for _, id := range entity.Children {
		children = append(children, kernel.Id(id))
	}
	return &domain.BoardEntity{
		ID:       kernel.Id(entity.ID.Hex()),
		Owner:    entity.Owner,
		Name:     entity.Name,
		Layout:   entity.Layout,
		Shared:   entity.Shared,
		Children: children,
	}, nil
}

func (r *boardDomainRepository) Persist(entity interface{}) (kernel.Id, error) {
	board, ok := entity.(*domain.BoardEntity)
	if !ok {
		return kernel.EmptyID, domain.ErrInvalidType
	}
	children := []string{}
	for _, id := range board.Children {
		children = append(children, string(id))
	}
	persistent := &BoardEntity{
		Owner:    board.Owner,
		Name:     board.Name,
		Layout:   board.Layout,
		Shared:   board.Shared,
		Children: children,
	}
	if board.ID.IsValid() {
		persistent.ID = bson.ObjectIdHex(string(board.ID))
		return board.ID, r.Repository.Update(r.ctx, persistent)
	} else {
		id, err := r.Repository.Create(r.ctx, persistent)
		if err == nil {
			return kernel.Id(id), nil
		}
		return kernel.EmptyID, err
	}
}
