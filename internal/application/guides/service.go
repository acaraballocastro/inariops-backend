package guidesapp

import (
	"fmt"
	"time"

	"inariops/internal/domain"
	guidedomain "inariops/internal/modules/guides"
	"inariops/internal/modules/guides/availabilities"
	languageguides "inariops/internal/modules/guides/language_guides"
	"inariops/internal/modules/guides/languages"
	zoneguides "inariops/internal/modules/guides/zone_guides"
	"inariops/internal/modules/guides/zones"
	"inariops/internal/modules/users"
	appErrors "inariops/internal/shared/errors"
	"inariops/internal/shared/logger"

	"github.com/google/uuid"
)

type Service struct {
	usersService             *users.Service
	usersRepository          *users.Repository
	guidesRepository         *guidedomain.Repository
	languagesRepository      *languages.Repository
	languageGuidesRepository *languageguides.Repository
	zonesRepository          *zones.Repository
	zoneGuidesRepository     *zoneguides.Repository
	availabilitiesRepository *availabilities.Repository
}

func NewService(
	usersService *users.Service,
	usersRepository *users.Repository,
	guidesRepository *guidedomain.Repository,
	languagesRepository *languages.Repository,
	languageGuidesRepository *languageguides.Repository,
	zonesRepository *zones.Repository,
	zoneGuidesRepository *zoneguides.Repository,
	availabilitiesRepository *availabilities.Repository,
) *Service {
	return &Service{
		usersService:             usersService,
		usersRepository:          usersRepository,
		guidesRepository:         guidesRepository,
		languagesRepository:      languagesRepository,
		languageGuidesRepository: languageGuidesRepository,
		zonesRepository:          zonesRepository,
		zoneGuidesRepository:     zoneGuidesRepository,
		availabilitiesRepository: availabilitiesRepository,
	}
}

