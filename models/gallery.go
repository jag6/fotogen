package models

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is used to tell the GalleryService where to store and locate images.
	// If not set, the GalleryService will default to using the "media" directory.
	ImagesDir string
}

func (gs *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}
	row := gs.DB.QueryRow(`
		SELECT title, user_id
		FROM galleries 
		WHERE id = $1;`, gallery.ID)
	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err = rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete gallery by id: %w", err)
	}
	err = os.RemoveAll(gs.galleryDir(id))
	if err != nil {
		return fmt.Errorf("delete gallery images: %w", err)
	}
	return nil
}

func (gs *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(gs.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, gs.extensions()) {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}
	return images, nil
}

func (gs *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(gs.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}
	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (gs *GalleryService) CreateImage(galleryID int, filename string, contents io.Reader) error {
	readBytes, err := checkContentType(contents, gs.imageContentTypes())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	err = checkExtension(filename, gs.extensions())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	galleryDir := gs.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d images directory: %w", galleryID, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	completeFile := io.MultiReader(
		bytes.NewBuffer(readBytes),
		contents,
	)
	_, err = io.Copy(dst, completeFile)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}
	return nil
}

func (gs *GalleryService) CreateImageViaURL(galleryID int, url string) error {
	filename := path.Base(url)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloading image: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading image: invalid status code %d", res.StatusCode)
	}
	return gs.CreateImage(galleryID, filename, res.Body)
}

func (gs *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := gs.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

func (gs *GalleryService) imageContentTypes() []string {
	return []string{"image/png", "image/jpg", "image/jpeg", "image/gif", "image/webp"}
}

func (gs *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}
}

func (gs *GalleryService) galleryDir(id int) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "media"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

// may move to users controller later
func (gs *GalleryService) ByUserIDForIndexRead(id int) (*User, error) {
	user := User{
		ID: id,
	}
	row := gs.DB.QueryRow(`
		SELECT username
		FROM users
		WHERE id = $1;`, user.ID)
	err := row.Scan(&user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return &user, nil
}
