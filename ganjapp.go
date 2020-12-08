package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kaigoh/ganjapp/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/kaigoh/ganjapp/middleware"
	"github.com/kaigoh/ganjapp/utilities"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/controllers"
)

// baseURL for SeedFinder database
var baseURL = "https://en.seedfinder.eu/"

var databaseURL = baseURL + "database/strains/alphabetical/"

// childURLs contains an array of URLs to append to the base URL
var childURLs = map[string]string{
	"numeric": "1234567890",
	"a":       "a-all",
	"b":       "b-all",
	"c":       "c-all",
	"d":       "d-all",
	"e":       "e-all",
	"f":       "f-all",
	"g":       "g-all",
	"h":       "h-all",
	"i":       "i-all",
	"j":       "j-all",
	"k":       "k-all",
	"l":       "l-all",
	"m":       "m-all",
	"n":       "n-all",
	"o":       "o-all",
	"p":       "p-all",
	"q":       "q-all",
	"r":       "r-all",
	"s":       "s-all",
	"t":       "t-all",
	"u":       "u-all",
	"v":       "v-all",
	"w":       "w-all",
	"x":       "x-all",
	"y":       "y-all",
	"z":       "z-all",
}

// App config
var appConfig = utilities.AppConfig

func main() {

	fmt.Println("--------------")
	fmt.Println("Ganjapp Server")
	fmt.Println("--------------")

	// Run the web server...
	go server()

	// Ensure we have an admin user...
	go checkAdminAccount()

	// Perform initial population of strain database...
	go scrapeDatabase()

	// Plus, create a ticker so that the scrape runs weekly...
	for range time.Tick(time.Hour * 168) {
		go scrapeDatabase()
	}

}

func checkAdminAccount() {
	// Do we need to create an initial admin account?
	if !models.HaveAdminAccount() {
		models.CreateInitialAdminAccount()
	}
}

// Server
func server() {

	r := gin.Default()

	store := cookie.NewStore([]byte(appConfig.CookieKey))
	r.Use(sessions.Sessions("ganjappsession", store))

	r.Use(helmet.Default())

	r.Delims("{!", "!}")
	r.LoadHTMLGlob(utilities.GetEnv("GANJAPP_TEMPLATE_PATH", "templates") + "/*.tmpl")

	/**
	 * Routes
	 */

	// Serve up static resources
	r.StaticFile("/favicon.ico", "./icons/favicon.ico")
	r.StaticFile("/site.webmanifest", "./web/site.webmanifest")
	r.StaticFile("/browserconfig.xml", "./web/browserconfig.xml")
	r.Static("/assets", "./assets")
	r.Static("/icons", "./icons")

	// Ruotes that handle logins etc.
	authentication := r.Group("/auth")
	{
		authentication.GET("/login", controllers.Login)
		authentication.POST("/login", controllers.Login)
		authentication.POST("/authenticate", controllers.Authenticate)
	}

	// Authenticated routes
	authenticated := r.Group("/", middleware.IsUserLoggedInMiddleware())
	{

		// Enable GZip compression...
		authenticated.Use(gzip.Gzip(gzip.DefaultCompression))

		// Catch-all route...
		authenticated.Any("/", controllers.Home)

		// Route for fetching photos from S3...
		authenticated.GET("/photos/:objectType/:imageUUID", controllers.Photo)

		// These are "fake" routes, for Vue SPA routing...
		authenticated.Any("/dashboard", controllers.Home)
		authenticated.Any("/events", controllers.Home)
		authenticated.Any("/environment/*environmentUUID", controllers.Home)

	}

	// API routes
	api := r.Group("/api", middleware.IsUserAuthenticatedMiddleware())
	{

		// Enable GZip compression...
		api.Use(gzip.Gzip(gzip.DefaultCompression))

		// Return a token the user can use for external API calls...
		// Note that we issue this token with NO expiry...
		api.GET("/gettoken", controllers.GetToken)

		// Fetch events for the logged in user...
		api.GET("/events", controllers.GetEvents)

		// Create a new environment...
		api.POST("/create/environment", controllers.CreateEnvironment)

		// Returns an array containing all the environments the user owns...
		api.GET("/environments", controllers.GetEnvironments)

		// Endpoint to allow an environment property to be updated...
		api.POST("/environment/:id/update/:property", controllers.UpdateEnvironment)

		// Endpoint to allow an environments extended data property to be updated...
		api.POST("/environment/:id/extendeddata/:property", controllers.UpdateEnvironmentExtendedData)

		// Endpoint that accepts image uploads for environments
		api.POST("/environment/:id/upload", controllers.UploadEnvironmentImage)

		// Endpoint that moves trees between environments...
		api.GET("/tree/:id/move/:environment", controllers.MoveTree)

		// Endpoint to allow tree properties to be updated...
		api.POST("/tree/:id/update/:property", controllers.UpdateTree)

		// Endpoint to allow a trees extended data property to be updated...
		api.POST("/tree/:id/extendeddata/:property", controllers.UpdateTreeExtendedData)

		// Endpoint that accepts image uploads for trees...
		api.POST("/tree/:id/upload", controllers.UploadTreeImage)

		// Endpoint that moves shrooms between environments...
		api.GET("/shroom/:id/move/:environment", controllers.MoveShroom)

		// Endpoint to allow shroom properties to be updated...
		api.POST("/shroom/:id/update/:property", controllers.UpdateShroom)

		// Endpoint to allow a shrooms extended data property to be updated...
		api.POST("/shroom/:id/extendeddata/:property", controllers.UpdateShroomExtendedData)

		// Endpoint that accepts image uploads for shrooms...
		api.POST("/shroom/:id/upload", controllers.UploadShroomImage)

	}

	// Push data routes
	// Note: Do NOT enable GZip compression on any of these routes!
	pushAPI := r.Group("/live", middleware.IsUserAuthenticatedMiddleware())
	{
		pushAPI.GET("/sse", controllers.SSE)
	}

	r.Run()

}