func (s *Service) GetGuideDetailByID(guideID string) (GuidesDetail, error) {
	allGuides, err := s.usersService.GetAllGuides()
	if err != nil {
		logger.Error("GetGuideDetailByID: failed to get all guides for guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	var guide *domain.GuideUser
	for _, g := range allGuides {
		if g.ID == guideID {
			guide = &g
			break
		}
	}

	if guide == nil {
		err := fmt.Errorf("guide not found")
		logger.Error("GetGuideDetailByID: guide not found for guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	languageGuides, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guide.ID)
	if err != nil {
		logger.Error("GetGuideDetailByID: failed to load language guides for guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	languagesList := make([]languages.Language, 0, len(languageGuides))

	for _, lg := range languageGuides {
		language, err := s.languagesRepository.GetLanguageByID(lg.LanguageID)
		if err != nil {
			logger.Error("GetGuideDetailByID: failed to load language id=%s for guideID=%s: %v", lg.LanguageID, guideID, err)
			return GuidesDetail{}, err
		}
		if language == nil {
			err := fmt.Errorf("language not found")
			logger.Error("GetGuideDetailByID: language not found id=%s for guideID=%s: %v", lg.LanguageID, guideID, err)
			return GuidesDetail{}, err
		}

		languagesList = append(languagesList, *language)
	}

	zoneGuides, err := s.zoneGuidesRepository.GetZoneGuideByGuideID(guide.ID)
	if err != nil {
		logger.Error("GetGuideDetailByID: failed to load zone guides for guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	zonesList := make([]zones.Zone, 0, len(zoneGuides))

	for _, zg := range zoneGuides {
		zone, err := s.zonesRepository.GetZoneByID(zg.ZoneID)
		if err != nil {
			logger.Error("GetGuideDetailByID: failed to load zone id=%s for guideID=%s: %v", zg.ZoneID, guideID, err)
			return GuidesDetail{}, err
		}
		if zone == nil {
			err := fmt.Errorf("zone not found")
			logger.Error("GetGuideDetailByID: zone not found id=%s for guideID=%s: %v", zg.ZoneID, guideID, err)
			return GuidesDetail{}, err
		}

		zonesList = append(zonesList, *zone)
	}

	return GuidesDetail{
		Guide:     *guide,
		Languages: languagesList,
		Zones:     zonesList,
	}, nil
}

func (s *Service) GetAllGuidesDetail() ([]GuidesDetail, error) {
	allGuides, err := s.usersService.GetAllGuides()
	if err != nil {
		logger.Error("GetAllGuidesDetail: failed to get all guides: %v", err)
		return nil, err
	}

	var guidesDetails []GuidesDetail

	for _, guide := range allGuides {
		languageGuides, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guide.ID)
		if err != nil {
			logger.Error("GetAllGuidesDetail: failed to load language guides for guideID=%s: %v", guide.ID, err)
			return nil, err
		}

		languagesList := make([]languages.Language, 0, len(languageGuides))

		for _, lg := range languageGuides {
			language, err := s.languagesRepository.GetLanguageByID(lg.LanguageID)
			if err != nil {
				logger.Error("GetAllGuidesDetail: failed to load language id=%s for guideID=%s: %v", lg.LanguageID, guide.ID, err)
				return nil, err
			}
			if language == nil {
				err := fmt.Errorf("language not found")
				logger.Error("GetAllGuidesDetail: language not found id=%s for guideID=%s: %v", lg.LanguageID, guide.ID, err)
				return nil, err
			}

			languagesList = append(languagesList, *language)
		}

		zoneguides, err := s.zoneGuidesRepository.GetZoneGuideByGuideID(guide.ID)
		if err != nil {
			logger.Error("GetAllGuidesDetail: failed to load zone guides for guideID=%s: %v", guide.ID, err)
			return nil, err
		}

		logger.Info("GetAllGuidesDetail: guideID=%s has %d zone guides", guide.ID, len(zoneguides))

		zonesList := make([]zones.Zone, 0, len(zoneguides))

		for _, zg := range zoneguides {
			zone, err := s.zonesRepository.GetZoneByID(zg.ZoneID)
			if err != nil {
				logger.Error("GetAllGuidesDetail: failed to load zone id=%s for guideID=%s: %v", zg.ZoneID, guide.ID, err)
				return nil, err
			}
			if zone == nil {
				err := fmt.Errorf("zone not found")
				logger.Error("GetAllGuidesDetail: zone not found id=%s for guideID=%s: %v", zg.ZoneID, guide.ID, err)
				return nil, err
			}

			zonesList = append(zonesList, *zone)
		}

		guidesDetails = append(guidesDetails, GuidesDetail{
			Guide:     guide,
			Languages: languagesList,
			Zones:     zonesList,
		})
	}

	return guidesDetails, nil
}

func (s *Service) AddLanguageToGuide(guideID, languageCode string) error {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		logger.Error("AddLanguageToGuide: failed to get guide id=%s: %v", guideID, err)
		return err
	}
	if guide.ID == "" {
		err := fmt.Errorf("guide not found")
		logger.Error("AddLanguageToGuide: guide not found id=%s: %v", guideID, err)
		return err
	}

	// Check if the language exists
	language, err := s.languagesRepository.GetLanguageByCode(languageCode)
	if err != nil {
		logger.Error("AddLanguageToGuide: failed to get language code=%s for guideID=%s: %v", languageCode, guideID, err)
		return err
	}
	if language == nil {
		err := fmt.Errorf("language not found")
		logger.Error("AddLanguageToGuide: language not found code=%s for guideID=%s: %v", languageCode, guideID, err)
		return err
	}

	// Check if the language-guide association already exists
	exists, err := s.languageGuidesRepository.IsLanguageGuideExists(guideID, language.ID)
	if err != nil {
		logger.Error("AddLanguageToGuide: failed to check association guideID=%s languageID=%s: %v", guideID, language.ID, err)
		return err
	}
	if exists {
		err := fmt.Errorf("language already associated with the guide")
		logger.Error("AddLanguageToGuide: association already exists guideID=%s languageID=%s: %v", guideID, language.ID, err)
		return err
	}

	languageGuide := &languageguides.LanguageGuide{
		GuideID:    guideID,
		LanguageID: language.ID,
	}

	_, err = s.languageGuidesRepository.CreateLanguageGuide(languageGuide)
	if err != nil {
		logger.Error("AddLanguageToGuide: failed to create association guideID=%s languageID=%s: %v", guideID, language.ID, err)
		return err
	}

	return nil
}

func (s *Service) RemoveLanguageFromGuide(guideID, languageCode string) error {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		logger.Error("RemoveLanguageFromGuide: failed to get guide id=%s: %v", guideID, err)
		return err
	}
	if guide.ID == "" {
		err := fmt.Errorf("guide not found")
		logger.Error("RemoveLanguageFromGuide: guide not found id=%s: %v", guideID, err)
		return err
	}

	// Check if the language exists
	language, err := s.languagesRepository.GetLanguageByCode(languageCode)
	if err != nil {
		logger.Error("RemoveLanguageFromGuide: failed to get language code=%s for guideID=%s: %v", languageCode, guideID, err)
		return err
	}
	if language == nil {
		err := fmt.Errorf("language not found")
		logger.Error("RemoveLanguageFromGuide: language not found code=%s for guideID=%s: %v", languageCode, guideID, err)
		return err
	}

	// Check if the language-guide association exists
	exists, err := s.languageGuidesRepository.IsLanguageGuideExists(guideID, language.ID)
	if err != nil {
		logger.Error("RemoveLanguageFromGuide: failed to check association guideID=%s languageID=%s: %v", guideID, language.ID, err)
		return err
	}
	if !exists {
		err := fmt.Errorf("language is not associated with the guide")
		logger.Error("RemoveLanguageFromGuide: association not found guideID=%s languageCode=%s: %v", guideID, languageCode, err)
		return err
	}

	err = s.languageGuidesRepository.DeleteLanguageGuide(guideID, language.ID)
	if err != nil {
		logger.Error("RemoveLanguageFromGuide: failed to delete association guideID=%s languageID=%s: %v", guideID, language.ID, err)
		return err
	}

	return nil
}

func (s *Service) GetLanguagesByGuideID(guideID string) ([]languages.Language, error) {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		logger.Error("GetLanguagesByGuideID: failed to get guide id=%s: %v", guideID, err)
		return nil, err
	}
	if guide.ID == "" {
		err := fmt.Errorf("guide not found")
		logger.Error("GetLanguagesByGuideID: guide not found id=%s: %v", guideID, err)
		return nil, err
	}

	languageGuides, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guideID)
	if err != nil {
		logger.Error("GetLanguagesByGuideID: failed to load language guides for guideID=%s: %v", guideID, err)
		return nil, err
	}

	languagesList := make([]languages.Language, 0, len(languageGuides))

	for _, lg := range languageGuides {
		language, err := s.languagesRepository.GetLanguageByID(lg.LanguageID)
		if err != nil {
			logger.Error("GetLanguagesByGuideID: failed to load language id=%s for guideID=%s: %v", lg.LanguageID, guideID, err)
			return nil, err
		}
		if language != nil {
			languagesList = append(languagesList, *language)
		}
	}

	return languagesList, nil
}

func (s *Service) CreateGuide(guide CreateGuideRequest) (GuidesDetail, error) {
	// Create the user
	newUser := &domain.User{
		ID:        uuid.New().String(),
		Name:      guide.Name,
		Email:     guide.Email,
		Phone:     guide.Phone,
		Role:      "GUIDE",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	err := s.usersRepository.CreateUser(*newUser)
	if err != nil {
		logger.Error("CreateGuide: failed to create user name=%s email=%s phone=%s: %v", guide.Name, guide.Email, guide.Phone, err)
		return GuidesDetail{}, err
	}

	// Create the guide
	newGuide := &domain.Guide{
		UserID:         newUser.ID,
		MaxToursPerDay: guide.MaxToursPerDay,
		CreatedAt:      time.Now(),
	}
	if err := s.guidesRepository.CreateGuide(*newGuide); err != nil {
		logger.Error("CreateGuide: failed to create guide for userID=%s: %v", newUser.ID, err)
		return GuidesDetail{}, err
	}

	createdGuide, err := s.guidesRepository.GetGuideByUserID(newUser.ID)
	if err != nil {
		logger.Error("CreateGuide: failed to get created guide for userID=%s: %v", newUser.ID, err)
		return GuidesDetail{}, err
	}

	// Add languages to the guide
	for _, lang := range guide.Languages {
		if err := s.AddLanguageToGuide(createdGuide.ID, lang.Code); err != nil {
			logger.Error("CreateGuide: failed to add language code=%s to guideID=%s: %v", lang.Code, createdGuide.ID, err)
			return GuidesDetail{}, err
		}
	}

	// Add zones to the guide
	for _, zone := range guide.Zones {
		if err := s.AddZoneToGuide(createdGuide.ID, zone.ID); err != nil {
			logger.Error("CreateGuide: failed to add zone id=%s to guideID=%s: %v", zone.ID, createdGuide.ID, err)
			return GuidesDetail{}, err
		}
	}

	// Fetch the created guide details
	createdGuideDetail, err := s.GetGuideDetailByID(createdGuide.ID)
	if err != nil {
		logger.Error("CreateGuide: failed to fetch created guide detail guideID=%s: %v", createdGuide.ID, err)
		return GuidesDetail{}, err
	}

	return createdGuideDetail, nil
}

func (s *Service) UpdateGuide(guideID string, guide UpdateGuideRequest) (GuidesDetail, error) {
	// Fetch the existing guide
	existingGuide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		logger.Error("UpdateGuide: failed to get guide id=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}
	if existingGuide.ID == "" {
		logger.Error("UpdateGuide: guide not found id=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	// Update the guide details
	existingGuide.MaxToursPerDay = guide.MaxToursPerDay
	if err := s.guidesRepository.UpdateGuide(existingGuide); err != nil {
		logger.Error("UpdateGuide: failed to update guide id=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	// Normalize the requested languages to codes so we can diff them against the DB state.
	requestedLanguageCodes := make(map[string]struct{}, len(guide.Languages))
	for _, lang := range guide.Languages {
		languageCode, err := s.resolveLanguageCode(lang)
		if err != nil {
			logger.Error("UpdateGuide: failed to resolve requested language for guideID=%s: %v", guideID, err)
			return GuidesDetail{}, err
		}
		requestedLanguageCodes[languageCode] = struct{}{}
	}

	// Remove languages that are no longer present in the update request.
	existingLanguages, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guideID)
	if err != nil {
		logger.Error("UpdateGuide: failed to load existing languages for guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	existingLanguageCodes := make(map[string]struct{}, len(existingLanguages))

	for _, existingLanguage := range existingLanguages {
		language, err := s.languagesRepository.GetLanguageByID(existingLanguage.LanguageID)
		if err != nil {
			logger.Error("UpdateGuide: failed to resolve existing language id=%s for guideID=%s: %v", existingLanguage.LanguageID, guideID, err)
			return GuidesDetail{}, err
		}
		if language == nil {
			logger.Error("UpdateGuide: existing language not found id=%s for guideID=%s: %v", existingLanguage.LanguageID, guideID, err)
			return GuidesDetail{}, err
		}

		existingLanguageCodes[language.Code] = struct{}{}

		if _, keep := requestedLanguageCodes[language.Code]; keep {
			continue
		}

		if err := s.RemoveLanguageFromGuide(guideID, language.Code); err != nil {
			logger.Error("UpdateGuide: failed to remove language code=%s from guideID=%s: %v", language.Code, guideID, err)
			return GuidesDetail{}, err
		}
	}

	// Add languages that are present in the request but not yet linked.
	for _, lang := range guide.Languages {
		languageCode, err := s.resolveLanguageCode(lang)
		if err != nil {
			logger.Error("UpdateGuide: failed to resolve requested language for guideID=%s: %v", guideID, err)
			return GuidesDetail{}, err
		}

		if _, exists := existingLanguageCodes[languageCode]; exists {
			continue
		}
		if err := s.AddLanguageToGuide(guideID, languageCode); err != nil {
			logger.Error("UpdateGuide: failed to add language code=%s to guideID=%s: %v", languageCode, guideID, err)
			return GuidesDetail{}, err
		}
	}

	// Remove zones that are no longer present in the update request.
	existingZones, err := s.zoneGuidesRepository.GetZoneGuideByGuideID(guideID)
	if err != nil {
		logger.Error("UpdateGuide: failed to load existing zones for guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	existingZoneIDs := make(map[string]struct{}, len(existingZones))

	for _, existingZone := range existingZones {
		existingZoneIDs[existingZone.ZoneID] = struct{}{}

		// Check if the existing zone is in the requested zones
		found := false
		for _, requestedZone := range guide.Zones {
			if existingZone.ZoneID == requestedZone.ID {
				found = true
				break
			}
		}

		if !found {
			if err := s.zoneGuidesRepository.DeleteZoneGuide(guideID, existingZone.ZoneID); err != nil {
				logger.Error("UpdateGuide: failed to remove zone id=%s from guideID=%s: %v", existingZone.ZoneID, guideID, err)
				return GuidesDetail{}, err
			}
		}
	}

	// Add zones that are present in the request but not yet linked.
	for _, requestedZone := range guide.Zones {
		if _, exists := existingZoneIDs[requestedZone.ID]; exists {
			continue
		}
		if err := s.AddZoneToGuide(guideID, requestedZone.ID); err != nil {
			logger.Error("UpdateGuide: failed to add zone id=%s to guideID=%s: %v", requestedZone.ID, guideID, err)
			return GuidesDetail{}, err
		}
	}

	// Fetch the updated guide details
	updatedGuideDetail, err := s.GetGuideDetailByID(guideID)
	if err != nil {
		logger.Error("UpdateGuide: failed to fetch updated guide detail guideID=%s: %v", guideID, err)
		return GuidesDetail{}, err
	}

	return updatedGuideDetail, nil
}

func (s *Service) resolveLanguageCode(language languages.Language) (string, error) {
	if language.Code != "" {
		return language.Code, nil
	}

	if language.ID == "" {
		return "", fmt.Errorf("language code or id is required")
	}

	resolvedLanguage, err := s.languagesRepository.GetLanguageByID(language.ID)
	if err != nil {
		return "", err
	}
	if resolvedLanguage == nil {
		return "", fmt.Errorf("language not found")
	}

	return resolvedLanguage.Code, nil
}

func (s *Service) AddZoneToGuide(guideID, zoneID string) error {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		logger.Error("AddZoneToGuide: failed to get guide id=%s: %v", guideID, err)
		return err
	}
	if guide.ID == "" {
		err := fmt.Errorf("guide not found")
		logger.Error("AddZoneToGuide: guide not found id=%s: %v", guideID, err)
		return err
	}

	// Check if the zone exists
	zone, err := s.zonesRepository.GetZoneByID(zoneID)
	if err != nil {
		logger.Error("AddZoneToGuide: failed to get zone id=%s for guideID=%s: %v", zoneID, guideID, err)
		return err
	}

	if zone == nil {
		err := fmt.Errorf("zone not found")
		logger.Error("AddZoneToGuide: zone not found id=%s for guideID=%s: %v", zoneID, guideID, err)
		return err
	}

	// Check if the zone-guide association already exists
	exists, err := s.zoneGuidesRepository.IsZoneGuideExists(guideID, zone.ID)
	if err != nil {
		logger.Error("AddZoneToGuide: failed to check association guideID=%s zoneID=%s: %v", guideID, zone.ID, err)
		return err
	}
	if exists {
		err := fmt.Errorf("zone already associated with the guide")
		logger.Error("AddZoneToGuide: association already exists guideID=%s zoneID=%s: %v", guideID, zone.ID, err)
		return err
	}

	zoneGuide := &zoneguides.ZoneGuide{
		GuideID: guideID,
		ZoneID:  zone.ID,
	}

	_, err = s.zoneGuidesRepository.CreateZoneGuide(zoneGuide)
	if err != nil {
		logger.Error("AddZoneToGuide: failed to create association guideID=%s zoneID=%s: %v", guideID, zone.ID, err)
		return err
	}

	return nil
}

// Availibilities

func (s *Service) CreateAvailability(availability availabilities.CreateAvailabilityRequest) (*availabilities.Availability, error) {
	startDate, endDate, err := validateAvailabilityDates(availability.StartDate, availability.EndDate)
	if err != nil {
		return nil, err
	}

	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(availability.GuideID)
	if err != nil {
		logger.Error("CreateAvailability: failed to get guide id=%s: %v", availability.GuideID, err)
		return nil, err
	}
	if guide.ID == "" {
		err := fmt.Errorf("guide not found")
		logger.Error("CreateAvailability: guide not found id=%s: %v", availability.GuideID, err)
		return nil, err
	}
	conflict, err := s.availabilitiesRepository.HasConflictingAvailability(availability.GuideID, startDate, endDate, nil)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, appErrors.ErrAvailabilityConflict
	}

	// Create the availability
	newAvailability := &availabilities.Availability{
		ID:        uuid.New(),
		GuideID:   availability.GuideID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    availability.Reason,
	}

	return s.availabilitiesRepository.CreateAvailability(newAvailability)
}

func (s *Service) GetAvailabilitiesByGuideID(guideID string) ([]*availabilities.Availability, error) {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		logger.Error("GetAvailabilitiesByGuideID: failed to get guide id=%s: %v", guideID, err)
		return nil, err
	}
	if guide.ID == "" {
		err := fmt.Errorf("guide not found")
		logger.Error("GetAvailabilitiesByGuideID: guide not found id=%s: %v", guideID, err)
		return nil, err
	}

	return s.availabilitiesRepository.GetAvailabilitiesByGuideID(guideID)
}

func (s *Service) DeleteAvailability(availabilityID string) error {
	availability, err := s.availabilitiesRepository.GetAvailabilityByID(availabilityID)
	if err != nil {
		logger.Error("deleteAvailability: failed to get availability id=%s: %v", availabilityID, err)
		return err
	}
	if availability == nil {
		err := fmt.Errorf("availability not found")
		logger.Error("deleteAvailability: availability not found id=%s: %v", availabilityID, err)
		return err
	}

	return s.availabilitiesRepository.DeleteAvailability(availabilityID)
}

func (s *Service) UpdateAvailability(guideID, availabilityID string, updatedAvailability availabilities.CreateAvailabilityRequest) (*availabilities.Availability, error) {
	startDate, endDate, err := validateAvailabilityDates(updatedAvailability.StartDate, updatedAvailability.EndDate)
	if err != nil {
		return nil, err
	}

	availability, err := s.availabilitiesRepository.GetAvailabilityByID(availabilityID)
	if err != nil {
		logger.Error("updateAvailability: failed to get availability id=%s: %v", availabilityID, err)
		return nil, err
	}
	if availability == nil {
		err := appErrors.ErrAvailabilityNotFound
		logger.Error("updateAvailability: availability not found id=%s: %v", availabilityID, err)
		return nil, err
	}
	if availability.GuideID != guideID {
		return nil, appErrors.ErrAvailabilityNotFound
	}
	conflict, err := s.availabilitiesRepository.HasConflictingAvailability(guideID, startDate, endDate, &availability.ID)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, appErrors.ErrAvailabilityConflict
	}

	availability.StartDate = startDate
	availability.EndDate = endDate
	availability.Reason = updatedAvailability.Reason

	if err := s.availabilitiesRepository.UpdateAvailability(availability); err != nil {
		logger.Error("updateAvailability: failed to update availability id=%s: %v", availabilityID, err)
		return nil, err
	}

	return availability, nil
}

func validateAvailabilityDates(start, end *time.Time) (*time.Time, *time.Time, error) {
	if start == nil || end == nil || end.Before(*start) {
		return nil, nil, appErrors.ErrInvalidAvailabilityDate
	}
	return start, end, nil
}
