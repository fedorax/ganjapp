package utilities

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/encrypt"
)

// S3Client returns a Minio S3 client
func S3Client() (*minio.Client, error) {
	// Use SSL to communicate with S3 endpoint?
	s3Https := true
	if GetEnv("GANJAPP_S3_HTTPS", "false") == "false" {
		s3Https = false
	}

	return minio.New(GetEnv("GANJAPP_S3_ENDPOINT", ""), &minio.Options{
		Creds:  credentials.NewStaticV4(GetEnv("GANJAPP_S3_ACCESS_KEY", ""), GetEnv("GANJAPP_S3_SECRET_KEY", ""), ""),
		Secure: s3Https,
	})
}

// S3CreateBucket attempts to create an S3 bucket
func S3CreateBucket(name string) (bool, error) {
	s3Client, err := S3Client()
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	err = s3Client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{Region: "us-east-1"})

	if err != nil {
		return false, err
	}

	return true, err
}

// S3BucketExists tests if the passed bucket name exists on the server
func S3BucketExists(name string) (bool, error) {
	s3Client, err := S3Client()
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	found, err := s3Client.BucketExists(context.Background(), name)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return found, err
}

// S3CreateBucketIfNotExists attempts to create an S3 bucket if it doesn't already exist
func S3CreateBucketIfNotExists(name string) (bool, error) {
	if exists, _ := S3BucketExists(name); !exists {
		return S3CreateBucket(name)
	}
	return true, nil
}

// S3UploadFile uploads the supplied file to an S3 bucket
func S3UploadFile(filePath string, remoteName string) (bool, error) {

	s3Client, err := S3Client()
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	// Open the local file that we will upload...
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}
	defer file.Close()

	// Get file stats...
	fstat, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	bucketName := GetEnv("GANJAPP_S3_BUCKET", "ganjapp")
	password := GetEnv("GANJAPP_S3_PASSPHRASE", "[!!error!!]")

	// Ensure the bucket exists...
	success, err := S3CreateBucketIfNotExists(bucketName)
	if err != nil || !success {
		return false, err
	}

	if password == "[!!error!!]" || password == "setme" {
		panic("S3 passphrase not set!")
	}

	// New SSE-C where the cryptographic key is derived from a password and the remoteName + bucketName as salt
	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketName+remoteName))

	// Encrypt file content and upload to the server
	n, err := s3Client.PutObject(context.Background(), bucketName, remoteName, file, fstat.Size(), minio.PutObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	log.Println("Uploaded ", remoteName, " (", filePath, ", ", n, ")")

	return true, nil
}

// S3GetFile attempts to fetch a file from an S3 bucket
func S3GetFile(file string) (reader *minio.Object, err error) {

	s3Client, err := S3Client()
	if err != nil {
		log.Fatalln(err)
		return
	}

	bucketName := GetEnv("GANJAPP_S3_BUCKET", "ganjapp")
	password := GetEnv("GANJAPP_S3_PASSPHRASE", "[!!error!!]")

	if password == "[!!error!!]" || password == "setme" {
		panic("S3 passphrase not set!")
	}

	// New SSE-C where the cryptographic key is derived from a password and the remoteName + bucketName as salt
	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketName+file))

	// Get the encrypted object
	r, e := s3Client.GetObject(context.Background(), bucketName, file, minio.GetObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		log.Fatalln(err)
	}

	defer reader.Close()

	return r, e

}

// S3DeleteFile removes an object from S3 storage
func S3DeleteFile(file string) (bool, error) {
	s3Client, err := S3Client()

	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	err = s3Client.RemoveObject(context.Background(), GetEnv("GANJAPP_S3_BUCKET", "ganjapp"), file, minio.RemoveObjectOptions{GovernanceBypass: true})

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, err

}
