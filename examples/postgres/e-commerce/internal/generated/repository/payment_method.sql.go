// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: payment_method.sql

package db

import (
	"context"
)

const getPaymentMethods = `-- name: GetPaymentMethods :many
SELECT payment_method_id, name, description
FROM payment_method
`

func (q *Queries) GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error) {
	rows, err := q.db.Query(ctx, getPaymentMethods)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PaymentMethod{}
	for rows.Next() {
		var i PaymentMethod
		if err := rows.Scan(&i.PaymentMethodID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
