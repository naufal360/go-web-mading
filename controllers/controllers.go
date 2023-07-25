package controllers

import (
	"jewepe/helpers"
	"jewepe/model"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var adminKey string = "adminkey"

// post
func (h *HttpServer) RegisterUsers(c *gin.Context) {
	// get form data
	username := c.PostForm("username")
	password := c.PostForm("password")

	// hashing password
	hassPash := helpers.HashPass(password)

	// uuid
	id := uuid.New().String()

	// bind data
	admin := model.Users{
		Id:       id,
		Username: username,
		Password: hassPash,
	}

	// gorm create user
	err := h.db.Debug().Create(&admin).Error
	// jika gagal
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": "Bad request",
		})
		return
	}

	// jika sukses
	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Success create admin",
	})
}

func (h *HttpServer) LoginUsers(c *gin.Context) {
	var admin model.Users

	// session
	session := sessions.Default(c)

	// get form data
	username := c.PostForm("username")
	password := c.PostForm("password")

	// find username
	if err := h.db.Debug().First(&admin, "username = ?", username).Error; err != nil {
		log.Println("account doesn't exist!")
		c.Redirect(http.StatusMovedPermanently, "./login")
		c.Abort()
		return
	}

	// check password
	compare := helpers.ComparePass(admin.Password, password)
	if compare == false {
		log.Println("wrong password!")
		c.Redirect(http.StatusMovedPermanently, "./login")
		c.Abort()
		return
	}

	// create session
	session.Set(adminKey, username)
	if err := session.Save(); err != nil {
		log.Println("succes create session!")
		c.Redirect(http.StatusMovedPermanently, "./login")
		c.Abort()
		return
	}

	log.Println("login succes!")
	c.Redirect(http.StatusFound, "./dashboard")
}

