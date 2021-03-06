package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/steven7/go-createmusic/go/config"
	"github.com/steven7/go-createmusic/go/controllers"
	"github.com/steven7/go-createmusic/go/middleware"
	"github.com/steven7/go-createmusic/go/models"
	"net/http"
	"os"
)

func main() {

	// Get development channel from environment variable set in docker file.
	// Either dev or prod.
	channel := os.Getenv("channel_env_var")
	// channel == "prod" || channel == "preprod"
	isProd := channel != "dev" // preprod is with local docker netwroking and aws db and s3

	boolPtr := flag.Bool("prod", isProd, "Provide this flag "+
		"in production. This ensures that a .config file is "+
		"provided before the application starts.")
	flag.Parse()

	fmt.Printf("___________________ we are in %s channel ___________________ %b\n", channel, isProd)

	cfg := config.LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	fmt.Println("trying with host ", dbCfg.Host)
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		// only log when not in prod
		models.WithLogMode(!cfg.IsProd()),
		// We want each of these services, but if we didn't need
		// one of them we could possibly skip that config func
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithTrack(),
		//models.WithImage(),
		//models.WithMusicFile(),
		models.WithFile(),
		models.WithOauth(),
	)

	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// not set up
	//mgCfg := cfg.Mailgun
	//emailer := email.NewClient(
	//	email.WithSender("ImageCloud.com Support", "support@"+mgCfg.Domain),
	//	email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	//)

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	// emailer is nil because we arent using it now
	usersC := controllers.NewUsers(services.User, nil)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)
	//tracksC := controllers.NewTracksController(services.Track, services.Image, services.MusicFile, r)
	tracksC := controllers.NewTracksController(services.Track, services.File, r)

	//
	//
	//  services.DestructiveReset()
	//
	// be careful with this ^^^^^

	//
	// Middleware
	// csrf middleware
	//

	corsMw := cors.New(cors.Options{
		AllowedHeaders: []string{"accept", "authorization", "content-type"},
		AllowedOrigins: []string{"http://localhost", "http://localhost:3000", "http://localhost:5000",
			"http://172.19.0.3:5000", "*"}, // * is for testing only not production
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	//r.Handle("/", staticC.Home).Methods("GET")
	//r.Handle("/contact", staticC.Contact).Methods("GET")
	//r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// This will assign the page to the nor found handler
	var h http.Handler = http.Handler(staticC.NotFound)
	r.NotFoundHandler = h

	// Gallery routes

	//r.Handle("/galleries",
	//	requireUserMw.ApplyFn(galleriesC.Index)).
	//	Methods("GET").
	//	Name(controllers.IndexGalleries)
	//r.Handle("/galleries/new",
	//	requireUserMw.Apply(galleriesC.New)).
	//	Methods("GET")
	//r.Handle("/galleries",
	//	requireUserMw.ApplyFn(galleriesC.Create)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}",
	//	galleriesC.Show).
	//	Methods("GET").
	//	Name(controllers.ShowGallery)
	//r.HandleFunc("/galleries/{id:[0-9]+}/edit",
	//	galleriesC.Edit).
	//	Methods("GET").
	//	Name(controllers.EditGallery)
	//r.HandleFunc("/galleries/{id:[0-9]+}/update",
	//	requireUserMw.ApplyFn(galleriesC.Update)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}/delete",
	//	requireUserMw.ApplyFn(galleriesC.Delete)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}/images",
	//	requireUserMw.ApplyFn(galleriesC.ImageUpload)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}/images/link",
	//	requireUserMw.ApplyFn(galleriesC.ImageViaLink)).
	//	Methods("POST")

	// tracks

	// view tracks
	r.Handle("/tracks",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.ApplyFn(tracksC.Index)).
		Methods("GET").
		Name(controllers.IndexTracks)

	// create new track
	r.Handle("/tracks/new",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.Apply(tracksC.ChooseTypeView)).
		Methods("GET")

	r.HandleFunc("/tracks/{id:[0-9]+}/play",
		tracksC.Play).
		Methods("GET").
		Name(controllers.PlayTrack)

	// create new // phase one
	r.Handle("/tracks/createlocal",
		requireUserMw.ApplyFn(tracksC.CreateLocal)).
		Methods("POST")

	// create with compose backend service
	r.Handle("/tracks/createWithComposeAI",
		requireUserMw.ApplyFn(tracksC.CreateWithComposeAI)).
		Methods("POST")

	r.Handle("/tracks/createWithDJ",
		requireUserMw.ApplyFn(tracksC.ChooseDJOptions)).
		Methods("POST")

	//
	r.Handle("/tracks/createWithDJWorking",
		requireUserMw.ApplyFn(tracksC.CreateWithDJWorking)).
		Methods("POST")
	r.Handle("/tracks/createWithDJComplete",
		requireUserMw.ApplyFn(tracksC.CreateWithDJComplete)).
		Methods("POST")

	// create new looks like edit
	// when create song is pressed
	r.Handle("/tracks/createlocalcomplete",
		requireUserMw.ApplyFn(tracksC.CreateLocalComplete)).
		Methods("POST")

	// edit existing
	r.Handle("/tracks/{id:[0-9]+}/editLocalTrack",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.ApplyFn(tracksC.EditLocal)).
		Methods("GET"). //
		Name(controllers.EditTrack)
	r.Handle("/tracks/{id:[0-9]+}/editDJCreatedTrack",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.ApplyFn(tracksC.EditDJ)).
		Methods("GET")

	//
	// create track
	//
	// upload
	r.HandleFunc("/tracks/{id:[0-9]+}/music",
		requireUserMw.ApplyFn(tracksC.MusicUpload)).
		Methods("POST")
	r.HandleFunc("/tracks/{id:[0-9]+}/images",
		requireUserMw.ApplyFn(tracksC.ImageUpload)).
		Methods("POST")

	// add to db -- when edit song pressed
	r.HandleFunc("/tracks/{id:[0-9]+}/create",
		requireUserMw.ApplyFn(tracksC.CreateLocalSongWithDB)).
		Methods("POST")
	// edit existing pressed
	r.HandleFunc("/tracks/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(tracksC.EditLocalSongComplete)).
		Methods("POST")

	// Image routes
	imageHandler := http.FileServer(http.Dir("./userfiles/tracks/"))
	r.PathPrefix("/userfiles/tracks/").Handler(http.StripPrefix("/userfiles/tracks/", imageHandler))

	// file routes
	//imageHandler := http.FileServer(http.Dir("./images/"))
	//r.PathPrefix("/images/").Handler(http.StripPrefix("/images/",imageHandler))

	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).
		Methods("POST")

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	//
	//
	// API routes
	//
	//
	usersCAPI := controllers.NewUsersAPI(services.User, nil)
	tracksCAPI := controllers.NewTracksAPI(services.Track, services.File, r, cfg) // Include config as parameter

	//
	//
	//  API routes
	//
	//

	//
	// Auth API
	r.HandleFunc("/api/auth/login", usersCAPI.AuthenticateWithAPI).Methods("POST")
	r.HandleFunc("/api/auth/signup", usersCAPI.CreateWithAPI).Methods("POST")

	//
	// Tracks API
	r.HandleFunc("/api/tracks/createlocal", tracksCAPI.CreateLocalWithAPI).Methods("POST")
	//r.HandleFunc("/api/tracks/createWithComposeAI", tracksCAPI.CreateWithComposeAI).Methods("POST")
	r.HandleFunc("/api/tracks/createWithComposeAI", tracksCAPI.CreateWithComposeAI).Methods("POST")
	//
	r.HandleFunc("/api/tracks/index", tracksCAPI.IndexWithAPI).Methods("POST")
	//r.HandleFunc("/api/tracks/one", tracksCAPI.GetTrackWithAPI).Methods("GET")
	r.HandleFunc("/api/tracks/one/coverimage", tracksCAPI.GetTrackCoverFileWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/one/musicfile", tracksCAPI.GetTrackMusicFileWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/update", tracksCAPI.UpdateTrackWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/delete", tracksCAPI.DeleteTrackWithAPI).Methods("POST")

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), corsMw.Handler(userMw.Apply(r)))

}
