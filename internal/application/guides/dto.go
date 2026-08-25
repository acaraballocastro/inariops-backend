package guidesapp

import (
	"inariops/internal/domain"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/guides/zones"
)

type GuidesDetail struct {
	Guide     domain.GuideUser     `json:"guide"`
	Languages []languages.Language `json:"languages"`
	Zones     []zones.Zone         `json:"zones"`
}

type CreateGuideRequest struct {
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	Phone          string               `json:"phone"`
	MaxToursPerDay int                  `json:"max_tours_per_day"`
	Languages      []languages.Language `json:"languages"`
	Zones          []zones.Zone         `json:"zones"`
}

type UpdateGuideRequest struct {
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	Phone          string               `json:"phone"`
	MaxToursPerDay int                  `json:"max_tours_per_day"`
	Languages      []languages.Language `json:"languages"`
	Zones          []zones.Zone         `json:"zones"`
}

type LanguageGuide struct {
	GuideID    string `json:"guide_id"`
	LanguageID string `json:"language_id"`
}

type ZoneGuide struct {
	GuideID string `json:"guide_id"`
	ZoneID  string `json:"zone_id"`
}
