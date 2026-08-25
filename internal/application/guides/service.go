package guidesapp

import (
	"fmt"
	"time"

	"inariops/internal/domain"
	guidedomain "inariops/internal/modules/guides"
	languageguides "inariops/internal/modules/guides/language_guides"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/users"
	"inariops/internal/shared/logger"

	"github.com/google/uuid"
)

type Service struct {
	usersService             *users.Service
	usersRepository          *users.Repository
	guidesRepository         *guidedomain.Repository
	languagesRepository      *languages.Repository
	languageGuidesRepository *languageguides.Repository
}

func NewService(
	usersService *users.Service,
	usersRepository *users.Repository,
	guidesRepository *guidedomain.Repository,
	languagesRepository *languages.Repository,
	languageGuidesRepository *languageguides.Repository,
) *Service {
	return &Service{
		usersService:             usersService,
		usersRepository:          usersRepository,
		guidesRepository:         guidesRepository,
		languagesRepository:      languagesRepository,
		languageGuidesRepository: languageGuidesRepository,
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

	return GuidesDetail{
		Guide:     *guide,
		Languages: languagesList,
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

		guidesDetails = append(guidesDetails, GuidesDetail{
			Guide:     guide,
			Languages: languagesList,
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

func (s *Service) xCreateGuide(guide CreateGuideRequest) (GuidesDetail, error) {
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
