package catalog

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite catalog database.
type Store struct {
	db *sql.DB
}

// OpenStore opens (creating if needed) the catalog database at path.
func OpenStore(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create store dir: %w", err)
	}
	dsn := path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(10000)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-20000)&_pragma=temp_store(MEMORY)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	db.SetMaxOpenConns(4)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping store: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS folders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL UNIQUE,
			display_name TEXT NOT NULL,
			is_default INTEGER NOT NULL DEFAULT 0,
			last_scanned INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS repos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			folder_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			path TEXT NOT NULL UNIQUE,
			category TEXT NOT NULL DEFAULT 'Uncategorized',
			readme_path TEXT NOT NULL DEFAULT '',
			readme_text TEXT NOT NULL DEFAULT '',
			UNIQUE(folder_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS scripts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id INTEGER NOT NULL,
			rel_path TEXT NOT NULL,
			name TEXT NOT NULL,
			ext TEXT NOT NULL DEFAULT '',
			lang TEXT NOT NULL DEFAULT 'text',
			size INTEGER NOT NULL DEFAULT 0,
			mtime INTEGER NOT NULL DEFAULT 0,
			abs_path TEXT NOT NULL,
			is_script INTEGER NOT NULL DEFAULT 1,
			snippet TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE INDEX IF NOT EXISTS idx_scripts_repo ON scripts(repo_id)`,
		`CREATE INDEX IF NOT EXISTS idx_scripts_abs ON scripts(abs_path)`,
		`CREATE INDEX IF NOT EXISTS idx_repos_folder ON repos(folder_id)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS scripts_fts USING fts5(content, tokenize='unicode61')`,
		`CREATE TABLE IF NOT EXISTS custom_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS custom_repos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			category TEXT NOT NULL DEFAULT 'Other',
			summary TEXT NOT NULL DEFAULT '',
			UNIQUE(name, category)
		)`,
		`CREATE TABLE IF NOT EXISTS custom_software (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			version TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT 'Other',
			download TEXT NOT NULL,
			notes TEXT NOT NULL DEFAULT '',
			winget_id TEXT NOT NULL DEFAULT '',
			UNIQUE(name)
		)`,
	}
	for _, q := range stmts {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}

// --- folders ---

func (s *Store) ListFolders() ([]*Folder, error) {
	rows, err := s.db.Query(`
		SELECT f.id, f.path, f.display_name, f.is_default, f.last_scanned,
			(SELECT count(*) FROM repos r WHERE r.folder_id = f.id),
			(SELECT count(*) FROM scripts sc JOIN repos r ON sc.repo_id = r.id WHERE r.folder_id = f.id)
		FROM folders f ORDER BY f.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Folder
	for rows.Next() {
		f := &Folder{}
		var ts int64
		if err := rows.Scan(&f.ID, &f.Path, &f.DisplayName, &f.IsDefault, &ts, &f.RepoCount, &f.ScriptCount); err != nil {
			return nil, err
		}
		if ts > 0 {
			f.LastScanned = time.Unix(ts, 0).Format("2006-01-02 15:04")
		}
		out = append(out, f)
	}
	if out == nil {
		out = []*Folder{}
	}
	return out, rows.Err()
}

func (s *Store) GetFolderByPath(path string) (*Folder, error) {
	row := s.db.QueryRow(`SELECT id, path, display_name, is_default, last_scanned FROM folders WHERE path = ?`, path)
	f := &Folder{}
	var ts int64
	if err := row.Scan(&f.ID, &f.Path, &f.DisplayName, &f.IsDefault, &ts); err != nil {
		return nil, err
	}
	if ts > 0 {
		f.LastScanned = time.Unix(ts, 0).Format("2006-01-02 15:04")
	}
	return f, nil
}

