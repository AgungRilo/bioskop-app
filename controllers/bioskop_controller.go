package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/AgungRilo/bioskop-app/config"
	"github.com/AgungRilo/bioskop-app/models"
	"github.com/gin-gonic/gin"
)

// parseID membaca parameter :id dan memastikan berupa angka positif
func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka positif"})
		return 0, false
	}
	return id, true
}

// bindBioskop membaca body JSON dan memvalidasi field wajib
func bindBioskop(c *gin.Context) (models.Bioskop, bool) {
	var input models.Bioskop
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data JSON tidak valid: " + err.Error()})
		return input, false
	}

	input.Nama = strings.TrimSpace(input.Nama)
	input.Lokasi = strings.TrimSpace(input.Lokasi)
	if input.Nama == "" || input.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
		return input, false
	}
	if input.Rating < 0 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating harus di antara 0 sampai 5"})
		return input, false
	}
	return input, true
}

// POST /bioskop
func CreateBioskop(c *gin.Context) {
	input, ok := bindBioskop(c)
	if !ok {
		return
	}

	err := config.DB.QueryRow(
		"INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id",
		input.Nama, input.Lokasi, input.Rating,
	).Scan(&input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "bioskop berhasil ditambahkan",
		"data":    input,
	})
}

// GET /bioskop
func GetAllBioskop(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, lokasi, COALESCE(rating, 0) FROM bioskop ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data: " + err.Error()})
		return
	}
	defer rows.Close()

	bioskops := []models.Bioskop{}
	for rows.Next() {
		var b models.Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membaca data: " + err.Error()})
			return
		}
		bioskops = append(bioskops, b)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membaca data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bioskops})
}

// GET /bioskop/:id
func GetBioskopByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var b models.Bioskop
	err := config.DB.QueryRow(
		"SELECT id, nama, lokasi, COALESCE(rating, 0) FROM bioskop WHERE id = $1", id,
	).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "bioskop dengan ID " + strconv.Itoa(id) + " tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": b})
}

// PUT /bioskop/:id
func UpdateBioskop(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	input, ok := bindBioskop(c)
	if !ok {
		return
	}

	err := config.DB.QueryRow(
		"UPDATE bioskop SET nama = $1, lokasi = $2, rating = $3 WHERE id = $4 RETURNING id",
		input.Nama, input.Lokasi, input.Rating, id,
	).Scan(&input.ID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "bioskop dengan ID " + strconv.Itoa(id) + " tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bioskop berhasil diperbarui",
		"data":    input,
	})
}

// DELETE /bioskop/:id
func DeleteBioskop(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	result, err := config.DB.Exec("DELETE FROM bioskop WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data: " + err.Error()})
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data: " + err.Error()})
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "bioskop dengan ID " + strconv.Itoa(id) + " tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bioskop berhasil dihapus"})
}
