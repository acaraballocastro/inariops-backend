package itinerary

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateItineraryItem(item ItineraryItem) error {
	_, err := r.db.Exec(`
		INSERT INTO itinerary_items (
			id, tour_day_id, order_index, start_time, end_time, type, place_id, activity_id, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		item.ID,
		item.TourDayID,
		item.OrderIndex,
		item.StartTime,
		item.EndTime,
		item.Type,
		item.PlaceID,
		item.ActivityID,
		item.Notes,
	)

	return err
}

func (r *Repository) GetItineraryItemsByTourDayID(tourDayID string) ([]ItineraryItem, error) {
	rows, err := r.db.Query(`
		SELECT id, tour_day_id, order_index, start_time, end_time, type, place_id, activity_id, notes
		FROM itinerary_items
		WHERE tour_day_id = $1
		ORDER BY order_index ASC
	`, tourDayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ItineraryItem
	for rows.Next() {
		var item ItineraryItem
		err := rows.Scan(
			&item.ID,
			&item.TourDayID,
			&item.OrderIndex,
			&item.StartTime,
			&item.EndTime,
			&item.Type,
			&item.PlaceID,
			&item.ActivityID,
			&item.Notes,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) GetItineraryItemByID(itemID string) (ItineraryItem, error) {
	item := ItineraryItem{}

	err := r.db.QueryRow(`
		SELECT id, tour_day_id, order_index, start_time, end_time, type, place_id, activity_id, notes
		FROM itinerary_items
		WHERE id = $1
	`, itemID).Scan(
		&item.ID,
		&item.TourDayID,
		&item.OrderIndex,
		&item.StartTime,
		&item.EndTime,
		&item.Type,
		&item.PlaceID,
		&item.ActivityID,
		&item.Notes,
	)

	if err != nil {
		return ItineraryItem{}, err
	}

	return item, nil
}

func (r *Repository) UpdateItineraryItem(item ItineraryItem) error {
	_, err := r.db.Exec(`
		UPDATE itinerary_items SET
			order_index = $1,
			start_time = $2,
			end_time = $3,
			type = $4,
			place_id = $5,
			activity_id = $6,
			notes = $7
		WHERE id = $8
	`,
		item.OrderIndex,
		item.StartTime,
		item.EndTime,
		item.Type,
		item.PlaceID,
		item.ActivityID,
		item.Notes,
		item.ID,
	)

	return err
}

func (r *Repository) DeleteItineraryItem(itemID string) error {
	_, err := r.db.Exec(`
		DELETE FROM itinerary_items WHERE id = $1
	`, itemID)
	return err
}

func (r *Repository) DeleteItineraryItemsByTourDayID(tourDayID string) error {
	_, err := r.db.Exec(`
		DELETE FROM itinerary_items WHERE tour_day_id = $1
	`, tourDayID)
	return err
}
