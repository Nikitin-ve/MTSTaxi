package taxi

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"trip_service/internal/app/trip_service/api/requests"
)

type PostgreSQL struct {
	DB *pgxpool.Pool
}

func (p *PostgreSQL) AddNewTrip(trip requests.ForDB) {
	pgxscan.Select(context.Background(), p.DB, &trip,
		`SELECT trip_id, source, type, datacontenttype, timer,
       offer_id, amount, currency, status, lat_from, lng_from, lat_to, lng_to FROM trip`)
}
