package apps

import (
	"database/sql"
	"mobile-connect/internal/apps/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	app entity.App,
) error {

	query := `
	INSERT INTO apps (
		id,
		filename,
		original_name,
		package_name,
		version_name,
		version_code,
		min_sdk,
		target_sdk,
		size,
		sha256
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		app.ID,
		app.FileName,
		app.OriginalName,
		app.PackageName,
		app.VersionName,
		app.VersionCode,
		app.MinSDK,
		app.TargetSDK,
		app.Size,
		app.SHA256,
	)

	return err
}

func (r *Repository) FindAll() (
	[]entity.App,
	error,
) {

	rows, err := r.db.Query(`
		SELECT
			id,
			filename,
			original_name,
			package_name,
			version_name,
			version_code,
			min_sdk,
			target_sdk,
			size,
			sha256,
			created_at
		FROM apps
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var apps []entity.App

	for rows.Next() {

		var app entity.App

		err := rows.Scan(
			&app.ID,
			&app.FileName,
			&app.OriginalName,
			&app.PackageName,
			&app.VersionName,
			&app.VersionCode,
			&app.MinSDK,
			&app.TargetSDK,
			&app.Size,
			&app.SHA256,
			&app.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		apps = append(
			apps,
			app,
		)
	}

	return apps, nil
}

func (r *Repository) FindByID(
	id string,
) (*entity.App, error) {

	var app entity.App

	err := r.db.QueryRow(`
		SELECT
			id,
			filename,
			original_name,
			package_name,
			version_name,
			version_code,
			min_sdk,
			target_sdk,
			size,
			sha256,
			created_at
		FROM apps
		WHERE id = ?
	`, id).Scan(
		&app.ID,
		&app.FileName,
		&app.OriginalName,
		&app.PackageName,
		&app.VersionName,
		&app.VersionCode,
		&app.MinSDK,
		&app.TargetSDK,
		&app.Size,
		&app.SHA256,
		&app.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *Repository) Delete(
	id string,
) error {

	_, err := r.db.Exec(
		"DELETE FROM apps WHERE id = ?",
		id,
	)

	return err
}
