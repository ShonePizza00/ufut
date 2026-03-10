package sqliteRepoInventory

import "context"

func (r *SQLiteRepo) ReserveItem(ctx context.Context, itemsIDs []string, quantities []int) ([]bool, error) {

}

func (r *SQLiteRepo) CancelItemReservation(ctx context.Context, itemsIDs []string) error {

}
