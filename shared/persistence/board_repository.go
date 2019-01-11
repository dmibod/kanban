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
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
	Children    []string      `bson:"children"`
}

// BoardVisitor type
type BoardVisitor func(*BoardEntity) error

// BoardRepository interface
type BoardRepository interface {
	db.RepositoryEntity
	db.Repository
	DomainRepository(context.Context) domain.Repository
	FindBoardByID(context.Context, kernel.ID) (*BoardEntity, error)
	FindBoards(context.Context, interface{}, BoardVisitor) error
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

func (r *boardRepository) FindBoardByID(ctx context.Context, id kernel.ID) (*BoardEntity, error) {
	if !id.IsValid() {
		return nil, domain.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id))
	if err != nil {
		return nil, err
	}

	board, ok := entity.(*BoardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	return board, nil
}

func (r *boardRepository) FindBoards(ctx context.Context, criteria interface{}, visitor BoardVisitor) error {
	return r.Repository.Find(ctx, criteria, func(entity interface{}) error {
		board, ok := entity.(*BoardEntity)
		if !ok {
			return ErrInvalidType
		}

		return visitor(board)
	})
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

func (r *boardDomainRepository) Fetch(id kernel.ID) (interface{}, error) {
	if !id.IsValid() {
		return nil, domain.ErrInvalidID
	}

	persistent, err := r.Repository.FindByID(r.ctx, string(id))
	if err != nil {
		return nil, err
	}

	entity, ok := persistent.(*BoardEntity)
	if !ok {
		return nil, domain.ErrInvalidType
	}

	return r.mapEntityToDomain(entity), nil
}

func (r *boardDomainRepository) Delete(id kernel.ID) (interface{}, error) {
	entity, err := r.Fetch(id)
	if err == nil {
		err = r.Repository.Remove(r.ctx, string(id))
	}

	return entity, err
}

func (r *boardDomainRepository) Persist(entity interface{}) (kernel.ID, error) {
	board, ok := entity.(domain.BoardEntity)
	if !ok {
		return kernel.EmptyID, domain.ErrInvalidType
	}

	persistent := r.mapDomainToEntity(&board)

	if board.ID.IsValid() {
		return board.ID, r.Repository.Update(r.ctx, persistent)
	}

	id, err := r.Repository.Create(r.ctx, persistent)
	if err != nil {
		return kernel.EmptyID, err
	}

	return kernel.ID(id), nil
}

func (r *boardDomainRepository) mapEntityToDomain(entity *BoardEntity) *domain.BoardEntity {
	children := []kernel.ID{}
	for _, id := range entity.Children {
		children = append(children, kernel.ID(id))
	}

	return &domain.BoardEntity{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
		Children:    children,
	}
}

func (r *boardDomainRepository) mapDomainToEntity(domainEntity *domain.BoardEntity) *BoardEntity {
	children := []string{}
	for _, id := range domainEntity.Children {
		children = append(children, string(id))
	}

	entity := &BoardEntity{
		Owner:       domainEntity.Owner,
		Name:        domainEntity.Name,
		Description: domainEntity.Description,
		Layout:      domainEntity.Layout,
		Shared:      domainEntity.Shared,
		Children:    children,
	}

	if domainEntity.ID.IsValid() {
		entity.ID = bson.ObjectIdHex(string(domainEntity.ID))
	}

	return entity
}