// Seedfinder.eu Scraping...

func scrapeDatabase() {
	// Do we need to scrape the database, or has it been done recently?
	l, err := models.GetObjectEventsByType("system", 0, "cannabis-strain-database-scraped", 0)
	if err != nil || len(l) == 0 || l[0].CreatedAt.Before(time.Now().AddDate(0, 0, -7)) {
		log.Println("Starting scrape of " + baseURL + "...")
		for _, v := range childURLs {
			url := databaseURL + v + "/"
			log.Println("Attempting to fetch data from " + url + "...")
			scrape(url)
		}
		// Log the scraping event to the database...
		models.LogObjectEvent(0, "system", 0, "cannabis-strain-database-scraped", "info", "seedfinder.eu cannabis strain database has been scraped", time.Now().String())
	} else {
		log.Println("Was going to start scrape of " + baseURL + ", but the database has been fetched recently, exiting...")
	}
}

// Scrape data from seedfinder.eu URL...
func scrape(url string) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "GanjappAPI/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf(url+" - Status code error: %d %s", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("table[id=\"cannabis-strain-table\"] tbody tr").Each(func(i int, s *goquery.Selection) {
		strainElem := s.Find("th a")
		strain := strainElem.Text()
		strainLinkRaw, _ := strainElem.Attr("href")
		strainLink := baseURL + strings.ReplaceAll(strainLinkRaw, "../", "")
		breeder := s.Find("td:nth-child(2)").Text()
		genetics, _ := s.Find("td:nth-child(3) img").Attr("title")
		environment, _ := s.Find("td:nth-child(4) img").Attr("title")
		floweringTime := s.Find("td:nth-child(5)").Text()
		seedGender, _ := s.Find("td:nth-child(6) img").Attr("title")

		// Try and find the strain and breeder in the database, creating a new entry if required...
		r, err := models.FindStrainAndBreeder(strain, breeder)

		if err != nil {
			// Create a new entry
			s := models.CannabisStrain{
				Breeder:       breeder,
				Strain:        strain,
				URL:           strainLink,
				Genetics:      genetics,
				Environment:   environment,
				FloweringTime: floweringTime,
				Gender:        seedGender,
			}
			models.DB.Create(&s)
		} else {
			// Update the existing entry
			r.Breeder = breeder
			r.Strain = strain
			r.URL = strainLink
			r.Genetics = genetics
			r.Environment = environment
			r.FloweringTime = floweringTime
			r.Gender = seedGender
			models.DB.Save(&r)
		}

	})

	time.Sleep(5 * time.Second)

}
