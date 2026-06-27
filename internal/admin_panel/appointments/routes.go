package appointments

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/user"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

func AppointmentRoutes(rg *gin.RouterGroup, s *Store) {
	protected := rg.Group("/")
	userStore := user.NewStore(s.db)
	protected.Use(user.WithJWTAuth(userStore))
	protected.GET("/appointments", s.ListAppointments)

}

func (s *Store) ListAppointments(c *gin.Context) {
	var appointments []store.Appointment

	// meds, err := medicationreminder.GetClientsWithCurrentMedications()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, med := range meds {
	// 	if med.MatterUid == "2y1ncrj8uh7qk8ye" {
	// 		log.Printf("med: %+v\n", med)
	// 	}
	// }

	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		log.Printf("page: %s", p)
	}
	if l := c.Query("limit"); l != "" {
		log.Printf("limit: %s", l)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64

	if err := s.db.Model(&store.Appointment{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count appointments",
		})
		return
	}

	if err := s.db.
		Limit(limit).
		Offset(offset).
		Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve appointments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": appointments,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}