func (h *HttpServer) LogoutUser(c *gin.Context) {
	// session
	session := sessions.Default(c)
	admin := session.Get(adminKey)
	log.Println("logging out user:", admin)

	// redirect jika session kosong
	if admin == nil {
		log.Println("Invalid session token")
		c.Redirect(http.StatusFound, "./login")
		c.Abort()
		return
	}
	session.Delete(adminKey)
	if err := session.Save(); err != nil {
		log.Println("Failed to save session:", err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "./")
}

func (h *HttpServer) CreateArtikel(c *gin.Context) {
	var artikel model.Artikels

	// get data form
	judul := c.PostForm("judul")
	isiArtikel := c.PostForm("isi_artikel")

	// uuid
	id := uuid.New().String()

	// fill model
	artikel = model.Artikels{
		Id:         id,
		Judul:      judul,
		IsiArtikel: isiArtikel,
	}

	// gorm create artikel
	err := h.db.Debug().Create(&artikel).Error

	// if failed
	if err != nil {
		log.Println("Failed to create artikel:", err)
		c.Redirect(http.StatusInternalServerError, "./create")
		c.Abort()
		return
	}

	// if success
	c.Redirect(http.StatusMovedPermanently, "./")
}

func (h *HttpServer) DeleteArtikel(c *gin.Context) {
	// get id param
	id := c.Param("id")

	// gorm delete artikel
	if err := h.db.Debug().Delete(&model.Artikels{}, "id = ?", id).Error; err != nil {
		log.Println("Failed to delete artikel:", err)
		return
	}
}

func (h *HttpServer) CreateKomentar(c *gin.Context) {
	var komentarModel model.Komentars

	// get param artikel id
	artikelId := c.Param("id")

	// get data form
	nama := c.PostForm("nama")
	email := c.PostForm("email")
	komentar := c.PostForm("isi_komentar")

	// uuid
	id := uuid.New().String()

	// fill model
	komentarModel = model.Komentars{
		Id:          id,
		Nama:        nama,
		Email:       email,
		IsiKomentar: komentar,
		ArtikelId:   artikelId,
	}

	// gorm create komentar
	err := h.db.Debug().Create(&komentarModel).Error

	// if failed
	if err != nil {
		log.Println("Failed to create komentar:", err)
		c.Redirect(http.StatusInternalServerError, "./create")
		c.Abort()
		return
	}

	// url
	url := "/jewepe/artikel/" + artikelId
	cleanUrl := strings.Trim(url, "")

	// if success
	c.Redirect(http.StatusMovedPermanently, cleanUrl)
}

func (h *HttpServer) DeleteKomentar(c *gin.Context) {
	// get id param
	id := c.Param("id")

	// gorm delete artikel
	if err := h.db.Debug().Delete(&model.Komentars{}, "id = ?", id).Error; err != nil {
		log.Println("Failed to delete komentar:", err)
		return
	}
}

// read
func (h *HttpServer) ViewLogin(c *gin.Context) {
	// check session session
	session := sessions.Default(c)
	adminSession := session.Get(adminKey)
	if adminSession != nil {
		c.Redirect(http.StatusTemporaryRedirect, "./dashboard")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "login.htm", gin.H{
		"title": "Login",
	})
}

func (h *HttpServer) ViewDashboard(c *gin.Context) {
	var artikel []model.Artikels
	var responses []model.ResponsArtikel
	var komentars []model.Komentars
	var count int

	// session
	session := sessions.Default(c)
	admin := session.Get(adminKey)

	// get form search
	search := c.Query("search")

	// logic search
	if search != "" {
		newSearch := "%" + search + "%"
		// get data search
		if err := h.db.Debug().Preload("Komentar").Where("judul LIKE ?", newSearch).Order("created_at DESC").Find(&artikel).Error; err != nil {
			log.Println("Failed to get specific artikel:", err)
			return
		}
	} else {
		// get all data
		if err := h.db.Debug().Preload("Komentar").Order("created_at DESC").Find(&artikel).Error; err != nil {
			log.Println("Failed to get all artikel:", err)
			return
		}
	}

	// get komentar
	countKomentar := h.db.Debug().Find(&komentars).RowsAffected

	// create custom response
	var sneakpeak string
	responses = make([]model.ResponsArtikel, len(artikel))
	for i, data := range artikel {
		// create sneakpeak
		if len(data.IsiArtikel) < 100 {
			sneakpeak = data.IsiArtikel
		} else {
			sneakpeak = data.IsiArtikel[0:100] + "..."
		}

		responses[i] = model.ResponsArtikel{
			Id:         data.Id,
			Judul:      data.Judul,
			Sneakpeak:  sneakpeak,
			IsiArtikel: data.IsiArtikel,
			CreatedAt:  data.CreatedAt,
			Komentar:   data.Komentar,
		}
	}

	// count artikel
	count = len(artikel)

	c.HTML(http.StatusOK, "dashboard.htm", gin.H{
		"title":         "Dashboard",
		"session":       admin,
		"artikel":       responses,
		"count":         count,
		"komentarCount": countKomentar,
	})
}

func (h *HttpServer) ViewIndex(c *gin.Context) {
	var artikel []model.Artikels
	var responses []model.ResponsArtikel
	var count int

	// get form search
	search := c.Query("search")

	// logic search
	if search != "" {
		newSearch := "%" + search + "%"
		// get data search
		if err := h.db.Debug().Preload("Komentar").Where("judul LIKE ?", newSearch).Order("created_at DESC").Find(&artikel).Error; err != nil {
			log.Println("Failed to get specific artikel:", err)
			return
		}
	} else {
		// get all data
		if err := h.db.Debug().Preload("Komentar").Order("created_at DESC").Find(&artikel).Error; err != nil {
			log.Println("Failed to get all artikel:", err)
			return
		}
	}

	// create custom response
	var sneakpeak string
	responses = make([]model.ResponsArtikel, len(artikel))
	for i, data := range artikel {
		// create sneakpeak
		if len(data.IsiArtikel) < 100 {
			sneakpeak = data.IsiArtikel
		} else {
			sneakpeak = data.IsiArtikel[0:100] + "..."
		}

		responses[i] = model.ResponsArtikel{
			Id:         data.Id,
			Judul:      data.Judul,
			Sneakpeak:  sneakpeak,
			IsiArtikel: data.IsiArtikel,
			CreatedAt:  data.CreatedAt,
			Komentar:   data.Komentar,
		}
	}

	// count artikel
	count = len(artikel)

	c.HTML(http.StatusOK, "index.htm", gin.H{
		"title":   "Home",
		"artikel": responses,
		"count":   count,
	})
}

func (h *HttpServer) ViewCreateArtikel(c *gin.Context) {
	c.HTML(http.StatusOK, "c_artikel.htm", gin.H{
		"title": "Create Artikel",
	})
}

func (h *HttpServer) ViewDetailArtikel(c *gin.Context) {
	var artikel model.Artikels
	var isLogin bool
	var count int

	// session
	session := sessions.Default(c)
	admin := session.Get(adminKey)
	if admin != nil {
		isLogin = true
	}

	// get param id
	id := c.Param("id")

	// find artikel
	if err := h.db.Debug().Preload("Komentar").First(&artikel, "id = ?", id).Error; err != nil {
		log.Println("Failed to get detail artikel:", err)
		return
	}

	// count komentar
	count = len(artikel.Komentar)

	// if success
	c.HTML(http.StatusOK, "r_artikel.htm", gin.H{
		"artikel":       artikel,
		"isLogin":       isLogin,
		"komentarCount": count,
	})
}

func (h *HttpServer) ViewCreateKomentar(c *gin.Context) {
	var artikel model.Artikels

	// session check
	session := sessions.Default(c)
	admin := session.Get(adminKey)

	// get artikel id
	artikelId := c.Param("id")

	// if admin, redirect
	if admin != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/jewepe")
		return
	}

	// find artikel
	if err := h.db.Debug().First(&artikel, "id = ?", artikelId).Error; err != nil {
		log.Println("Failed to get detail artikel:", err)
		return
	}

	log.Println(artikel.Judul)

	// if success
	c.HTML(http.StatusOK, "komentar.htm", gin.H{
		"title":        "Create Komentar",
		"artikelId":    artikelId,
		"judulArtikel": artikel.Judul,
	})
}
