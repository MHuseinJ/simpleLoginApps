package repository

import (
	"context"
	"fmt"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
