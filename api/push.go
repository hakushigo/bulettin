package handler

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func Push(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// retrieve data and store it in a variables (so easy)
	_ = r.ParseForm()
	title, content := r.FormValue("title"), r.FormValue("content")
	imgthumbdata, imgthumbhead, err := r.FormFile("thumbnail")

	// UPLOAD THE IMAGE TO THE R2
	// Get S3 credentials
	s3uri := os.Getenv("S3_CONNECT_URI")
	s3idkey := os.Getenv("S3_ID")
	s3secret := os.Getenv("S3_SECRET")

	// Do the job
	// connect
	s3client, err := minio.New(s3uri, &minio.Options{
		Creds:  credentials.NewStaticV4(s3idkey, s3secret, ""),
		Secure: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	// upload the file
	// special : generate random number
	prefix := strconv.Itoa(rand.Int())
	_, err = s3client.PutObject(ctx, "bulettin", prefix+imgthumbhead.Filename, imgthumbdata, imgthumbhead.Size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})

	if err != nil {
		log.Fatal(err)
	}

	// send it to MONGODB
	// retrive it's configuration :<
	mongouri := os.Getenv("MONGODB_CONNECT_URI")
	database := "bulletin"
	collection := "posts"

	// and the job
	// Connect
	mclient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))

	if err != nil {
		log.Fatal(err)
	}

	type doctemplate struct {
		title        string
		descriptive  string
		thumbnailuri string
	}

	_, err = mclient.Database(database).Collection(collection).InsertOne(ctx, doctemplate{
		title:        title,
		descriptive:  content,
		thumbnailuri: "https://bulletin.pool.owo.my.id/" + prefix + imgthumbhead.Filename,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer func(mclient *mongo.Client, ctx context.Context) {
		err := mclient.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}(mclient, ctx)

	w.WriteHeader(200)
	_, _ = w.Write([]byte("OK"))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
