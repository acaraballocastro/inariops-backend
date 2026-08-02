package guidesapp

import (
	"fmt"

	"inariops/internal/domain"
	guidedomain "inariops/internal/modules/guides"
	languageguides "inariops/internal/modules/guides/language_guides"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/users"
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
		return GuidesDetail{}, fmt.Errorf("guide not found")
	}

	languageGuides, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guide.ID)
	if err != nil {
		return GuidesDetail{}, err
	}

	languagesList := make([]languages.Language, 0, len(languageGuides))

	for _, lg := range languageGuides {
		language, err := s.languagesRepository.GetLanguageByID(lg.LanguageID)
		if err != nil {
			return GuidesDetail{}, err
		}
		if language == nil {
			return GuidesDetail{}, fmt.Errorf("language not found")
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
		return nil, err
	}

	var guidesDetails []GuidesDetail

	for _, guide := range allGuides {
		languageGuides, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guide.ID)
		if err != nil {
			return nil, err
		}

		languagesList := make([]languages.Language, 0, len(languageGuides))

		for _, lg := range languageGuides {
			language, err := s.languagesRepository.GetLanguageByID(lg.LanguageID)
			if err != nil {
				return nil, err
			}
			if language == nil {
				return nil, fmt.Errorf("language not found")
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
		return err
	}
	if guide.ID == "" {
		return fmt.Errorf("guide not found")
	}

	// Check if the language exists
	language, err := s.languagesRepository.GetLanguageByCode(languageCode)
	if err != nil {
		return err
	}
	if language == nil {
		return fmt.Errorf("language not found")
	}

	// Check if the language-guide association already exists
	exists, err := s.languageGuidesRepository.IsLanguageGuideExists(guideID, language.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("language already associated with the guide")
	}

	languageGuide := &languageguides.LanguageGuide{
		GuideID:    guideID,
		LanguageID: language.ID,
	}

	_, err = s.languageGuidesRepository.CreateLanguageGuide(languageGuide)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveLanguageFromGuide(guideID, languageCode string) error {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		return err
	}
	if guide.ID == "" {
		return fmt.Errorf("guide not found")
	}

	// Check if the language exists
	language, err := s.languagesRepository.GetLanguageByCode(languageCode)
	if err != nil {
		return err
	}
	if language == nil {
		return fmt.Errorf("language not found")
	}

	// Check if the language-guide association exists
	exists, err := s.languageGuidesRepository.IsLanguageGuideExists(guideID, languageCode)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("language is not associated with the guide")
	}

	err = s.languageGuidesRepository.DeleteLanguageGuide(guideID, language.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetLanguagesByGuideID(guideID string) ([]languages.Language, error) {
	// Check if the guide exists
	guide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		return nil, err
	}
	if guide.ID == "" {
		return nil, fmt.Errorf("guide not found")
	}

	languageGuides, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guideID)
	if err != nil {
		return nil, err
	}

	languagesList := make([]languages.Language, 0, len(languageGuides))

	for _, lg := range languageGuides {
		language, err := s.languagesRepository.GetLanguageByID(lg.LanguageID)
		if err != nil {
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
	user, err := s.usersService.CreateUser(guide.Name, guide.Email, guide.Phone, "GUIDE")

	if err != nil {
		return GuidesDetail{}, err
	}

	// Create the guide
	newGuide := &domain.Guide{
		UserID:         user.User.ID,
		MaxToursPerDay: guide.MaxToursPerDay,
	}
	if err := s.guidesRepository.CreateGuide(*newGuide); err != nil {
		return GuidesDetail{}, err
	}

	// Add languages to the guide
	for _, lang := range guide.Languages {
		if err := s.AddLanguageToGuide(newGuide.ID, lang.Code); err != nil {
			return GuidesDetail{}, err
		}
	}

	// Fetch the created guide details
	createdGuideDetail, err := s.GetGuideDetailByID(newGuide.ID)
	if err != nil {
		return GuidesDetail{}, err
	}

	return createdGuideDetail, nil
}

func (s *Service) UpdateGuide(guideID string, guide UpdateGuideRequest) (GuidesDetail, error) {
	// Fetch the existing guide
	existingGuide, err := s.guidesRepository.GetGuideByID(guideID)
	if err != nil {
		return GuidesDetail{}, err
	}
	if existingGuide.ID == "" {
		return GuidesDetail{}, fmt.Errorf("guide not found")
	}

	// Update the guide details
	existingGuide.MaxToursPerDay = guide.MaxToursPerDay
	if err := s.guidesRepository.UpdateGuide(existingGuide); err != nil {
		return GuidesDetail{}, err
	}

	// Remove all existing languages associated with the guide
	existingLanguages, err := s.languageGuidesRepository.GetLanguageGuideByGuideID(guideID)
	if err != nil {
		return GuidesDetail{}, err
	}
	for _, lang := range existingLanguages {
		if err := s.RemoveLanguageFromGuide(guideID, lang.LanguageID); err != nil {
			return GuidesDetail{}, err
		}
	}

	// Add the new languages to the guide
	for _, lang := range guide.Languages {
		if err := s.AddLanguageToGuide(guideID, lang.Code); err != nil {
			return GuidesDetail{}, err
		}
	}

	// Fetch the updated guide details
	updatedGuideDetail, err := s.GetGuideDetailByID(guideID)
	if err != nil {
		return GuidesDetail{}, err
	}

	return updatedGuideDetail, nil
}
