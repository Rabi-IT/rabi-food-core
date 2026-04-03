package database

import (
	"sync"

	"gorm.io/gorm"
)

// Paginate performs pagination on the given GORM query.
func Paginate(query *gorm.DB, count *int64, data any, paginate PaginateInput) error {
	var wg sync.WaitGroup
	countErr := error(nil)
	countSession := query.Session(&gorm.Session{})
	wg.Go(func() {
		result := countSession.Count(count)
		countErr = result.Error
	})

	paginateErr := error(nil)
	paginateSession := query.Session(&gorm.Session{})
	wg.Go(func() {
		result := paginateSession.Limit(paginate.PageSize).Offset(paginate.Offset()).Find(data)
		paginateErr = result.Error
	})

	wg.Wait()
	if countErr != nil {
		return countErr
	}

	if paginateErr != nil {
		return paginateErr
	}

	return nil
}
