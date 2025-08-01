package locationrepo

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/pkg/ddd"

	"github.com/google/uuid"
)

// DomainToDTO converts Location domain model to LocationDTO
func DomainToDTO(l *location.Location) LocationDTO {
	return LocationDTO{
		ID:          l.ID().String(),
		Name:        l.Name,
		Latitude:    l.Coordinate.Latitude(),
		Longitude:   l.Coordinate.Longitude(),
		Address:     l.Coordinate.GetAddress(),
		Description: l.Description,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

// DtoToDomain converts LocationDTO to Location domain model
func DtoToDomain(dto LocationDTO) (*location.Location, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return nil, err
	}

	coordinate, err := kernel.NewGeoCoordinate(dto.Latitude, dto.Longitude, dto.Address)
	if err != nil {
		return nil, err
	}

	// Получаем адрес из координат
	var locationAddress string
	if coordinate.GetAddress() != nil {
		locationAddress = *coordinate.GetAddress()
	}

	l := &location.Location{
		BaseAggregate: ddd.NewBaseAggregate(id),
		Name:          dto.Name,
		Coordinate:    coordinate,
		Address:       locationAddress,
		Description:   dto.Description,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}

	return l, nil
}
