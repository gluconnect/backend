package endpoints

import (
	"encoding/binary"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/stanleymw/glucose/models"
	"github.com/stanleymw/glucose/password"
)

func Auth(db *bolt.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := db.View(func(tx *bolt.Tx) error {
			users := tx.Bucket([]byte("users"))
			if users == nil {
				return nil
			}

			var req models.LoginRequest
			err := ctx.ShouldBindBodyWithJSON(&req)

			log.Printf("New login request: %s", req)
			if err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return err
			}

			real := users.Bucket([]byte(req.Email))
			if real == nil {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return nil
			}

			realUser := models.DecodeStruct[models.User](real.Get([]byte("data")))
			if realUser.Password != password.Hash([]byte(req.Password)) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return nil
			}

			ctx.Set("requestEmail", realUser.Email)
			return nil
		})
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func Register(db *bolt.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := db.Update(func(tx *bolt.Tx) error {
			users, err := tx.CreateBucketIfNotExists([]byte("users"))
			if err != nil {
				return err
			}

			var req models.LoginRequest
			err = ctx.ShouldBindBodyWithJSON(&req)
			if err != nil {
				return err
			}

			userBuck, err := users.CreateBucket([]byte(req.Email))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Account Already Exists!")
				return err
			}

			hashed := password.Hash([]byte(req.Password))
			userBuck.Put([]byte("data"), models.EncodeStruct[models.User](models.User{Email: req.Email, Password: hashed}))
			ctx.JSON(http.StatusOK, map[string]string{})
			return nil
		}); err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func AddReading(db *bolt.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := db.Update(func(tx *bolt.Tx) error {
			users, err := tx.CreateBucketIfNotExists([]byte("users"))
			if err != nil {
				return err
			}

			myEmail := ctx.GetString("requestEmail")

			userBucket := users.Bucket([]byte(myEmail))
			if userBucket == nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return nil
			}

			readings, err := userBucket.CreateBucketIfNotExists([]byte("readings"))
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return err
			}

			var reading models.GlucoseReading
			err = ctx.ShouldBindBodyWithJSON(&reading)
			if err != nil {
				log.Println(err.Error())
				ctx.AbortWithStatus(http.StatusBadRequest)
				return err
			}

			timeStampBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(timeStampBytes, uint64(reading.Timestamp.UnixMilli()))

			readings.Put(timeStampBytes, models.EncodeStruct[models.GlucoseReading](reading))
			return nil
		}); err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func GetReadings(db *bolt.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := db.Update(func(tx *bolt.Tx) error {
			users, err := tx.CreateBucketIfNotExists([]byte("users"))
			if err != nil {
				return err
			}

			myEmail := ctx.GetString("requestEmail")

			userBucket := users.Bucket([]byte(myEmail))
			if userBucket == nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return nil
			}

			readings, err := userBucket.CreateBucketIfNotExists([]byte("readings"))
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return err
			}

			ret := make(map[uint64]models.GlucoseReading)

			readings.ForEach(func(k, v []byte) error {
				reading := models.DecodeStruct[models.GlucoseReading](v)
				ret[binary.LittleEndian.Uint64(k)] = reading
				return nil
			})

			log.Println(ret)
			return nil
		}); err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}
