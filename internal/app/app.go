package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"toucan/internal/content"
	"toucan/internal/courses"
	"toucan/internal/database"
	"toucan/internal/enrollments"
	"toucan/internal/identity"
	"toucan/internal/sections"
	"toucan/internal/storage"
	"toucan/internal/uploads"
	"toucan/internal/users"
)

type App struct {
	Handler http.Handler
	Close   func() error
}

func New(dbCfg database.Config, storageCfg storage.Config, identityCfg identity.Config, logger *log.Logger) (*App, error) {
	db, err := database.Open(dbCfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := database.EnsureSchema(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ensure schema: %w", err)
	}

	store, err := initStorage(storageCfg)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init storage: %w", err)
	}

	auth, err := identity.NewAuthenticator(identityCfg)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init identity: %w", err)
	}

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)

	// Hook JIT provisioning into the authenticator
	auth.SyncUser = func(ctx context.Context, principal identity.Principal) error {
		_, err := userService.EnsureUser(ctx, principal)
		return err
	}

	enrollmentRepo := enrollments.NewRepository(db)
	enrollmentService := enrollments.NewService(enrollmentRepo, userService)
	courseRepo := courses.NewRepository(db)
	courseService := courses.NewService(courseRepo, userService, enrollmentService)
	sectionRepo := sections.NewRepository(db)
	sectionService := sections.NewService(sectionRepo, courseService)
	contentRepo := content.NewRepository(db)
	contentService := content.NewService(contentRepo, sectionService)

	return &App{
		Handler: buildHandler(logger, courseService, enrollmentService, sectionService, contentService, userService, store, auth),
		Close:   db.Close,
	}, nil
}

func initStorage(cfg storage.Config) (storage.Store, error) {
	switch cfg.Driver {
	case storage.BlobDriverS3:
		if cfg.S3Bucket == "" {
			return nil, fmt.Errorf("s3 bucket is required when using s3 driver")
		}
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.S3Region))
		if err != nil {
			return nil, fmt.Errorf("load aws config: %w", err)
		}
		return storage.NewS3Store(s3.NewFromConfig(awsCfg), cfg.S3Bucket), nil

	case storage.BlobDriverLocal:
		return storage.NewLocalStore(cfg.LocalPath, "/uploads")

	default:
		return nil, fmt.Errorf("unsupported blob driver %q", cfg.Driver)
	}
}

func buildHandler(
	logger *log.Logger,
	courseService *courses.Service,
	enrollmentService *enrollments.Service,
	sectionService *sections.Service,
	contentService *content.Service,
	userService users.Service,
	store storage.Store,
	auth *identity.Authenticator,
) http.Handler {
	courseHandler := courses.NewHandler(courseService, enrollmentService, logger)
	enrollmentHandler := enrollments.NewHandler(enrollmentService)
	sectionHandler := sections.NewHandler(sectionService)
	contentHandler := content.NewHandler(contentService)
	userHandler := users.NewHandler(userService)
	uploadHandler := uploads.NewHandler(store)
	requireAuth := auth.Middleware

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", courseHandler.HandleHealth)

	// API routes
	mux.Handle("GET /api/v1/courses", requireAuth(http.HandlerFunc(courseHandler.HandleListCourses)))
	mux.Handle("GET /api/v1/me/courses", requireAuth(http.HandlerFunc(courseHandler.HandleListMyCourses)))
	mux.Handle("POST /api/v1/courses", requireAuth(http.HandlerFunc(courseHandler.HandleCreateCourse)))
	mux.Handle("GET /api/v1/courses/{id}", requireAuth(http.HandlerFunc(courseHandler.HandleGetCourse)))
	mux.Handle("GET /api/v1/courses/{id}/member/me", requireAuth(http.HandlerFunc(courseHandler.HandleGetMyMembership)))
	mux.Handle("PUT /api/v1/courses/{id}", requireAuth(http.HandlerFunc(courseHandler.HandleUpdateCourse)))
	mux.Handle("DELETE /api/v1/courses/{id}", requireAuth(http.HandlerFunc(courseHandler.HandleDeleteCourse)))
	mux.Handle("POST /api/v1/courses/{id}/publish", requireAuth(http.HandlerFunc(courseHandler.HandlePublishCourse)))
	mux.Handle("POST /api/v1/courses/{id}/archive", requireAuth(http.HandlerFunc(courseHandler.HandleArchiveCourse)))

	mux.Handle("POST /api/v1/enrollments", requireAuth(http.HandlerFunc(enrollmentHandler.HandleCreateEnrollment)))
	mux.Handle("DELETE /api/v1/enrollments", requireAuth(http.HandlerFunc(enrollmentHandler.HandleDeleteEnrollment)))
	mux.Handle("GET /api/v1/courses/{id}/enrollments", requireAuth(http.HandlerFunc(enrollmentHandler.HandleListCourseEnrollments)))

	mux.Handle("GET /api/v1/sections", requireAuth(http.HandlerFunc(sectionHandler.HandleListSections)))
	mux.Handle("POST /api/v1/sections", requireAuth(http.HandlerFunc(sectionHandler.HandleCreateSection)))
	mux.Handle("GET /api/v1/sections/{id}", requireAuth(http.HandlerFunc(sectionHandler.HandleGetSection)))
	mux.Handle("PUT /api/v1/sections/{id}", requireAuth(http.HandlerFunc(sectionHandler.HandleUpdateSection)))
	mux.Handle("DELETE /api/v1/sections/{id}", requireAuth(http.HandlerFunc(sectionHandler.HandleDeleteSection)))

	mux.Handle("GET /api/v1/content", requireAuth(http.HandlerFunc(contentHandler.HandleListContent)))
	mux.Handle("POST /api/v1/content", requireAuth(http.HandlerFunc(contentHandler.HandleCreateContent)))
	mux.Handle("GET /api/v1/content/{id}", requireAuth(http.HandlerFunc(contentHandler.HandleGetContent)))
	mux.Handle("PUT /api/v1/content/{id}", requireAuth(http.HandlerFunc(contentHandler.HandleUpdateContent)))
	mux.Handle("DELETE /api/v1/content/{id}", requireAuth(http.HandlerFunc(contentHandler.HandleDeleteContent)))

	mux.Handle("GET /api/v1/users", requireAuth(http.HandlerFunc(userHandler.HandleListUsers)))
	mux.Handle("POST /api/v1/users", requireAuth(http.HandlerFunc(userHandler.HandleCreateUser)))
	mux.Handle("GET /api/v1/users/{id}", requireAuth(http.HandlerFunc(userHandler.HandleGetUser)))

	mux.Handle("POST /api/v1/uploads/presign", requireAuth(http.HandlerFunc(uploadHandler.HandlePresign)))
	mux.Handle("POST /api/v1/uploads", requireAuth(http.HandlerFunc(uploadHandler.HandleUpload)))
	mux.Handle("GET /api/v1/uploads/{key}", requireAuth(http.HandlerFunc(uploadHandler.HandleDownload)))

	// Serve static files from public/ folder (build output of frontend)
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request is for an API or appears to be a static asset with an extension,
		// let the FileServer handle it.
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.Contains(r.URL.Path, ".") {
			fs.ServeHTTP(w, r)
			return
		}
		// Otherwise, fallback to index.html to support SPA client-side routing.
		http.ServeFile(w, r, "public/index.html")
	}))

	return courseHandler.LoggingMiddleware(mux)
}
