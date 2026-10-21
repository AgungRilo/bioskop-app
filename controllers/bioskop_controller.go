package controllers

import (
	"net/http"
	"strings"

	"github.com/AgungRilo/bioskop-app/config"
	"github.com/AgungRilo/bioskop-app/models"
	"github.com/gin-gonic/gin"
)

// POST /bioskop
func CreateBioskop(c *gin.Context) {
	var input models.Bioskop
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data JSON tidak valid: " + err.Error()})
		return
	}

	input.Nama = strings.TrimSpace(input.Nama)
	input.Lokasi = strings.TrimSpace(input.Lokasi)
	if input.Nama == "" || input.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
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
	rows, err := config.DB.Query("SELECT id, nama, lokasi, rating FROM bioskop ORDER BY id")
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

	c.JSON(http.StatusOK, gin.H{"data": bioskops})
}