func (s *Store) AddFolder(path, displayName string, isDefault bool) (*Folder, error) {
	var id int64
	var exists int
	err := s.db.QueryRow(`SELECT id FROM folders WHERE path = ?`, path).Scan(&id)
	switch {
	case err == nil:
		exists = 1
	case err == sql.ErrNoRows:
		res, err := s.db.Exec(`INSERT INTO folders (path, display_name, is_default) VALUES (?, ?, ?)`, path, displayName, isDefault)
		if err != nil {
			return nil, err
		}
		id, _ = res.LastInsertId()
	default:
		return nil, err
	}
	_ = exists
	f, err := s.GetFolderByPath(path)
	if err != nil {
		return nil, err
	}
	f.ID = id
	return f, nil
}

func (s *Store) RemoveFolder(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repoIDs, err := s.repoIDsForFolder(tx, id)
	if err != nil {
		return err
	}
	for _, rid := range repoIDs {
		if err := s.deleteRepoData(tx, rid); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM folders WHERE id = ?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) SetFolderScanned(id int64, ts int64) error {
	_, err := s.db.Exec(`UPDATE folders SET last_scanned = ? WHERE id = ?`, ts, id)
	return err
}

// --- repos ---

func (s *Store) repoIDsForFolder(tx *sql.Tx, folderID int64) ([]int64, error) {
	rows, err := tx.Query(`SELECT id FROM repos WHERE folder_id = ?`, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Store) deleteRepoData(tx *sql.Tx, repoID int64) error {
	rows, err := tx.Query(`SELECT id FROM scripts WHERE repo_id = ?`, repoID)
	if err != nil {
		return err
	}
	var scriptIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		scriptIDs = append(scriptIDs, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range scriptIDs {
		if _, err := tx.Exec(`DELETE FROM scripts_fts WHERE rowid = ?`, id); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM scripts WHERE repo_id = ?`, repoID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM repos WHERE id = ?`, repoID); err != nil {
		return err
	}
	return nil
}

// upsertRepoTx creates or updates a repo inside the given transaction.
func upsertRepoTx(tx *sql.Tx, folderID int64, name, path, category, readmePath, readmeText string) (int64, error) {
	var id int64
	err := tx.QueryRow(`SELECT id FROM repos WHERE path = ?`, path).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := tx.Exec(`
			INSERT INTO repos (folder_id, name, path, category, readme_path, readme_text)
			VALUES (?, ?, ?, ?, ?, ?)`, folderID, name, path, category, readmePath, readmeText)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	if err != nil {
		return 0, err
	}
	if _, err := tx.Exec(`
		UPDATE repos SET folder_id = ?, name = ?, category = ?, readme_path = ?, readme_text = ?
		WHERE id = ?`, folderID, name, category, readmePath, readmeText, id); err != nil {
		return 0, err
	}
	return id, nil
}

// ListRepoIDs returns all repo ids for a folder.
func (s *Store) ListRepoIDs(folderID int64) ([]int64, error) {
	return listRepoIDsTx(s.db, folderID)
}

func listRepoIDsTx(q interface {
	Query(query string, args ...any) (*sql.Rows, error)
}, folderID int64) ([]int64, error) {
	rows, err := q.Query(`SELECT id FROM repos WHERE folder_id = ?`, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// --- scripts ---

// upsertScriptTx inserts or updates a script inside the given transaction.
// If the file is unchanged (same size and mtime), it is skipped.
// Content is only stored/indexed when storeContent is true.
func upsertScriptTx(tx *sql.Tx, repoID int64, relPath, absPath, name, ext, lang string, size, mtime int64, isScript bool, snippet, content string, storeContent bool) (int64, error) {
	var id int64
	var oldSize, oldMTime int64
	err := tx.QueryRow(`SELECT id, size, mtime FROM scripts WHERE abs_path = ?`, absPath).Scan(&id, &oldSize, &oldMTime)
	if err == sql.ErrNoRows {
		res, err := tx.Exec(`
			INSERT INTO scripts (repo_id, rel_path, name, ext, lang, size, mtime, abs_path, is_script, snippet, content)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			repoID, relPath, name, ext, lang, size, mtime, absPath, isScript, snippet, content)
		if err != nil {
			return 0, err
		}
		id, _ = res.LastInsertId()
		if storeContent && content != "" {
			if _, err := tx.Exec(`INSERT INTO scripts_fts (rowid, content) VALUES (?, ?)`, id, content); err != nil {
				return 0, err
			}
		}
		return id, nil
	}
	if err != nil {
		return 0, err
	}
	if oldSize == size && oldMTime == mtime {
		return id, nil // unchanged
	}
	if _, err := tx.Exec(`
		UPDATE scripts SET repo_id = ?, rel_path = ?, name = ?, ext = ?, lang = ?, size = ?, mtime = ?,
			abs_path = ?, is_script = ?, snippet = ?, content = ?
		WHERE id = ?`,
		repoID, relPath, name, ext, lang, size, mtime, absPath, isScript, snippet, content, id); err != nil {
		return 0, err
	}
	if storeContent {
		if _, err := tx.Exec(`INSERT OR REPLACE INTO scripts_fts (rowid, content) VALUES (?, ?)`, id, content); err != nil {
			return 0, err
		}
	} else {
		if _, err := tx.Exec(`DELETE FROM scripts_fts WHERE rowid = ?`, id); err != nil {
			return 0, err
		}
	}
	return id, nil
}

// deleteScriptsNotInTx removes indexed scripts for a repo that no longer exist on disk.
func deleteScriptsNotInTx(tx *sql.Tx, repoID int64, seenAbsPaths map[string]bool) (int, error) {
	rows, err := tx.Query(`SELECT id, abs_path FROM scripts WHERE repo_id = ?`, repoID)
	if err != nil {
		return 0, err
	}
	var toDelete []int64
	for rows.Next() {
		var id int64
		var p string
		if err := rows.Scan(&id, &p); err != nil {
			rows.Close()
			return 0, err
		}
		if !seenAbsPaths[p] {
			toDelete = append(toDelete, id)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for _, id := range toDelete {
		if _, err := tx.Exec(`DELETE FROM scripts_fts WHERE rowid = ?`, id); err != nil {
			return 0, err
		}
		if _, err := tx.Exec(`DELETE FROM scripts WHERE id = ?`, id); err != nil {
			return 0, err
		}
	}
	return len(toDelete), nil
}

// --- queries ---

func (s *Store) CatalogView() (*CatalogView, error) {
	// Empty slices (not nil) so JSON marshals [] instead of null — the
	// frontend reads .length on these.
	v := &CatalogView{Folders: []*Folder{}, Categories: []*Category{}, Repos: []*Repo{}}
	folders, err := s.ListFolders()
	if err != nil {
		return nil, err
	}
	v.Folders = folders

	rows, err := s.db.Query(`
		SELECT r.id, r.folder_id, f.path, r.name, r.path, r.category, r.readme_path,
			(SELECT count(*) FROM scripts sc WHERE sc.repo_id = r.id),
			(SELECT count(*) FROM scripts sc WHERE sc.repo_id = r.id AND sc.is_script = 1)
		FROM repos r JOIN folders f ON f.id = r.folder_id
		ORDER BY r.category, r.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	catCount := map[string]int{}
	for rows.Next() {
		r := &Repo{}
		if err := rows.Scan(&r.ID, &r.FolderID, &r.FolderPath, &r.Name, &r.Path, &r.Category, &r.ReadmePath, &r.FileCount, &r.ScriptCount); err != nil {
			return nil, err
		}
		v.Repos = append(v.Repos, r)
		catCount[r.Category]++
		v.TotalScripts += r.ScriptCount
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, name := range AllCategories() {
		if c, ok := catCount[name]; ok {
			v.Categories = append(v.Categories, &Category{Name: name, Count: c})
		}
	}
	return v, nil
}

func (s *Store) RepoDetail(repoID int64) (*RepoDetail, error) {
	r := &Repo{}
	err := s.db.QueryRow(`
		SELECT r.id, r.folder_id, f.path, r.name, r.path, r.category, r.readme_path, r.readme_text,
			(SELECT count(*) FROM scripts sc WHERE sc.repo_id = r.id),
			(SELECT count(*) FROM scripts sc WHERE sc.repo_id = r.id AND sc.is_script = 1)
		FROM repos r JOIN folders f ON f.id = r.folder_id WHERE r.id = ?`, repoID).
		Scan(&r.ID, &r.FolderID, &r.FolderPath, &r.Name, &r.Path, &r.Category, &r.ReadmePath, &r.ReadmeText, &r.FileCount, &r.ScriptCount)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, repo_id, rel_path, abs_path, name, ext, lang, size, mtime, is_script, snippet
		FROM scripts WHERE repo_id = ? ORDER BY is_script DESC, rel_path LIMIT 3000`, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scripts []*ScriptFile
	for rows.Next() {
		sc := &ScriptFile{}
		if err := rows.Scan(&sc.ID, &sc.RepoID, &sc.RelPath, &sc.AbsPath, &sc.Name, &sc.Ext, &sc.Lang, &sc.Size, &sc.MTime, &sc.IsScript, &sc.Snippet); err != nil {
			return nil, err
		}
		sc.Repo = r.Name
		scripts = append(scripts, sc)
	}
	if scripts == nil {
		scripts = []*ScriptFile{}
	}
	return &RepoDetail{Repo: r, Scripts: scripts}, rows.Err()
}

// ScriptContent returns the text of a script. The real file is read from disk
// (scripts always ship next to the app); the DB value is a capped fallback.
func (s *Store) ScriptContent(path string) (string, error) {
	if data, err := os.ReadFile(path); err == nil {
		return string(data), nil
	}
	var content string
	err := s.db.QueryRow(`SELECT content FROM scripts WHERE abs_path = ?`, path).Scan(&content)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *Store) ScriptByPath(path string) (*ScriptFile, error) {
	sc := &ScriptFile{}
	err := s.db.QueryRow(`
		SELECT sc.id, sc.repo_id, sc.rel_path, sc.abs_path, sc.name, sc.ext, sc.lang, sc.size, sc.mtime, sc.is_script, sc.snippet
		FROM scripts sc WHERE sc.abs_path = ?`, path).
		Scan(&sc.ID, &sc.RepoID, &sc.RelPath, &sc.AbsPath, &sc.Name, &sc.Ext, &sc.Lang, &sc.Size, &sc.MTime, &sc.IsScript, &sc.Snippet)
	if err != nil {
		return nil, err
	}
	_ = s.db.QueryRow(`SELECT name FROM repos WHERE id = ?`, sc.RepoID).Scan(&sc.Repo)
	return sc, nil
}

// SearchName searches repo names and script names with LIKE.
func (s *Store) SearchName(q string) (*SearchResult, error) {
	q = "%" + strings.TrimSpace(q) + "%"
	res := &SearchResult{Repos: []*Repo{}, Scripts: []*ScriptFile{}}
	rows, err := s.db.Query(`
		SELECT r.id, r.folder_id, f.path, r.name, r.path, r.category, r.readme_path,
			(SELECT count(*) FROM scripts sc WHERE sc.repo_id = r.id),
			(SELECT count(*) FROM scripts sc WHERE sc.repo_id = r.id AND sc.is_script = 1)
		FROM repos r JOIN folders f ON f.id = r.folder_id
		WHERE r.name LIKE ? OR r.category LIKE ?
		ORDER BY r.name LIMIT 100`, q, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r := &Repo{}
		if err := rows.Scan(&r.ID, &r.FolderID, &r.FolderPath, &r.Name, &r.Path, &r.Category, &r.ReadmePath, &r.FileCount, &r.ScriptCount); err != nil {
			rows.Close()
			return nil, err
		}
		res.Repos = append(res.Repos, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = s.db.Query(`
		SELECT sc.id, sc.repo_id, r.name, sc.rel_path, sc.abs_path, sc.name, sc.ext, sc.lang, sc.size, sc.mtime, sc.is_script, sc.snippet
		FROM scripts sc JOIN repos r ON r.id = sc.repo_id
		WHERE sc.name LIKE ?
		ORDER BY sc.name LIMIT 200`, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		sc := &ScriptFile{}
		if err := rows.Scan(&sc.ID, &sc.RepoID, &sc.Repo, &sc.RelPath, &sc.AbsPath, &sc.Name, &sc.Ext, &sc.Lang, &sc.Size, &sc.MTime, &sc.IsScript, &sc.Snippet); err != nil {
			rows.Close()
			return nil, err
		}
		res.Scripts = append(res.Scripts, sc)
	}
	return res, rows.Err()
}

// SearchContent uses FTS5 to find scripts whose indexed content matches.
func (s *Store) SearchContent(q string) ([]*ScriptFile, error) {
	match, ok := ftsMatch(q)
	if !ok {
		return []*ScriptFile{}, nil
	}
	rows, err := s.db.Query(`
		SELECT sc.id, sc.repo_id, r.name, sc.rel_path, sc.abs_path, sc.name, sc.ext, sc.lang, sc.size, sc.mtime, sc.is_script, sc.snippet
		FROM scripts_fts f JOIN scripts sc ON sc.id = f.rowid JOIN repos r ON r.id = sc.repo_id
		WHERE scripts_fts MATCH ?
		ORDER BY rank LIMIT 200`, match)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ScriptFile
	for rows.Next() {
		sc := &ScriptFile{}
		if err := rows.Scan(&sc.ID, &sc.RepoID, &sc.Repo, &sc.RelPath, &sc.AbsPath, &sc.Name, &sc.Ext, &sc.Lang, &sc.Size, &sc.MTime, &sc.IsScript, &sc.Snippet); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	if out == nil {
		out = []*ScriptFile{}
	}
	return out, rows.Err()
}

// ftsMatch converts a user query into a safe FTS5 MATCH expression.
func ftsMatch(q string) (string, bool) {
	fields := strings.Fields(strings.TrimSpace(q))
	if len(fields) == 0 {
		return "", false
	}
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		if strings.ContainsAny(f, `"*:()`) {
			f = strings.Trim(f, `"*:()[]{}^~-+`)
		}
		f = strings.ReplaceAll(f, `"`, `""`)
		if f == "" {
			continue
		}
		parts = append(parts, `"`+f+`"`)
	}
	if len(parts) == 0 {
		return "", false
	}
	return strings.Join(parts, " AND "), true
}

// ListScripts returns script files filtered by optional language and name query.
func (s *Store) ListScripts(lang, q string) ([]*ScriptFile, error) {
	sqlQ := `SELECT sc.id, sc.repo_id, r.name, sc.rel_path, sc.abs_path, sc.name, sc.ext, sc.lang, sc.size, sc.mtime, sc.is_script, sc.snippet
		FROM scripts sc JOIN repos r ON r.id = sc.repo_id WHERE sc.is_script = 1`
	var args []any
	if lang != "" && lang != "all" {
		sqlQ += ` AND sc.lang = ?`
		args = append(args, lang)
	}
	if q != "" {
		sqlQ += ` AND (sc.name LIKE ? OR sc.rel_path LIKE ?)`
		like := "%" + q + "%"
		args = append(args, like, like)
	}
	sqlQ += ` ORDER BY sc.name LIMIT 500`
	rows, err := s.db.Query(sqlQ, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*ScriptFile{}
	for rows.Next() {
		sc := &ScriptFile{}
		if err := rows.Scan(&sc.ID, &sc.RepoID, &sc.Repo, &sc.RelPath, &sc.AbsPath, &sc.Name, &sc.Ext, &sc.Lang, &sc.Size, &sc.MTime, &sc.IsScript, &sc.Snippet); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

// Totals returns total repo and script counts.
func (s *Store) Totals() (repos, scripts int, err error) {
	if err := s.db.QueryRow(`SELECT count(*) FROM repos`).Scan(&repos); err != nil {
		return 0, 0, err
	}
	if err := s.db.QueryRow(`SELECT count(*) FROM scripts WHERE is_script = 1`).Scan(&scripts); err != nil {
		return 0, 0, err
	}
	return repos, scripts, nil
}

// --- Custom Categories ---

// CustomCategory represents a user-created category.
type CustomCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListCustomCategories returns all user-created categories.
func (s *Store) ListCustomCategories() ([]*CustomCategory, error) {
	rows, err := s.db.Query(`SELECT id, name FROM custom_categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*CustomCategory
	for rows.Next() {
		c := &CustomCategory{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []*CustomCategory{}
	}
	return out, rows.Err()
}

// CreateCustomCategory adds a new category.
func (s *Store) CreateCustomCategory(name string) (*CustomCategory, error) {
	res, err := s.db.Exec(`INSERT INTO custom_categories (name) VALUES (?)`, name)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &CustomCategory{ID: id, Name: name}, nil
}

// RenameCustomCategory renames a category and updates all repos using it.
func (s *Store) RenameCustomCategory(id int64, newName string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var oldName string
	if err := tx.QueryRow(`SELECT name FROM custom_categories WHERE id = ?`, id).Scan(&oldName); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE custom_categories SET name = ? WHERE id = ?`, newName, id); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE custom_repos SET category = ? WHERE category = ?`, newName, oldName); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteCustomCategory removes a category and all repos assigned to it.
func (s *Store) DeleteCustomCategory(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var name string
	if err := tx.QueryRow(`SELECT name FROM custom_categories WHERE id = ?`, id).Scan(&name); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM custom_repos WHERE category = ?`, name); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM custom_categories WHERE id = ?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

// --- Custom Repos ---

// CustomRepo represents a user-added GitHub repository.
type CustomRepo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

// ListCustomRepos returns all user-added repos.
func (s *Store) ListCustomRepos() ([]*CustomRepo, error) {
	rows, err := s.db.Query(`SELECT id, name, url, category, summary FROM custom_repos ORDER BY category, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*CustomRepo
	for rows.Next() {
		r := &CustomRepo{}
		if err := rows.Scan(&r.ID, &r.Name, &r.URL, &r.Category, &r.Summary); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if out == nil {
		out = []*CustomRepo{}
	}
	return out, rows.Err()
}

// AddCustomRepo adds a user-provided GitHub repository.
func (s *Store) AddCustomRepo(name, url, category, summary string) (*CustomRepo, error) {
	res, err := s.db.Exec(
		`INSERT INTO custom_repos (name, url, category, summary) VALUES (?, ?, ?, ?)`,
		name, url, category, summary,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &CustomRepo{ID: id, Name: name, URL: url, Category: category, Summary: summary}, nil
}

// RemoveCustomRepo deletes a user-added repo by ID.
func (s *Store) RemoveCustomRepo(id int64) error {
	_, err := s.db.Exec(`DELETE FROM custom_repos WHERE id = ?`, id)
	return err
}

// --- Custom Software ---

// CustomSoftware is a user-added software item.
type CustomSoftware struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Category string `json:"category"`
	Download string `json:"download"`
	Notes    string `json:"notes"`
	WingetID string `json:"wingetId"`
}

// ListCustomSoftware returns all user-added software items.
func (s *Store) ListCustomSoftware() ([]*CustomSoftware, error) {
	rows, err := s.db.Query(`SELECT id, name, version, category, download, notes, winget_id FROM custom_software ORDER BY category, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*CustomSoftware
	for rows.Next() {
		sw := &CustomSoftware{}
		if err := rows.Scan(&sw.ID, &sw.Name, &sw.Version, &sw.Category, &sw.Download, &sw.Notes, &sw.WingetID); err != nil {
			return nil, err
		}
		out = append(out, sw)
	}
	if out == nil {
		out = []*CustomSoftware{}
	}
	return out, rows.Err()
}

// AddCustomSoftware adds a user-defined software item.
func (s *Store) AddCustomSoftware(name, version, category, download, notes, wingetID string) (*CustomSoftware, error) {
	res, err := s.db.Exec(
		`INSERT INTO custom_software (name, version, category, download, notes, winget_id) VALUES (?, ?, ?, ?, ?, ?)`,
		name, version, category, download, notes, wingetID,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &CustomSoftware{ID: id, Name: name, Version: version, Category: category, Download: download, Notes: notes, WingetID: wingetID}, nil
}

// RemoveCustomSoftware deletes a user-added software item by ID.
func (s *Store) RemoveCustomSoftware(id int64) error {
	_, err := s.db.Exec(`DELETE FROM custom_software WHERE id = ?`, id)
	return err
}
