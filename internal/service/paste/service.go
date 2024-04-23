package paste

import (
	"context"
	"github.com/artshirshov/gastebin/internal/mapper"
	model "github.com/artshirshov/gastebin/internal/model/paste"
	repo "github.com/artshirshov/gastebin/internal/repository/paste"
	hasher "github.com/artshirshov/gastebin/pkg/hash"
	"github.com/artshirshov/gastebin/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

type Service interface {
	GetPasteByHash(ctx context.Context, hash string) (model.ResponseDto, error)

	CreatePaste(ctx context.Context, reqDto model.RequestDto) (model.ResponseDto, error)

	UpdatePaste(ctx context.Context, hash string, reqDto model.RequestDto) (model.ResponseDto, error)

	DeletePaste(ctx context.Context, hash string) error
}

type service struct {
	repository repo.Repository
}

func NewService(repository repo.Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetPasteByHash(
	ctx context.Context,
	hash string,
) (model.ResponseDto, error) {
	if validateErr := validateHash(hash); validateErr != nil {
		return model.ResponseDto{}, errors.New("")
	}

	entity, getPasteError := s.repository.GetPasteByHash(ctx, hash)
	if getPasteError != nil {
		return model.ResponseDto{}, errors.New("paste not found")
	}

	respDto, typeError := mapper.EntityToResponseDto(entity)
	if typeError != nil {
		return model.ResponseDto{}, errors.New("mapping failed")
	}

	return respDto, nil
}

func (s *service) CreatePaste(
	ctx context.Context,
	reqDto model.RequestDto,
) (model.ResponseDto, error) {
	remoteAddr := ctx.Value("remoteAddr").(string)
	if remoteAddr == "" {
		return model.ResponseDto{}, errors.New("remote address is missing or invalid")
	}
	hash := hasher.GenerateHash(remoteAddr)
	newEntity := mapper.PasteDtoToEntity(hash, reqDto)

	entity, createPasteErr := s.repository.CreatePaste(ctx, newEntity)
	if createPasteErr != nil {
		logger.Log.With(zap.Error(createPasteErr)).Error("error during save paste in db")
		return model.ResponseDto{}, errors.New("")
	}

	respDto, typeError := mapper.EntityToResponseDto(entity)
	if typeError != nil {
		return model.ResponseDto{}, errors.New("mapping failed")
	}

	return respDto, nil
}

func (s *service) UpdatePaste(
	ctx context.Context,
	hash string,
	reqDto model.RequestDto,
) (model.ResponseDto, error) {
	if validateErr := validateHash(hash); validateErr != nil {
		return model.ResponseDto{}, errors.New("mapping failed")
	}

	dbEntity, getErr := s.repository.GetPasteByHash(ctx, hash)
	if getErr != nil {
		return model.ResponseDto{}, errors.New("paste not found")
	}
	mapper.UpdateFromDTO(&dbEntity, reqDto)

	entity, updatePasteErr := s.repository.UpdatePaste(ctx, dbEntity)
	if updatePasteErr != nil {
		return model.ResponseDto{}, errors.New("failed to update paste")
	}

	respDto, typeError := mapper.EntityToResponseDto(entity)
	if typeError != nil {
		return model.ResponseDto{}, errors.New("mapping failed")
	}

	return respDto, nil
}

func (s *service) DeletePaste(
	ctx context.Context,
	hash string,
) error {
	deleteErr := s.repository.DeletePaste(ctx, hash)
	if deleteErr != nil {
		return errors.New("failed to delete paste")
	}
	return nil
}

func validateHash(hash string) error {
	// Игнорируем статичные файлы которые запрашивает браузер
	if strings.Contains(hash, ".") {
		return errors.New("hash not include point letter")
	}
	return nil
}
