// Package page defines HTTP handlers for the web application.
package page

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/models"
	"github.com/minhtranin/ct/internal/render"
	view "github.com/minhtranin/ct/internal/view/components"
	"github.com/minhtranin/ct/internal/view/layouts"
)

func Story(f *fiber.Ctx) error {
	// Get user from session
	userID, err := handler.GetSession(f)
	if err != nil {
		return f.Redirect("/signin")
	}

	// Fetch user from database
	db := handler.GetDB()
	usersCollection := db.Database("ct").Collection("users")

	var user models.User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = usersCollection.FindOne(f.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return f.Redirect("/signin")
	}

	return render.HTML(
		f,
		layouts.Dashboard("Story", view.StoryPage(), false, user.Email),
	)
}
