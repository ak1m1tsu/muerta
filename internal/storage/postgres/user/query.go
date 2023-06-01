package user

const (
	countUsers = `
		SELECT COUNT(*)
		FROM users 
		WHERE deleted_at IS NULL 
			AND name ILIKE $1
	`
	createShelfLife = `
		WITH inserted AS (
			INSERT INTO shelf_lives (id_user, id_product, id_storage, id_measure, quantity, purchase_date, end_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, id_product, id_storage, id_measure, quantity, purchase_date, end_date
		)
		SELECT 
			i.id, p.name, s.name, m.name
		FROM inserted i
		JOIN products p ON i.id_product = p.id
		JOIN storages s ON i.id_storage = s.id
		JOIN measures m ON i.id_measure = m.id
		WHERE p.deleted_at IS NULL AND 
			s.deleted_at IS NULL
		LIMIT 1
	`
	deleteShelfLife = `
		UPDATE shelf_lives
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id_user = $1 AND id = $2
	`
	findShelfLife = `
		SELECT 
			sl.id, sl.id_product, sl.id_storage, sl.id_measure, 
			sl.quantity, sl.purchase_date, sl.end_date,
			p.name, s.name, m.name
		FROM shelf_lives sl
		JOIN products p ON sl.id_product = p.id
		JOIN storages s ON sl.id_storage = s.id
		JOIN measures m ON sl.id_measure = m.id
		WHERE sl.id_user = $1 AND 
			sl.id = $2 AND
			sl.deleted_at IS NULL
		LIMIT 1
	`
	findShelfLives = `
		SELECT 
			sl.id, sl.id_product, sl.id_storage, sl.id_measure, 
			sl.quantity, sl.purchase_date, sl.end_date,
			p.name, s.name, m.name
		FROM shelf_lives sl
		JOIN products p ON sl.id_product = p.id
		JOIN storages s ON sl.id_storage = s.id
		JOIN measures m ON sl.id_measure = m.id
		WHERE sl.id_user = $1 AND 
			sl.deleted_at IS NULL
		ORDER BY sl.end_date DESC
	`
	restoreShelfLife = `
		WITH updated AS (
			UPDATE shelf_lives
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id_user = $1 AND id = $2
			RETURNING id, id_product, id_storage, id_measure, quantity, purchase_date, end_date
		)
		SELECT 
			u.id, u.id_product, u.id_storage, u.id_measure, 
			u.quantity, u.purchase_date, u.end_date,
			p.name, s.name, m.name
		FROM updated u
		JOIN products p ON u.id_product = p.id
		JOIN storages s ON u.id_storage = s.id
		JOIN measures m ON u.id_measure = m.id
		WHERE p.deleted_at IS NULL AND
			s.deleted_at IS NULL AND
		LIMIT 1
	`
	updateShelfLife = `
		WITH updated AS (
			UPDATE shelf_lives
			SET id_product = $3,
				id_storage = $4,
				id_measure = $5,
				quantity = $6,
				purchase_date = $7,
				end_date = $8,
				updated_at = NOW()
			WHERE id_user = $1 AND id = $2
			RETURNING id_product, id_storage, id_measure, quantity, purchase_date, end_date
		)
		SELECT 
			u.id_product, u.id_storage, u.id_measure, 
			u.quantity, u.purchase_date, u.end_date,
			p.name, s.name, m.name
		FROM updated u
		JOIN products p ON u.id_product = p.id
		JOIN storages s ON u.id_storage = s.id
		JOIN measures m ON u.id_measure = m.id
		WHERE p.deleted_at IS NULL AND
			s.deleted_at IS NULL AND
		LIMIT 1
	`
	addVault = `
		WITH inserted AS (
			INSERT INTO users_storages (id_user, id_storage)
			VALUES ($1, $2)
			RETURNING id_user, id_storage
		)
		SELECT s.id, s.name, s.temperature, s.humidity, st.id, st.name
		FROM storages s
		JOIN storages_types st ON s.id_type = st.id
		JOIN inserted i ON i.id_storage = s.id
		WHERE i.id_user = $1 AND 
			s.id = i.id_storage AND
			s.deleted_at IS NULL
		LIMIT 1
	`
	removeVault = `
		DELETE FROM users_storages
		WHERE id_user = $1 AND id_storage = $2
	`
	findVaults = `
		SELECT s.id, s.name, st.name, s.temperature, s.humidity
		FROM storages s
		JOIN storages_types st ON s.id_type = st.id
		JOIN users_storages us ON us.id_storage = s.id
		WHERE us.id_user = $1
	`
	findRoles = `
		SELECT r.id, r.name
		FROM roles r
		JOIN users_roles ur ON ur.id_role = r.id
		WHERE ur.id_user = $1 AND r.deleted_at IS NULL
	`
	restoreUser = `
		UPDATE users
		SET deleted_at = NULL
		SET updated_at = NOW()
		WHERE id = $1
	`
	deleteUser = `
		UPDATE users
		SET deleted_at = NOW()
		SET updated_at = NOW()
		WHERE id = $1
	`
	updateUser = `
		UPDATE users
		SET name = $1,
			updated_at = NOW()
		WHERE id = $2
	`
	createUser = `
		INSERT INTO users 
			(name, salt)
		VALUES
			($1, $2)
		RETURNING id
	`
	createPassword = `
		INSERT INTO passwords (passhash)
		VALUES ($1)
	`
	findPassword = `
		SELECT passhash
		FROM passwords
		WHERE passhash = $1
		LIMIT 1
	`
	findUsers = `
		SELECT id, name, created_at
		FROM users
		WHERE name LIKE $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3
	`
	findUserByName = `
		SELECT id, name, salt, created_at
		FROM users
		WHERE name = $1
		LIMIT 1
	`
	updateSetting = `
		UPDATE users_settings
		SET value = $2
		WHERE id_user = $1 AND id_setting = $3
	`
	findSetting = `
		SELECT s.name, us.value, sc.name FROM settings s
		JOIN users_settings us ON s.id = us.id_setting
		JOIN settings_categories sc ON s.id_category = sc.id
		WHERE us.id_user = $1
		LIMIT 1
	`
	findSettings = `
		SELECT s.id, s.name, us.value, sc.name
		FROM settings s
		JOIN users_settings us ON s.id = us.id_setting
		JOIN settings_categories sc ON s.id_category = sc.id
		WHERE us.id_user = $1
	`
)
