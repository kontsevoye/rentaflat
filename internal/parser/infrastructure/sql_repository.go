package infrastructure

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/kontsevoye/rentaflat/internal/parser/domain"
	"time"
)

func NewSqlRepository(connection *sqlx.DB, flatFactory domain.FlatFactory) *SqlRepository {
	return &SqlRepository{connection, flatFactory}
}

type SqlRepository struct {
	connection  *sqlx.DB
	flatFactory domain.FlatFactory
}

type rawFlat struct {
	Id          string    `db:"id"`
	ServiceId   string    `db:"service_id"`
	Url         string    `db:"url"`
	PhotoUrls   []byte    `db:"photo_urls"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Area        uint      `db:"area"`
	Rooms       uint      `db:"rooms"`
	Floor       uint      `db:"floor"`
	Price       uint      `db:"price"`
	ContactName string    `db:"contact_name"`
	Phone       string    `db:"phone"`
	IsAgency    bool      `db:"is_agency"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
}

func (f rawFlat) isInitialized() bool {
	return f.Id != ""
}

func (f rawFlat) unmarshalPhotoUrls() ([]string, error) {
	var photoUrlsRaw []string
	err := json.Unmarshal(f.PhotoUrls, &photoUrlsRaw)

	return photoUrlsRaw, err
}

func (s *SqlRepository) Add(flat domain.Flat) error {
	photos, err := json.Marshal(flat.PhotoUrlsAsStrings())
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"id":           flat.Id(),
		"service_id":   flat.ServiceId(),
		"url":          flat.Url().String(),
		"photo_urls":   photos,
		"title":        flat.Title(),
		"description":  flat.Description(),
		"area":         flat.Area(),
		"rooms":        flat.Rooms(),
		"floor":        flat.Floor(),
		"price":        flat.Price(),
		"contact_name": flat.ContactName(),
		"phone":        flat.Phone(),
		"is_agency":    flat.IsAgency(),
		"published_at": flat.PublishedAt(),
		"created_at":   flat.CreatedAt(),
	}
	_, err = s.connection.NamedExec(
		`INSERT INTO public.flats (
                   id,
                   service_id,
                   url,
                   photo_urls,
                   title,
                   description,
                   area,
                   rooms,
                   floor,
                   price,
                   contact_name,
                   phone,
                   is_agency,
                   published_at,
                   created_at
			   ) VALUES (
				   :id,
                   :service_id,
                   :url,
                   :photo_urls,
                   :title,
                   :description,
                   :area,
                   :rooms,
                   :floor,
                   :price,
                   :contact_name,
                   :phone,
                   :is_agency,
                   :published_at,
                   :created_at
			   )`,
		data,
	)

	return err
}

func (s *SqlRepository) FindByUrl(url string) (domain.Flat, error) {
	f := rawFlat{}
	err := s.connection.Get(&f, "SELECT * FROM public.flats WHERE url = $1", url)
	if err != nil {
		return domain.Flat{}, err
	}

	if f.isInitialized() {
		urls, err := f.unmarshalPhotoUrls()
		if err != nil {
			return domain.Flat{}, err
		}

		return s.flatFactory.LoadFlat(
			f.Id,
			f.ServiceId,
			f.Url,
			urls,
			f.Title,
			f.Description,
			f.Area,
			f.Rooms,
			f.Floor,
			f.Price,
			f.ContactName,
			f.Phone,
			f.IsAgency,
			f.PublishedAt,
			f.CreatedAt,
		)
	} else {
		return domain.Flat{}, domain.ErrFlatNotFoundWithId(url)
	}
}

func (s *SqlRepository) Has(url string) (bool, error) {
	var has bool
	err := s.connection.Get(&has, "SELECT COUNT(*) > 0 as has FROM public.flats WHERE url = $1", url)

	return has, err
}

func (s *SqlRepository) FindLatest() (domain.Flat, error) {
	f := rawFlat{}
	err := s.connection.Get(&f, "SELECT * FROM public.flats ORDER BY published_at DESC LIMIT 1")
	if err != nil {
		return domain.Flat{}, err
	}

	if f.isInitialized() {
		urls, err := f.unmarshalPhotoUrls()
		if err != nil {
			return domain.Flat{}, err
		}

		return s.flatFactory.LoadFlat(
			f.Id,
			f.ServiceId,
			f.Url,
			urls,
			f.Title,
			f.Description,
			f.Area,
			f.Rooms,
			f.Floor,
			f.Price,
			f.ContactName,
			f.Phone,
			f.IsAgency,
			f.PublishedAt,
			f.CreatedAt,
		)
	} else {
		return domain.Flat{}, domain.ErrFlatNotFound
	}
}
